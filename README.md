# zipkin-replay
Replay Zipkin traces for debug and demo

A small utility for replaying Zipkin v1 traces captured with the [usurp](https://github.com/prydin/usurp) tool. It will try to parse every record as an array of Zipkin v1 
spans and ignore anything it can't parse. All timestamps are adjusted to count from the time the tool was launched, minus an 
optional offset. All IDs are remapped to make sure traces are unique even if the tool is run multiple times against the same 
dataset. Also, namespace suffixes (as they appear in Istio traces) can be remapped.

## Usage

Installing:
```
go get github.com/prydin/zipkin-replay
go build -o zreplay
```

## Usage
```
Usage of ./zreplay:
  -file string
    	Input file captured with usurp
  -forever
    	Run realtime with delays according to original timing
  -gz
    	Input is gzipped
  -log.format value
    	Set the log target and format. Example: "logger:syslog?appname=bob&local=7" or "logger:stdout?json=true" (default "logger:stderr")
  -log.level value
    	Only log messages with the given severity or above. Valid levels: [debug, info, warn, error, fatal] (default "info")
  -namespace string
    	Namespace substitution (default "default")
  -offset duration
    	Time offset in minutes
  -realtime
    	Run realtime with delays according to original timing
  -target string
    	The target URL
```
## Example
```
./zreplay -file data/zipkin.out -target http://localhost:9411/api/v1/spans -namespace foobar
```
