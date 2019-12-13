package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
)

var usage = `greeter

Example service which greets a client using a configurable prefix.
`

var (
	Verbose = flag.Bool("verbose", false, "Print verbose logs. [VERBOSE]")
	Prefix  = flag.String("prefix", "Hello", "Greet the world with the given prefix. [PREFIX]")
	Listen  = flag.String("listen", "0.0.0.0:8080", "Listen address (interface and port). [LISTEN]")
	Version = flag.Bool("version", false, "Print the version information and exit without starting. [VERSION]")

	// Populated during build:
	brch string // active branch
	date string // date stamp
	vers string // version tag or "tip"
	hash string // active commit hash (no "dirty" flag)
)

func main() {
	// Override default help (usage text)
	// with the simple message and flag details.
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, usage)
		flag.PrintDefaults()
	}
	// Parse environment and command-line flags,
	// resulting in values populated in the following
	// order of precidence:
	// Flag defaults -> Environment Variables -> CLI flags
	parseEnv()
	flag.Parse()

	// If Verbose was enabled, print effective config.
	if *Verbose {
		printConfig()
	}
	// If Version was requested, print the verbose
	// binary information populated at build-time and
	// immediately exit.
	if *Version {
		fmt.Println(verboseVersion())
		return
	}

	// Run the service, exiting non-zero if any error
	// is generated.
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func parseEnv() {
	parseBool("VERSION", Version)
	parseBool("VERBOSE", Verbose)
	parseString("PREFIX", Prefix)
	parseString("LISTEN", Listen)
}

func printConfig() {
	fmt.Printf("VERBOSE=%v\n", *Verbose)
	fmt.Printf("LISTEN=%v\n", *Listen)
	fmt.Printf("PREFIX=%v\n", *Prefix)
}

func verboseVersion() string {
	return fmt.Sprintf("%s-%s-%s-%s", brch, date, vers, hash)
}

func run() (err error) {
	s := http.Server{
		Addr:           *Listen,
		Handler:        http.HandlerFunc(handle),
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   5 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	if *Verbose {
		fmt.Fprintf(os.Stderr, "greeter listening at %v\n", *Listen)
	}
	return s.ListenAndServe()
}

func handle(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("content-type", "text/plain")
	fmt.Fprintf(res, "%v World!\n", *Prefix)
}

func parseBool(key string, value *bool) {
	if val, ok := os.LookupEnv(key); ok {
		b, err := strconv.ParseBool(val)
		if err != nil {
			panic(err)
		}
		*value = b
	}
}

func parseString(key string, value *string) {
	if val, ok := os.LookupEnv(key); ok {
		*value = val
	}
}
