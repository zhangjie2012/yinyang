TARGET=./bin/yinyang

VERSION=`git rev-parse HEAD`
BUILD_TIME=`date +'%Y-%m-%d %H:%M:%S'`
BRANCH=`git rev-parse --abbrev-ref HEAD`
LDFLAGS=-ldflags "-X 'main.Branch=${BRANCH}' -X 'main.Version=${VERSION}' -X 'main.BuildTime=${BUILD_TIME}'"

build-linux:
	GOOS=linux GOARCH=amd64 go build ${LDFLAGS} -o ${TARGET} ./cmd/main.go

build-local:
	go build ${LDFLAGS} -o ${TARGET} ./cmd/main.go

run-test:
	go test -count=1 -v ./...
