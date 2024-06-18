#!/bin/sh
echo "checking for protoc"
panic() {
    echo "error: $1"
    exit 1
}
which -s protoc || panic "protoc is not installed"

getProtocGenGo() {
    echo "protoc-gen-go is not installed, installing..."
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest || panic "failed to install protoc-gen-go"
}
echo "checking for protoc-gen-go"
which -s protoc-gen-go || getProtocGenGo

echo "compiling protobuf..."
find ./**/*.proto -type f -exec echo "compiling {}" \; -exec protoc -I=. --go_out=. "{}" \;
echo "done"
