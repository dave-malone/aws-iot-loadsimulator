rm -f build/worker-handler*

GOOS=linux go build -o ./build/worker-handler ./cmd/lambda/worker/main.go
pushd build
zip worker-handler.zip worker-handler ../certs/*
popd
