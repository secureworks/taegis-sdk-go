package middleware

import (
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/secureworks/tdr-sdk-go/log"
	"github.com/pkg/errors"
)

func TestHTTPBaseContext(t *testing.T) {
	srv := httptest.NewUnstartedServer(nil)
	noop := log.Noop()

	var c io.Closer
	srv.Config, c = NewHTTPServer(noop, 0)
	defer c.Close()

	//seems like a pointless test... but coverage!
	srv.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if log.LoggerFromCtx(r.Context()) == nil {
			t.Fatal("Nil logger in request scoped context")
		}
	})
	srv.Start()

	resp, err := srv.Client().Get(srv.URL)
	if err != nil {
		t.Fatalf("Failed to make http request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Non-200 status from http.Server: %d", resp.StatusCode)
	}
}

func TestHTTPRequestMiddleware(t *testing.T) {
	ml := mLog{log.Noop()}
	mid := NewHTTPRequestMiddleware(ml, 0)

	meta := "meta"
	data := "data"
	msg := "hello world"
	method := http.MethodGet
	path := "/foobar"
	code := http.StatusCreated
	var entry log.Entry //hacky but whatever

	handler := mid(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		entry = log.EntryFromCtx(r.Context())
		entry.WithStr(meta, data).Msg(msg)
		w.WriteHeader(code)
	}))

	req := httptest.NewRequest(method, path, nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != code {
		t.Fatalf("Wrong code from http handler: %d", code)
	}

	ent := entry.(*mEntry)
	if !ent.async || !ent.sent || ent.lvl != log.Level(0) || ent.msg != msg {
		t.Fatal("Entry fields incorrect")
	}

	if _, ok := ent.vals[log.ReqDuration]; !ok {
		t.Fatal("Request duration does not exist in log entry")
	}

	delete(ent.vals, log.ReqDuration)
	deepEqual := reflect.DeepEqual(ent.vals, map[string]interface{}{
		meta:               data,
		"http_method":      method,
		"http_path":        path,
		"http_remote_addr": req.RemoteAddr,
	})

	if !deepEqual {
		t.Fatalf("Unequal values in log entry: %v", ent.vals)
	}
}

func TestHTTPRequestMiddlewarePanic(t *testing.T) {
	ml := mLog{log.Noop()}
	mid := NewHTTPRequestMiddleware(ml, 0)

	pv := "this is fine"
	var entry log.Entry
	handler := mid(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		entry = log.EntryFromCtx(r.Context())
		panic(pv)
	}))

	req := httptest.NewRequest(http.MethodGet, "/path", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	ent := entry.(*mEntry)

	if v, ok := ent.vals[log.PanicValue].(string); !ok || v != pv {
		t.Fatalf("log.PanicValue not what was expected: %v", v)
	}

	if st, ok := ent.vals[log.PanicStack].(errors.StackTrace); !ok || len(st) == 0 {
		t.Fatalf("log.PanicStack not type that was expected: %T", ent.vals[log.PanicStack])
	}

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("Unexpected status code from panic/recover: %d", rec.Code)
	}
}

//some mock types with certain methods shadowed

type mLog struct {
	log.Logger
}

func (l mLog) Entry(lvl log.Level) log.Entry {
	return &mEntry{
		Entry: l.Logger.Entry(lvl),
		lvl:   lvl,
		vals:  make(map[string]interface{}),
	}
}

type mEntry struct {
	log.Entry
	async bool
	sent  bool
	lvl   log.Level
	msg   string
	vals  map[string]interface{}
}

func (m *mEntry) WithStr(key string, strs ...string) log.Entry {
	m.vals[key] = strings.Join(strs, "")
	return m
}

func (m *mEntry) WithFields(fields map[string]interface{}) log.Entry {
	for k, v := range fields {
		m.vals[k] = v
	}

	return m
}

func (m *mEntry) Error() log.Entry {
	return m
}

func (m *mEntry) Msg(msg string) {
	m.msg = msg
}

func (m *mEntry) Async() log.Entry {
	m.async = true
	return m
}

func (m *mEntry) Send() {
	m.sent = true
}
