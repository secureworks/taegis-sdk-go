package logrus

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/makasim/sentryhook"
	"github.com/sirupsen/logrus"

	"github.com/secureworks/taegis-sdk-go/log"
	"github.com/secureworks/taegis-sdk-go/log/internal/common"
)

func init() {
	log.Register("logrus", newLogger)
}

func newLogger(conf *log.Config, opts ...log.Option) (log.Logger, error) {
	lrus := logrus.New()
	lrus.SetNoLock()
	if conf.Output == nil {
		conf.Output = os.Stderr
	}
	lrus.SetOutput(conf.Output)
	lrus.SetLevel(lvlToLogrus(conf.Level))
	if conf.Format == log.JSONFormat {
		jsonF := &logrus.JSONFormatter{
			PrettyPrint: conf.LocalDevel,
		}

		lrus.SetFormatter(jsonF)
	}

	if conf.EnableErrStack {
		lrus.AddHook(errHook{})
	}

	if conf.Sentry.DSN != "" {
		tp := sentry.NewHTTPSyncTransport()
		tp.Timeout = time.Second * 15

		opts := sentry.ClientOptions{
			Dsn:              conf.Sentry.DSN,
			Release:          conf.Sentry.Release,
			Environment:      conf.Sentry.Env,
			ServerName:       conf.Sentry.Server,
			Debug:            conf.Sentry.Debug,
			AttachStacktrace: conf.EnableErrStack,
			Transport:        tp,
		}

		if err := common.InitSentry(opts); err != nil {
			return nil, err
		}

		lrusLvls := make([]logrus.Level, 0, len(conf.Sentry.Levels))
		for _, lvl := range conf.Sentry.Levels {
			lrusLvls = append(lrusLvls, lvlToLogrus(lvl))
		}

		lrus.AddHook(sentryhook.New(lrusLvls, sentryhook.WithConverter(sentryConverter)))
	}

	logger := &logger{lg: lrus, errStack: conf.EnableErrStack}
	for _, opt := range opts {
		if err := opt(logger); err != nil {
			return nil, err
		}
	}

	return logger, nil
}

type logger struct {
	lg       *logrus.Logger
	errStack bool
}

func (l *logger) newEntry(lvl logrus.Level) *entry {
	return &entry{
		ent:      logrus.NewEntry(l.lg),
		errStack: l.errStack,
		lvl:      lvl,
	}
}

func (l *logger) Entry(lvl log.Level) log.Entry { return l.newEntry(lvlToLogrus(lvl)) }

func (l *logger) Trace() log.Entry { return l.newEntry(logrus.TraceLevel) }
func (l *logger) Debug() log.Entry { return l.newEntry(logrus.DebugLevel) }
func (l *logger) Info() log.Entry  { return l.newEntry(logrus.InfoLevel) }
func (l *logger) Warn() log.Entry  { return l.newEntry(logrus.WarnLevel) }
func (l *logger) Error() log.Entry { return l.newEntry(logrus.ErrorLevel) }
func (l *logger) Panic() log.Entry { return l.newEntry(logrus.PanicLevel) }
func (l *logger) Fatal() log.Entry { return l.newEntry(logrus.FatalLevel) }

func (l *logger) WithError(err error) log.Entry {
	ent := l.Error()
	return ent.WithError(err)
}

func (l *logger) WithField(key string, val interface{}) log.Entry {
	ent := l.Entry(0)
	return ent.WithField(key, val)
}

func (l *logger) WithFields(fields map[string]interface{}) log.Entry {
	ent := l.Entry(0)
	return ent.WithFields(fields)
}

func lvlToLogrus(lvl log.Level) logrus.Level {
	switch lvl {
	case log.TRACE:
		return logrus.TraceLevel
	case log.DEBUG:
		return logrus.DebugLevel
	case log.INFO:
		return logrus.InfoLevel
	case log.WARN:
		return logrus.WarnLevel
	case log.ERROR:
		return logrus.ErrorLevel
	case log.PANIC:
		return logrus.PanicLevel
	case log.FATAL:
		return logrus.FatalLevel
	default:
		return logrus.InfoLevel
	}
}

func (l *logger) WriteCloser(lvl log.Level) io.WriteCloser {
	return l.lg.WriterLevel(lvlToLogrus(lvl))
}

func (l *logger) GetLogger() interface{} {
	return l.lg
}

func (l *logger) SetLogger(iface interface{}) {
	if lg, ok := iface.(*logrus.Logger); ok {
		l.lg = lg
	}
}

type entry struct {
	ent      *logrus.Entry
	lvl      logrus.Level
	async    bool
	errStack bool
	msg      string
}

