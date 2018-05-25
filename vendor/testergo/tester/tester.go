package tester

import (
	"errors"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/fsnotify/fsnotify"
)

const (
	StatusPass  = "pass"
	StatusFail  = "fail"
	StatusPanic = "panic"
)

type Tester struct {
	Dir     string
	watcher *fsnotify.Watcher

	running  bool
	Status   string
	Response []*Resp
}

type Resp struct {
	Func   string
	Status string
	Log    []string
}

func NewTester(dir string) (*Tester, error) {
	fi, err := os.Stat(dir)
	if err != nil {
		return nil, err
	}

	if !fi.IsDir() {
		return nil, errors.New("this is not directory")
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	return &Tester{Dir: dir, watcher: watcher}, nil
}

func (t *Tester) Start() (chan bool, error) {
	t.RunTests()

	c := make(chan bool)
	go t.events(c)

	return c, t.watcher.Add(t.Dir)
}

func (t *Tester) Close() error {
	t.running = false
	return t.watcher.Close()
}

func (t *Tester) events(c chan bool) {
	t.running = true

	for t.running {
		select {
		case <-t.watcher.Events:
			c <- false
			t.RunTests()
			c <- true

		case err := <-t.watcher.Errors:
			log.Println("error:", err)
		}
	}
}

func (t *Tester) RunTests() {
	out, err := exec.Command("go", "test", "-v", t.Dir).CombinedOutput()
	if err != nil {
		t.Status = StatusFail
	} else {
		t.Status = StatusPass
	}

	t.Response = t.Response[:0]
	chunks := regexp.MustCompile("(?m)^=== RUN").Split(string(out), -1)
	for _, chunk := range chunks {
		chunk = strings.Trim(chunk, " ")
		if len(chunk) == 0 {
			continue
		}

		resp := t.parseResponse(chunk)
		if resp.Status == StatusPanic {
			t.Status = StatusPanic
		}

		t.Response = append(t.Response, resp)
	}
}

func (t *Tester) parseResponse(text string) *Resp {
	var resp = new(Resp)

	for _, line := range strings.Split(text, "\n") {
		if len(line) == 0 {
			continue
		}

		if resp.Func == "" {
			resp.Func = line
			continue
		}

		if strings.HasPrefix(line, "--- PASS") {
			resp.Status = StatusPass
			continue
		}

		if strings.HasPrefix(line, "--- FAIL") {
			resp.Status = StatusFail
			continue
		}

		if strings.HasPrefix(line, "panic:") {
			resp.Status = StatusPanic
		}

		if strings.HasPrefix(line, "PASS") || strings.HasPrefix(line, "FAIL") {
			break
		}

		resp.Log = append(resp.Log, line)
	}

	return resp
}
