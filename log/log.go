/*
Package log provides a unified interface for loggers such as logrus or zerolog, along with Sentry support and other custom options.
Look to the design doc here for more information: https://code.8labs.io/project/designs/-/tree/platform-ctpx-sdk-logging/platform/ctpx-sdk-logging
*/
package log

import (
	"context"
	"io"
	"strings"
	"time"
)

//LoggerFormat is the base type for logging formats supported by this package
type LoggerFormat int

const (
	//ImplementationDefaultFormat leaves the format up to the logger implementation default.
	ImplementationDefaultFormat LoggerFormat = iota - 1
	//JSONFormat is the default (zero value) format for loggers registered with this package.
	JSONFormat
)

//IsValid checks if a logger format is valid.
func (l LoggerFormat) IsValid() bool {
	switch l {
	case ImplementationDefaultFormat, JSONFormat:
		return true
	default:
		return false
	}
}

//Level is the base type for logging levels supported by this package.
type Level int

//LevelFromString parses str and returns the closest level.
//If one isn't found the default level is returned.
func LevelFromString(str string) (lvl Level) {
	switch strings.ToUpper(str) {
	case "TRACE":
		lvl = TRACE
	case "DEBUG":
		lvl = DEBUG
	case "INFO":
		lvl = INFO
	case "WARN":
		lvl = WARN
	case "ERROR":
		lvl = ERROR
	case "PANIC":
		lvl = PANIC
	case "FATAL":
		lvl = FATAL
	}

	//default case isn't needed, default is determined by enum zero value
	return
}

//IsValid checks if the current level is valid relative to known values.
func (l Level) IsValid() bool {
	return l >= TRACE && l <= FATAL
}

//IsEnabled checks if the level l is enabled relative to en.
func (l Level) IsEnabled(en Level) bool {
	return l.IsValid() && en.IsValid() && l >= en
}

const (
	//TRACE Level.
	TRACE Level = iota + -2
	//DEBUG Level.
	DEBUG
	//INFO Level; this is the default (zero value).
	INFO
	//WARN Level.
	WARN
	//ERROR Level.
	ERROR
	//PANIC Level; note, depending on usage this will cause the logger to panic.
	PANIC
	//FATAL Level; note, depending on usage this will cause the logger to force a program exit.
	FATAL
)

//AllLevels is a convenience function returning all levels as a slice.
func AllLevels() []Level {
	return []Level{
		TRACE,
		DEBUG,
		INFO,
		WARN,
		ERROR,
		PANIC,
		FATAL,
	}
}

type ctxKey int

const (
	//LoggerKey is the key value to use with context.Context for Logger put and retrieval.
	LoggerKey ctxKey = iota + 1
	//EntryKey is the key value to use with context.Context for Logger put and retrieval.
	EntryKey
)

//CtxWithLogger returns a context with Logger l as its value.
func CtxWithLogger(ctx context.Context, l Logger) context.Context {
	return context.WithValue(ctx, LoggerKey, l)
}

//LoggerFromCtx returns the Logger in ctx, or nil if none exists.
func LoggerFromCtx(ctx context.Context) Logger {
	l, _ := ctx.Value(LoggerKey).(Logger)
	return l
}

//CtxWithEntry returns a context with Entry e as its value.
func CtxWithEntry(ctx context.Context, e Entry) context.Context {
	return context.WithValue(ctx, EntryKey, e)
}

//EntryFromCtx returns the Entry in ctx, or nil if none exists.
func EntryFromCtx(ctx context.Context) Entry {
	e, _ := ctx.Value(EntryKey).(Entry)
	return e
}

//TODO: Other common fields here instead of sdks/go/common ?

const (
	//ReqDuration is a key for Logger data.
	ReqDuration = `request_duration`
	//PanicStack is a key for Logger data.
	PanicStack = `panic_stack`
	//PanicValue is a key for Logger data.
	PanicValue = `panic_value`
	//CallerField is a key for Logger data.
	CallerField = `caller`
	//StackField is a key for Logger data.
	StackField = `stack`
	//XRequestID is a common field and header from API calls.
	XRequestID = `X-Request-Id`
	//XTraceID is a common field and header from API calls.
	XTraceID = `X-Trace-Id`
	//XSpanID is a common field and header from API calls.
	XSpanID = `X-Span-Id`
)

//Logger is the minimum interface loggers should implement when used with CTPx packages.
type Logger interface {
	//WithError inserts the given error into a new Entry and returns the Entry.
	WithError(err error) Entry
	//WithField inserts key & val into a new Entry and returns the Entry.
	WithField(key string, val interface{}) Entry
	//WithFields inserts the given set of fields into a new Entry and returns the Entry.
	WithFields(fields map[string]interface{}) Entry

	//Entry returns a new Entry at the provided log level.
	Entry(Level) Entry
	//Trace returns a new Entry at TRACE level.
	Trace() Entry
	//Debug returns a new Entry at DEBUG level.
	Debug() Entry
	//Info returns a new Entry at INFO level.
	Info() Entry
	//Warn returns a new Entry at WARN level.
	Warn() Entry
	//Error returns a new Entry at ERROR level.
	Error() Entry
	//Panic returns a new Entry at PANIC level.
	//Implementations should panic once the final message for the Entry is logged.
	Panic() Entry
	//Fatal returns a new Entry at FATAL level.
	//Implementations should exit non-zero once the final message for the Entry is logged.
	Fatal() Entry

	//WriteCloser returns an io.Writer that when written to writes logs at the given level.
	//It is the callers responsibility to call Close when finished.
	//This is particularly useful for redirecting the output of other loggers or even Readers
	//with the help of io.TeeReader.
	WriteCloser(Level) io.WriteCloser
}

//Entry is the primary interface by which individual log entries are made.
type Entry interface {
	//Async flips the current Entry to be asynchronous, or back if called more than once.
	//If set to asynchronous an Entry implementation should not log its final message until Send is called.
	Async() Entry
	//Caller embeds a caller value into the existing Entry.
	//A caller value is a filepath followed by line number.
	//Skip determines the number of additional stack frames to ascend when
	//determining the value. By default the caller of the method is the value used
	//and skip does not need to be supplied in that case.
	//Caller may be called multiple times on the Entry to build a stack or execution trace.
	Caller(skip ...int) Entry

	WithError(errs ...error) Entry
	WithField(key string, val interface{}) Entry
	WithFields(fields map[string]interface{}) Entry
	WithBool(key string, bls ...bool) Entry
	WithDur(key string, durs ...time.Duration) Entry
	WithInt(key string, is ...int) Entry
	WithUint(key string, us ...uint) Entry
	WithStr(key string, strs ...string) Entry

	//WithTime adds the respective time values to the Entry at the given key.
	//Note that many loggers add a "time" key automatically and time formatting
	//may be dependant on configuration or logger choice.
	WithTime(key string, ts ...time.Time) Entry

	Trace() Entry
	Debug() Entry
	Info() Entry
	Warn() Entry
	Error() Entry
	Panic() Entry
	Fatal() Entry

	//Msgf formats and sets the final log message for this Entry.
	//It will also send the message if Async has not been set.
	Msgf(string, ...interface{})
	//Msg sets the final log message for this Entry.
	//It will also send the message if Async has not been set.
	Msg(msg string)
	//Send sends the final log entry.
	//This interface does not define the behavior of calling this method more than once.
	Send()
}

//UnderlyingLogger is an escape hatch allowing Loggers registered with this package
//the option to return their underlying implementation, as well as reset it.
//Note this is currently required for CustomOption's to work.
type UnderlyingLogger interface {
	GetLogger() interface{}
	SetLogger(interface{})
}
