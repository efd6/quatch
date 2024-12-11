package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"quamina.net/go/quamina"
)

func main() {
	pattern := flag.String("p", "", "JSON pattern to search for (required)")
	file := flag.String("i", "", "file to search (stdin if empty)")
	l := flag.String("l", "", "label for match (incompatible with v)")
	v := flag.String("v", "", "label for non-match (incompatible with l)")
	flag.Parse()
	if *pattern == "" || (*v != "" && *l != "") {
		flag.Usage()
		os.Exit(2)
	}

	q, err := quamina.New()
	if err != nil {
		log.Fatalf("failed to create matcher: %v", err)
	}
	err = q.AddPattern("query", *pattern)
	if err != nil {
		log.Fatalf("failed to add pattern: %v", err)
	}

	f := os.Stdin
	if *file != "" {
		f, err = os.Open(*file)
		if err != nil {
			log.Fatalf("failed to open file: %v", err)
		}
		defer f.Close()
	}
	var buf bytes.Buffer
	_, err = io.Copy(&buf, f)
	if err != nil {
		log.Fatalf("failed to read data: %v", err)
	}

	m, err := q.MatchesForEvent(buf.Bytes())
	if err != nil {
		log.Fatalf("failed match search: %v", err)
	}
	if len(m) == 0 {
		if *v != "" {
			fmt.Println(*v)
		}
		os.Exit(4)
	}
	if *l != "" {
		fmt.Println(*l)
	}
}
