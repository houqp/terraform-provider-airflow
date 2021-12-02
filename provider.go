package main

import (
	"context"
	"fmt"
	"log"
	"net/url"

	"github.com/apache/airflow-client-go/airflow"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type ProviderConfig struct {
	ApiClient   *airflow.APIClient
	AuthContext context.Context
}

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"base_endpoint": {
				Type:     schema.TypeString,
				Required: true,
			},
			// username and password are used for API basic auth
			"username": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The username to use for API basic authentication",
			},
			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The password to use for API basic authentication",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"airflow_variable": resourceVariable(),
			"airflow_user":     resourceUser(),
		},
		ConfigureFunc: func(p *schema.ResourceData) (interface{}, error) {
			endpoint := p.Get("base_endpoint").(string)
			u, err := url.Parse(endpoint)
			if err != nil {
				return nil, fmt.Errorf("invalid base_endpoint: %w", err)
			}

			authCtx := context.Background()

			if username, ok := p.GetOk("username"); ok {
				var password interface{}
				if password, ok = p.GetOk("password"); !ok {
					return nil, fmt.Errorf("found username for basic auth, but password not specified")
				}
				log.Printf("[DEBUG] Using API Basic Auth")

				cred := airflow.BasicAuth{
					UserName: username.(string),
					Password: password.(string),
				}
				authCtx = context.WithValue(authCtx, airflow.ContextBasicAuth, cred)
			}

			return ProviderConfig{
				ApiClient: airflow.NewAPIClient(&airflow.Configuration{
					Scheme: u.Scheme,
					Host:   u.Host,
				}),
				AuthContext: authCtx,
			}, nil
		},
	}
}
