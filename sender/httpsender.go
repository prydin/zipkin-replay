package sender

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/prydin/wf-replay/reader"
	"net/http"
)

type HTTPSender struct {
	url string
	client http.Client
}

func NewHTTPSender(url string) *HTTPSender{
	return &HTTPSender{ url: url }
}

func (h* HTTPSender) Send(r *reader.Record) error {
	buf, err := json.Marshal(r.Span)
	if err != nil {
		return err
	}
	rq, err := http.NewRequest("POST", h.url, bytes.NewReader(buf))
	if err != nil {
		return err
	}
	rq.Header = r.Headers
	resp, err := h.client.Do(rq)
	if err != nil {
		return err
	}
	if resp.StatusCode > 299 {
		return fmt.Errorf("server returned error code: %d %s", resp.StatusCode, resp.Status)
	}
	return nil
}