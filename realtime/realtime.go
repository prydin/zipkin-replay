package realtime

import (
	"bufio"
	"github.com/prydin/zipkin-replay/reader"
	"github.com/prydin/zipkin-replay/sender"
	"github.com/prydin/zipkin-replay/transform"
	"io"
	"time"
)

type Realtime struct {
	stop chan struct{}
	rdr *bufio.Reader
	stopped bool
}

func NewRealtime(rdr *bufio.Reader) (*Realtime, error) {
	return &Realtime{
		stop: make(chan struct{}),
		rdr: rdr,
	}, nil
}

func (rt *Realtime) Stop() {
	close(rt.stop)
}

func (rt *Realtime) RunForever(target, namespace string) error {
	for !rt.stopped {
		err := rt.Run(target, namespace)
		if err != io.EOF {
			return err
		}
	}
	return nil
}

func (rt *Realtime) Run(target, namespace string) error {
	snd := sender.NewHTTPSender(target)
	tr, err := transform.NewTransformer()
	if err != nil {
		return err
	}
	dr := reader.NewDumpReader(rt.rdr)
	for {
		r, err := dr.Next()
		if err != nil {
			return err
		}
		if r == nil {
			return nil
		}
		if len(r.Span) == 0 {
			continue // Shouldn't happen, but handle it gracefully
		}
		tr.Transform(&r.Span, namespace, 0)
		d := time.Duration((r.Span[0].Timestamp * 1000) - time.Now().UnixNano())
		if d < 0 {
			d = 0
		}
		tick := time.After(d)
		select {
		case <-rt.stop:
			rt.stopped = true
			return nil
		case <-tick:
		}
		err = snd.Send(r)
		if err != nil {
			return err
		}
	}
	return nil
}

