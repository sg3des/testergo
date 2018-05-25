package templates

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"testergo/tester"
)

func JS(w io.Writer, addr string) {
	w.Write([]byte(`<script type="application/javascript" defer>`))
	fmt.Fprintf(w, `new rattle.NewConnection("ws://%v/ws", true);`, addr)
	w.Write([]byte(`function favicon(name) {
	document.getElementById("favicon").href = "/assets/"+name+".png";
}
</script>`))
}

func Index(w io.Writer, addr string) {
	w.Write([]byte(`<!DOCTYPE html>
<html>
<head>
	<meta charset="utf-8">
	<title>Testergo</title>
	<link rel="stylesheet"    href="/assets/main.css">
	<link id='favicon' rel="shortcut icon" href="/assets/pass.png" type="image/png">
	<script src="/assets/rattle.js"></script>`))
	JS(w, addr)
	w.Write([]byte(`</head>
<body>
	<header></header>
	<main></main>
</body>
</html>`))
}

func Tests(t *tester.Tester) []byte {
	w := new(bytes.Buffer)
	fmt.Fprintf(w, `	<nav>%v</nav>`, t.Dir)
	for _, r := range t.Response {
		fmt.Fprintf(w, `		<div class='func %v'>`, r.Status)
		fmt.Fprintf(w, `			<div class='name'>%v</div>`, r.Func)
		w.Write([]byte(`			<div class='log'>`))
		for _, s := range r.Log {
			fmt.Fprintf(w, `					<p>%v</p>`, s)
		}
		w.Write([]byte(`			</div>
		</div>`))
	}
	return w.Bytes()
}

func Status(status string) []byte {
	w := new(bytes.Buffer)
	fmt.Fprintf(w, `<div class='status %v`, status)
	fmt.Fprintf(w, `'>%v</div>`, strings.ToTitle(status))
	return w.Bytes()
}
