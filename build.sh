
export GOPATH=~/workspace/GolangHelper:$GOPATH
go build -ldflags "-s -w" -o "bin/GolangHelper" "src/gentwolf/main.go"
