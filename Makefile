APPNAME = kplcloud
BIN = $(GOPATH)/bin
GOCMD = go
GOBUILD = $(GOCMD) build
GORUN = $(GOCMD) run
BINARY_UNIX = $(BIN)/$(APPNAME)
GOPROXY = https://goproxy.cn
PID = .pid

start:
	./$(APPNAME) & echo $$! > $(PID) 2>1 &

restart:
	@echo restart the app...
	@kill `cat $(PID)` || true
	./$(APPNAME) ./app.cfg & echo $$! > $(PID) 2>1 &

install:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOINSTALL) -v

stop:
	@kill `cat $(PID)` || true

build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) ./main.go -o $(BINARY_UNIX) -v

run:
	GOPROXY=$(GOPROXY) GO111MODULE=on $(GORUN) ./cmd/main.go start -p :8080 -c app.dev.cfg -a dev.$(APPNAME)