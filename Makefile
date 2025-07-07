.PHONY: build clean test install deps lint format

# 变量定义
BINARY_NAME=xfs-quota-kit
BUILD_DIR=build
CMD_DIR=cmd/xfs-quota-kit
VERSION=$(shell git describe --tags --always --dirty)
LDFLAGS=-ldflags "-X main.version=$(VERSION)"

# 默认目标
all: deps lint test build

# 安装依赖
deps:
	go mod download
	go mod tidy

# 构建
build:
	mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) ./$(CMD_DIR)

# 构建多平台版本
build-all:
	mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 ./$(CMD_DIR)
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 ./$(CMD_DIR)

# 安装到系统
install: build
	sudo cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/

# 运行测试
test:
	go test -v -race -coverprofile=coverage.out ./...

# 代码检查
lint:
	@which golangci-lint > /dev/null || (echo "Installing golangci-lint..." && go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
	golangci-lint run

# 代码格式化
format:
	go fmt ./...
	goimports -w .

# 清理
clean:
	rm -rf $(BUILD_DIR)
	go clean

# 开发模式运行
dev:
	go run ./$(CMD_DIR) --config configs/dev.yaml

# 生成文档
docs:
	@which godoc > /dev/null || (echo "Installing godoc..." && go install golang.org/x/tools/cmd/godoc@latest)
	echo "Visit http://localhost:6060/pkg/github.com/xfs-quota-kit/ for documentation"
	godoc -http=:6060

# Docker相关
docker-build:
	docker build -t $(BINARY_NAME):$(VERSION) .

docker-run:
	docker run --rm -it --privileged -v /mnt:/mnt $(BINARY_NAME):$(VERSION) 