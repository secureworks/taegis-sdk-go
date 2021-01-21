package log

import (
	"fmt"
	"io"
	"net/url"
	"os"
	"reflect"
	"strings"
)

//EnvKey is a publicly documented string type for environment lookups performed for DefaultConfig.
//It is otherwise unspecial.
type EnvKey string

const (
	//LogLevel is the EnvKey used for looking up the log level.
	LogLevel EnvKey = "LOG_LEVEL"
	//ErrStack is the EnvKey used for looking up error/stack trace logging; Value should be true || True || TRUE.
	ErrStack EnvKey = "ERROR_STACK"
	//SentryDSN is the EnvKey used for looking up the Sentry project DNS; Empty value disables Sentry.
	SentryDSN EnvKey = "SENTRY_DSN"
	//SentryLevels is the EnvKey used for looking up which log levels will be sent to Sentry; Values should be comma separated.
	SentryLevels EnvKey = "SENTRY_LEVELS"
	//Environment is the EnvKey used for looking up the current deployment environment; Values commonly dev || prod.
	Environment EnvKey = "ENVIRONMENT"
)

func (ek EnvKey) String() string {
	return string(ek)
}

//Config defines common logger configuration options.
type Config struct {
	//Level is the level at which returned Logger's will be considered enabled.
	//For example, setting WARN then logging and sending a Debug entry will cause
	//the entry to not be logged.
	Level Level

	//LocalDevel, may be used by some Logger's for local debugging changes.
	LocalDevel bool

	//Format is the format the Logger should log in.
	Format LoggerFormat

	//EnableErrStack enables error stack gathering and logging.
	EnableErrStack bool

	//Output is the io.Writer the Logger will write messages to.
	Output io.Writer

	//Sentry is a sub-config type for configurating Sentry if desired.
	//No other portion of this struct is considered if DSN is not set and valid.
	Sentry struct {
		//DSN is the Sentry DSN.
		DSN string
		//Release is the program release or revision.
		Release string
		//Env is the deployment environment; "prod", "dev", etc.
		Env string
		//Server is the server or hostname.
		Server string
		//Levels are the log levels that will trigger an event to be sent to Sentry.
		Levels []Level
		//Debug is a passthrough for Sentry debugging.
		Debug bool
	}
}

//DefaultConfig returns a Config instance with sane defaults.
//env is a callback for looking up EnvKeys, it is set to os.Getenv if nil.
//Fields and values returned by this function can be altered.
func DefaultConfig(env func(string) string) *Config {
	/*
		TODO: Configuration is still an issue, as discussed during design phase.
		If we use consol + vault, most of our values should become env vars.
		This is a work in progress.
	*/
	conf := new(Config)

	if env == nil {
		env = os.Getenv
	}

	if lvlStr := env(LogLevel.String()); lvlStr != "" {
		conf.Level = LevelFromString(lvlStr)
	}

	if errStackStr := env(ErrStack.String()); errStackStr != "" {
		conf.EnableErrStack = strings.ToUpper(errStackStr) == "TRUE"
	}

	sentryDSN := env(SentryDSN.String())

	if sentryDSN != "" {
		lvls := []Level{FATAL, PANIC, ERROR}

		split := strings.Split(env(SentryLevels.String()), ",")
		if len(split) > 0 && split[0] != "" {
			lvlSet := make(map[Level]bool, len(split))
			for _, lvl := range split {
				lvlSet[LevelFromString(lvl)] = true
			}

			lvls = make([]Level, 0, len(lvlSet))
			for lvl := range lvlSet {
				lvls = append(lvls, lvl)
			}
		}

		host, _ := os.Hostname()

		conf.Sentry.DSN = sentryDSN
		conf.Sentry.Levels = lvls
		// TODO: fix this
		//conf.Sentry.Release = version.Get().Revision
		conf.Sentry.Server = host
		conf.Sentry.Env = env(Environment.String())
	}

	if _, err := url.Parse(sentryDSN); err != nil {
		conf.Sentry.DSN = ""
	}

	conf.Output = os.Stderr
	return conf
}

//NewLogger is a function type for Logger implemenations to register themselves.
type NewLogger func(*Config, ...Option) (Logger, error)

var (
	setup = make(map[string]NewLogger, 2)
)

//Register registers the provided NewLogger function under the given name for use with Open.
//Note, this method is not concurreny safe, nil NewLoggers or duplicate registration
//will cause a panic.
func Register(name string, nl NewLogger) {
	if _, ok := setup[name]; ok || nl == nil {
		panic(fmt.Errorf("log: %s already registered with logging package", name))
	}

	setup[name] = nl
}

