package sender

import (
	"bufio"
	"github.com/prydin/wf-replay/reader"
	"github.com/prydin/wf-replay/transform"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestSenderOne(t *testing.T) {
	f, err := os.Open("../data/zipkin.one")
	assert.NoError(t, err, "Error accessing test data")
	s := NewHTTPSender("http://testproxy:9411/api/v1/spans")
	rdr := bufio.NewReader(f)
	dr := reader.NewDumpReader(rdr)
	tr, err := transform.NewTransformer()
	assert.NoError(t, err)

	for {
		r, err := dr.Next()
		assert.NoError(t, err, "Error advancing to next record")
		if r == nil {
			break
		}
		tr.Transform(&r.Span, "foobar", 0)
		assert.NoError(t, s.Send(r))
	}
}

func TestSenderMany(t *testing.T) {
	f, err := os.Open("../data/zipkin.medium")
	assert.NoError(t, err, "Error accessing test data")
	s := NewHTTPSender("http://testproxy:9411/api/v1/spans")
	rdr := bufio.NewReader(f)
	dr := reader.NewDumpReader(rdr)
	tr, err := transform.NewTransformer()
	assert.NoError(t, err)

	for {
		r, err := dr.Next()
		assert.NoError(t, err, "Error advancing to next record")
		if r == nil {
			break
		}
		tr.Transform(&r.Span, "foobar", 0)
		assert.NoError(t, s.Send(r))
	}
}