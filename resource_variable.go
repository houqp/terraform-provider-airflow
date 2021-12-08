package main

import (
	"fmt"

	"github.com/apache/airflow-client-go/airflow"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceVariableCreate(d *schema.ResourceData, m interface{}) error {
	pcfg := m.(ProviderConfig)
	client := pcfg.ApiClient
	key := d.Get("key").(string)
	value := d.Get("value").(string)

	v := airflow.Variable{
		Key:   &key,
		Value: &value,
	}

	req := client.VariableApi.PostVariables(pcfg.AuthContext)
	_, _, err := req.Variable(v).Execute()
	if err != nil {
		return err
	}

	d.SetId(key)

	return resourceVariableRead(d, m)
}

func resourceVariableRead(d *schema.ResourceData, m interface{}) error {
	pcfg := m.(ProviderConfig)
	client := pcfg.ApiClient
	key := d.Id()

	req := client.VariableApi.GetVariable(pcfg.AuthContext, key)
	variable, resp, err := req.Execute()
	if resp != nil && resp.StatusCode == 404 {
		d.SetId("")
		return nil
	}
	if err != nil {
		return fmt.Errorf("failed to get variable `%s` from Airflow: %w", key, err)
	}

	d.Set("key", variable.Key)
	d.Set("value", variable.Value)
	return nil
}

func resourceVariableUpdate(d *schema.ResourceData, m interface{}) error {
	pcfg := m.(ProviderConfig)
	client := pcfg.ApiClient
	key := d.Get("key").(string)
	value := d.Get("value").(string)

	v := airflow.Variable{
		Key:   &key,
		Value: &value,
	}

	req := client.VariableApi.PatchVariable(pcfg.AuthContext, key)
	_, _, err := req.Variable(v).Execute()
	if err != nil {
		return err
	}

	d.SetId(key)
	return resourceVariableRead(d, m)
}

func resourceVariableDelete(d *schema.ResourceData, m interface{}) error {
	pcfg := m.(ProviderConfig)
	client := pcfg.ApiClient
	key := d.Get("key").(string)
	req := client.VariableApi.DeleteVariable(pcfg.AuthContext, key)
	_, err := req.Execute()
	return err
}

func resourceVariableImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	if err := resourceVariableRead(d, m); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}

func resourceVariable() *schema.Resource {
	return &schema.Resource{
		Create: resourceVariableCreate,
		Read:   resourceVariableRead,
		Update: resourceVariableUpdate,
		Delete: resourceVariableDelete,
		Importer: &schema.ResourceImporter{
			State: resourceVariableImport,
		},
		Schema: map[string]*schema.Schema{
			"key": {
				Type:     schema.TypeString,
				Required: true,
			},
			"value": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}
