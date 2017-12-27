package ovm

import (
	"github.com/dbgeek/go-ovm-helper/ovmHelper"
	"github.com/hashicorp/terraform/helper/schema"
)

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
	keymapName         string
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

func resourceOvmVm() *schema.Resource {
	return &schema.Resource{
		Create: resourceOvmVmCreate,
		Read:   resourceOvmVmRead,
		Delete: resourceOvmVmDelete, /*
			Update: resourceOvmCheckUpdate,
			Delete: resourceOvmCheckDelete,
			Importer: &schema.ResourceImporter{
				State: resourceOvmCheckImporter,
			},*/

		Schema: map[string]*schema.Schema{
			/*			"Id": &schema.Schema{
						Type:     schema.TypeString,
						Required: true,
						ForceNew: false,
					},*/
			"repositoryid": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"serverpoolid": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"cpucount": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			/*			"hugepagesenabled": &schema.Schema{
									Type:     schema.TypeBool,
									Required: false,
									ForceNew: false,
								},
						"memory": &schema.Schema{
							Type:     schema.TypeInt,
							Required: false,
							ForceNew: true,
						},*/
			"vmdomaintype": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func checkForResource(d *schema.ResourceData) (ovmHelper.Vm, error) {

	vmParams := &ovmHelper.Vm{}

	// required
	if v, ok := d.GetOk("repositoryid"); ok {
		vmParams.RepositoryId = ovmHelper.Id{Value: v.(string),
			Type: "com.oracle.ovm.mgr.ws.model.Repository"}
	}

	if v, ok := d.GetOk("serverpoolid"); ok {
		vmParams.ServerPoolId = ovmHelper.Id{Value: v.(string),
			Type: "com.oracle.ovm.mgr.ws.model.ServerPool"}
	}

	if v, ok := d.GetOk("vmdomaintype"); ok {
		vmParams.VmDomainType = v.(string)
	}

	//Optinal parameters
	if v, ok := d.GetOk("cpucount"); ok {
		vmParams.CpuCount = v.(int)
	}

	return *vmParams, nil
}

func resourceOvmVmCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ovmHelper.Client)

	vm, err := checkForResource(d)
	if err != nil {
		return err
	}

	//	log.Printf("[DEBUG] Check create configuration: %#v, %#v", d.Get("name"), d.Get("hostname"))

	v, err := client.Vms.CreateVm(vm)
	if err != nil {
		return err
	}

	d.SetId(*v)

	return nil
}

func resourceOvmVmRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceOvmVmDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
