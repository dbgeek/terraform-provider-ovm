# terraform-provider-ovm #

# This repo is not maintained anymore
## I can not maintaine this repo anymore as I do not have any lab/play environment

### Dependencies ###

This project is a [terraform](http://www.terraform.io/) provider for [OVM](https://www.oracle.com/virtualization/oracle-vm-server-for-x86/index.html)

You should have a working Go environment setup.  If not check out the Go [getting started](http://golang.org/doc/install) guide.

This relies on the [go-ovm-helper](https://github.com/dbgeek/go-ovm-helper) library. To
get that: `go get github.com/dbgeek/go-pingdom/go-ovm-helper/ovmHelper`.

You'll also need the libraries from terraform.  Check out those docs under [plugin basics](http://www.terraform.io/docs/plugins/basics.html)

### Build ###

Run `go install github.com/dbgeek/terraform-provider-ovm`

### Install ###

Add the following to `$HOME/.terraformrc`

```
providers {
    ovm = "$GOPATH/bin/terraform-provider-ovm"
}
```

## Usage ##

**provider.ovm: new or changed plugin executable**

```
export TF_SKIP_PROVIDER_VERIFY=1
```

**Configure the Provider**

***Configure in TF configuration***

```
provider "ovm" {
  user       = "${var.ovm_username}"
  password   = "${var.ovm_password}"
  entrypoint = "${var.ovm_endpoint}"
}
```

***Configure in environment***

Set username(`OVM_USERNAME`) and password(`OVM_PASSWORD`) and endpoint(`OVM_ENDPOINT`) in environment
```
provider "ovm" {}
```

**Basic vm provision**

Create one vm and create two virtual disks and mapp them to the vm
```
resource "ovm_vm" "vm1" {
  name          = "vm1"
  repositoryid  = "${var.vm_repositoryid}"
  serverpoolid  = "${var.serverpoolid}"
  vmdomaintype  = "XEN_HVM"
  cpucount      = 2
  cpucountlimit = 2
  memory        = 512 //MB
}

resource "ovm_vd" "vm1_virtualdisk" {
  count        = 2
  name         = "vm1_vd${count.index}"
  sparse       = true
  shareable    = false
  repositoryid = "${var.vd_repositoryid}"
  size         = 104857600 //bytes
}

resource "ovm_vdm" "vm1_vdm" {
  count       = 2
  vmid        = "${ovm_vm.vm1.id}"
  vdid        = "${element(ovm_vd.vm1.*.id, count.index)}"
  name        = "vm1_vdm_2${count.index}"
  slot        = "${count.index}"
  description = "Virtual disk mapping for vm1 and vm1_vdm_2${count.index}"
}
```

**Create VM from a Template**

```
//Creating VmCloneCustomizer
resource "ovm_vmcd" "oe7_tmpl_cst" {
  vmid        = "0004fb000014000014feb8708c34fc0f"
  name        = "oe7_tmpl_cst"
  description = "Desc oel7 cust"
}

//Defining Vm Clone Storage Mapping
resource "ovm_vmcsm" "oel7_vmclonestoragemapping" {
  vmdiskmappingid     = "0004fb0000130000f8e1fa844def645e"
  vmclonedefinitionid = "${ovm_vmcd.oe7_tmpl_cst.id}"
  repositoryid        = "0004fb00000300003a68daf22a32ebc5"
  name                = "oel_cust_storage"
  clonetype           = "SPARSE_COPY"
}

//Defining Vm Clone Network Mappings.
resource "ovm_vmcnm" "oel7_vmclonenetworkmapping" {
  networkid           = "${var.networkid}"
  vmclonedefinitionid = "${ovm_vmcd.oe7_tmpl_cst.id}"
  virtualnicid        = "${var.virtualnicid}"
  name                = "oel_cust_network"
}

resource "ovm_vm" "cloneoel7" {
  name                = "cloneoel7Vm"
  repositoryid        = "${var.vm_repositoryid}" //Where vm.cfg should be stored
  serverpoolid        = "${var.serverpoolid}"
  vmdomaintype        = "XEN_HVM"
  clonevmid           = "${var.template_vmid}"
  vmclonedefinitionid = "${ovm_vmcd.oel7_cust.id}"

  sendmessages {
    "com.oracle.linux.network.hostname"    = "cloneoel7vm"
    "com.oracle.linux.network.device.0"    = "eth0"
    "com.oracle.linux.network.bootproto.0" = "dhcp"
    "com.oracle.linux.network.onboot.0"    = "yes"
    "com.oracle.linux.root-password"       = "Welcome!"
  }
 
  depends_on          = ["ovm_vmcnm.oel7_cust_vmcnm", "ovm_vmcsm.oel7_cust_vmcsm"]
}

resource "ovm_vd" "clonevmvd" {
  count        = 2                           //nr of vritual disk to create
  name         = "clonedvm${count.index}"
  sparse       = true
  shareable    = false
  repositoryid = "${var.vd_repositoryid}"
  size         = 104857600 //bytes
}

//Mapping the Virtual Disk to the vm
resource "ovm_vdm" "clonevmvdm" {
  count       = 2
  vmid        = "${ovm_vm.cloneoel7.id}"
  vdid        = "${element(ovm_vd.clonevmvd.*.id, count.index)}"
  name        = "clonevmvdm${count.index +1}"
  slot        = "${count.index + 1}"          //The template has one disk that already attached to slot 0
  description = "Extra disk that get attached to the vm"
}
```
