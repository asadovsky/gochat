export GOPATH=${HOME}/dev/gochat

${GOPATH}/tools/lint.sh

cd src/main
go run main.go

go get code.google.com/p/go.net/websocket
