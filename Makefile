build:
	go build -o terraform-provider-airflow

docs:
	go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs
