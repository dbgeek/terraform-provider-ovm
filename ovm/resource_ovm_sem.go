package ovm

import (
	"log"

	"github.com/dbgeek/go-ovm-helper/ovmHelper"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceOvmSem() *schema.Resource {
	return &schema.Resource{
		Create: resourceOvmSemCreate,
		Read:   resourceOvmSemRead,
		Delete: resourceOvmSemDelete,

		//		Update: resourceOvmVmdUpdate,
		/*			Importer: &schema.ResourceImporter{
					State: resourceOvmCheckImporter,
				},*/

		Schema: map[string]*schema.Schema{
			"vmid": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"wwid": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"slot": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
		},
	}
}

func checkForResourceSem(d *schema.ResourceData) (ovmHelper.Vdm, error) {

	vdmParams := &ovmHelper.Vdm{}

	// required
	if v, ok := d.GetOk("vmid"); ok {
		vdmParams.VmId = &ovmHelper.Id{Value: v.(string),
			Type: "com.oracle.ovm.mgr.ws.model.Vm"}
	}
	if v, ok := d.GetOk("vdid"); ok {
		vdmParams.VirtualDiskId = &ovmHelper.Id{Value: v.(string),
			Type: "com.oracle.ovm.mgr.ws.model.VirtualDisk"}
	}
	if v, ok := d.GetOk("slot"); ok {
		vdmParams.DiskTarget = v.(int)
		log.Printf("[DEBUG] Slot: %v DiskTarget: %v", v.(int), vdmParams.DiskTarget)
	}
	//optional
	if v, ok := d.GetOk("description"); ok {
		vdmParams.Description = v.(string)
	}
	if v, ok := d.GetOk("name"); ok {
		vdmParams.Name = v.(string)
	}

	return *vdmParams, nil
}

func resourceOvmSemRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ovmHelper.Client)

	vdm, _ := client.Vdms.Read(d.Get("vmid").(string), d.Id())

	if vdm == nil {
		d.SetId("")
		return nil
	}

	d.Set("vmid", vdm.VmId)
	d.Set("wwid", vdm.VirtualDiskId)
	d.Set("slot", vdm.DiskTarget)
	d.Set("description", vdm.Description)
	d.Set("name", vdm.Name)
	return nil
}

func resourceOvmSemCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ovmHelper.Client)

	vdm, err := checkForResourceSem(d)
	if err != nil {
		return err
	}
	log.Printf("[INFO] Creating vdm for vmid: %v, vdid: %v, slot: %v", vdm.VmId.Value, vdm.VirtualDiskId.Value, vdm.DiskTarget)

	v, err := client.Vdms.Create(vdm)
	if err != nil {
		return err
	}

	d.SetId(*v)

	return nil
}

func resourceOvmSemDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ovmHelper.Client)
	//log.Printf("[INFO] Deleting Vdm: %v", d.Id())

	err := client.Vdms.Delete(d.Get("vmid").(string), d.Id())
	if err != nil {
		return err
	}
	return nil
}
