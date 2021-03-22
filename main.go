package main

import (
	"bufio"
	"context"
	"encoding/json"
	"flag"
	"log"
	"os"

	"cloud.google.com/go/logging"
	"github.com/google/uuid"
)

var (
	pid  = flag.String("p", "", "gcp project id")
	tid  = flag.String("id", uuid.New().String(), "test id for log attribute")
	name = flag.String("name", "go-test-log", "logName for Cloud Logging")
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
		logger.Log(logging.Entry{Payload: v})
	}
	if err := scanner.Err(); err != nil {
		log.Fatal("reading standard input failed. %#v", err)
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
