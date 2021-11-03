package main

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceConnection() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceConnectionRead,
		Schema: map[string]*schema.Schema{
			"connection_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"conn_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"host": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"login": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"schema": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"port": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"password": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"extra": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceConnectionRead(d *schema.ResourceData, m interface{}) error {
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
	d.Set("password", connection.GetPassword())
	d.Set("extra", connection.GetExtra())

	return nil
}
