package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"airflow_variable": resourceVariable(),
		},
		ConfigureFunc: func(*schema.ResourceData) (interface{}, error) {
			return &AirflowClient{}, nil
		},
	}
}
