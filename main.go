package main

import (
	"bufio"
	"context"
	"encoding/json"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strconv"

	"cloud.google.com/go/logging"
	"github.com/google/uuid"

	"github.com/nktks/go-test-cloud-logging-exporter/junitxml"
)

var (
	pid      = flag.String("p", "", "gcp project id")
	tid      = flag.String("id", uuid.New().String(), "test id for log attribute")
	name     = flag.String("name", "go-test-log", "logName for Cloud Logging")
	top      = flag.Int("top", 50, "logging target top number sorted by elapsed")
	junitXML = flag.String("junitxml", "", "gotestsum --junitxml file path")
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
	tests := []*Payload{}
	packages := []*Payload{}
	if *junitXML == "" {
		t, p, err := scanGoTestPayloads(os.Stdin, *tid)
		if err != nil {
			log.Fatalf("reading standard input failed. %#v", err)
		}
		tests = t
		packages = p
	} else {
		b, err := ioutil.ReadFile(*junitXML)
		if err != nil {
			log.Fatalf("Failed open file: %v", err)
		}
		j, err := junitxml.Unmarshal(b)
		if err != nil {
			log.Fatalf("Failed parse junit xml: %v", err)
		}
		t, p, err := convertJunitPayloads(j, *tid)
		if err != nil {
			log.Fatalf("Failed convert junit xml to payloads. %#v", err)
		}
		tests = t
		packages = p
	}
	for _, v := range packages {
		logger.Log(logging.Entry{Payload: v})
	}
	sort.SliceStable(tests, func(i, j int) bool { return tests[i].Elapsed > tests[j].Elapsed })
	for i, v := range tests {
		if i >= *top {
			break
		}
		logger.Log(logging.Entry{Payload: v})
	}
	log.Printf("tid %s export completed.", *tid)

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

func scanGoTestPayloads(i io.Reader, tid string) ([]*Payload, []*Payload, error) {
	scanner := bufio.NewScanner(i)
	tests := []*Payload{}
	packages := []*Payload{}
	for scanner.Scan() {
		v := &Payload{}
		err := json.Unmarshal(scanner.Bytes(), v)
		if err != nil {
			log.Printf("json unmarshal failed. %#v", err)
			continue
		}
		if v.Action != "pass" && v.Action != "fail" {
			continue
		}
		v.TestID = tid
		if v.Test == "" {
			packages = append(packages, v)
		} else {
			tests = append(tests, v)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, nil, err
	}
	return tests, packages, nil
}
func convertJunitPayloads(j *junitxml.JUnitTestSuites, tid string) ([]*Payload, []*Payload, error) {
	tests := []*Payload{}
	packages := []*Payload{}
	for _, s := range j.Suites {
		var action string
		if s.Failures > 0 {
			action = "fail"
		} else {
			action = "pass"
		}
		var elapsed float64
		if f, err := strconv.ParseFloat(s.Time, 64); err != nil {
			return nil, nil, err
		} else {
			elapsed = f
		}

		packages = append(packages, &Payload{
			Time:    "",
			Action:  action,
			Package: s.Name,
			Test:    "",
			Elapsed: elapsed,
			TestID:  tid,
		})
		for _, c := range s.TestCases {
			var caction string
			if c.Failure != nil {
				caction = "fail"
			} else {
				caction = "pass"
			}
			var celapsed float64
			if f, err := strconv.ParseFloat(s.Time, 64); err != nil {
				return nil, nil, err
			} else {
				celapsed = f
			}
			tests = append(tests, &Payload{
				Time:    "",
				Action:  caction,
				Package: c.Classname,
				Test:    c.Name,
				Elapsed: celapsed,
				TestID:  tid,
			})
		}
	}
	return tests, packages, nil

}
