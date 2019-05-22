package transform

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/prydin/wf-replay/reader"
	"github.com/stretchr/testify/assert"
	"os"
	"regexp"
	"strings"
	"testing"
	"time"
)

func TestNamespace(t *testing.T) {
	f, err := os.Open("../data/zipkin.one")
	assert.NoError(t, err, "Error accessing test data")
	rdr := bufio.NewReader(f)
	dr := reader.NewDumpReader(rdr)

	tr, err := NewTransformer()
	assert.NoError(t, err)

	nameRe := regexp.MustCompile(".+\\.foobar\\.svc.cluster.local:.+")
	svcRe := regexp.MustCompile(".+\\.foobar")
	for {
		r, err := dr.Next()
		assert.NoError(t, err)
		if r == nil {
			break
		}
		tr.Transform(&r.Span, "foobar", 0)
		buf, err := json.Marshal(&r.Span)
		dst := bytes.NewBufferString("")
		json.Indent(dst, buf, "", "    ")
		fmt.Println(dst.String())
		for _, s := range r.Span {
			//fmt.Printf("%s\n%s\n%d\n\n", s.ID, s.ParentID, len(s.ID))
			assert.NotEqual(t, s.ParentID, s.ID, "Parent ID must not equal span ID")
			if strings.Index(s.Name, "async") != 0 {
				assert.Regexp(t, nameRe, s.Name, "Name didn't match")
			}
			for _, a := range s.Annotations {
				if strings.Index(a.Endpoint.ServiceName, "async") != 0 {
					assert.Regexp(t, svcRe, a.Endpoint.ServiceName, "Name didn't match")
				}
			}
		}
	}

}

func TestTime(t *testing.T) {
	f, err := os.Open("../data/zipkin.out")
	assert.NoError(t, err, "Error accessing test data")
	rdr := bufio.NewReader(f)
	dr := reader.NewDumpReader(rdr)

	tr, err := NewTransformer()
	assert.NoError(t, err)
	for {
		r, err := dr.Next()
		assert.NoError(t, err)
		tr.Transform(&r.Span, "foobar", 0)
		s := r.Span
		assert.InDelta(t, time.Now().UnixNano()/1000, s[0].Timestamp, 100)
	}
}