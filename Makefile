GOPATH=$(shell go env GOPATH)

run: install-deps
	@$(GOPATH)/bin/gin --appPort=80 --bin='app-bin' --immediate --buildArgs='-v -x -mod=vendor -buildvcs=false' run main.go  

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

generate: install-deps	
	@export PATH="$(GOPATH)/bin:$(PATH)"; \
	 go mod tidy; \
	 go mod vendor; \
	 go generate

build:
	@go build -v -mod=vendor -buildvcs=false -o ./app-bin; \
	 chmod a+x ./app-bin

install:
	@go build -v -mod=vendor -buildvcs=false -o ./EchoPilot; \
	 chmod a+x ./EchoPilot && mv ./EchoPilot $(GOPATH)/bin/