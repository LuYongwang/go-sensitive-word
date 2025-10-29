.PHONY: lint lint-fix test build clean help

# 默认目标
.DEFAULT_GOAL := help

# Go 相关变量
GOCMD = go
GOBUILD = $(GOCMD) build
GOTEST = $(GOCMD) test
GOGET = $(GOCMD) get

# golangci-lint 配置
GOLANGCI_LINT = golangci-lint
LINT_ENABLED = errcheck,govet,staticcheck,unused,misspell,ineffassign
LINT_TIMEOUT = 5m

help: ## 显示帮助信息
	@echo "可用命令:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

lint: ## 运行 golangci-lint 代码检查
	@echo "🔍 运行代码检查..."
	@$(GOLANGCI_LINT) run --no-config \
		--enable=$(LINT_ENABLED) \
		--timeout=$(LINT_TIMEOUT) \
		./...

lint-fix: ## 自动修复可修复的 lint 问题
	@echo "🔧 自动修复代码问题..."
	@$(GOLANGCI_LINT) run --no-config \
		--enable=$(LINT_ENABLED) \
		--timeout=$(LINT_TIMEOUT) \
		--fix \
		./...

test: ## 运行测试
	@echo "🧪 运行测试..."
	@$(GOTEST) -v ./...

test-coverage: ## 运行测试并生成覆盖率报告
	@echo "📊 生成测试覆盖率报告..."
	@$(GOTEST) -v -coverprofile=coverage.out ./...
	@$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "✅ 覆盖率报告已生成: coverage.html"

build: ## 构建项目
	@echo "🔨 构建项目..."
	@$(GOBUILD) -v ./...

build-examples: ## 构建示例程序
	@echo "🔨 构建示例程序..."
	@$(GOBUILD) -o examples/demo examples/demo.go
	@$(GOBUILD) -o examples/dynamic/main examples/dynamic/main.go

clean: ## 清理构建产物
	@echo "🧹 清理构建产物..."
	@rm -f examples/demo examples/dynamic/main
	@rm -f coverage.out coverage.html
	@$(GOCMD) clean

install-lint: ## 安装 golangci-lint
	@echo "📦 安装 golangci-lint..."
	@$(GOCMD) install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest
	@echo "✅ golangci-lint 安装完成"

check: lint test ## 运行代码检查和测试（CI 使用）

ci: check ## CI 流程（lint + test）

all: lint test build ## 运行所有检查、测试并构建

