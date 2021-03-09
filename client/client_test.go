package client

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/http/httptrace"
	"os"
	"sync/atomic"
	"testing"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/secureworks/taegis-sdk-go/log"
	_ "github.com/secureworks/taegis-sdk-go/log/logrus"

	"github.com/stretchr/testify/assert"
)

var _ logrus.Hook = make(logHook, 1)

type logHook chan *logrus.Entry

func (l logHook) Levels() []logrus.Level {
	return logrus.AllLevels
}
func (l logHook) Fire(e *logrus.Entry) error {
	select {
	case l <- e:
	default:
		return errors.New("failed sending message from hook")
	}

	return nil
}

func loggerWithHook(t *testing.T) (log.Logger, logHook) {
	logger, err := log.Open("logrus", &log.Config{Level: log.DEBUG})
	assert.NoError(t, err)
	hook := make(logHook, 1)
	logger.(log.UnderlyingLogger).GetLogger().(*logrus.Logger).Hooks.Add(hook)
	return logger, hook
}

func TestLogger(t *testing.T) {
	logger, hook := loggerWithHook(t)
	c := NewClient(WithLogger(logger))
	c.Logger.Info().Msg("test")
	e := <-hook
	assert.Equal(t, "test", e.Message)
}

func TestCurlDebug(t *testing.T) {
	logger, hook := loggerWithHook(t)
	c := NewClient(WithLogger(logger))

	c.Logger = logger
	fakeHandler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(``))
	}
	server := httptest.NewServer(http.HandlerFunc(fakeHandler))
	defer server.Close()

	req, _ := http.NewRequest(http.MethodPost, server.URL, nil)
	os.Setenv("CURL_DEBUG", "true")
	defer os.Unsetenv("CURL_DEBUG")
	resp, err := c.Do(req)
	assert.Nil(t, err)
	defer resp.Body.Close()
	e := <-hook

	assert.Equal(t, logrus.Fields{"command": "curl -X 'POST' -H 'Content-Type: application/json' '" + server.URL + "'"}, e.Data)
}

func TestBearer(t *testing.T) {
	c := NewClient(WithBearerToken("test"))
	fakeHandler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, r.Header.Get("Authorization"), "Bearer test")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(``))
	}
	server := httptest.NewServer(http.HandlerFunc(fakeHandler))
	defer server.Close()

	req, _ := http.NewRequest(http.MethodPost, server.URL, nil)
	resp, err := c.Do(req)
	assert.Nil(t, err)
	assert.NotNil(t, resp)
	resp.Body.Close()
}

func TestTimeouts(t *testing.T) {
	c := NewClient(WithHTTPTimeout(10 * time.Millisecond))
	fakeHandler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		w.WriteHeader(http.StatusOK)
		time.Sleep(200 * time.Millisecond)
		_, _ = w.Write([]byte(``))
	}
	server := httptest.NewServer(http.HandlerFunc(fakeHandler))
	defer server.Close()

	req, _ := http.NewRequest(http.MethodPost, server.URL, nil)
	_, err := c.Do(req) //nolint: bodyclose not testing the body just the error response
	assert.Contains(t, err.Error(), "Timeout")
}

type transport struct {
	getCount, gotCount, rdTripCount uint64

	logf func(string, ...interface{})
	rt   http.RoundTripper
}

func (tr *transport) GetConn(hostPort string) {
	inc(&tr.getCount)
	tr.logf("GetConn: %s", hostPort)
}

func (tr *transport) GotConn(gci httptrace.GotConnInfo) {
	inc(&tr.gotCount)
	tr.logf("GotConn: Reused (%v), WasIdle (%v), IdleTime (%v)", gci.Reused, gci.WasIdle, gci.IdleTime)
}

func (tr *transport) RoundTrip(req *http.Request) (*http.Response, error) {
	inc(&tr.rdTripCount)
	return tr.rt.RoundTrip(req)
}

func (tr *transport) counts() (uint64, uint64, uint64) {
	return load(&tr.getCount), load(&tr.gotCount), load(&tr.rdTripCount)
}

func inc(i *uint64) uint64 {
	return atomic.AddUint64(i, 1)
}

func load(i *uint64) uint64 {
	return atomic.LoadUint64(i)
}

func TestCustomHTTPClient(t *testing.T) {
	fakeHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.Copy(w, r.Body)
	})

	srv := httptest.NewServer(fakeHandler)
	defer srv.Close()

	cli := srv.Client()
	myt := &transport{
		rt:   cli.Transport,
		logf: func(_ string, _ ...interface{}) {},
	}

	cli.Transport = myt
	trace := &httptrace.ClientTrace{
		GetConn: myt.GetConn,
		GotConn: myt.GotConn,
	}

	if testing.Verbose() {
		myt.logf = t.Logf
	}

	c := NewClient(WithHTTPClient(cli))

	const contents = "Hello World!\n"
	buf := new(bytes.Buffer)
	fmt.Fprint(buf, contents)

	req, err := http.NewRequestWithContext(httptrace.WithClientTrace(context.Background(), trace),
		http.MethodPost, srv.URL, buf)
	assert.Nil(t, err)

	resp, err := c.Do(req)
	assert.Nil(t, err)
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	assert.Nil(t, err)
	assert.Equal(t, contents, string(body))

	get, got, rd := myt.counts()
	one := uint64(1)
	assert.Equal(t, get, one)
	assert.Equal(t, got, one)
	assert.Equal(t, rd, one)
}
