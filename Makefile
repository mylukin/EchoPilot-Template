IMAGE_NAME := EchoPilot/app-api
IMAGE_TAG := $(shell date +%Y%m%d)
GOPATH=$(shell go env GOPATH)
APP_NAME={APP_NAME}
GIN_PORT=$(shell if [ -f .port ]; then cat .port; else echo 3000; fi)
APP_PORT=$$(( $(GIN_PORT) + 1 ))

.PHONY: run
run: install-deps
	@$(GOPATH)/bin/gin --port=$(GIN_PORT) --appPort=$(APP_PORT) --bin='app-bin' --immediate --buildArgs='-v -x -mod=readonly -buildvcs=false' run main.go  

.PHONY: install-deps
install-deps:
	@ls $(GOPATH)/bin/gin > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		echo "install gin ..."; \
		go install -mod=mod github.com/codegangsta/gin; \
	fi; \

	@ls $(GOPATH)/bin/codetool > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		echo "install codetool ..."; \
		go install -mod=mod github.com/mylukin/EchoPilot/codetool; \
	fi; \
	
	@ls $(GOPATH)/bin/easyi18n > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		echo "install easyi18n ..."; \
		go install -mod=mod github.com/mylukin/easy-i18n/easyi18n; \
	fi; \

.PHONY: generate
generate: install-deps
	@export PATH="$(GOPATH)/bin:$(PATH)"; \
	 go mod tidy; \
	 go mod vendor; \
	 go generate

.PHONY: app-build
app-build:
	@GOOS=linux GOARCH=amd64 go build -v -mod=readonly -buildvcs=false -o ./docker/app-bin; \
	 chmod a+x ./docker/app-bin

.PHONY: install
install:
	@go build -v -mod=readonly -buildvcs=false -o ./$(APP_NAME); \
	 chmod a+x ./$(APP_NAME) && mv ./$(APP_NAME) $(GOPATH)/bin/

.PHONY: build
build: app-build
	docker build --platform=linux/amd64 -t $(IMAGE_NAME):$(IMAGE_TAG) ./docker
	docker tag $(IMAGE_NAME):$(IMAGE_TAG) $(IMAGE_NAME):latest

.PHONY: publish
publish:
	docker push $(IMAGE_NAME):$(IMAGE_TAG)
	docker push $(IMAGE_NAME):latest

.PHONY: release
release: build publish