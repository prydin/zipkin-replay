package main

import (
	"bufio"
	"flag"
	"log"
	"github.com/prydin/wf-replay/reader"
	"github.com/prydin/wf-replay/sender"
	"github.com/prydin/wf-replay/transform"
	"os"
	"time"
)

func main() {
	targetPtr := flag.String("target", "", "The target URL")
	filePtr := flag.String("file", "", "Input file captured with usurp")
	nsPtr := flag.String("namespace", "default", "Namespace substitution")
	offsetPtr := flag.Int("offset", 0, "Time offset in minutes")
	flag.Parse()

	if *filePtr == "" || *targetPtr == "" {
		log.Fatal("Both -file and -target must be specified")
	}

	f, err := os.Open(*filePtr)
	if err != nil {
		log.Fatalf("Could not open input file: %s\n", err)
	}
	s := sender.NewHTTPSender(*targetPtr)
	rdr := bufio.NewReader(f)
	dr := reader.NewDumpReader(rdr)
	tr, err := transform.NewTransformer()
	if err != nil {
		log.Fatalf("Could not open input file: %s\n", err)
	}

	for {
		r, err := dr.Next()
		if err != nil {
			log.Fatalf("Error parsing input: %s\n", err)
		}
		if r == nil {
			break
		}
		tr.Transform(&r.Span, *nsPtr, time.Duration(*offsetPtr) * time.Minute)
		if err := s.Send(r); err != nil {
			log.Fatalf("Error posting traces: %s\n", err)
		}
	}
}
