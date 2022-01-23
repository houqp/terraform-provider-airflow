package main

import (
	"fmt"

	"github.com/apache/airflow-client-go/airflow"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceDag() *schema.Resource {
	return &schema.Resource{
		Create: resourceDagUpdate,
		Read:   resourceDagRead,
		Update: resourceDagUpdate,
		Delete: resourceDagDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"dag_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"file_token": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"fileloc": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_active": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"is_paused": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"is_subdag": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"root_dag_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceDagUpdate(d *schema.ResourceData, m interface{}) error {
	pcfg := m.(ProviderConfig)
	client := pcfg.ApiClient

	dagId := d.Get("dag_id").(string)
	dagApi := client.DAGApi
	dag := *airflow.NewDAG()
	dag.SetIsPaused(d.Get("is_paused").(bool))

	_, res, err := dagApi.PatchDag(pcfg.AuthContext, dagId).DAG(dag).Execute()
	if res.StatusCode != 200 {
		return fmt.Errorf("failed to update DAG `%s` from Airflow: %w", dagId, err)
	}
	d.SetId(dagId)

	return resourceDagRead(d, m)
}

func resourceDagRead(d *schema.ResourceData, m interface{}) error {
	pcfg := m.(ProviderConfig)
	client := pcfg.ApiClient

	DAG, resp, err := client.DAGApi.GetDag(pcfg.AuthContext, d.Id()).Execute()
	if resp != nil && resp.StatusCode == 404 {
		d.SetId("")
		return nil
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("failed to get DAG `%s` from Airflow: %w", d.Id(), err)
	}

	d.Set("dag_id", DAG.DagId)
	d.Set("is_paused", DAG.IsPaused.Get())
	d.Set("is_active", DAG.IsActive.Get())
	d.Set("is_subdag", DAG.IsSubdag)
	d.Set("description", DAG.Description.Get())
	d.Set("file_token", DAG.FileToken)
	d.Set("fileloc", DAG.Fileloc)
	d.Set("root_dag_id", DAG.RootDagId.Get())

	return nil
}

func resourceDagDelete(d *schema.ResourceData, m interface{}) error {
	// pcfg := m.(ProviderConfig)
	// client := pcfg.ApiClient.DAGApi

	// resp, err := client.DeleteDag(pcfg.AuthContext, d.Id()).Execute()
	// if err != nil {
	// 	return fmt.Errorf("failed to delete DAG `%s` from Airflow: %w", d.Id(), err)
	// }

	// if resp != nil && resp.StatusCode == 404 {
	// 	return nil
	// }

	return nil
}
