# =============================================================================
# {APP_NAME} - Makefile
# =============================================================================

# 变量定义
GOPATH := $(shell go env GOPATH)
APP_NAME := {APP_NAME}
APP_NAME_LOWER := {APP_NAME_LOWER}
DOCKER_IMAGE := ghcr.io/{PACKAGE_NAME}:latest
GIN_PORT := $(shell if [ -f .port ]; then cat .port; else echo 3000; fi)
APP_PORT := $$(( $(GIN_PORT) + 1 ))
DOCKER_COMPOSE := docker compose
PROJECT_DIR := $(shell pwd)

# 默认目标
.DEFAULT_GOAL := help

# =============================================================================
# 帮助信息
# =============================================================================

.PHONY: help
help: ## 显示帮助信息
	@echo "$(APP_NAME) Makefile Commands:"
	@echo ""
	@echo "\033[33m开发相关命令:\033[0m"
	@grep -E '^[a-zA-Z_-]+:.*?## \[dev\]' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## \\[dev\\] "}; {printf "\033[36m  %-18s\033[0m %s\n", $$1, $$2}'
	@echo ""
	@echo "\033[33m构建相关命令:\033[0m"
	@grep -E '^[a-zA-Z_-]+:.*?## \[build\]' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## \\[build\\] "}; {printf "\033[36m  %-18s\033[0m %s\n", $$1, $$2}'
	@echo ""
	@echo "\033[33mDocker 生产环境:\033[0m"
	@grep -E '^[a-zA-Z_-]+:.*?## \[docker\]' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## \\[docker\\] "}; {printf "\033[36m  %-18s\033[0m %s\n", $$1, $$2}'
	@echo ""
	@echo "\033[33mDocker 应用命令:\033[0m"
	@grep -E '^[a-zA-Z_-]+:.*?## \[docker-app\]' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## \\[docker-app\\] "}; {printf "\033[36m  %-18s\033[0m %s\n", $$1, $$2}'
	@echo ""
	@echo "\033[33m常驻进程管理:\033[0m"
	@grep -E '^[a-zA-Z_-]+:.*?## \[workers\]' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## \\[workers\\] "}; {printf "\033[36m  %-18s\033[0m %s\n", $$1, $$2}'
	@echo ""
	@echo "\033[33mCron 管理命令:\033[0m"
	@grep -E '^[a-zA-Z_-]+:.*?## \[cron\]' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## \\[cron\\] "}; {printf "\033[36m  %-18s\033[0m %s\n", $$1, $$2}'
	@echo ""
	@echo "\033[33m工具相关命令:\033[0m"
	@grep -E '^[a-zA-Z_-]+:.*?## \[tool\]' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## \\[tool\\] "}; {printf "\033[36m  %-18s\033[0m %s\n", $$1, $$2}'
	@echo ""
	@echo "\033[33m其他命令:\033[0m"
	@grep -E '^[a-zA-Z_-]+:.*?## \[other\]' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## \\[other\\] "}; {printf "\033[36m  %-18s\033[0m %s\n", $$1, $$2}'
	@echo ""

# =============================================================================
# 开发相关命令
# =============================================================================

.PHONY: run
run: install-deps ## [dev] 运行开发服务器
	@echo "Starting development server on port $(GIN_PORT)..."
	@$(GOPATH)/bin/gin --port=$(GIN_PORT) --appPort=$(APP_PORT) --bin='app-bin' --immediate --buildArgs='-v -x -mod=readonly -buildvcs=false' run main.go

.PHONY: install-deps
install-deps: ## [dev] 安装开发依赖工具
	@echo "Checking and installing dependencies..."
	@ls $(GOPATH)/bin/gin > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		echo "Installing gin..."; \
		go install -mod=mod github.com/codegangsta/gin; \
	fi
	@ls $(GOPATH)/bin/easyi18n > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		echo "Installing easyi18n..."; \
		go install -mod=mod github.com/mylukin/easy-i18n/easyi18n; \
	fi
	@ls $(GOPATH)/bin/translator > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		echo "Installing translator..."; \
		go install -mod=mod github.com/mylukin/translator; \
	fi
	@echo "Dependencies check completed."

