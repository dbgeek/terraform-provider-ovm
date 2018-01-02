package ovm

import (
	"fmt"
	"log"

	"github.com/dbgeek/go-ovm-helper/ovmHelper"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceOvmVd() *schema.Resource {
	return &schema.Resource{
		Create: resourceOvmVdCreate,
		Read:   resourceOvmVdRead,
		Delete: resourceOvmVdDelete,
		Update: resourceOvmVdUpdate,
		/*			Importer: &schema.ResourceImporter{
					State: resourceOvmCheckImporter,
				},*/

		Schema: map[string]*schema.Schema{
			"repositoryid": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"size": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: false,
			},
			"sparse": &schema.Schema{
				Type:     schema.TypeBool,
				Required: true,
				ForceNew: false,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: false,
				ForceNew: false,
				Optional: true,
				Computed: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Required: false,
				ForceNew: false,
				Optional: true,
				Computed: true,
			},
			"shareable": &schema.Schema{
				Type:     schema.TypeBool,
				Required: false,
				ForceNew: false,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func checkForResourceVd(d *schema.ResourceData) (ovmHelper.Vd, error) {

	vdParams := &ovmHelper.Vd{}

	// required
	if v, ok := d.GetOk("repositoryid"); ok {
		vdParams.RepositoryId = &ovmHelper.Id{Value: v.(string),
			Type: "com.oracle.ovm.mgr.ws.model.Repository"}
	}
	if v, ok := d.GetOk("size"); ok {
		vdParams.Size = v.(int)
	}
	//Optional
	if v, ok := d.GetOk("name"); ok {
		vdParams.Name = v.(string)
	}
	if v, ok := d.GetOk("description"); ok {
		vdParams.Description = v.(string)
	}
	if v, ok := d.GetOk("shareable"); ok {
		vdParams.Shareable = v.(bool)
	}

	return *vdParams, nil
}

func resourceOvmVdRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ovmHelper.Client)

	vd, _ := client.Vds.Read(d.Id())

	if vd == nil {
		d.SetId("")
		fmt.Println("Not find any vm")
		return nil
	}

	d.Set("repositoryid", vd.RepositoryId)
	d.Set("name", vd.Name)
	d.Set("size", vd.Size)
	d.Set("description", vd.Description)
	d.Set("shareable", vd.Shareable)
	return nil
}

func resourceOvmVdCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ovmHelper.Client)

	vd, err := checkForResourceVd(d)
	if err != nil {
		return err
	}

	v, err := client.Vds.Create(d.Get("repositoryid").(string), d.Get("sparse").(bool), vd)
	if err != nil {
		return err
	}

	d.SetId(*v)

	return nil
}

func resourceOvmVdUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ovmHelper.Client)
	vd, err := checkForResourceVd(d)
	if err != nil {
		return err
	}
	err = client.Vds.Update(d.Id(), vd)
	if err != nil {
		return err
	}
	return nil
}

func resourceOvmVdDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ovmHelper.Client)
	log.Printf("[INFO] Deleting Vd: %v", d.Id())
	err := client.Vds.Delete(d.Get("repositoryid").(string), d.Id())
	if err != nil {
		return err
	}
	return nil
}
