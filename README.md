# go-test-cloud-logging-exporter

Export `go test -json` results to Google Cloud Logging for monitoring.

# How to use.

## Install
```
go get -u github.com/nakatamixi/go-test-cloud-logging-exporter
```

## Usage
```
go-test-cloud-logging-exporter -h
Usage of go-test-cloud-logging-exporter:
  -id string
    	test id for log attribute (default "4542a761-e459-4d51-9b6e-e49794e6945c")
  -name string
    	logName for Cloud Logging (default "go-test-log")
  -p string
    	gcp project id
  -top int
    	logging target top number sorted by elapsed (default 50)
```
## Run
```
go test -json ./... | GOOGLE_APPLICATION_CREDENTIALS=/path/to/cred.json go-test-cloud-logging-exporter -p your-gcp-project
```

You need `roles/logging.logWriter` permission to your account.
