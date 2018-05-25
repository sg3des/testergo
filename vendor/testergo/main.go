package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/sg3des/argum"
	"github.com/sg3des/rattle"

	"testergo/assets"
	"testergo/templates"
	"testergo/tester"
)

var args struct {
	Dir     string `argum:"pos" help:"path to working directory"`
	Address string `argum:"--address" help:"listening address" default:":8000"`
}

var t *Testergo

func init() {
	argum.MustParse(&args)
	log.SetFlags(log.Lshortfile)
}

func main() {
	err := initTestergo(args.Dir, args.Address)
	if err != nil {
		log.Fatal(err)
	}
}

func initTestergo(dir, addr string) (err error) {
	t = &Testergo{
		wsaddr: addrToWS(addr),
	}

	t.tester, err = tester.NewTester(args.Dir)
	if err != nil {
		return
	}

	t.event, err = t.tester.Start()
	if err != nil {
		return
	}

	http.Handle("/ws", rattle.SetControllers(t))
	rattle.SetOnConnect(t.onConnect)

	http.HandleFunc("/", t.Index)
	http.HandleFunc("/assets/", t.Assets)

	fmt.Println("listen:", addr)
	return http.ListenAndServe(addr, nil)
}

func addrToWS(addr string) string {
	if len(addr) == 0 {
		return addr
	}

	if addr[0] == ':' {
		return "127.0.0.1" + addr
	}

	return addr
}

type Testergo struct {
	wsaddr string

	tester *tester.Tester
	event  chan bool
}

//
// HTTP handlers
//

func (t *Testergo) Index(w http.ResponseWriter, r *http.Request) {
	templates.Index(w, t.wsaddr)
}

//Assets serve static files stored to go code how bind-data
func (t *Testergo) Assets(w http.ResponseWriter, r *http.Request) {
	assetname := r.URL.Path[1:]

	//lookup assets file
	fi, err := assets.AssetInfo(assetname)
	if err != nil {
		http.Error(w, fmt.Sprintf("File %s not found", r.URL.Path), 404)
		return
	}

	//check modified date
	modSince, err := time.Parse(time.RFC1123, r.Header.Get("If-Modified-Since"))
	if err == nil && fi.ModTime().Before(modSince) {
		w.WriteHeader(http.StatusNotModified)
		return
	}

	//restore assets from bind data
	data, err := assets.Asset(assetname)
	if err != nil {
		http.Error(w, fmt.Sprintf("File %s not found", r.URL.Path), http.StatusNotFound)
		return
	}

	//server file
	http.ServeContent(w, r, r.URL.Path, time.Now(), bytes.NewReader(data))
}

//
// WEBSOCKET
//

func (*Testergo) State(r *rattle.Conn) {
	r.NewMessage("favicon", []byte(t.tester.Status)).Send()
	r.NewMessage("=header", templates.Status(t.tester.Status)).Send()
	r.NewMessage("=main", templates.Tests(t.tester)).Send()
}

func (*Testergo) onConnect(r *rattle.Conn) {
	t.State(r)

	for {
		<-t.event
		t.State(r)
	}
}
