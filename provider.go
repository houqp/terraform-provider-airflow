package main

import (
	"fmt"
	"log"
	"net/url"
	"path"

	"github.com/apache/airflow-client-go/airflow"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"base_endpoint": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"airflow_variable": resourceVariable(),
		},
		ConfigureFunc: func(p *schema.ResourceData) (interface{}, error) {
			endpoint := p.Get("base_endpoint").(string)
			u, err := url.Parse(endpoint)
			if err != nil {
				return nil, fmt.Errorf("invalid base_endpoint: %w", err)
			}

			basePath := path.Join(u.Path + "/api/v1")
			log.Printf("[DEBUG] Using API prefix: %s", basePath)

			return airflow.NewAPIClient(&airflow.Configuration{
				Scheme:   u.Scheme,
				Host:     u.Host,
				BasePath: basePath,
			}), nil
		},
	}
}
