#!/bin/bash
# Get the command parameter, default to "start"
command=${1:-start}

# Determine the project root directory (parent directory of the script)
PROJECT_ROOT="$( cd "$( dirname "${BASH_SOURCE[0]}" )/.." && pwd )"

# Set environment variables
export HOSPITAL_API_ENVIRONMENT="Development"
export HOSPITAL_API_PORT="8080"

export HOSPITAL_API_MONGODB_USERNAME="root"
export HOSPITAL_API_MONGODB_PASSWORD="root"

# Define a helper function for docker compose using our compose file
mongo() {
    echo "mongo command $@"
    docker compose --file "${PROJECT_ROOT}/deployments/docker-compose/compose.yaml" "$@"
    echo "mongo running with command: $@"
}

# Execute the corresponding command
case "$command" in
  start)
    # "Try" block:
    # Bring up MongoDB containers in detached mode.
    mongo up --detach
    echo "mongo probably started"
    # Set a trap that will run when the script exits (normal or error),
    # simulating a "finally" block to ensure MongoDB is shut down.
    trap 'mongo down' EXIT

    # Run the API service.
    go run "${PROJECT_ROOT}/cmd/hospital-api-service"
    # When this script exits (even due to an error), the trap will execute
    # and shut down the MongoDB containers.
    ;;
  openapi)
    docker run --rm -ti -v "${PROJECT_ROOT}":/local openapitools/openapi-generator-cli generate -c /local/scripts/generator-cfg.yaml
    ;;
  docker)
    docker build -t sempresionus/hospital-wl-api:local-build -f ${PROJECT_ROOT}/build/docker/Dockerfile .
    ;;
  test)
    go test -v ./...
    ;;
  mongo)
    # Shift removes the "mongo" command so the rest of the arguments go to our helper.
    shift
    mongo "$@"
    ;;
  *)
    echo "Unknown command: $command" >&2
    exit 1
    ;;
esac
