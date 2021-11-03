package main

import (
	"context"
	"fmt"
	"net/url"

	"github.com/apache/airflow-client-go/airflow"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type ProviderConfig struct {
	ApiClient   *airflow.APIClient
	AuthContext context.Context
}

func AirflowProvider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"base_endpoint": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("AIRFLOW_BASE_ENDPOINT", nil),
			},
			"oauth2_token": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The oauth to use for API authentication",
				DefaultFunc: schema.EnvDefaultFunc("AIRFLOW_OAUTH2_TOKEN", nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"airflow_connection": resourceConnection(),
			"airflow_variable":   resourceVariable(),
			"airflow_pool":       resourcePool(),
			"airflow_role":       resourceRole(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	endpoint := d.Get("base_endpoint").(string)
	u, err := url.Parse(endpoint)
	if err != nil {
		return nil, fmt.Errorf("invalid base_endpoint: %w", err)
	}

	authCtx := context.Background()

	authCtx = context.WithValue(authCtx, airflow.ContextAccessToken, d.Get("oauth2_token"))

	return ProviderConfig{
		ApiClient: airflow.NewAPIClient(&airflow.Configuration{
			Scheme: u.Scheme,
			Host:   u.Host,
			Debug:  true,
			Servers: airflow.ServerConfigurations{
				{
					URL:         "/api/v1",
					Description: "Apache Airflow Stable API.",
				},
			},
		}),
		AuthContext: authCtx,
	}, nil
}