.PHONY: generate
generate: install-deps ## [dev] 生成代码和依赖
	@echo "Generating code and tidying modules..."
	@export PATH="$(GOPATH)/bin:$(PATH)"; \
	 go mod tidy; \
	 go mod vendor; \
	 go generate

# =============================================================================
# 构建相关命令
# =============================================================================

.PHONY: app-build
app-build: ## [build] 构建 Linux 应用程序
	@echo "Building application for Linux..."
	@GOOS=linux GOARCH=amd64 go build -v -mod=readonly -buildvcs=false -o ./docker/app-bin
	@chmod a+x ./docker/app-bin
	@echo "Build completed: ./docker/app-bin"

.PHONY: install
install: ## [build] 构建并安装到 GOPATH/bin
	@echo "Building and installing $(APP_NAME)..."
	@go build -v -mod=readonly -buildvcs=false -o ./$(APP_NAME)
	@chmod a+x ./$(APP_NAME) && mv ./$(APP_NAME) $(GOPATH)/bin/
	@echo "$(APP_NAME) installed to $(GOPATH)/bin/"

.PHONY: clean
clean: ## [build] 清理构建文件
	@echo "Cleaning build artifacts..."
	@rm -f ./docker/app-bin
	@rm -f ./$(APP_NAME)
	@rm -rf ./vendor
	@echo "Cleanup completed."

# =============================================================================
# Docker 相关命令
# =============================================================================

.PHONY: docker-pull
docker-pull: ## [docker] 拉取 GitHub 镜像
	@echo "Pulling Docker image from GitHub Container Registry..."
	@docker pull $(DOCKER_IMAGE)
	@echo "Docker image pulled successfully."

.PHONY: docker-network
docker-network: ## [docker] 创建 Docker 网络
	@echo "Creating Docker network if not exists..."
	@docker network ls | grep $(APP_NAME_LOWER)-network >/dev/null 2>&1 || docker network create $(APP_NAME_LOWER)-network
	@echo "Docker network ready."

.PHONY: docker-up
docker-up: docker-pull ## [docker] 启动 API 服务
	@echo "Starting API service..."
	@$(DOCKER_COMPOSE) up -d
	@echo "API service started at http://localhost:3000"

.PHONY: docker-down
docker-down: ## [docker] 停止 API 服务
	@echo "Stopping API service..."
	@$(DOCKER_COMPOSE) down
	@echo "API service stopped."

.PHONY: docker-restart
docker-restart: docker-pull ## [docker] 重启 API 服务
	@echo "Restarting API service..."
	@$(DOCKER_COMPOSE) down
	@$(DOCKER_COMPOSE) up -d
	@echo "API service restarted."

.PHONY: docker-logs
docker-logs: ## [docker-app] 查看 API 日志
	@$(DOCKER_COMPOSE) logs -f -n 1000 api

.PHONY: docker-exec
docker-exec: ## [docker-app] 进入 API 容器
	@$(DOCKER_COMPOSE) exec api bash

.PHONY: docker-run
docker-run: ## [docker-app] 一次性执行命令 (用法: make docker-run CMD=命令，或 make docker-run)
	@if [ -z "$(CMD)" ]; then \
		echo "Running container without command..."; \
		$(DOCKER_COMPOSE) --profile tools run --rm runner bash; \
	else \
		echo "Running command: $(CMD)"; \
		$(DOCKER_COMPOSE) --profile tools run --rm runner /go/app/app-bin $(CMD); \
	fi

.PHONY: docker-test
docker-test: ## [docker-app] 执行 tdotme 爬虫
	@make docker-run CMD=test

.PHONY: docker-clean
docker-clean: ## [docker-app] 清理资源
	@echo "Cleaning Docker resources..."
	@$(DOCKER_COMPOSE) down --remove-orphans
	@echo "Docker cleanup completed."

# =============================================================================
# 常驻进程管理命令
# =============================================================================

.PHONY: workers-up
workers-up: docker-network ## [workers] 启动常驻进程容器
	@echo "Starting workers container..."
	@$(DOCKER_COMPOSE) up -d workers
	@echo "Workers container started successfully."
	@echo "Use 'make workers-status' to check status."

