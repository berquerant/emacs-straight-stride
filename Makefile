GOMOD = go mod
GOBUILD = go build -trimpath
GOTEST = go test -cover -race

BIN = dist/stride
CMD = "./cmd/stride"

.PHONY: $(BIN)
$(BIN):
	$(GOBUILD) -o $@ $(CMD)

.PHONY: test
test:
	$(GOTEST) ./...

.PHONY: init
init:
	$(GOMOD) tidy

.PHONY: vuln
vuln:
	go tool govulncheck ./...

.PHONY: vet
vet:
	go vet ./...
