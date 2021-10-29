package main

import (
	"fmt"

	"github.com/apache/airflow-client-go/airflow"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceConnectionCreate(d *schema.ResourceData, m interface{}) error {
	pcfg := m.(ProviderConfig)
	client := pcfg.ApiClient
	connId := d.Get("connection_id").(string)
	connType := d.Get("conn_type").(string)

	conn := airflow.Connection{
		ConnectionId: &connId,
		ConnType:     &connType,
	}
	connApi := client.ConnectionApi

	_, _, err := connApi.PostConnection(pcfg.AuthContext).Connection(conn).Execute()
	if err != nil {
		return fmt.Errorf("failed to create connection `%s` from Airflow: %w", connId, err)
	}
	d.SetId(connId)
	return resourceConnectionRead(d, m)
}

func resourceConnectionRead(d *schema.ResourceData, m interface{}) error {
	pcfg := m.(ProviderConfig)
	client := pcfg.ApiClient
	connId := d.Get("connection_id").(string)
	connection, resp, err := client.ConnectionApi.GetConnection(pcfg.AuthContext, connId).Execute()
	if resp != nil && resp.StatusCode == 404 {
		d.SetId("")
		return nil
	}
	if err != nil {
		return fmt.Errorf("failed to get connection `%s` from Airflow: %w", connId, err)
	}

	d.Set("connection_id", connection.ConnectionId)
	d.Set("conn_type", connection.ConnType)
	return nil
}

func resourceConnectionUpdate(d *schema.ResourceData, m interface{}) error {
	pcfg := m.(ProviderConfig)
	client := pcfg.ApiClient
	connId := d.Get("connection_id").(string)
	connType := d.Get("conn_type").(string)

	conn := airflow.Connection{
		ConnectionId: &connId,
		ConnType:     &connType,
	}
	_, _, err := client.ConnectionApi.PatchConnection(pcfg.AuthContext, connId).Connection(conn).Execute()
	if err != nil {
		return fmt.Errorf("failed to update connection `%s` from Airflow: %w", connId, err)
	}
	d.SetId(connId)
	return resourceConnectionRead(d, m)
}

func resourceConnectionDelete(d *schema.ResourceData, m interface{}) error {
	pcfg := m.(ProviderConfig)
	client := pcfg.ApiClient
	connId := d.Get("connection_id").(string)
	_, err := client.ConnectionApi.DeleteConnection(pcfg.AuthContext, connId).Execute()
	if err != nil {
		return fmt.Errorf("failed to delete connection `%s` from Airflow: %w", connId, err)
	}

	return nil
}

func resourceConnection() *schema.Resource {
	return &schema.Resource{
		Create: resourceConnectionCreate,
		Read:   resourceConnectionRead,
		Update: resourceConnectionUpdate,
		Delete: resourceConnectionDelete,

		Schema: map[string]*schema.Schema{
			"connection_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"conn_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"host": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"login": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"schema": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"port": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"password": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"extra": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}
