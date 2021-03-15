package zerolog

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/getsentry/sentry-go"
	"github.com/rs/zerolog"

	"github.com/secureworks/taegis-sdk-go/log"
	"github.com/secureworks/taegis-sdk-go/log/internal/common"
)

/*
An unfortunate but necessary type, as zerolog's hook interface is next to useless for extracting data from.
Older sentry writers used raven-go instead of sentry-go and zerolog authors themselves
seem to prefer using io.Writer instead of the hook interface. So here we are.

https://github.com/rs/zerolog/blob/72acd6cfe8bbbf5c52bfc805a3889c6941499c95/journald/journald.go#L37
https://github.com/rs/zerolog/blob/72acd6cfe8bbbf5c52bfc805a3889c6941499c95/console.go#L86
https://github.com/rs/zerolog/issues/93
https://gist.github.com/asdine/f821abe6189a04250ae61b77a3048bd9
*/
type sentryWriter struct {
	lvlField []byte
	hub      *sentry.Hub
	lvlSet   map[zerolog.Level]bool
}

func newSentryWriter(lvls ...log.Level) *sentryWriter {
	lvlSet := make(map[zerolog.Level]bool, len(lvls))
	for _, lvl := range lvls {
		lvlSet[lvlToZerolog(lvl)] = true
	}

	return &sentryWriter{
		lvlField: []byte(`"` + zerolog.LevelFieldName + `"`),
		hub:      sentry.CurrentHub(),
		lvlSet:   lvlSet,
	}
}

func (sw *sentryWriter) Write(p []byte) (n int, err error) {
	n = len(p)
	zlvl, ok := sw.checkLvl(p)
	if !ok {
		return
	}

	dat := make(map[string]interface{})
	json.Unmarshal(p, &dat)

	delete(dat, zerolog.LevelFieldName)
	if len(dat) == 0 {
		return
	}

	event := sentry.NewEvent()
	event.Level = zlvlToSentry(zlvl)

	if msg, ok := dat[zerolog.MessageFieldName].(string); ok {
		event.Message = msg
		delete(dat, zerolog.MessageFieldName)
	}

	//we should just never send "Exception"'s to Sentry on principle
	var exe *sentry.Exception

	if iface, ok := dat[zerolog.ErrorFieldName]; ok {
		es, ok := iface.(string)
		if !ok {
			es = fmt.Sprintf("%v", iface)
		}

		delete(dat, zerolog.ErrorFieldName)
		exe = &sentry.Exception{
			Value: es,
		}
	}

	if iface, ok := dat[zerolog.ErrorStackFieldName].([]interface{}); ok {
		frames := common.ParseFrames(iface...)

		if len(frames) > 0 {
			delete(dat, zerolog.ErrorStackFieldName)
			if exe == nil {
				exe = new(sentry.Exception)
			}

			exe.Stacktrace = &sentry.Stacktrace{
				Frames: frames,
			}
		}
	}

	if exe != nil {
		event.Exception = append(event.Exception, *exe)
		exe = nil
	}

	if iface, ok := dat[log.PanicValue]; ok {
		pv, ok := iface.(string)
		if !ok {
			pv = fmt.Sprintf("%v", pv)
		}

		delete(dat, log.PanicValue)
		exe = &sentry.Exception{
			Value: pv,
		}
	}

	if iface, ok := dat[log.PanicStack].([]interface{}); ok {
		frames := common.ParseFrames(iface...)

		if len(frames) > 0 {
			delete(dat, log.PanicStack)
			if exe == nil {
				exe = new(sentry.Exception)
			}

			exe.Stacktrace = &sentry.Stacktrace{
				Frames: frames,
			}
		}
	}

	if exe != nil {
		event.Exception = append(event.Exception, *exe)
		exe = nil
	}

	//additional values as "Extra" vs "Tags" vs "Breadcrumbs"
	event.Extra = dat
	sw.hub.CaptureEvent(event)

	return
}

func (sw *sentryWriter) checkLvl(p []byte) (zlvl zerolog.Level, shouldSend bool) {
	if p == nil || len(p) < len(sw.lvlField) {
		return
	}

	//we have it on authority that the level field will be in the latter part of the message
	i := bytes.Index(p[len(p)/2:], sw.lvlField)
	if i == -1 {
		//try the full slice, wasn't in the latter half (or was split between halfs)
		i = bytes.Index(p, sw.lvlField)
	} else {
		//add back the len we skipped
		i += (len(p) / 2)
	}

	//still don't have level field?
	if i == -1 {
		return
	}

	//hardcore math
	startingOffset := i + len(sw.lvlField) + 2
	if startingOffset >= len(p) {
		return
	}

	i = bytes.IndexByte(p[startingOffset:], '"')
	if i == -1 {
		return
	}

	i += startingOffset
	lvl, err := zerolog.ParseLevel(string(p[startingOffset:i]))

	return lvl, err == nil && sw.lvlSet[lvl]
}

func zlvlToSentry(lvl zerolog.Level) sentry.Level {
	switch lvl {
	case zerolog.TraceLevel, zerolog.DebugLevel:
		return sentry.LevelDebug
	case zerolog.InfoLevel:
		return sentry.LevelInfo
	case zerolog.WarnLevel:
		return sentry.LevelWarning
	case zerolog.ErrorLevel, zerolog.PanicLevel:
		return sentry.LevelError
	case zerolog.FatalLevel:
		return sentry.LevelFatal
	default:
		return sentry.LevelWarning
	}
}
