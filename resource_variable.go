package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceVariableCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*AirflowClient)
	key := d.Get("key").(string)
	err := client.SetVariable(key, d.Get("value").(string))
	if err != nil {
		return err
	}
	d.SetId(key)
	return resourceVariableRead(d, m)
}

func resourceVariableRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*AirflowClient)
	key := d.Get("key").(string)
	value := client.ReadVariable(key)
	if value == "" {
		d.SetId("")
		return nil
	}

	d.Set("key", key)
	d.Set("value", value)
	return nil
}

func resourceVariableUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceVariableCreate(d, m)
}

func resourceVariableDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*AirflowClient)
	key := d.Get("key").(string)
	return client.DeleteVariable(key)
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
