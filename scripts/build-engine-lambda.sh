rm -f build/engine-handler*

GOOS=linux go build -o ./build/engine-handler ./cmd/lambda/engine/main.go
pushd build
zip engine-handler.zip engine-handler
popd
