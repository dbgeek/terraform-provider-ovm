package ovm

import (
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"user": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},

			"password": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},

			"entrypoint": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"ovm_vm": resourceOvmVm(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func ProviderConfigure(d *schema.ResourceData) (interface{}, error) {
	var config Config
	configRaw := d.Get("").(map[string]interface{})
	if err := mapstructure.Decode(configRaw, &config); err != nil {
		return nil, err
	}

	log.Println("[INFO] Initializing OVM client")
	return config.Client()
}
