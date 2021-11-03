package main

import (
	"fmt"

	"github.com/apache/airflow-client-go/airflow"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceConnection() *schema.Resource {
	return &schema.Resource{
		Create: resourceConnectionCreate,
		Read:   resourceConnectionRead,
		Update: resourceConnectionUpdate,
		Delete: resourceConnectionDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"connection_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IsPortNumberOrZero,
			},
			"password": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"extra": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceConnectionCreate(d *schema.ResourceData, m interface{}) error {
	pcfg := m.(ProviderConfig)
	client := pcfg.ApiClient
	connId := d.Get("connection_id").(string)
	connType := d.Get("conn_type").(string)

	conn := airflow.Connection{
		ConnectionId: &connId,
		ConnType:     &connType,
	}

	if v, ok := d.GetOk("host"); ok {
		conn.SetHost(v.(string))
	}

	if v, ok := d.GetOk("login"); ok {
		conn.SetLogin(v.(string))
	}

	if v, ok := d.GetOk("schema"); ok {
		conn.SetSchema(v.(string))
	}

	if v, ok := d.GetOk("port"); ok {
		conn.SetPort(int32(v.(int)))
	}

	conn.SetPassword(d.Get("password").(string))

	if v, ok := d.GetOk("extra"); ok {
		conn.SetExtra(v.(string))
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
	connection, resp, err := client.ConnectionApi.GetConnection(pcfg.AuthContext, d.Id()).Execute()
	if resp != nil && resp.StatusCode == 404 {
		d.SetId("")
		return nil
	}
	if err != nil {
		return fmt.Errorf("failed to get connection `%s` from Airflow: %w", d.Id(), err)
	}

	d.Set("connection_id", connection.GetConnectionId())
	d.Set("conn_type", connection.GetConnType())
	d.Set("host", connection.GetHost())
	d.Set("login", connection.GetLogin())
	d.Set("schema", connection.GetSchema())
	d.Set("port", connection.GetPort())
	d.Set("extra", connection.GetExtra())

	if v, ok := connection.GetPasswordOk(); ok {
		d.Set("password", v)
	} else if v, ok := d.GetOk("password"); ok {
		d.Set("password", v)
	}

	return nil
}

func resourceConnectionUpdate(d *schema.ResourceData, m interface{}) error {
	pcfg := m.(ProviderConfig)
	client := pcfg.ApiClient
	connId := d.Id()
	connType := d.Get("conn_type").(string)

	conn := airflow.Connection{
		ConnectionId: &connId,
		ConnType:     &connType,
	}

	if v, ok := d.GetOk("host"); ok {
		conn.SetHost(v.(string))
	}

	if v, ok := d.GetOk("login"); ok {
		conn.SetLogin(v.(string))
	}

	if v, ok := d.GetOk("schema"); ok {
		conn.SetSchema(v.(string))
	}

	if v, ok := d.GetOk("port"); ok {
		conn.SetPort(int32(v.(int)))
	}

	if v, ok := d.GetOk("password"); ok && v.(string) != "" {
		conn.SetPassword(v.(string))
	}

	if v, ok := d.GetOk("extra"); ok {
		conn.SetExtra(v.(string))
	}

	_, _, err := client.ConnectionApi.PatchConnection(pcfg.AuthContext, connId).Connection(conn).Execute()
	if err != nil {
		return fmt.Errorf("failed to update connection `%s` from Airflow: %w", connId, err)
	}

	return resourceConnectionRead(d, m)
}

func resourceConnectionDelete(d *schema.ResourceData, m interface{}) error {
	pcfg := m.(ProviderConfig)
	client := pcfg.ApiClient

	resp, err := client.ConnectionApi.DeleteConnection(pcfg.AuthContext, d.Id()).Execute()
	if err != nil {
		return fmt.Errorf("failed to delete connection `%s` from Airflow: %w", d.Id(), err)
	}

	if resp != nil && resp.StatusCode == 404 {
		return nil
	}

	return nil
}
