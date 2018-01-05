package ovm

import (
	"github.com/dbgeek/go-ovm-helper/ovmHelper"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceOvmVmcnm() *schema.Resource {
	return &schema.Resource{
		Create: resourceOvmVmcnmCreate,
		Read:   resourceOvmVmcnmRead,
		Delete: resourceOvmVmcnmDelete,

		//		Update: resourceOvmVmdUpdate,
		/*			Importer: &schema.ResourceImporter{
					State: resourceOvmCheckImporter,
				},*/

		Schema: map[string]*schema.Schema{
			"networkid": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vmclonedefinitionid": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"virtualnicid": &schema.Schema{
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

func checkForResourceVmcnm(d *schema.ResourceData) (ovmHelper.Vmcnm, error) {

	vmcnmParams := &ovmHelper.Vmcnm{}

	// required
	if v, ok := d.GetOk("networkid"); ok {
		vmcnmParams.NetworkId = &ovmHelper.Id{Value: v.(string),
			Type: "com.oracle.ovm.mgr.ws.model.Network"}
	}
	if v, ok := d.GetOk("vmclonedefinitionid"); ok {
		vmcnmParams.VmCloneDefinitionId = &ovmHelper.Id{Value: v.(string),
			Type: "com.oracle.ovm.mgr.ws.model.VmCloneDefinition"}
	}
	if v, ok := d.GetOk("virtualnicid"); ok {
		vmcnmParams.VirtualNicId = &ovmHelper.Id{Value: v.(string),
			Type: "com.oracle.ovm.mgr.ws.model.VirtualNic"}
	}
	if v, ok := d.GetOk("name"); ok {
		vmcnmParams.Name = v.(string)
	}
	return *vmcnmParams, nil
}

func resourceOvmVmcnmRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ovmHelper.Client)

	vmcnm, _ := client.Vmcnms.Read(d.Id())

	if vmcnm == nil {
		d.SetId("")
		return nil
	}

	d.Set("networkid", vmcnm.NetworkId)
	d.Set("vmclonedefinitionid", vmcnm.VmCloneDefinitionId)
	d.Set("virtualnicid", vmcnm.VirtualNicId)
	d.Set("name", vmcnm.Name)
	return nil
}

func resourceOvmVmcnmCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ovmHelper.Client)

	vmcnm, err := checkForResourceVmcnm(d)
	if err != nil {
		return err
	}
	//log.Printf("[INFO] Creating vdm for vmid: %v, vdid: %v, slot: %v", vdm.VmId.Value, vdm.VirtualDiskId.Value, vdm.DiskTarget)

	v, err := client.Vmcnms.Create(vmcnm.VmCloneDefinitionId.Value, vmcnm)
	if err != nil {
		return err
	}

	d.SetId(*v)

	return nil
}

func resourceOvmVmcnmDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ovmHelper.Client)
	//log.Printf("[INFO] Deleting Vdm: %v", d.Id())

	err := client.Vmcnms.Delete(d.Get("vmclonedefinitionid").(string), d.Id())
	if err != nil {
		return err
	}
	return nil
}
