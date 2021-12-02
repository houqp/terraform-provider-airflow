package main

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAcc_Variable(t *testing.T) {
	value1 := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	value2 := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	key := "test_key"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: providers(),
		Steps: []resource.TestStep{
			{
				Config: variableConfig(value1, key),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("airflow_variable.test", "value", value1),
					resource.TestCheckResourceAttr("airflow_variable.test", "key", key),
				),
			},
			{
				Config: variableConfig(value2, key),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("airflow_variable.test", "value", value2),
					resource.TestCheckResourceAttr("airflow_variable.test", "key", key),
				),
			},
		},
	})
}

func variableConfig(value, key string) string {
	return fmt.Sprintf(`
	resource "airflow_variable" "test" {
		value = "%s"
		key = "%s"
	}
	`, value, key)
}
