build:
	go build -o terraform-provider-airflow

docs:
	go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs

vendor:
	go mod vendor
	go mod tidy

test-acceptance:
	@AIRFLOW_API_USERNAME=airflow \
	AIRFLOW_API_PASSWORD=airflow \
	AIRFLOW_BASE_ENDPOINT=http://127.0.0.1:8080 \
	TF_ACC=1 $(go_test) go test  -v -cover $(shell go list ./... | grep -v vendor)

start-airflow:
	docker run -d -e AIRFLOW__API__AUTH_BACKEND=airflow.api.auth.backend.basic_auth \
	-p 127.0.0.1:8080:8080 --name=airflow apache/airflow:2.2.2 standalone

stop-airflow:
	docker stop airflow; docker rm airflow
