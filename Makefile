.PHONY: run
run:
	go run cmd/auth/main.go

.PHONY: lint
lint:
	golangci-lint run --timeout 5m --config .golangci.yaml -v ./...

.PHONY: build
build:  .build

.build:
	go mod download && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
		-v -o ./bin/auth-service$(shell go env GOEXE) ./cmd/auth/main.go

.PHONY: generate
generate: .generate-install-buf .generate-go .generate-finalize-go

.generate-install-buf:
	@ command -v buf 2>&1 > /dev/null || (echo "Install buf" && \
    		curl -sSL0 https://github.com/bufbuild/buf/releases/download/$(BUF_VERSION)/buf-$(OS_NAME)-$(OS_ARCH)$(shell go env GOEXE) --create-dirs -o "$(BUF_EXE)" && \
    		chmod +x "$(BUF_EXE)")

.generate-go:
	$(BUF_EXE) generate

.generate-finalize-go:
	mv pkg/$(SERVICE_NAME)/gitlab.com/$(SERVICE_PATH)/$(SERVICE_NAME)/* pkg/$(SERVICE_NAME)
	rm -rf pkg/$(SERVICE_NAME)/gitlab.com/
	cd pkg/$(SERVICE_NAME) && ls go.mod || (go mod init gitlab.com/$(SERVICE_PATH)/pkg/$(SERVICE_NAME) && go mod tidy)

# ----------------------------------------------------------------

.PHONY: deps
deps: deps-go

.PHONY: deps-go
deps-go:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.27.1
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1.0
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.5.0
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.5.0
	go install github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger@latest

.PHONY: docs
docs:
	swag init --parseDependency --parseInternal --dir ./internal/adapters/rest --generalInfo api.go --output ./api/swagger/public --parseDepth 4
