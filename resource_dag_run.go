package main

import (
	"fmt"
	"strings"

	"github.com/apache/airflow-client-go/airflow"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceDagRun() *schema.Resource {
	return &schema.Resource{
		Create: resourceDagRunCreate,
		Read:   resourceDagRunRead,
		Delete: resourceDagRunDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"dag_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"dag_run_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"conf": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceDagRunCreate(d *schema.ResourceData, m interface{}) error {
	pcfg := m.(ProviderConfig)
	client := pcfg.ApiClient.DAGRunApi

	dagId := d.Get("dag_id").(string)
	// dagRunId := d.Get("dag_run_id").(string)

	dagRun := *airflow.NewDAGRunWithDefaults()

	if v, ok := d.GetOk("dag_run_id"); ok {
		dagRun.SetDagRunId(v.(string))
	}

	if v, ok := d.GetOk("conf"); ok {
		dagRun.SetConf(v.(map[string]interface{}))
	}

	res, _, err := client.PostDagRun(pcfg.AuthContext, dagId).DAGRun(dagRun).Execute()
	if err != nil {
		return fmt.Errorf("failed to create dagRunId `%s` from Airflow: %w", dagId, err)
	}
	d.SetId(fmt.Sprintf("%s:%s", dagId, *res.DagRunId.Get()))

	return resourceDagRunRead(d, m)
}

func resourceDagRunRead(d *schema.ResourceData, m interface{}) error {
	pcfg := m.(ProviderConfig)
	client := pcfg.ApiClient.DAGRunApi

	dagId, dagRunId, err := airflowDagRunId(d.Id())
	if err != nil {
		return err
	}

	dagRun, resp, err := client.GetDagRun(pcfg.AuthContext, dagId, dagRunId).Execute()
	if resp != nil && resp.StatusCode == 404 {
		d.SetId("")
		return nil
	}
	if err != nil {
		return fmt.Errorf("failed to get dagRunId `%s` from Airflow: %w", d.Id(), err)
	}

	d.Set("dag_id", dagRun.DagId)
	d.Set("dag_run_id", dagRun.DagRunId.Get())
	d.Set("conf", dagRun.Conf)
	d.Set("state", dagRun.State)

	return nil
}

func resourceDagRunDelete(d *schema.ResourceData, m interface{}) error {
	pcfg := m.(ProviderConfig)
	client := pcfg.ApiClient.DAGRunApi

	dagId, dagRunId, err := airflowDagRunId(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.DeleteDagRun(pcfg.AuthContext, dagId, dagRunId).Execute()
	if err != nil {
		return fmt.Errorf("failed to delete dagRunId `%s` from Airflow: %w", d.Id(), err)
	}

	if resp != nil && resp.StatusCode == 404 {
		return nil
	}

	return nil
}

func airflowDagRunId(id string) (string, string, error) {
	parts := strings.SplitN(id, ":", 2)

	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf("unexpected format of ID (%s), expected DAG-ID:DAG-RUN-ID", id)
	}

	return parts[0], parts[1], nil
}
