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
	if v := os.Getenv("AIRFLOW_OAUTH2_TOKEN"); v == "" {
		t.Fatal("AIRFLOW_OAUTH2_TOKEN must be set for acceptance tests")
	}
	if v := os.Getenv("AIRFLOW_BASE_ENDPOINT"); v == "" {
		t.Fatal("AIRFLOW_BASE_ENDPOINT must be set for acceptance tests")
	}
}
