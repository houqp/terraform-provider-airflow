package main

import (
	"fmt"

	"github.com/apache/airflow-client-go/airflow"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceRole() *schema.Resource {
	return &schema.Resource{
		Create: resourceRoleCreate,
		Read:   resourceRoleRead,
		Update: resourceRoleUpdate,
		Delete: resourceRoleDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"action": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"action": {
							Type:     schema.TypeString,
							Required: true,
						},
						"resource": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
		},
	}
}

func resourceRoleCreate(d *schema.ResourceData, m interface{}) error {
	pcfg := m.(ProviderConfig)
	client := pcfg.ApiClient

	name := d.Get("name").(string)
	varApi := client.RoleApi
	role := airflow.Role{
		Name: &name,
	}

	if v, ok := d.GetOk("action"); ok && v.(*schema.Set).Len() > 0 {
		actions := expandAirflowRoleActions(d.Get("action").(*schema.Set).List())
		role.Actions = &actions
	}

	_, _, err := varApi.PostRole(pcfg.AuthContext).Role(role).Execute()
	if err != nil {
		return fmt.Errorf("failed to create role `%s` from Airflow: %w", name, err)
	}
	d.SetId(name)

	return resourceRoleRead(d, m)
}

func resourceRoleRead(d *schema.ResourceData, m interface{}) error {
	pcfg := m.(ProviderConfig)
	client := pcfg.ApiClient

	role, resp, err := client.RoleApi.GetRole(pcfg.AuthContext, d.Id()).Execute()
	if resp != nil && resp.StatusCode == 404 {
		d.SetId("")
		return nil
	}
	if err != nil {
		return fmt.Errorf("failed to get role `%s` from Airflow: %w", d.Id(), err)
	}

	d.Set("name", role.Name)
	if err := d.Set("action", flattenAirflowRoleActions(*role.Actions)); err != nil {
		return fmt.Errorf("error setting action: %w", err)
	}

	return nil
}

func resourceRoleUpdate(d *schema.ResourceData, m interface{}) error {
	pcfg := m.(ProviderConfig)
	client := pcfg.ApiClient

	name := d.Id()
	actions := expandAirflowRoleActions(d.Get("action").(*schema.Set).List())
	role := airflow.Role{
		Name:    &name,
		Actions: &actions,
	}

	_, _, err := client.RoleApi.PatchRole(pcfg.AuthContext, name).Role(role).Execute()
	if err != nil {
		return fmt.Errorf("failed to update role `%s` from Airflow: %w", name, err)
	}

	return resourceRoleRead(d, m)
}

func resourceRoleDelete(d *schema.ResourceData, m interface{}) error {
	pcfg := m.(ProviderConfig)
	client := pcfg.ApiClient

	resp, err := client.RoleApi.DeleteRole(pcfg.AuthContext, d.Id()).Execute()
	if err != nil {
		return fmt.Errorf("failed to delete role `%s` from Airflow: %w", d.Id(), err)
	}

	if resp != nil && resp.StatusCode == 404 {
		return nil
	}

	return nil
}

func expandAirflowRoleActions(tfList []interface{}) []airflow.ActionResource {
	if len(tfList) == 0 {
		return nil
	}

	apiObjects := make([]airflow.ActionResource, 0)

	for _, tfMapRaw := range tfList {
		tfMap, ok := tfMapRaw.(map[string]interface{})

		if !ok {
			continue
		}

		action := tfMap["action"].(string)
		resource := tfMap["resource"].(string)

		apiObject := airflow.ActionResource{
			Action: &airflow.Action{
				Name: &action,
			},
			Resource: &airflow.Resource{
				Name: &resource,
			},
		}
		apiObjects = append(apiObjects, apiObject)
	}

	return apiObjects
}

func flattenAirflowRoleActions(apiObjects []airflow.ActionResource) []interface{} {
	if len(apiObjects) == 0 {
		return nil
	}

	var tfList []interface{}

	for _, apiObject := range apiObjects {
		tfList = append(tfList, map[string]interface{}{
			"action":   apiObject.Action.Name,
			"resource": apiObject.Resource.Name,
		})
	}

	return tfList
}
