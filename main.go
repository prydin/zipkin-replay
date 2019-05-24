package main

import (
	"bufio"
	"compress/gzip"
	"flag"
	"github.com/prydin/zipkin-replay/reader"
	"github.com/prydin/zipkin-replay/realtime"
	"github.com/prydin/zipkin-replay/sender"
	"github.com/prydin/zipkin-replay/transform"
	"log"
	"os"
)

func main() {
	targetPtr := flag.String("target", "", "The target URL")
	filePtr := flag.String("file", "", "Input file captured with usurp")
	nsPtr := flag.String("namespace", "default", "Namespace substitution")
	offsetPtr := flag.Duration("offset", 0, "Time offset in minutes")
	rtPtr := flag.Bool("realtime", false, "Run realtime with delays according to original timing")
	fePtr := flag.Bool("forever", false, "Run realtime with delays according to original timing")
	gzPtr := flag.Bool("gz", false, "Input is gzipped")

	flag.Parse()

	if *filePtr == "" || *targetPtr == "" {
		log.Fatal("Both -file and -target must be specified")
	}

	f, err := os.Open(*filePtr)
	if err != nil {
		log.Fatalf("Could not open input file: %s\n", err)
	}
	defer f.Close()

	var rdr *bufio.Reader
	if *gzPtr {
		gzr, err := gzip.NewReader(f)
		if err != nil {
			log.Fatalf("Error initializing gzip: %s\n", err)
		}
		rdr = bufio.NewReader(gzr)
	} else {
		bufio.NewReader(f)
	}

	if *rtPtr {
		rt, err := realtime.NewRealtime(rdr)
		if err != nil {
			log.Fatal(err)
		}
		if *fePtr {
			rt.RunForever(*targetPtr, *nsPtr)
		} else {
			rt.Run(*targetPtr, *nsPtr)
		}
		rt.Run(*targetPtr, *nsPtr)
	} else {
		s := sender.NewHTTPSender(*targetPtr)
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
			tr.Transform(&r.Span, *nsPtr, *offsetPtr)
			if err := s.Send(r); err != nil {
				log.Fatalf("Error posting traces: %s\n", err)
			}
		}
	}
}
