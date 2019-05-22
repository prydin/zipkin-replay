package transform

import (
	"github.com/google/uuid"
	lru "github.com/hashicorp/golang-lru"
	"github.com/prometheus/common/log"
	"github.com/prydin/wf-replay/model"
	"strings"
	"time"
)

type Transformer struct {
	idCache *lru.Cache
	startOffset int64
}

func NewTransformer() (*Transformer, error) {
	a := &Transformer{ }
	var err error
	a.idCache, err = lru.New(1000)
	if err != nil {
		return nil, err
	}
	return a, nil
}

func (t *Transformer) Transform(span *model.SpanList, namespace string, offset time.Duration) {
	s := *span
	for i := range s {
		s[i].TraceID = t.getID(s[i].TraceID, 32)
		if s[i].ParentID != "" {
			s[i].ParentID = t.getID(s[i].ParentID, 16)
		}
		s[i].ID = t.getID(s[i].ID, 16)
		for j := range s[i].BinaryAnnotations {
			if s[i].BinaryAnnotations[j].Key == "guid:x-request-id" {
				s[i].BinaryAnnotations[j].Value = t.getID(s[i].BinaryAnnotations[j].Value, 32)
			}
		}
	}
	if namespace != "" {
		t.transformNamespace(span, namespace)
	}
	t.adjustTimestamps(span, offset)
}

func (t *Transformer) getID(id string, length int) string {
	ifc, ok := t.idCache.Get(id)
	if !ok {
		u, err := uuid.NewRandom()
		if err != nil {
			log.Fatal(err)
		}
		s := strings.ReplaceAll(u.String(), "-", "")
		s = s[len(s)-length:]
		t.idCache.Add(id, s)
		return s
	}
	return ifc.(string)
}

func (t *Transformer) transformNamespace(span *model.SpanList, namespace string) {
	s := *span
	for i := range s {
		s[i].Name = strings.Replace(s[i].Name, "default.", namespace + ".", 1)
		for j := range s[i].Annotations {
			s[i].Annotations[j].Endpoint.ServiceName = strings.Replace(s[i].Annotations[j].Endpoint.ServiceName,
				".default", "." + namespace, 1)
		}
	}
}

func (t *Transformer) adjustTimestamps(span *model.SpanList, offset time.Duration) {
	s := *span
	if t.startOffset == 0 {
		start := toTimestamp(s[0].Timestamp)
		t.startOffset = time.Since(start).Nanoseconds() / 1000
	}
	o := t.startOffset - offset.Nanoseconds() / 1000
	for i := range s {
		s[i].Timestamp += o
		for j := range s[i].Annotations {
			s[i].Annotations[j].Timestamp += o
		}
	}
}

func toTimestamp(t int64) time.Time {
	return time.Unix(t/1000000, 1000*(t%1000000))
}
