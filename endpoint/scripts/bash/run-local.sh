#!/bin/bash
app_id=endpoint
dapr_port=3502

if [[ ! -d components ]]; then
  mkdir components
fi

cat > ./components/pubsub.yaml <<EOF
apiVersion: dapr.io/v1alpha1
kind: Component
metadata:
  name: reports
spec:
  type: pubsub.redis
  version: v1
  metadata:
  - name: redisHost
    value: localhost:6379
  - name: redisPassword
    value: ""

EOF

dapr run \
  --app-id $app_id \
  --app-protocol grpc \
  --dapr-grpc-port $dapr_port \
  --components-path ./components -- go run main.go

rm -rf ./components
