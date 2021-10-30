package main

import (
	"fmt"

	"github.com/apache/airflow-client-go/airflow"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourcePool() *schema.Resource {
	return &schema.Resource{
		Create: resourcePoolCreate,
		Read:   resourcePoolRead,
		Update: resourcePoolUpdate,
		Delete: resourcePoolDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"slots": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"occupied_slots": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"used_slots": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"queued_slots": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"open_slots": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourcePoolCreate(d *schema.ResourceData, m interface{}) error {
	pcfg := m.(ProviderConfig)
	client := pcfg.ApiClient

	name := d.Get("name").(string)
	slots := int32(d.Get("slots").(int))
	varApi := client.PoolApi

	pool := airflow.Pool{
		Name:  &name,
		Slots: &slots,
	}

	_, _, err := varApi.PostPool(pcfg.AuthContext).Pool(pool).Execute()
	if err != nil {
		return fmt.Errorf("failed to create pool `%s` from Airflow: %w", name, err)
	}
	d.SetId(name)

	return resourcePoolRead(d, m)
}

func resourcePoolRead(d *schema.ResourceData, m interface{}) error {
	pcfg := m.(ProviderConfig)
	client := pcfg.ApiClient

	pool, resp, err := client.PoolApi.GetPool(pcfg.AuthContext, d.Id()).Execute()
	if resp != nil && resp.StatusCode == 404 {
		d.SetId("")
		return nil
	}
	if err != nil {
		return fmt.Errorf("failed to get pool `%s` from Airflow: %w", d.Id(), err)
	}

	d.Set("name", pool.Name)
	d.Set("slots", pool.Slots)
	d.Set("occupied_slots", pool.OccupiedSlots)
	d.Set("queued_slots", pool.QueuedSlots)
	d.Set("open_slots", pool.OpenSlots)
	d.Set("used_slots", pool.UsedSlots)

	return nil
}

func resourcePoolUpdate(d *schema.ResourceData, m interface{}) error {
	pcfg := m.(ProviderConfig)
	client := pcfg.ApiClient

	slots := int32(d.Get("slots").(int))
	name := d.Id()

	pool := airflow.Pool{
		Name:  &name,
		Slots: &slots,
	}

	_, _, err := client.PoolApi.PatchPool(pcfg.AuthContext, name).Pool(pool).Execute()
	if err != nil {
		return fmt.Errorf("failed to update pool `%s` from Airflow: %w", name, err)
	}

	return resourcePoolRead(d, m)
}

func resourcePoolDelete(d *schema.ResourceData, m interface{}) error {
	pcfg := m.(ProviderConfig)
	client := pcfg.ApiClient

	resp, err := client.PoolApi.DeletePool(pcfg.AuthContext, d.Id()).Execute()
	if err != nil {
		return fmt.Errorf("failed to delete pool `%s` from Airflow: %w", d.Id(), err)
	}

	if resp != nil && resp.StatusCode == 404 {
		return nil
	}

	return nil
}
