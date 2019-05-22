# zipkin-replay
Replay Zipkin traces for debug and demo

A small utility for replaying Zipkin v1 traces captured with the usurp tool. It will try to parse every record as a Zipkin v1 
trace and ignore anything it can't parse. All timestamps are adjusted to count from the time the tool was launched, minus an 
optional offset. All IDs are remapped to make sure traces are unique even if the tool is run multiple times against the same 
dataset. Also, namespace suffixes (as they appear in Istio traces) can be remapped.

Work in progress.