//Open returns a new instance of the selected Logger with config and options.
func Open(name string, conf *Config, opts ...Option) (Logger, error) {
	nl, ok := setup[name]
	if !ok {
		return nil, fmt.Errorf("log: No logger by name (%s)", name)
	}

	if conf == nil {
		conf = DefaultConfig(nil)
	}

	return nl(conf, opts...)
}

//Option is a function type that accepts an interface value and returns an error.
type Option func(interface{}) error

func noopOption(_ interface{}) error { return nil }

func errOption(err error) Option {
	return func(_ interface{}) error {
		return err
	}
}

//CustomOption takes the name of a method on the UnderlyingLogger of a Logger implementation
//as well as a value, and returns an Option. name is the case sensitive name of a method while
//val should be a single value needed as input to the named method. If several values are needed as input
//then val should be a function that accepts no input and returns values to be used as input to the
//named method. If val is a function and returns an error as its only or last value it will be checked
//and returned if non nil, otherwise remaining values are fed to the named method.
//A nil val is valid so long as the named method expects nil input or no input.
//If the named method returns an instance of itself it will be set back as the new UnderlyingLogger.
//If the named method returns an error that error will be checked and returned.
//Look to the CustomOption tests and package level example for basic usage.
func CustomOption(name string, val interface{}) Option {
	if name == "" {
		return noopOption
	}

	valFunc, err := getReflectVals(val)
	if err != nil {
		return errOption(err)
	}

	return func(topLogger interface{}) (err error) {
		ul, ok := topLogger.(UnderlyingLogger)
		if !ok {
			return fmt.Errorf("log: Logger type (%T) does not support the UnderlyingLogger interface", topLogger)
		}

		defer func() {
			pv := recover()

			if pv == nil || err != nil {
				return
			}

			if e, ok := pv.(error); ok {
				err = e
			} else {
				err = fmt.Errorf("log: Panic caught in CustomOption for %s: %v", name, pv)
			}
		}()

		//get logger and check if nil, if it is this was
		//intentional by implementation and we should just return
		logger := ul.GetLogger()
		if logger == nil {
			return
		}

		logVal := reflect.ValueOf(logger)
		if !logVal.IsValid() {
			return
		}

		methVal := logVal.MethodByName(name)
		if !methVal.IsValid() {
			return
		}

		var wasError bool
		vals := valFunc()
		//check if last val is error or error that is nil
		if l := len(vals); l > 0 {
			err, wasError = valueToError(vals[l-1])
			if err != nil {
				return
			}

			//remove value
			if wasError {
				vals = vals[:l-1]
			}
		}

		//from this point we have the method we want to call and
		//the values with which to call it
		//we could check if each input matches what is expected but
		//instead we'll just call the method and rely on the defer recover
		//above to stop us from calling something wrong
		out := methVal.Call(vals)

		le := len(out)
		if le == 0 {
			return
		}

		err, wasError = valueToError(out[le-1])
		if err != nil {
			return
		}
		if wasError {
			out = out[:le-1]
			le = len(out)
		}

		if le == 0 {
			return
		}

		//if one of the remaining types in the output is the same
		//as the underlying logger, set it back as the underlying logger
		//this is common in methods/funcs that chain configuration
		logType := logVal.Type()
		for _, val := range out {
			if !val.IsValid() {
				continue
			}

			/*
				Since logType is an interface this never seems correct, but according to
				https://golang.org/pkg/reflect/#Type
				"Type values are comparable, such as with the == operator, so they can be used as map keys.
				Two Type values are equal if they represent identical types."
			*/
			if val.Type() == logType {
				ul.SetLogger(val.Interface())
				break
			}
		}

		return
	}
}

func getReflectVals(val interface{}) (func() []reflect.Value, error) {
	reflval := reflect.ValueOf(val)
	if !reflval.IsValid() {
		return func() []reflect.Value { return []reflect.Value{} }, nil
	}

	typ := reflval.Type()
	if typ.Kind() != reflect.Func {
		return func() []reflect.Value { return []reflect.Value{reflval} }, nil
	}

	//we know it's a func, check to make sure it doesn't take args
	if typ.NumIn() > 0 {
		name := typ.Name()
		if name == "" {
			name = "anon func"
		}

		return nil, fmt.Errorf("log: Function value (%s) expects inputs", name)
	}

	return func() []reflect.Value { return reflval.Call([]reflect.Value{}) }, nil
}

var errInterface = reflect.TypeOf((*error)(nil)).Elem()

//checks and converts a reflect.Value to an error if appropriate.
//nolint
func valueToError(val reflect.Value) (err error, wasError bool) {
	if !val.IsValid() {
		return
	}

	if val.Kind() != reflect.Interface || !val.Type().Implements(errInterface) {
		return
	}

	wasError = true
	if val.IsNil() {
		return
	}

	err = val.Elem().Interface().(error)
	return
}
