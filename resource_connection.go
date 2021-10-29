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
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return old == ""
				},
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
		val := v.(string)
		conn.Host = *airflow.NewNullableString(&val)
	}

	if v, ok := d.GetOk("login"); ok {
		val := v.(string)
		conn.Login = *airflow.NewNullableString(&val)
	}

	if v, ok := d.GetOk("schema"); ok {
		val := v.(string)
		conn.Schema = *airflow.NewNullableString(&val)
	}

	if v, ok := d.GetOk("port"); ok {
		val := int32(v.(int))
		conn.Port = *airflow.NewNullableInt32(&val)
	}

	if v, ok := d.GetOk("password"); ok {
		val := v.(string)
		conn.Password = &val
	}

	if v, ok := d.GetOk("extra"); ok {
		val := v.(string)
		conn.Extra = *airflow.NewNullableString(&val)
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

	d.Set("connection_id", connection.ConnectionId)
	d.Set("conn_type", connection.ConnType)
	d.Set("host", connection.Host.Get())
	d.Set("login", connection.Login.Get())
	d.Set("schema", connection.Schema.Get())
	d.Set("port", connection.Port.Get())
	d.Set("password", connection.Password)
	d.Set("extra", connection.Extra.Get())

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
		val := v.(string)
		conn.Host = *airflow.NewNullableString(&val)
	}

	if v, ok := d.GetOk("login"); ok {
		val := v.(string)
		conn.Login = *airflow.NewNullableString(&val)
	}

	if v, ok := d.GetOk("schema"); ok {
		val := v.(string)
		conn.Schema = *airflow.NewNullableString(&val)
	}

	if v, ok := d.GetOk("port"); ok {
		val := int32(v.(int))
		conn.Port = *airflow.NewNullableInt32(&val)
	}

	if v, ok := d.GetOk("password"); ok {
		val := v.(string)
		conn.Password = &val
	}

	if v, ok := d.GetOk("extra"); ok {
		val := v.(string)
		conn.Extra = *airflow.NewNullableString(&val)
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
