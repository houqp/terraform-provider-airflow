package main

import (
	"context"
	"fmt"
	"log"
	"net/url"

	"github.com/apache/airflow-client-go/airflow"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

type ProviderConfig struct {
	ApiClient   *airflow.APIClient
	AuthContext context.Context
}

func AirflowProvider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"base_endpoint": {
				Type:         schema.TypeString,
				Required:     true,
				DefaultFunc:  schema.EnvDefaultFunc("AIRFLOW_BASE_ENDPOINT", nil),
				ValidateFunc: validation.IsURLWithHTTPorHTTPS,
			},
			"oauth2_token": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "The oauth to use for API authentication",
				DefaultFunc:   schema.EnvDefaultFunc("AIRFLOW_OAUTH2_TOKEN", nil),
				ConflictsWith: []string{"username", "password"},
			},
			"username": {
				Type:          schema.TypeString,
				DefaultFunc:   schema.EnvDefaultFunc("AIRFLOW_API_USERNAME", nil),
				Optional:      true,
				Description:   "The username to use for API basic authentication",
				RequiredWith:  []string{"password"},
				ConflictsWith: []string{"oauth2_token"},
			},
			"password": {
				Type:          schema.TypeString,
				DefaultFunc:   schema.EnvDefaultFunc("AIRFLOW_API_PASSWORD", nil),
				Optional:      true,
				Description:   "The password to use for API basic authentication",
				RequiredWith:  []string{"username"},
				ConflictsWith: []string{"oauth2_token"},
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"airflow_connection": resourceConnection(),
			"airflow_dag":        resourceDag(),
			"airflow_dag_run":    resourceDagRun(),
			"airflow_variable":   resourceVariable(),
			"airflow_pool":       resourcePool(),
			"airflow_role":       resourceRole(),
			"airflow_user":       resourceUser(),
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
	if v, ok := d.GetOk("oauth2_token"); ok {
		authCtx = context.WithValue(authCtx, airflow.ContextAccessToken, v)
	}

	if username, ok := d.GetOk("username"); ok {
		var password interface{}
		if password, ok = d.GetOk("password"); !ok {
			return nil, fmt.Errorf("found username for basic auth, but password not specified")
		}
		log.Printf("[DEBUG] Using API Basic Auth")

		cred := airflow.BasicAuth{
			UserName: username.(string),
			Password: password.(string),
		}
		authCtx = context.WithValue(authCtx, airflow.ContextBasicAuth, cred)
	}

	clientConf := &airflow.Configuration{
		Scheme: u.Scheme,
		Host:   u.Host,
		Debug:  true,
		Servers: airflow.ServerConfigurations{
			{
				URL:         "/api/v1",
				Description: "Apache Airflow Stable API.",
			},
		},
	}

	return ProviderConfig{
		ApiClient:   airflow.NewAPIClient(clientConf),
		AuthContext: authCtx,
	}, nil
}
