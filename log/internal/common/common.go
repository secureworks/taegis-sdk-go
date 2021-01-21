package common

import (
	"fmt"
	"path/filepath"
	"sync/atomic"

	"github.com/getsentry/sentry-go"
	"github.com/pkg/errors"
)

//StackTracer is one common interface for extracting stack information from types
//Sentry and several other packages check for this
type StackTracer interface {
	StackTrace() errors.StackTrace
}

//WithStackTrace checks err for implementing StackTracer and returns it if it does.
//Otherwise it'll wrap err in an error type that implements StackTracer and return both.
//If nil is passed then nil is returned.
func WithStackTrace(err error) (st StackTracer, e error) {
	if err == nil {
		return
	}

	e = err
	var ok bool

	st, ok = err.(StackTracer)
	if !ok {
		e = errors.WithStack(err)
		st = e.(StackTracer)
	}

	return
}

var nonRentrant int32

//InitSentry is a shared function for initializing the global or "CurrentHub" for Sentry.
//Unfortunately while Sentry supports several hubs, the default and most commonly used one is "CurrentHub",
//a package level variable that cannot be setup concurrently (it is data race free but not race condition free).
//Ideally each logger instance would have their own hub, but users will likely expect this to setup Sentry "in total".
//For now just do what we can so that multiple logger instances use the same thing if desired.
func InitSentry(opts sentry.ClientOptions) (err error) {
	if !atomic.CompareAndSwapInt32(&nonRentrant, 0, 1) {
		return
	}

	defer func() {
		if err != nil {
			//allow another call if this one failed
			atomic.StoreInt32(&nonRentrant, 0)
		}
	}()

	//TODO: Reader, do we setup fakesentry? Allows for interception/debugging of data sent to sentry
	//if opts.Debug && setupFakeSentry != nil

	err = sentry.Init(opts)
	return
}

//ParseFrames returns a slice of sentry.Frames for string values produced by a StackTracer.
//It accepts interfaces as it is meant to be used with JSON marshaling, otherwise call ParseFrame directly.
func ParseFrames(vals ...interface{}) (frames []sentry.Frame) {
	if len(vals) == 0 {
		return
	}

	frames = make([]sentry.Frame, 0, len(vals))
	for _, v := range vals {
		s, ok := v.(string)
		if !ok {
			break
		}

		frames = append(frames, ParseFrame(s))
	}

	return
}

//ParseFrame parses a single sentry.Frame from a string produced by a StackTracer.
func ParseFrame(str string) sentry.Frame {
	fnName, file, lineNo := parseFrameStr(str)

	return sentry.Frame{
		Function: fnName,
		Filename: filepath.Base(file), //this will become "." if file is empty
		AbsPath:  file,
		Lineno:   lineNo,
	}
}

func parseFrameStr(frame string) (fnName, file string, lineNo int) {
	//if this fails because stack was "unknown" then just fnName becomes unknown
	//https://github.com/pkg/errors/blob/614d223910a179a466c1767a985424175c39b465/stack.go#L93
	fmt.Sscanf(frame, "%s %s:%d", &fnName, &file, &lineNo)
	return
}
