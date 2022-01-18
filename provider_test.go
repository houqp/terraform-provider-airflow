package main

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var testAccProviders map[string]*schema.Provider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = AirflowProvider()
	testAccProviders = map[string]*schema.Provider{
		"airflow": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := AirflowProvider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ *schema.Provider = AirflowProvider()
}

func testAccPreCheck(t *testing.T) {
	_, tokenOk := os.LookupEnv("AIRFLOW_OAUTH2_TOKEN")
	_, userOk := os.LookupEnv("AIRFLOW_API_USERNAME")
	_, passOk := os.LookupEnv("AIRFLOW_API_PASSWORD")

	if tokenOk && !(userOk || passOk) {
		t.Fatal("AIRFLOW_OAUTH2_TOKEN OR AIRFLOW_API_USERNAME/AIRFLOW_API_PASSWORD must be set for acceptance tests")
	}

	if v := os.Getenv("AIRFLOW_BASE_ENDPOINT"); v == "" {
		t.Fatal("AIRFLOW_BASE_ENDPOINT must be set for acceptance tests")
	}
}
