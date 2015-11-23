package main

import (
	"github.com/hashicorp/terraform/helper/schema"
	//"log"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"username": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				//Default:     "",
				Description: descriptions["username"],
			},
			"password": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				//Default:     "",
				Description: descriptions["password"],
			},
			"management_ip": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				//Default:     "",
				Description: descriptions["management_ip"],
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"bigip_node": resourceNode(),
			"bigip_pool": resourcePool(),
			"bigip_pool_member": resourcePoolMember(),
			"bigip_vserver": resourceVServer(),
		},
		ConfigureFunc: providerConfigure,
	}
}


var descriptions map[string]string

func init() {
	descriptions = map[string]string{
		"management_ip": "The management IP address of the F5 BIGIP device",
		"username": "The username to use for RestAPI calls",
		"password": "The password to use for RestAPI calls",
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := Config{
		RestUsername:        d.Get("username").(string),
		RestPassword:        d.Get("password").(string),
		RestIP:              d.Get("management_ip").(string),
	}
	return config.Client()
}

