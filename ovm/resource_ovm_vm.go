package ovm

import "github.com/hashicorp/terraform/helper/schema"

type commonVmParams struct {
	repositoryId int64
	serverPoolId int64
	//The following fields of the new Vm are optional:
	bootOrder          []string
	cpuCount           int
	cpuCountLimit      int
	cpuPriority        int
	cpuUtilizationCap  int
	highAvailability   bool
	hugePagesEnabled   bool
	keymapName         sting
	memory             int
	memoryLimit        int
	networkInstallPath string
	osType             string
	serverId           int64
	vmDomainType       string
	vmMouseType        string
	vmRunState         string
	vmStartPolicy      string
}

func resourcePingdomCheck() *schema.Resource {
	return &schema.Resource{
		Create: resourceOvmCheckCreate,
		Read:   resourceOvmCheckRead,
		Update: resourceOvmCheckUpdate,
		Delete: resourceOvmCheckDelete,
		Importer: &schema.ResourceImporter{
			State: resourceOvmCheckImporter,
		},

		Schema: map[string]*schema.Schema{
			"repositoryId": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: false,
			},
			"serverPoolId": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: false,
			},
		},
	}
}
