package log_test

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/secureworks/taegis-sdk-go/log"
	"github.com/secureworks/taegis-sdk-go/log/middleware"
)

var (
	//these are UNUSED and only for making this example compile
	e http.Handler
)

func Example() {
	//if config is passed as nil to log.Open then this is used
	config := log.DefaultConfig(nil)

	//various other options
	//config.EnableErrStack = true
	//config.LocalDevel = true

	//optional custom options, passing nothing for this is fine, look to log.CustomOption for more info
	//var opts []log.Option
	//opts = append(opts, log.CustomOption("Sample", mySamplerInterfaceType))

	//Zerolog implementation underscore imported
	logger, err := log.Open("zerolog", config)
	//logger, err := log.Open("zerolog", config, opts...)
	if err != nil {
		//handle error
		fmt.Fprintf(os.Stderr, "Unable to setup logger: %v", err)
		os.Exit(1)
	}

	//where e is your base handler, echo and buffalo have adapters
	handler := middleware.NewHTTPRequestMiddleware(logger, 0)(e)

	//...
	//setup logger for use by http.Server or anything that can use an io.Writer
	//wc := logger.WriteCloser(log.WARN)
	//defer wc.Close()

	errc := make(chan error)
	listenAddr := ":8080"

	//or create new http.Server with logger as base context for use by request handlers
	//and ErrorLog set to write to logger at WARN level
	srv, c := middleware.NewHTTPServer(logger, log.WARN)
	defer c.Close()

	//more options you should be setting
	srv.Addr = listenAddr
	srv.Handler = handler
	srv.WriteTimeout = time.Second * 15
	srv.ReadTimeout = time.Second * 10
	logger.Info().Msgf("Starting server on: %s", listenAddr)
	go func() { errc <- srv.ListenAndServe() }()

	//block on sig term || int || etc, whatever else
}
