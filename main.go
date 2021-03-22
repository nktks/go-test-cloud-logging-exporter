package main

import (
	"bufio"
	"context"
	"encoding/json"
	"flag"
	"log"
	"os"
	"sort"

	"cloud.google.com/go/logging"
	"github.com/google/uuid"
)

var (
	pid  = flag.String("p", "", "gcp project id")
	tid  = flag.String("id", uuid.New().String(), "test id for log attribute")
	name = flag.String("name", "go-test-log", "logName for Cloud Logging")
	top  = flag.Int("top", 50, "logging target top number sorted by elapsed")
)

type Payload struct {
	Time    string
	Action  string
	Package string
	Test    string
	Elapsed float64
	TestID  string
}

func main() {
	flag.Parse()
	projectID := fixProjectID(*pid)
	ctx := context.Background()
	client, err := logging.NewClient(ctx, projectID)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	logger := client.Logger(*name)
	scanner := bufio.NewScanner(os.Stdin)
	tests := []*Payload{}
	packages := []*Payload{}
	for scanner.Scan() {
		v := &Payload{}
		err := json.Unmarshal(scanner.Bytes(), v)
		if err != nil {
			log.Fatal("json unmarshal failed. %#v", err)
		}
		if v.Action != "pass" && v.Action != "fail" {
			continue
		}
		v.TestID = *tid
		if v.Test == "" {
			packages = append(packages, v)
		} else {
			tests = append(tests, v)
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal("reading standard input failed. %#v", err)
	}
	sort.SliceStable(tests, func(i, j int) bool { return tests[i].Elapsed > tests[j].Elapsed })
	for i, v := range tests {
		if i >= *top {
			return
		}
		logger.Log(logging.Entry{Payload: v})
	}
	for _, v := range packages {
		logger.Log(logging.Entry{Payload: v})
	}

}

func fixProjectID(i string) string {
	if i == "" {
		e := os.Getenv("GOOGLE_CLOUD_PROJECT")
		if e == "" {
			log.Fatal("need -p flag or GOOGLE_CLOUD_PROJECT env for specify project.")
			return ""
		}
		return e
	}
	return i
}
