#!/bin/bash
app_id=worker
app_port=3001
dapr_port=3501

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

cat > ./components/bindings.yaml <<EOF
apiVersion: dapr.io/v1alpha1
kind: Component
metadata:
  name: reports-output
spec:
  type: bindings.redis
  version: v1
  metadata:
  - name: redisHost
    value: localhost:6379
  - name: redisPassword
    value: ""

EOF

dapr run \
  --app-id $app_id \
  --app-port $app_port \
  --app-protocol grpc \
  --dapr-grpc-port $dapr_port \
  --components-path ./components -- go run main.go

rm -rf ./components
