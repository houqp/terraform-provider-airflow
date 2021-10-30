TEST?=$$(go list ./...)
ACCTEST_PARALLELISM?=20

build:
	go build -o terraform-provider-airflow

testacc:
	TF_ACC=1 go test $(TEST) -parallel $(ACCTEST_PARALLELISM) -v $(TESTARGS) -timeout 5m

test:
	go test	