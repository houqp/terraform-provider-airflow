package main

import (
	"fmt"

	"github.com/apache/airflow-client-go/airflow"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceUserCreate(d *schema.ResourceData, m interface{}) error {
	pcfg := m.(ProviderConfig)
	client := pcfg.ApiClient
	username := d.Get("username").(string)
	email := d.Get("email").(string)
	firstName := d.Get("first_name").(string)
	lastName := d.Get("last_name").(string)
	password := d.Get("password").(string)
	roles := stringListToRoles(d.Get("roles").(*schema.Set).List())

	u := airflow.User{
		Username:  &username,
		Email:     &email,
		FirstName: &firstName,
		LastName:  &lastName,
		Password:  &password,
		Roles:     &roles,
	}

	req := client.UserApi.PostUser(pcfg.AuthContext)
	req.User(u)
	_, _, err := req.Execute()
	if err != nil {
		return err
	}

	d.SetId(username)

	return resourceUserRead(d, m)
}

func resourceUserRead(d *schema.ResourceData, m interface{}) error {
	pcfg := m.(ProviderConfig)
	client := pcfg.ApiClient
	username := d.Get("username").(string)

	req := client.UserApi.GetUser(pcfg.AuthContext, username)

	user, resp, err := req.Execute()
	if resp != nil && resp.StatusCode == 404 {
		d.SetId("")
		return nil
	}
	if err != nil {
		return fmt.Errorf("failed to get user `%s` from Airflow: %w", username, err)
	}

	d.Set("username", user.Username)
	d.Set("email", user.Email)
	d.Set("first_name", user.FirstName)
	d.Set("last_name", user.LastName)
	roles := rolesToStringList(*user.Roles)
	d.Set("roles", roles)

	return nil
}

func resourceUserUpdate(d *schema.ResourceData, m interface{}) error {
	pcfg := m.(ProviderConfig)
	client := pcfg.ApiClient
	username := d.Get("username").(string)
	email := d.Get("email").(string)
	firstName := d.Get("first_name").(string)
	lastName := d.Get("last_name").(string)
	password := d.Get("password").(string)
	roles := stringListToRoles(d.Get("roles").(*schema.Set).List())

	u := airflow.User{
		Username:  &username,
		Email:     &email,
		FirstName: &firstName,
		LastName:  &lastName,
		Password:  &password,
		Roles:     &roles,
	}

	req := client.UserApi.PatchUser(pcfg.AuthContext, username)
	req.User(u)
	_, _, err := req.Execute()
	if err != nil {
		return err
	}

	d.SetId(username)

	return resourceUserRead(d, m)
}

func resourceUserDelete(d *schema.ResourceData, m interface{}) error {
	pcfg := m.(ProviderConfig)
	client := pcfg.ApiClient
	username := d.Get("username").(string)
	req := client.UserApi.DeleteUser(pcfg.AuthContext, username)
	_, err := req.Execute()
	return err
}

func resourceUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceUserCreate,
		Read:   resourceUserRead,
		Update: resourceUserUpdate,
		Delete: resourceUserDelete,
		Schema: map[string]*schema.Schema{
			"username": {
				Type:     schema.TypeString,
				Required: true,
			},
			"email": {
				Type:     schema.TypeString,
				Required: true,
			},
			"first_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"last_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"roles": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "**WARNING:** this will put the password in the terraform state file. Use carefully.",
			},
		},
	}
}

func stringListToRoles(list []interface{}) []airflow.UserCollectionItemRoles {
	vs := make([]airflow.UserCollectionItemRoles, 0, len(list))
	for _, v := range list {
		val, ok := v.(string)
		if ok && val != "" {
			vs = append(vs, airflow.UserCollectionItemRoles{Name: &val})
		}
	}
	return vs
}

func rolesToStringList(roles []airflow.UserCollectionItemRoles) []string {
	vs := make([]string, 0, len(roles))
	for _, v := range roles {
		vs = append(vs, *v.Name)
	}
	return vs
}
