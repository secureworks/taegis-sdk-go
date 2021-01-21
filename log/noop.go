package log

import (
	"io"
	"time"
)

func init() {
	Register("noop", func(_ *Config, _ ...Option) (Logger, error) {
		return Noop(), nil
	})
}

//Noop returns an implementation of Logger that does nothing when its methods are called.
//This can also be retrieved using "noop" with "Open", though in that case any configuration is ignored.
func Noop() Logger {
	return noopL{}
}

////log impl
type noopL struct{}

func (noopL) WithError(_ error) Entry                   { return noopE{} }
func (noopL) WithField(_ string, _ interface{}) Entry   { return noopE{} }
func (noopL) WithFields(_ map[string]interface{}) Entry { return noopE{} }

func (noopL) Entry(_ Level) Entry { return noopE{} }
func (noopL) Trace() Entry        { return noopE{} }
func (noopL) Debug() Entry        { return noopE{} }
func (noopL) Info() Entry         { return noopE{} }
func (noopL) Warn() Entry         { return noopE{} }
func (noopL) Error() Entry        { return noopE{} }
func (noopL) Panic() Entry        { return noopE{} }
func (noopL) Fatal() Entry        { return noopE{} }

func (np noopL) WriteCloser(_ Level) io.WriteCloser { return np }
func (noopL) Write(p []byte) (int, error)           { return len(p), nil }
func (noopL) Close() error                          { return nil }

////

////entry impl
type noopE struct{}

func (np noopE) Async() Entry          { return np }
func (np noopE) Caller(_ ...int) Entry { return np }

func (np noopE) WithError(_ ...error) Entry                 { return np }
func (np noopE) WithField(_ string, _ interface{}) Entry    { return np }
func (np noopE) WithFields(_ map[string]interface{}) Entry  { return np }
func (np noopE) WithBool(_ string, _ ...bool) Entry         { return np }
func (np noopE) WithDur(_ string, _ ...time.Duration) Entry { return np }
func (np noopE) WithInt(_ string, _ ...int) Entry           { return np }
func (np noopE) WithUint(_ string, _ ...uint) Entry         { return np }
func (np noopE) WithStr(_ string, _ ...string) Entry        { return np }
func (np noopE) WithTime(_ string, _ ...time.Time) Entry    { return np }

func (np noopE) Trace() Entry { return np }
func (np noopE) Debug() Entry { return np }
func (np noopE) Info() Entry  { return np }
func (np noopE) Warn() Entry  { return np }
func (np noopE) Error() Entry { return np }
func (np noopE) Panic() Entry { return np }
func (np noopE) Fatal() Entry { return np }

func (noopE) Msgf(_ string, _ ...interface{}) {}
func (noopE) Msg(_ string)                    {}
func (noopE) Send()                           {}

////
