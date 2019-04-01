
GOOS=linux go build ./cmd/lambda/worker/main.go ./build/worker-handler
zip ./build/worker-handler.zip ./cmd/lambda/worker/worker-handler ./certs/*

aws lambda create-function \
  --region us-east-1 \
  --function-name iot_simulator_worker \
  --memory 2048 \
  --timeout 300 \
  --role arn:aws:iam::068311527115:role/service-role/iot_simulator_worker-role-xmccwpnk \
  --runtime go1.x \
  --zip-file fileb://build/worker-handler.zip \
  --handler worker-handler \
  --environment Variables={MAX_CONCURRENT_CLIENTS=1000}
