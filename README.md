# simple-httpd
An http server that serves up a static directory on localhost

## Installation 

```
go install github.com/efixler/simple-httpd@latest
```

## Usage

```
% simple-httpd -h
Usage: 
        simple-httpd [flags] 

Flags:
 
  -h
        Show this help message
  -dir value
        Directory to serve
        Environment: SIMPLE_HTTPD_DIRECTORY (default .)
  -host value
        Host to run the server on [* for all interfaces]
        Environment: SIMPLE_HTTPD_HOST (default 127.0.0.1)
  -log-level value
        Set the log level [debug|error|info|warn]
        Environment: SIMPLE_HTTPD_LOG_LEVEL
  -port value
        Port to run the server on
        Environment: SIMPLE_HTTPD_PORT (default 4411)
```

## Building

```
% make help

Usage:
  make 
  build            build the binaries, to the build/ folder (default target)
  clean            clean the build directory
  cognitive        run the cognitive complexity checker
  help             show this help message
```
