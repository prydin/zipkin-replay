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

Running:
```
./zreplay -file data/zipkin.out -target http://localhost:9411/api/v1/spans -namespace foobar
```
