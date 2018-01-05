package ovm

import (
	"github.com/dbgeek/go-ovm-helper/ovmHelper"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceOvmVmcsm() *schema.Resource {
	return &schema.Resource{
		Create: resourceOvmVmcsmCreate,
		Read:   resourceOvmVmcsmRead,
		Delete: resourceOvmVmcsmDelete,

		//		Update: resourceOvmVmdUpdate,
		/*			Importer: &schema.ResourceImporter{
					State: resourceOvmCheckImporter,
				},*/

		Schema: map[string]*schema.Schema{
			"vmdiskmappingid": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vmclonedefinitionid": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"repositoryid": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"clonetype": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func checkForResourceVmcsm(d *schema.ResourceData) (ovmHelper.Vmcsm, error) {

	vmcsmParams := &ovmHelper.Vmcsm{}

	// required
	if v, ok := d.GetOk("vmdiskmappingid"); ok {
		vmcsmParams.VmDiskMappingId = &ovmHelper.Id{Value: v.(string),
			Type: "com.oracle.ovm.mgr.ws.model.VmDiskMapping"}
	}
	if v, ok := d.GetOk("vmclonedefinitionid"); ok {
		vmcsmParams.VmCloneDefinitionId = &ovmHelper.Id{Value: v.(string),
			Type: "com.oracle.ovm.mgr.ws.model.VmCloneDefinition"}
	}
	if v, ok := d.GetOk("repositoryid"); ok {
		vmcsmParams.RepositoryId = &ovmHelper.Id{Value: v.(string),
			Type: "com.oracle.ovm.mgr.ws.model.Repository"}
	}
	if v, ok := d.GetOk("clonetype"); ok {
		vmcsmParams.CloneType = v.(string)
	}
	if v, ok := d.GetOk("name"); ok {
		vmcsmParams.Name = v.(string)
	}
	return *vmcsmParams, nil
}

func resourceOvmVmcsmRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ovmHelper.Client)

	vmcsm, _ := client.Vmcsms.Read(d.Id())

	if vmcsm == nil {
		d.SetId("")
		return nil
	}

	d.Set("vmdiskmappingid", vmcsm.VmDiskMappingId.Value)
	d.Set("vmclonedefinitionid", vmcsm.VmCloneDefinitionId)
	d.Set("repositoryid", vmcsm.RepositoryId.Value)
	d.Set("clonetype", vmcsm.CloneType)
	d.Set("name", vmcsm.Name)
	return nil
}

func resourceOvmVmcsmCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ovmHelper.Client)

	vmcsm, err := checkForResourceVmcsm(d)
	if err != nil {
		return err
	}
	//log.Printf("[INFO] Creating vdm for vmid: %v, vdid: %v, slot: %v", vdm.VmId.Value, vdm.VirtualDiskId.Value, vdm.DiskTarget)

	v, err := client.Vmcsms.Create(vmcsm.VmCloneDefinitionId.Value, vmcsm)
	if err != nil {
		return err
	}

	d.SetId(*v)

	return nil
}

func resourceOvmVmcsmDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ovmHelper.Client)
	//log.Printf("[INFO] Deleting Vdm: %v", d.Id())

	err := client.Vmcsms.Delete(d.Get("vmclonedefinitionid").(string), d.Id())
	if err != nil {
		return err
	}
	return nil
}
