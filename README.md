


## Usage ##

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


