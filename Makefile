PROTO_OUT=sdk/go/proto
BUF_VERSION=1.32.1
BUF_BIN=bin/buf

all: generate

generate: $(BUF_BIN)
	$(BUF_BIN) generate

$(BUF_BIN):
	@mkdir -p bin
	@echo "Installing buf $(BUF_VERSION) locally..."
	@curl -sSL "https://github.com/bufbuild/buf/releases/download/v$(BUF_VERSION)/buf-$$(uname -s)-$$(uname -m)" -o $(BUF_BIN)
	@chmod +x $(BUF_BIN)
	@echo "buf installed to $(BUF_BIN)"

tidy:
	go mod tidy

clean:
	rm -rf $(PROTO_OUT)

clean-all: clean
	rm -rf bin

test:
	go test ./...

buf: $(BUF_BIN)

buf-lint: $(BUF_BIN)
	$(BUF_BIN) lint

lint: $(BUF_BIN)
	golangci-lint run
	$(BUF_BIN) lint

validate: test lint

depend:
	@echo "Installing Go development tools..."
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@v2.3.1
	@go install golang.org/x/tools/cmd/goimports@latest
	@go install github.com/fatih/gomodifytags@latest
	@go install github.com/josharian/impl@latest
	@go install github.com/cweill/gotests/gotests@latest
	@go install github.com/golang/mock/mockgen@latest
	@go install github.com/axw/gocov/gocov@latest
	@go install github.com/AlekSi/gocov-xml@latest
	@go install github.com/tebeka/go2xunit@latest
	@echo "Go development tools installed successfully!"