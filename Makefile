TEST?=$$(go list ./...)

build:
	go build -o terraform-provider-airflow

testacc:
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 5m