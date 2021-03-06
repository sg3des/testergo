package templates

import (
	"bytes"
	"fmt"
	"io"
	"testergo/tester"
)

func Index(w io.Writer, addr string) {
	w.Write([]byte(`<!DOCTYPE html>
<html>
<head>
	<meta charset="utf-8">
	<title>Testergo</title>
	<link rel="stylesheet"    href="/assets/main.css">
	<link id='favicon' rel="shortcut icon" href="/assets/pass.png" type="image/png">
	<script src="/assets/main.js"></script>
	<script src="/assets/rattle.js"></script>
	<script type="application/javascript" defer>`))
	fmt.Fprintf(w, `		var r = new rattle.NewConnection("ws://%v/ws", true);`, addr)
	w.Write([]byte(`	</script>
</head>
<body>
	<header>
		<div id='status'></div>
		<div id='loading'></div>
	</header>
	<main></main>
</body>
</html>`))
}

func Tests(t *tester.Tester) []byte {
	w := new(bytes.Buffer)
	w.Write([]byte(`<table id='toolbar'><tr>
	<td>
		<b id='reload' class='btn' onclick='r.send("reload")'>re</b>
	<td width='100%'>`))
	fmt.Fprintf(w, `		<input id='wd' class='btn' onchange='r.send("changewd", this.value)' value='%v'>`, t.Dir)
	w.Write([]byte(`</table>`))
	for _, r := range t.Response {
		fmt.Fprintf(w, `	<div class='func %v'>`, r.Status)
		fmt.Fprintf(w, `		<div class='name'>%v</div>`, r.Func)
		w.Write([]byte(`		<div class='log'>`))
		for _, s := range r.Log {
			fmt.Fprintf(w, `				<p>%v</p>`, s)
		}
		w.Write([]byte(`		</div>
	</div>`))
	}
	return w.Bytes()
}
