#!/bin/bash

resource_group=$RESOURCE_GROUP_NAME
environment=$CONTAINER_APP_ENVIRONMENT_NAME
identity_name=$IDENTITY_NAME
registry=$REGISTRY_FQDN
servicebus_namespace=$SERVICEBUS_NAMESPACE

endpoint_name="endpoint"
endpoint_version=1.0.0
endpoint_port=3000
endpoint_api_keys=""

worker_name="worker"
worker_version=1.0.0
worker_port=3001
worker_scale_rule_type="queue"

for arg in "$@"
do
  case $arg in
    --resource-group)
      shift
      resource_group=$1
      shift
      ;;
    --environment)
      shift
      environment=$1
      shift
      ;;
    --identity)
      shift
      identity_name=$1
      shift
      ;;
    --registry)
      shift
      registry=$1
      shift
      ;;
    --servicebus-namespace)
      shift
      servicebus_namespace=$1
      shift
      ;;
    --servicebus-namespace-authorization-rule)
      shift
      servicebus_namespace_authorization_rule=$1
      shift
      ;;
    --worker-scale-rule-type)
      shift
      worker_scale_rule_type=$1
      shift
      ;;
    --endpoint-version)
      shift
      endpoint_version=$1
      shift
      ;;
    --worker-version)
      shift
      worker_version=$1
      shift
      ;;
    --endpoint-api-keys)
      shift
      endpoint_api_keys=$1
      shift
      ;;
  esac
done

if [[ -z "$resource_group" ]]; then
  echo "A resource group name must be provided."
  exit 1
fi

if [[ -z "$environment" ]]; then
  echo "A container app environment name must be provided."
  exit 1
fi

if [[ -z "$identity_name" ]]; then
  echo "An identity name must be provided."
  exit 1
fi

if [[ -z "$registry" ]]; then
  echo "A registry FQDN must be provided."
  exit 1
fi

if [[ -z "$servicebus_namespace" ]]; then
  echo "A servicebus namespace name must be provided."
fi

endpoint_image=$registry/$endpoint_name:$endpoint_version
worker_image=$registry/$worker_name:$worker_version

# Get identity and set client ID and resource ID.
identity=$(az identity show --resource-group $resource_group --name $identity_name)
identity_resource_id=$(echo $identity | jq -r .id)
identity_client_id=$(echo $identity | jq -r .clientId)

worker_scale_rule=(
   "--scale-rule-type azure-servicebus"
   "--scale-rule-auth connection=servicebus-connection-string"
)

case $worker_scale_rule_type in
  "queue")
    worker_scale_rule+=("--scale-rule-name queue-scale-rule" "--scale-rule-metadata namespace=$servicebus_namespace queueName=create messageCount=1")
    servicebus_namespace_authorization_rule="queue-scaling"
    ;;
  "topic")
    worker_scale_rule+=("--scale-rule-name topic-scale-rule" "--scale-rule-metadata namespace=$servicebus_namespace subscriptionName=$worker_name topicName=create")
    servicebus_namespace_authorization_rule="topic-scaling"
    ;;
  *)
    echo "Invalid worker scale rule type: $worker_scale_rule_type"
    exit 1
    ;;
esac

# Get servicebus primary connection string.
servicebus_connection_string=$(az servicebus namespace authorization-rule keys list \
  --resource-group $resource_group \
  --namespace-name $servicebus_namespace \
  --name $servicebus_namespace_authorization_rule \
  --out tsv --query primaryConnectionString)

# Deploy endpoint container app.
az containerapp create \
  --resource-group $resource_group \
  --environment $environment \
  --name $endpoint_name \
  --container-name $endpoint_name \
  --user-assigned $identity_resource_id \
  --enable-dapr \
  --dal \
  --dapr-app-id $endpoint_name \
  --dapr-app-port $endpoint_port \
  --dapr-app-protocol grpc \
  --registry-server $registry \
  --registry-identity $identity_resource_id \
  --image $endpoint_image \
  --cpu 0.25 \
  --memory 0.5Gi \
  --min-replicas 0 \
  --max-replicas 3 \
  --ingress external \
  --target-port $endpoint_port \
  --secrets \
      endpoint-security-keys=$endpoint_api_keys \
  --env-vars \
      ENDPOINT_SECURITY_KEYS=secretref:endpoint-security-keys \
      DAPR_CLIENT_TIMEOUT_SECONDS=15 \
  --scale-rule-name http-scale-rule \
  --scale-rule-http-concurrency 50

# Deploy worker container app.
az containerapp create \
  --resource-group $resource_group \
  --environment $environment \
  --name $worker_name \
  --container-name $worker_name \
  --user-assigned $identity_resource_id \
  --enable-dapr \
  --dal \
  --dapr-app-id $worker_name \
  --dapr-app-port $worker_port \
  --dapr-app-protocol grpc \
  --registry-server $registry \
  --registry-identity $identity_resource_id \
  --image $worker_image \
  --cpu 0.25 \
  --memory 0.5Gi \
  --min-replicas 0 \
  --max-replicas 3 \
  --secrets \
      servicebus-connection-string=$servicebus_connection_string \
  --env-vars \
      DAPR_CLIENT_TIMEOUT_SECONDS=15 \
   ${worker_scale_rule[@]}
