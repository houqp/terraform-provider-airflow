package main

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAirflowDagRun_basic(t *testing.T) {
	dagId := "example_bash_operator"

	resourceName := "airflow_dag_run.test"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAirflowDagRunCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAirflowDagRunConfigBasic(dagId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "dag_id", dagId),
					resource.TestCheckResourceAttrSet(resourceName, "dag_run_id"),
					resource.TestCheckResourceAttr(resourceName, "conf.%", "0"),
					resource.TestCheckResourceAttrSet(resourceName, "state"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAirflowDagRun_dagRunId(t *testing.T) {
	dagRunId := acctest.RandomWithPrefix("tf-acc-test")
	dagId := "example_bash_operator"

	resourceName := "airflow_dag_run.test"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAirflowDagRunCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAirflowDagRunConfigRunId(dagId, dagRunId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "dag_id", dagId),
					resource.TestCheckResourceAttr(resourceName, "dag_run_id", dagRunId),
					resource.TestCheckResourceAttr(resourceName, "conf.%", "0"),
					resource.TestCheckResourceAttrSet(resourceName, "state"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAirflowDagRun_conf(t *testing.T) {
	dagId := "example_bash_operator"

	resourceName := "airflow_dag_run.test"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAirflowDagRunCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAirflowDagRunConfigConf(dagId),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "dag_id", dagId),
					resource.TestCheckResourceAttrSet(resourceName, "dag_run_id"),
					resource.TestCheckResourceAttr(resourceName, "conf.%", "1"),
					resource.TestCheckResourceAttr(resourceName, fmt.Sprintf("conf.%s", dagId), dagId),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckAirflowDagRunCheckDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(ProviderConfig)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "airflow_dag_run" {
			continue
		}

		dagId, dagRunId, err := airflowDagRunId(rs.Primary.ID)
		if err != nil {
			return err
		}

		dagRun, res, err := client.ApiClient.DAGRunApi.GetDagRun(client.AuthContext, dagId, dagRunId).Execute()
		if err == nil {
			if *dagRun.DagRunId.Get() == dagRunId {
				return fmt.Errorf("Airflow DagRun (%s) still exists.", rs.Primary.ID)
			}
		}

		if res != nil && res.StatusCode == 404 {
			continue
		}
	}

	return nil
}

func testAccAirflowDagRunConfigBasic(dagId string) string {
	return fmt.Sprintf(`
resource "airflow_dag_run" "test" {
  dag_id     = %[1]q
}
`, dagId)
}

func testAccAirflowDagRunConfigRunId(dagId, dagRunId string) string {
	return fmt.Sprintf(`
resource "airflow_dag_run" "test" {
  dag_id     = %[1]q
  dag_run_id = %[2]q
}
`, dagId, dagRunId)
}

func testAccAirflowDagRunConfigConf(dagId string) string {
	return fmt.Sprintf(`
resource "airflow_dag_run" "test" {
  dag_id     = %[1]q

  conf = {
    %[1]q = %[1]q
  }  
}
`, dagId)
}
