# terraform-provider-ovm #

### Dependencies ###

You should have a working Go environment setup.  If not check out the Go [getting started](http://golang.org/doc/install) guide.

This relies on the [go-pingdom](https://github.com/dbgeek/go-ovm-helper) library. To
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


