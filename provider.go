package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"variables_output_path": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"airflow_variable": resourceVariable(),
		},
		ConfigureFunc: func(p *schema.ResourceData) (interface{}, error) {
			return &AirflowClient{
				VariablesOutputPath: p.Get("variables_output_path").(string),
			}, nil
		},
	}
}
