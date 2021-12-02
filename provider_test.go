package main

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var testAccProviders map[string]*schema.Provider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider()
	testAccProviders = map[string]*schema.Provider{
		"airflow": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func testAccPreCheck(t *testing.T) {
	if err := os.Getenv("AIRFLOW_BASE_ENDPOINT"); err == "" {
		t.Fatal("AIRFLOW_BASE_ENDPOINT must be set for acceptance tests")
	}
	if err := os.Getenv("AIRFLOW_API_USERNAME"); err == "" {
		t.Fatal("AIRFLOW_API_USERNAME must be set for acceptance tests")
	}
	if err := os.Getenv("AIRFLOW_API_PASSWORD"); err == "" {
		t.Fatal("AIRFLOW_API_PASSWORD must be set for acceptance tests")
	}
}
