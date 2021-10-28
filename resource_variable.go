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
	val := d.Get("value").(string)
	varApi := client.VariableApi

	_, _, err := varApi.PostVariables(pcfg.AuthContext).Variable(airflow.Variable{
		Key:   &key,
		Value: &val,
	}).Execute()
	if err != nil {
		return err
	}
	d.SetId(key)
	return resourceVariableRead(d, m)
}

func resourceVariableRead(d *schema.ResourceData, m interface{}) error {
	pcfg := m.(ProviderConfig)
	client := pcfg.ApiClient
	key := d.Get("key").(string)
	variable, resp, err := client.VariableApi.GetVariable(pcfg.AuthContext, key).Execute()
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
	val := d.Get("value").(string)
	key := d.Get("key").(string)
	_, _, err := client.VariableApi.PatchVariable(pcfg.AuthContext, key).Variable(airflow.Variable{
		Key:   &key,
		Value: &val,
	}).Execute()
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
	_, err := client.VariableApi.DeleteVariable(pcfg.AuthContext, key).Execute()
	return err
}

func resourceVariable() *schema.Resource {
	return &schema.Resource{
		Create: resourceVariableCreate,
		Read:   resourceVariableRead,
		Update: resourceVariableUpdate,
		Delete: resourceVariableDelete,

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
