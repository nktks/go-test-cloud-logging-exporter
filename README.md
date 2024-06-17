# go-test-cloud-logging-exporter

Export `go test -json` results to Google Cloud Logging for monitoring.

# How to use.

## Install
```
go get -u github.com/nktks/go-test-cloud-logging-exporter
```

## Usage
```
go-test-cloud-logging-exporter -h
Usage of go-test-cloud-logging-exporter:
  -id string
    	test id for log attribute (default "fbcabe2b-5196-475d-afdb-eb726165ad3a")
  -junitxml string
    	gotestsum --junitxml file path
  -name string
    	logName for Cloud Logging (default "go-test-log")
  -p string
    	gcp project id
  -top int
    	logging target top number sorted by elapsed (default 50)
```
## Run
by `go test -json`
```
go test -json ./... | GOOGLE_APPLICATION_CREDENTIALS=/path/to/cred.json go-test-cloud-logging-exporter -p your-gcp-project
```
by `gotestsum --junitxml=junit.xml`
```
gotestsum --junitfile=junit.xml ./...
GOOGLE_APPLICATION_CREDENTIALS=/path/to/cred.json go-test-cloud-logging-exporter -p your-gcp-project -junitxml junit.xml
```

You need `roles/logging.logWriter` permission to your account.
