package ovm

import (
	"fmt"
	"log"

	"github.com/dbgeek/go-ovm-helper/ovmHelper"
	"github.com/hashicorp/terraform/helper/schema"
)

type tfVmCfg struct {
	networkId string
}

func resourceOvmVm() *schema.Resource {
	return &schema.Resource{
		Create: resourceOvmVmCreate,
		Read:   resourceOvmVmRead,
		Delete: resourceOvmVmDelete,
		Update: resourceOvmVmUpdate,
		/*			Importer: &schema.ResourceImporter{
					State: resourceOvmCheckImporter,
				},*/

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: false,
				ForceNew: false,
				Optional: true,
				Computed: true,
			},
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
				Optional: true,
				Computed: true,
				Required: false,
				ForceNew: false,
			},
			"cpucountlimit": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				Required: false,
				ForceNew: false,
			},
			"hugepagesenabled": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
				Required: false,
				ForceNew: false,
			},
			"memory": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				Required: false,
				ForceNew: false,
			},
			"vmdomaintype": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Required: false,
				ForceNew: false,
			},
			"clonevmid": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Required: false,
				ForceNew: true,
			},
			"vmclonedefinitionid": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Required: false,
				ForceNew: true,
			},
			"networkid": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Required: false,
				ForceNew: true,
			},
			"sendmessages": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
			},
			"startvm": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
		},
	}
}

func checkForResource(d *schema.ResourceData) (ovmHelper.Vm, ovmHelper.CfgVm, error) {

	vmParams := &ovmHelper.Vm{}
	tfVmCfgParams := &ovmHelper.CfgVm{}

	// required
	if v, ok := d.GetOk("repositoryid"); ok {
		vmParams.RepositoryId = &ovmHelper.Id{Value: v.(string),
			Type: "com.oracle.ovm.mgr.ws.model.Repository"}
	}

	if v, ok := d.GetOk("serverpoolid"); ok {
		vmParams.ServerPoolId = &ovmHelper.Id{Value: v.(string),
			Type: "com.oracle.ovm.mgr.ws.model.ServerPool"}
	}

	if v, ok := d.GetOk("vmdomaintype"); ok {
		vmParams.VmDomainType = v.(string)
	}

	//Optinal parameters
	if v, ok := d.GetOk("cpucount"); ok {
		vmParams.CpuCount = v.(int)
	}
	if v, ok := d.GetOk("cpucountlimit"); ok {
		vmParams.CpuCountLimit = v.(int)
	}
	if v, ok := d.GetOk("name"); ok {
		vmParams.Name = v.(string)
	}
	if v, ok := d.GetOk("hugepagesenabled"); ok {
		vmParams.HugePagesEnabled = v.(bool)
	}
	if v, ok := d.GetOk("memory"); ok {
		vmParams.Memory = v.(int)
	}
	if v, ok := d.GetOk("networkid"); ok {
		tfVmCfgParams.NetworkId = v.(string)
	}

	if v, ok := d.GetOk("sendmessages"); ok {
		sendmessages, rootPassword := sendmessagesFromMap(v.(map[string]interface{}))
		tfVmCfgParams.SendMessages = sendmessages
		tfVmCfgParams.RootPassword = rootPassword
	}

	return *vmParams, *tfVmCfgParams, nil
}

func resourceOvmVmCreate(d *schema.ResourceData, meta interface{}) error {
	var v *string

	client := meta.(*ovmHelper.Client)

	vm, tfVmCfgParams, err := checkForResource(d)
	if err != nil {
		return err
	}

	if d.Get("clonevmid").(string) == "" {
		v, err = client.Vms.CreateVm(vm, tfVmCfgParams)
		if err != nil {
			return err
		}
	} else {
		v, err = client.Vms.CloneVm(d.Get("clonevmid").(string), d.Get("vmclonedefinitionid").(string), vm, tfVmCfgParams)
		if err != nil {
			return err
		}

	}

	if d.Get("startvm").(bool) {
		err = client.Vms.Start(*v)
		if err != nil {
			return err
		}
	}

	d.SetId(*v)

	return nil
}

func resourceOvmVmRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ovmHelper.Client)

	vm, _ := client.Vms.Read(d.Id())

	if vm == nil {
		d.SetId("")
		fmt.Println("Not find any vm")
		return nil
	}

	d.Set("repositoryid", vm.RepositoryId)
	d.Set("serverpoolid", vm.ServerPoolId)
	d.Set("vmdomaintype", vm.VmDomainType)
	d.Set("cpucount", vm.CpuCount)
	return nil
}

func resourceOvmVmDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ovmHelper.Client)
	log.Printf("[INFO] Deleting Vm: %v", d.Id())
	err := client.Vms.DeleteVm(d.Id())
	if err != nil {
		return err
	}
	return nil
}

func resourceOvmVmUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ovmHelper.Client)
	vm, _, err := checkForResource(d)
	if err != nil {
		return err
	}
	err = client.Vms.UpdateVm(d.Id(), vm)
	if err != nil {
		return err
	}
	return nil
}

func sendmessagesFromMap(m map[string]interface{}) (*[]ovmHelper.KeyValuePair, *[]ovmHelper.KeyValuePair) {

	result := make([]ovmHelper.KeyValuePair, 0, len(m))
	password := make([]ovmHelper.KeyValuePair, 0, len(m))
	for k, v := range m {
		t := ovmHelper.KeyValuePair{
			Key:   k,
			Value: v.(string),
		}
		if k == "com.oracle.linux.root-password" {
			password = append(password, t)
		} else {
			result = append(result, t)
		}
	}

	return &result, &password
}
