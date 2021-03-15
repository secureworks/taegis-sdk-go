package zerolog

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/rs/zerolog"

	"github.com/secureworks/taegis-sdk-go/log"
	"github.com/secureworks/taegis-sdk-go/log/internal/common"
)

func init() {
	/*
		These really should be in newLogger below, but they're package level vars in zerolog (bad)
		which means no matter where we put them they have a chance to race, they'll just race less here.
		We also can't do much about the fact that each instance of zerolog.Logger will use these even if
		we don't want
	*/
	zerolog.ErrorStackFieldName = log.StackField
	zerolog.ErrorStackMarshaler = func(err error) interface{} {
		st, _ := common.WithStackTrace(err)
		return st.StackTrace()
	}

	log.Register("zerolog", newLogger)
}

func newLogger(conf *log.Config, opts ...log.Option) (log.Logger, error) {
	zlvl := lvlToZerolog(conf.Level)
	logger := &logger{
		errStack: conf.EnableErrStack,
		lvl:      zlvl,
	}

	output := conf.Output
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

		output = io.MultiWriter(output, newSentryWriter(conf.Sentry.Levels...))
	}

	zlog := zerolog.New(output).Level(zlvl)
	logger.lg = &zlog

	for _, opt := range opts {
		if err := opt(logger); err != nil {
			return nil, err
		}
	}

	return logger, nil
}

type logger struct {
	lg       *zerolog.Logger
	lvl      zerolog.Level
	errStack bool
}

func (l *logger) notValid() bool {
	return l == nil || l.lg == nil
}

func (l *logger) WithError(err error) log.Entry {
	return l.Error().WithError(err)
}

func (l *logger) WithField(key string, val interface{}) log.Entry {
	return l.Entry(0).WithField(key, val)
}

func (l *logger) WithFields(fields map[string]interface{}) log.Entry {
	return l.Entry(0).WithFields(fields)
}

func (l *logger) newEntry(lvl zerolog.Level) log.Entry {
	if l.notValid() {
		return l.DisabledEntry()
	}

	//we have to use NoLevel or we can't change them after the fact
	//https://github.com/rs/zerolog/blob/7825d863376faee2723fc99c061c538bd80812c8/log.go#L419
	//https://github.com/rs/zerolog/pull/255
	//Our own *entry type will write the level as needed
	//TODO: Only using NoLevel silently breaks zerolog.Hook interface
	ent := l.lg.WithLevel(zerolog.NoLevel)

	if l.errStack {
		ent = ent.Stack()
	}

	return &entry{
		ent:    ent,
		caller: make([]string, 0, 1),
		loglvl: l.lvl,
		lvl:    lvl,
	}
}

func lvlToZerolog(lvl log.Level) zerolog.Level {
	switch lvl {
	case log.TRACE:
		return zerolog.TraceLevel
	case log.DEBUG:
		return zerolog.DebugLevel
	case log.INFO:
		return zerolog.InfoLevel
	case log.WARN:
		return zerolog.WarnLevel
	case log.ERROR:
		return zerolog.ErrorLevel
	case log.PANIC:
		return zerolog.PanicLevel
	case log.FATAL:
		return zerolog.FatalLevel
	default:
		return zerolog.InfoLevel
	}
}

func (l *logger) Entry(lvl log.Level) log.Entry { return l.newEntry(lvlToZerolog(lvl)) }

func (l *logger) Trace() log.Entry { return l.newEntry(zerolog.TraceLevel) }
func (l *logger) Debug() log.Entry { return l.newEntry(zerolog.DebugLevel) }
func (l *logger) Info() log.Entry  { return l.newEntry(zerolog.InfoLevel) }
func (l *logger) Warn() log.Entry  { return l.newEntry(zerolog.WarnLevel) }
func (l *logger) Error() log.Entry { return l.newEntry(zerolog.ErrorLevel) }
func (l *logger) Panic() log.Entry { return l.newEntry(zerolog.PanicLevel) }
func (l *logger) Fatal() log.Entry { return l.newEntry(zerolog.FatalLevel) }

//TODO: could refactor this if useful anywhere else
type writeLevelCloser struct {
	log log.Logger
	lvl log.Level
}

//An implementation of zerolog.Logger.Write method but with level support
func (wlc writeLevelCloser) Write(p []byte) (n int, err error) {
	n = len(p)
	if n > 0 && p[n-1] == '\n' {
		// Trim CR added by stdlog.
		p = p[0 : n-1]
	}
	wlc.log.Entry(wlc.lvl).Msg(string(p))
	return
}

func (wlc writeLevelCloser) Close() error {
	return nil
}

func (l *logger) WriteCloser(lvl log.Level) io.WriteCloser {
	return writeLevelCloser{log: l, lvl: lvl}
}

//An assertable method/interface if someone wants to disable zerolog events at runtime
func (l *logger) DisabledEntry() log.Entry {
	return (*entry)(nil)
}

func (l *logger) GetLogger() interface{} {
	if l.notValid() {
		return nil
	}

	return l.lg
}

func (l *logger) SetLogger(iface interface{}) {
	if lg, ok := iface.(*zerolog.Logger); ok && !l.notValid() {
		l.lg = lg
	}
}

type entry struct {
	ent    *zerolog.Event
	caller []string
	msg    string
	async  bool
	loglvl zerolog.Level
	lvl    zerolog.Level
}

func (e *entry) notValid() bool {
	return e == nil || e.ent == nil
}

func (e *entry) enabled() bool {
	return !e.notValid() && e.lvl >= e.loglvl
}

func (e *entry) Async() log.Entry {
	if e.notValid() {
		return e
	}

	e.async = !e.async
	return e
}

