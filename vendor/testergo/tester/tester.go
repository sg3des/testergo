package tester

import (
	"errors"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
)

type State int

const (
	StateTesting State = 0
	StateDone    State = 1
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

func (t *Tester) Start() (chan State, error) {
	t.RunTests()

	webchan := make(chan State)
	fschan := make(chan State)

	go t.events(fschan)
	go t.channels(webchan, fschan)

	return webchan, t.watcher.Add(t.Dir)
}

func (t *Tester) SetWD(dir string) error {
	if err := t.watcher.Remove(t.Dir); err != nil {
		return err
	}

	t.Dir = dir
	t.RunTests()

	return t.watcher.Add(dir)
}

func (t *Tester) Close() error {
	t.running = false
	return t.watcher.Close()
}

func (t *Tester) events(fschan chan State) {
	t.running = true

	for t.running {
		select {
		case <-t.watcher.Events:
			fschan <- StateDone

		case err := <-t.watcher.Errors:
			log.Println("error:", err)
		}
	}
}

func (t *Tester) channels(webchan, fschan chan State) {
	dur := 200 * time.Millisecond
	timer := time.NewTimer(dur)

	for t.running {
		select {
		case <-fschan:
			webchan <- StateTesting
			timer.Reset(dur)

		case <-timer.C:
			t.RunTests()
			webchan <- StateDone
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
