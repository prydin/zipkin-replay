package model

type Span struct {
	TraceID     string `json:"traceId"`
	Name        string `json:"name"`
	ID          string `json:"id"`
	ParentID    string `json:"parentId,omitempty"`
	Timestamp   int64  `json:"timestamp"`
	Duration    int64    `json:"duration,omitempty"`
	Annotations []struct {
		Timestamp int64  `json:"timestamp"`
		Value     string `json:"value"`
		Endpoint  struct {
			Ipv4        string `json:"ipv4"`
			Port        int    `json:"port"`
			ServiceName string `json:"serviceName"`
		} `json:"endpoint"`
	} `json:"annotations"`
	BinaryAnnotations []struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	} `json:"binaryAnnotations"`
}

type SpanList []Span
