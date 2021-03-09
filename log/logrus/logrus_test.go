package logrus

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/VerticalOps/fakesentry"
	"github.com/getsentry/sentry-go"
	. "github.com/smartystreets/goconvey/convey"

	"github.com/secureworks/taegis-sdk-go/log"
)

func TestLogrus(t *testing.T) {
	Convey("LogrusRoot", t, func(c C) {
		testLogrusRoot(t, c)
	})
}

func testLogrusRoot(t *testing.T, c C) {
	sentrySrv := fakesentry.NewUnstartedServer()

	bufc := make(chan []byte, 1)
	opt := fakesentry.AsMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jb, ok := fakesentry.FromRequest(r)
		c.So(ok, ShouldBeTrue)

		if testing.Verbose() {
			buf := new(bytes.Buffer)
			json.Indent(buf, jb, "", "  ")
			fmt.Fprintf(os.Stderr, "Sentry Body: \n%s\n", buf.Bytes())
		}

		select {
		case bufc <- jb:
		default:
		}
	}))

	sentrySrv.Server = &http.Server{Handler: fakesentry.NewHandler(opt)}
	go sentrySrv.Serve(sentrySrv.Listener())
	defer sentrySrv.Close()

	conf := log.DefaultConfig(func(ev string) string {
		if ev == log.SentryDSN.String() {
			return `http://thisis:myfakeauth@localhost/1`
		}
		return os.Getenv(ev)
	})

	conf.EnableErrStack = true
	if !testing.Verbose() {
		conf.Output = ioutil.Discard
	}

	logger, err := log.Open("logrus", conf)
	So(err, ShouldBeNil)

	//wouldn't normally do it this way, good enough for test though
	copts := sentry.CurrentHub().Client().Options()
	copts.HTTPTransport = sentrySrv.Transport()
	sc, err := sentry.NewClient(copts)
	So(err, ShouldBeNil)
	sentry.CurrentHub().BindClient(sc)

	testLoggedData := func() {
		const errValue, metaValue, msgValue = "oh nooo", "data", "testSentryWriter!"

		logger.WithError(errors.New(errValue)).WithStr("meta", metaValue).Msg(msgValue)

		timer := time.NewTimer(time.Second * 2)
		var dat []byte

		select {
		case <-timer.C:
			t.Fatal("Failed to get logger data")
		case dat = <-bufc:
			timer.Stop()
		}

		event := new(sentry.Event)
		err = json.Unmarshal(dat, event)
		So(err, ShouldBeNil)

		So(event.Message, ShouldEqual, msgValue)
		So(event.Extra["meta"], ShouldEqual, metaValue)
		So(len(event.Exception), ShouldEqual, 1)
		So(event.Exception[0].Value, ShouldEqual, errValue)
		So(len(event.Exception[0].Stacktrace.Frames), ShouldBeGreaterThan, 0)
	}

	Convey("LoggedData", testLoggedData)
}

func TestLogrus_newLogger(t *testing.T) {
	t.Run("Nothing is logged for debug when at INFO", func(t *testing.T) {
		msg := "Debug message here"
		buf := make([]byte, 0, 100)
		out := bytes.NewBuffer(buf)
		cfg := log.DefaultConfig(nil)
		cfg.Output = out
		logger, err := newLogger(cfg)
		assert.Nil(t, err)
		// Log something
		logger.Debug().Msg(msg)
		assert.Nil(t, err)
		data, err := ioutil.ReadAll(out)
		assert.Nil(t, err)
		assert.Equal(t, len(data), 0)
	})
	t.Run("Setting debug level", func(t *testing.T) {
		msg := "Debug message here"
		buf := make([]byte, 0, 100)
		out := bytes.NewBuffer(buf)
		cfg := log.DefaultConfig(nil)
		cfg.Level = log.DEBUG
		cfg.Output = out
		logger, err := newLogger(cfg)
		assert.Nil(t, err)
		// Log something
		logger.Debug().Msg(msg)
		assert.Nil(t, err)
		data, err := ioutil.ReadAll(out)
		assert.Nil(t, err)
		assert.Contains(t, string(data), msg)
	})
	t.Run("Configuration with nil output", func(t *testing.T) {
		cfg := log.DefaultConfig(nil)
		cfg.Output = nil
		logger, err := newLogger(cfg)
		assert.Nil(t, err)
		assert.NotNil(t, logger)
	})
}
