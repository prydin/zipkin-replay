package reader

import (
	"bufio"
	"encoding/json"
	"github.com/prydin/wf-replay/model"
	"io"
	"net/http"
	"strings"
)

type Record struct {
	Headers http.Header
	Span model.SpanList
}

type DumpReader struct {
	in *bufio.Reader
}

func NewDumpReader(rdr *bufio.Reader) *DumpReader {
	return &DumpReader{rdr}
}

func (d *DumpReader) Next() (*Record, error) {
	for {
		r := Record{}

		// Skip to start
		for {
			s, err := d.in.ReadString('\n')
			if err != nil {
				return handleEOF(nil, err)
			}
			if s == "POST /api/v1/spans\n" {
				break
			}
		}

		// Read headers
		r.Headers = http.Header{}
		for {
			s, err := d.in.ReadString('\n')
			if err != nil {
				return handleEOF(nil, err)
			}
			if s == "\n" {
				break
			}
			h := strings.Split(strings.TrimSpace(s), ": ")
			r.Headers.Add(h[0], h[1])
		}

		// Read data until we find a "---". This should be our Zipkin payload
		payload := ""
		for {
			s, err := d.in.ReadString('\n')
			if err != nil {
				return handleEOF(nil, err)
			}
			if s == "---\n" {
				break
			}
			payload += s
		}

		// Parse the span
		r.Span = make(model.SpanList, 0)
		err := json.Unmarshal([]byte(payload), &r.Span)
		if err == nil {
			return &r, nil
		} else {
			//log.Printf("Skipped record due to error: %s", err)
		}
	}
}

func handleEOF(r *Record, err error) (*Record, error) {
	if err == io.EOF {
		return nil, nil
	}
	return nil, err
}