func (e *entry) Caller(skip ...int) log.Entry {
	if e.notValid() {
		return e
	}

	sk := 1
	if len(skip) > 0 {
		sk += skip[0]
	}

	//originally it was planned to use zerologs Caller implementation
	//but the interface was changed during the design phase to allow
	//Caller to be called multiple times which zerolog won't do without
	//adding dup fields
	_, file, line, ok := runtime.Caller(sk)
	if !ok {
		return e
	}

	//TODO: Use zerolog.CallerMarshalFunc ?
	e.caller = append(e.caller, fmt.Sprintf("%s:%d", file, line))

	return e
}

func (e *entry) WithError(errs ...error) log.Entry {
	le := len(errs)
	if e.notValid() || le == 0 {
		return e
	}

	if le == 1 {
		e.ent = e.ent.Err(errs[0])
	} else {
		e.ent = e.ent.Errs(zerolog.ErrorFieldName, errs)
	}

	return e
}

func (e *entry) WithField(key string, val interface{}) log.Entry {
	if e.notValid() {
		return e
	}

	e.ent = e.ent.Interface(key, val)

	return e
}

func (e *entry) WithFields(fields map[string]interface{}) log.Entry {
	if e.notValid() || len(fields) == 0 {
		return e
	}

	e.ent = e.ent.Fields(fields)

	return e
}

func (e *entry) WithBool(key string, bls ...bool) log.Entry {
	lb := len(bls)
	if e.notValid() || lb == 0 {
		return e
	}

	if lb == 1 {
		e.ent = e.ent.Bool(key, bls[0])
	} else {
		e.ent = e.ent.Bools(key, bls)
	}

	return e
}

func (e *entry) WithDur(key string, durs ...time.Duration) log.Entry {
	ld := len(durs)
	if e.notValid() || ld == 0 {
		return e
	}

	if ld == 1 {
		e.ent = e.ent.Dur(key, durs[0])
	} else {
		e.ent = e.ent.Durs(key, durs)
	}

	return e
}

func (e *entry) WithInt(key string, is ...int) log.Entry {
	li := len(is)
	if e.notValid() || li == 0 {
		return e
	}

	if li == 1 {
		e.ent = e.ent.Int(key, is[0])
	} else {
		e.ent = e.ent.Ints(key, is)
	}

	return e
}

func (e *entry) WithUint(key string, us ...uint) log.Entry {
	lu := len(us)
	if e.notValid() || lu == 0 {
		return e
	}

	if lu == 1 {
		e.ent = e.ent.Uint(key, us[0])
	} else {
		e.ent = e.ent.Uints(key, us)
	}

	return e
}

func (e *entry) WithStr(key string, strs ...string) log.Entry {
	ls := len(strs)
	if e.notValid() || ls == 0 {
		return e
	}

	if ls == 1 {
		e.ent = e.ent.Str(key, strs[0])
	} else {
		e.ent = e.ent.Strs(key, strs)
	}

	return e
}

func (e *entry) WithTime(key string, ts ...time.Time) log.Entry {
	lt := len(ts)
	if e.notValid() || lt == 0 {
		return e
	}

	if lt == 1 {
		e.ent = e.ent.Time(key, ts[0])
	} else {
		e.ent = e.ent.Times(key, ts)
	}

	return e
}

func (e *entry) setLevel(lvl zerolog.Level) log.Entry {
	if e.notValid() {
		return e
	}

	e.lvl = lvl
	return e
}

func (e *entry) Trace() log.Entry { return e.setLevel(zerolog.TraceLevel) }
func (e *entry) Debug() log.Entry { return e.setLevel(zerolog.DebugLevel) }
func (e *entry) Info() log.Entry  { return e.setLevel(zerolog.InfoLevel) }
func (e *entry) Warn() log.Entry  { return e.setLevel(zerolog.WarnLevel) }
func (e *entry) Error() log.Entry { return e.setLevel(zerolog.ErrorLevel) }
func (e *entry) Panic() log.Entry { return e.setLevel(zerolog.PanicLevel) }
func (e *entry) Fatal() log.Entry { return e.setLevel(zerolog.FatalLevel) }

func (e *entry) Msgf(format string, vals ...interface{}) {
	e.Msg(fmt.Sprintf(format, vals...))
}

func (e *entry) Msg(msg string) {
	if e.notValid() {
		return
	}

	e.msg = msg

	if !e.async {
		e.Send()
	}
}

func (e *entry) Send() {
	if e != nil && e.ent != nil {
		defer func() {
			putEvent(e.ent)
			e.ent = nil
		}()
	}

	if !e.enabled() {
		return
	}

	if len(e.caller) > 0 {
		e.ent = e.ent.Strs(log.CallerField, e.caller)
	}

	e.ent = e.ent.Str(zerolog.LevelFieldName, zerolog.LevelFieldMarshalFunc(e.lvl))
	e.ent.Msg(e.msg)

	switch e.lvl {
	case zerolog.PanicLevel:
		panic(e.msg)
	case zerolog.FatalLevel:
		os.Exit(1)
	}
}

//An assertable method/interface if someone wants to disable zerolog events at runtime
func (e *entry) DisabledEntry() log.Entry {
	if e.notValid() {
		return e
	}

	//this will disable all other methods
	if e.ent != nil {
		putEvent(e.ent)
		e.ent = nil
	}

	return e
}

func (e *entry) GetLogger() interface{} {
	if e.notValid() {
		return nil
	}

	return e.ent
}

func (e *entry) SetLogger(l interface{}) {
	if ent, ok := l.(*zerolog.Event); ok && !e.notValid() {
		e.ent = ent
	}
}