.PHONY: workers-down
workers-down: ## [workers] 停止常驻进程容器
	@echo "Stopping workers container..."
	@$(DOCKER_COMPOSE) stop workers
	@echo "Workers container stopped."

.PHONY: workers-restart
workers-restart: ## [workers] 重启常驻进程容器
	@echo "Restarting workers container..."
	@$(DOCKER_COMPOSE) restart workers
	@echo "Workers container restarted."

.PHONY: workers-logs
workers-logs: ## [workers] 查看常驻进程日志
	@$(DOCKER_COMPOSE) logs -f -n 1000 workers

.PHONY: workers-status
workers-status: ## [workers] 查看常驻进程状态
	@echo "Workers container status:"
	@$(DOCKER_COMPOSE) ps workers

.PHONY: workers-exec
workers-exec: ## [workers] 进入常驻进程容器
	@$(DOCKER_COMPOSE) exec workers bash

# =============================================================================
# Cron 管理命令
# =============================================================================

.PHONY: install-cron
install-cron: ## [cron] 安装/更新定时任务到系统 cron
	@echo "Installing/updating cron jobs to system..."
	@echo "Project directory: $(PROJECT_DIR)"
	@# 创建日志目录
	@sudo mkdir -p /var/log/$(APP_NAME_LOWER)
	@sudo chown $(USER):$(USER) /var/log/$(APP_NAME_LOWER)
	@# 生成实际的 cron.d 文件到项目目录
	@sed -e 's|PROJECT_DIR|$(PROJECT_DIR)|g' -e 's|USER_NAME|$(USER)|g' crontab.system.conf > $(APP_NAME_LOWER).cron
	@# 删除现有的文件（如果存在）
	@sudo rm -f /etc/cron.d/$(APP_NAME_LOWER)
	@# 复制文件
	@sudo cp "$(PROJECT_DIR)/$(APP_NAME_LOWER).cron" /etc/cron.d/$(APP_NAME_LOWER)
	@# 设置正确的文件权限
	@sudo chown root:root /etc/cron.d/$(APP_NAME_LOWER)
	@sudo chmod 644 /etc/cron.d/$(APP_NAME_LOWER)
	@# 重启 cron 服务
	@sudo systemctl restart cron
	@echo "Cron jobs installed/updated successfully!"
	@echo "Configuration: $(PROJECT_DIR)/$(APP_NAME_LOWER).cron -> /etc/cron.d/$(APP_NAME_LOWER)"
	@echo "Log files: /var/log/$(APP_NAME_LOWER)/"
	@echo "Cron service restarted."

.PHONY: logs-cron
logs-cron: ## [cron] 查看定时任务日志
	@echo "$(APP_NAME) cron job logs:"
	@echo ""
	@if [ -d /var/log/$(APP_NAME_LOWER) ]; then \
		for log in /var/log/$(APP_NAME_LOWER)/*.log; do \
			if [ -f "$$log" ]; then \
				echo "=== $$(basename $$log) ==="; \
				tail -20 "$$log" 2>/dev/null || echo "No logs yet"; \
				echo ""; \
			fi; \
		done; \
	else \
		echo "Log directory /var/log/$(APP_NAME_LOWER) not found. Run 'make install-cron' first."; \
	fi

# =============================================================================
# 工具相关命令
# =============================================================================

.PHONY: test
test: ## [tool] 运行测试
	@echo "Running tests..."
	@go test -v ./...

.PHONY: fmt
fmt: ## [tool] 格式化代码
	@echo "Formatting code..."
	@go fmt ./...

.PHONY: vet
vet: ## [tool] 运行 go vet 检查
	@echo "Running go vet..."
	@go vet ./...

.PHONY: mod-tidy
mod-tidy: ## [tool] 整理 Go modules
	@echo "Tidying modules..."
	@go mod tidy

# =============================================================================
# 其他命令
# =============================================================================

.PHONY: info
info: ## [other] 显示项目信息
	@echo "Project Information:"
	@echo "  App Name: $(APP_NAME)"
	@echo "  Gin Port: $(GIN_PORT)"
	@echo "  App Port: $(APP_PORT)"
	@echo "  Go Path: $(GOPATH)"
	@echo "  Project Dir: $(PROJECT_DIR)"
