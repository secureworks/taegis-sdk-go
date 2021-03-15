package zerolog

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

	"github.com/secureworks/taegis-sdk-go/log"

	"github.com/VerticalOps/fakesentry"
	"github.com/getsentry/sentry-go"
	"github.com/rs/zerolog"
	. "github.com/smartystreets/goconvey/convey"
)

func TestSentryWriter(t *testing.T) {
	Convey("SentryWriterRoot", t, func(c C) {
		testSentryWriter(t, c)
	})
}

func testSentryWriter(t *testing.T, c C) {
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

	testWriterLvl := func() {
		w := newSentryWriter(log.WARN)

		//doesn't need to be totally valid json
		//place level at the beginning to fall back for bytes.Index usage
		dat := []byte(`"level":"warn","foo":"bar","hello":"dfdfdf bad json`)

		lvl, shouldLog := w.checkLvl(dat)
		So(lvl, ShouldEqual, zerolog.WarnLevel)
		So(shouldLog, ShouldBeTrue)
	}

	testSentryWriter := func() {
		logger, err := log.Open("zerolog", conf)
		So(err, ShouldBeNil)

		//wouldn't normally do it this way, good enough for test though
		copts := sentry.CurrentHub().Client().Options()
		copts.HTTPTransport = sentrySrv.Transport()
		sc, err := sentry.NewClient(copts)
		So(err, ShouldBeNil)
		sentry.CurrentHub().BindClient(sc)

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

	Convey("WriterLevelCheck", testWriterLvl)
	Convey("SentryWriterData", testSentryWriter)
}
