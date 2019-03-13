deps:
	go get -u golang.org/x/lint/golint

# Lint the go code. Note: golint doesn't support vendor folder exclusion so we use find to filter it out
lint: deps
	@echo "Using vet to check for common mistakes..."
	@go vet ./...
	@echo "Checking style with golint..."
	@find . -type d -not -path "./vendor*" -exec golint {} \;

vendor:
	dep ensure -vendor-only

# Build in the real Google cloud
cb:
	gcloud builds submit --config cloudbuild.yaml .

# Build  using your local environment
cbl:
	cloud-build-local --config cloudbuild.yaml --dryrun=false .
