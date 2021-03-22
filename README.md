# go-test-cloud-logging-exporter

Export `go test -json` results to Google Cloud Logging for monitoring.

# How to use.

install
```
go get -u github.com/nakatamixi/go-test-cloud-logging-exporter
```

run
```
go test -json ./... | GOOGLE_APPLICATION_CREDENTIALS=/path/to/cred.json go-test-cloud-logging-exporter -p your-gcp-project
```

You need `roles/logging.logWriter` permission to your account.
