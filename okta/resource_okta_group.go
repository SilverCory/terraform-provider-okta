package okta

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/okta/okta-sdk-golang/v2/okta"
)

func resourceGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGroupCreate,
		ReadContext:   resourceGroupRead,
		UpdateContext: resourceGroupUpdate,
		DeleteContext: resourceGroupDelete,
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
				importID := strings.Split(d.Id(), "/")
				if len(importID) == 1 {
					return []*schema.ResourceData{d}, nil
				}
				if len(importID) > 2 {
					return nil, errors.New("invalid format used for import ID, format must be 'group_id' or 'group_id/skip_users'")
				}
				d.SetId(importID[0])
				if !isValidSkipArg(importID[1]) {
					return nil, fmt.Errorf("'%s' is invalid value to be used as part of import ID, it can only be 'skip_users'", importID[1])
				}
				_ = d.Set(importID[1], true)
				return []*schema.ResourceData{d}, nil
			},
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Group name",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Group description",
			},
			"users": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Users associated with the group. This can also be done per user.",
				Deprecated:  "The `users` field is now deprecated for the resource `okta_group`, please replace all uses of this with: `okta_group_memberships`",
			},
			"skip_users": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Ignore users sync. This is a temporary solution until 'users' field is supported in this resource",
				Default:     false,
			},
		},
	}
}

func resourceGroupCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	logger(m).Info("creating group", "name", d.Get("name").(string))
	group := buildGroup(d)
	responseGroup, _, err := getOktaClientFromMetadata(m).Group.CreateGroup(ctx, *group)
	if err != nil {
		return diag.Errorf("failed to create group: %v", err)
	}
	d.SetId(responseGroup.Id)
	err = updateGroupUsers(ctx, d, m)
	if err != nil {
		return diag.Errorf("failed to update group users on group create: %v", err)
	}
	return resourceGroupRead(ctx, d, m)
}

func resourceGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	logger(m).Info("reading group", "id", d.Id(), "name", d.Get("name").(string))
	g, resp, err := getOktaClientFromMetadata(m).Group.GetGroup(ctx, d.Id())
	if err := suppressErrorOn404(resp, err); err != nil {
		return diag.Errorf("failed to get group: %v", err)
	}
	if g == nil {
		d.SetId("")
		return nil
	}
	_ = d.Set("name", g.Profile.Name)
	_ = d.Set("description", g.Profile.Description)
	err = syncGroupUsers(ctx, d, m)
	if err != nil {
		return diag.Errorf("failed to get group users: %v", err)
	}
	return nil
}

func resourceGroupUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	logger(m).Info("updating group", "id", d.Id(), "name", d.Get("name").(string))
	group := buildGroup(d)
	_, _, err := getOktaClientFromMetadata(m).Group.UpdateGroup(ctx, d.Id(), *group)
	if err != nil {
		return diag.Errorf("failed to update group: %v", err)
	}
	err = updateGroupUsers(ctx, d, m)
	if err != nil {
		return diag.Errorf("failed to update group users on group update: %v", err)
	}
	return resourceGroupRead(ctx, d, m)
}

func resourceGroupDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	logger(m).Info("deleting group", "id", d.Id(), "name", d.Get("name").(string))
	_, err := getOktaClientFromMetadata(m).Group.DeleteGroup(ctx, d.Id())
	if err != nil {
		return diag.Errorf("failed to delete group: %v", err)
	}
	return nil
}

func syncGroupUsers(ctx context.Context, d *schema.ResourceData, m interface{}) error {
	// temp solution until 'users' field is supported
	if d.Get("skip_users").(bool) {
		return nil
	}
	userIDList, err := listGroupUserIDs(ctx, m, d.Id())
	if err != nil {
		return err
	}
	return d.Set("users", convertStringSliceToSet(userIDList))
}

func updateGroupUsers(ctx context.Context, d *schema.ResourceData, m interface{}) error {
	if !d.HasChange("users") {
		return nil
	}
	// temp solution until 'users' field is supported
	if d.Get("skip_users").(bool) {
		return nil
	}
	client := getOktaClientFromMetadata(m)
	oldGM, newGM := d.GetChange("users")
	oldSet := oldGM.(*schema.Set)
	newSet := newGM.(*schema.Set)
	usersToAdd := convertInterfaceArrToStringArr(newSet.Difference(oldSet).List())
	usersToRemove := convertInterfaceArrToStringArr(oldSet.Difference(newSet).List())
	err := addUserToGroups(ctx, client, d.Id(), usersToAdd)
	if err != nil {
		return err
	}
	return removeUserFromGroups(ctx, client, d.Id(), usersToRemove)
}

func containsUser(users []*okta.User, id string) bool {
	for _, user := range users {
		if user.Id == id {
			return true
		}
	}
	return false
}

func buildGroup(d *schema.ResourceData) *okta.Group {
	return &okta.Group{
		Profile: &okta.GroupProfile{
			Name:        d.Get("name").(string),
			Description: d.Get("description").(string),
		},
	}
}