func (e *entry) Caller(skip ...int) log.Entry {
	if e == nil {
		return e
	}

	sk := 1
	if len(skip) > 0 {
		sk += skip[0]
	}

	_, file, line, ok := runtime.Caller(sk)
	if !ok {
		return e
	}

	//not normal logrus, append to existing field
	//even if nil, won't panic
	cls, _ := e.ent.Data[log.CallerField].([]string)
	cls = append(cls, fmt.Sprintf("%s:%d", file, line))
	e.ent.Data[log.CallerField] = cls

	return e
}

type multiError struct {
	errs []error
}

func (me multiError) Error() string {
	sb := new(strings.Builder)
	sb.Grow(len(me.errs) * 32)

	for _, e := range me.errs {
		fmt.Fprintf(sb, "%v\n", e)
	}

	return sb.String()
}

func (e *entry) WithError(errs ...error) log.Entry {
	if len(errs) == 0 || e == nil {
		return e
	}

	err := errs[0]
	if len(errs) > 1 {
		err = multiError{errs}
	}

	if e.errStack {
		_, err = common.WithStackTrace(err)
	}

	return e.WithField(logrus.ErrorKey, err)
}

func (e *entry) WithField(key string, val interface{}) log.Entry {
	//this defer relies on the fact that defer's args are eval'd when defer is called
	//not when the defer'd function is run
	defer releaseEntry(e.ent.Logger, e.ent)

	e.ent = e.ent.WithField(key, val)
	return e
}

func (e *entry) WithFields(fields map[string]interface{}) log.Entry {
	defer releaseEntry(e.ent.Logger, e.ent)

	e.ent = e.ent.WithFields(fields)
	return e
}

//boy generics would help with these multi-value/variadic methods 😤
//not gonna use reflection cause that would be too slow

func (e *entry) WithBool(key string, bls ...bool) log.Entry {
	if e == nil || len(bls) == 0 {
		return e
	}

	var i interface{} = bls[0]
	if len(bls) > 1 {
		i = bls
	}

	return e.WithField(key, i)
}

func (e *entry) WithDur(key string, durs ...time.Duration) log.Entry {
	if e == nil || len(durs) == 0 {
		return e
	}

	var i interface{} = durs[0]
	if len(durs) > 1 {
		i = durs
	}

	return e.WithField(key, i)
}

func (e *entry) WithInt(key string, is ...int) log.Entry {
	if e == nil || len(is) == 0 {
		return e
	}

	var i interface{} = is[0]
	if len(is) > 1 {
		i = is
	}

	return e.WithField(key, i)
}

func (e *entry) WithUint(key string, us ...uint) log.Entry {
	if e == nil || len(us) == 0 {
		return e
	}

	var i interface{} = us[0]
	if len(us) > 1 {
		i = us
	}

	return e.WithField(key, i)
}

func (e *entry) WithStr(key string, strs ...string) log.Entry {
	if e == nil || len(strs) == 0 {
		return e
	}

	//string allocates when placed into empty interface 🙁
	var i interface{} = strs[0]
	if len(strs) > 1 {
		i = strs
	}

	return e.WithField(key, i)
}

func (e *entry) WithTime(key string, ts ...time.Time) log.Entry {
	if e == nil || len(ts) == 0 {
		return e
	}

	var i interface{} = ts[0]
	if len(ts) > 1 {
		i = ts
	}

	//avoid using WithTime here from logrus as we don't want to unnecessarily override time value
	return e.WithField(key, i)
}

func (e *entry) Msgf(format string, vals ...interface{}) {
	e.Msg(fmt.Sprintf(format, vals...))
}

func (e *entry) Msg(msg string) {
	e.msg = msg

	if !e.async {
		e.Send()
	}
}

func (e *entry) Async() log.Entry {
	e.async = !e.async
	return e
}

func (e *entry) Send() {
	if e == nil || e.ent == nil {
		return
	}

	defer releaseEntry(e.ent.Logger, e.ent)

	switch e.lvl {
	case logrus.PanicLevel:
		e.ent.Panic(e.msg)
	case logrus.FatalLevel:
		e.ent.Fatal(e.msg)
	default:
		e.ent.Log(e.lvl, e.msg)
	}

	e.ent = nil
}

func (e *entry) Trace() log.Entry { e.lvl = logrus.TraceLevel; return e }
func (e *entry) Debug() log.Entry { e.lvl = logrus.DebugLevel; return e }
func (e *entry) Info() log.Entry  { e.lvl = logrus.InfoLevel; return e }
func (e *entry) Warn() log.Entry  { e.lvl = logrus.WarnLevel; return e }
func (e *entry) Error() log.Entry { e.lvl = logrus.ErrorLevel; return e }
func (e *entry) Panic() log.Entry { e.lvl = logrus.PanicLevel; return e }
func (e *entry) Fatal() log.Entry { e.lvl = logrus.FatalLevel; return e }

func (e *entry) GetLogger() interface{} {
	return e.ent
}

func (e *entry) SetLogger(l interface{}) {
	if ent, ok := l.(*logrus.Entry); ok {
		e.ent = ent
	}
}
