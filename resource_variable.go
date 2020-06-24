package main

import (
	"context"
	"fmt"

	"github.com/apache/airflow-client-go/airflow"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceVariableCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*airflow.APIClient)
	key := d.Get("key").(string)
	_, _, err := client.VariableApi.PostVariables(context.TODO(), airflow.Variable{
		Key:   key,
		Value: d.Get("value").(string),
	})
	if err != nil {
		return err
	}
	d.SetId(key)
	return resourceVariableRead(d, m)
}

func resourceVariableRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*airflow.APIClient)
	key := d.Get("key").(string)
	variable, resp, err := client.VariableApi.GetVariable(context.TODO(), key)
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
	client := m.(*airflow.APIClient)
	key := d.Get("key").(string)
	_, _, err := client.VariableApi.PatchVariable(context.TODO(), key, airflow.Variable{
		Key:   key,
		Value: d.Get("value").(string),
	}, nil)
	if err != nil {
		return err
	}
	d.SetId(key)
	return resourceVariableRead(d, m)
}

func resourceVariableDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*airflow.APIClient)
	key := d.Get("key").(string)
	_, err := client.VariableApi.DeleteVariable(context.TODO(), key)
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
