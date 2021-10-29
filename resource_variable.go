package main

import (
	"fmt"

	"github.com/apache/airflow-client-go/airflow"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceVariable() *schema.Resource {
	return &schema.Resource{
		Create: resourceVariableCreate,
		Read:   resourceVariableRead,
		Update: resourceVariableUpdate,
		Delete: resourceVariableDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"key": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"value": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceVariableCreate(d *schema.ResourceData, m interface{}) error {
	pcfg := m.(ProviderConfig)
	client := pcfg.ApiClient
	key := d.Get("key").(string)
	val := d.Get("value").(string)
	varApi := client.VariableApi

	_, _, err := varApi.PostVariables(pcfg.AuthContext).Variable(airflow.Variable{
		Key:   &key,
		Value: &val,
	}).Execute()
	if err != nil {
		return fmt.Errorf("failed to create variable `%s` from Airflow: %w", key, err)
	}
	d.SetId(key)

	return resourceVariableRead(d, m)
}

func resourceVariableRead(d *schema.ResourceData, m interface{}) error {
	pcfg := m.(ProviderConfig)
	client := pcfg.ApiClient

	variable, resp, err := client.VariableApi.GetVariable(pcfg.AuthContext, d.Id()).Execute()
	if resp != nil && resp.StatusCode == 404 {
		d.SetId("")
		return nil
	}
	if err != nil {
		return fmt.Errorf("failed to get variable `%s` from Airflow: %w", d.Id(), err)
	}

	d.Set("key", variable.Key)
	d.Set("value", variable.Value)

	return nil
}

func resourceVariableUpdate(d *schema.ResourceData, m interface{}) error {
	pcfg := m.(ProviderConfig)
	client := pcfg.ApiClient

	val := d.Get("value").(string)
	key := d.Id()
	_, _, err := client.VariableApi.PatchVariable(pcfg.AuthContext, key).Variable(airflow.Variable{
		Key:   &key,
		Value: &val,
	}).Execute()
	if err != nil {
		return fmt.Errorf("failed to update variable `%s` from Airflow: %w", key, err)
	}

	return resourceVariableRead(d, m)
}

func resourceVariableDelete(d *schema.ResourceData, m interface{}) error {
	pcfg := m.(ProviderConfig)
	client := pcfg.ApiClient

	resp, err := client.VariableApi.DeleteVariable(pcfg.AuthContext, d.Id()).Execute()
	if err != nil {
		return fmt.Errorf("failed to delete variable `%s` from Airflow: %w", d.Id(), err)
	}

	if resp != nil && resp.StatusCode == 404 {
		return nil
	}

	return nil
}
