package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/efixler/envflags"
	"github.com/efixler/webutil/middleware"
)

const (
	DefaultPort = 4411
)

var (
	flags   flag.FlagSet
	host    *envflags.Value[string]
	port    *envflags.Value[int]
	dir     *envflags.Value[string]
	noCache *envflags.Value[bool]
)

func main() {
	assertDirectoryExists(dir.Get())

	fs := http.StripPrefix("/", http.FileServer(http.Dir(dir.Get()))).ServeHTTP
	if noCache.Get() {
		fs = middleware.Chain(fs, middleware.NoCache)
	}

	http.HandleFunc("/", fs)

	slog.Info("simple-httpd starting up", "dir", dir.Get(), "host:port", fmt.Sprintf("%s:%d", host.Get(), port.Get()))
	host := host.Get()
	if host == "*" {
		host = ""
	}

	if err := http.ListenAndServe(fmt.Sprintf("%s:%d", host, port.Get()), nil); err != http.ErrServerClosed {
		slog.Error("simple-httpd unable to start", "error", err)
		os.Exit(1)
	}
}

func assertDirectoryExists(path string) {
	stat, err := os.Stat(path)
	defer func() {
		if r := recover(); r != nil {
			slog.Error("simple-httpd unable to start", "error", r)
			os.Exit(1)
		}
	}()
	switch err {
	case nil:
		if !stat.IsDir() {
			panic(fmt.Errorf("%s is not a directory", path))
		}
	default:
		if os.IsNotExist(err) {
			panic(fmt.Errorf("directory %s does not exist", path))
		} else {
			panic(err)
		}
	}
}

func init() {
	flags.Init("", flag.ExitOnError)
	flags.Usage = usage
	envflags.EnvPrefix = "SIMPLE_HTTPD_"
	host = envflags.NewString("HOST", "127.0.0.1")
	host.AddTo(&flags, "host", "Host to run the server on [* for all interfaces]")
	port = envflags.NewInt("PORT", DefaultPort)
	port.AddTo(&flags, "port", "Port to run the server on")
	dir = envflags.NewString("DIRECTORY", ".")
	dir.AddTo(&flags, "dir", "Directory to serve")
	noCache = envflags.NewBool("NO_CACHE", false)
	noCache.AddTo(&flags, "no-cache", "Send no-cache headers on every request")
	logLevel := envflags.NewLogLevel("LOG_LEVEL", slog.LevelInfo)
	logLevel.AddTo(&flags, "log-level", "Set the log level [debug|error|info|warn]")
	flags.Parse(os.Args[1:])
	ll := logLevel.Get()
	logger := slog.New(slog.NewTextHandler(
		os.Stderr,
		&slog.HandlerOptions{
			Level: ll,
		},
	))
	slog.SetDefault(logger)
}

func usage() {
	fmt.Println(`Usage: 
	simple-httpd [flags] 

Flags:
 
  -h	
  	Show this help message`)
	flags.PrintDefaults()
}
