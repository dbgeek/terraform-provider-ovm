package main

import (
	"github.com/hashicorp/terraform/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() ovm.Provider {
			return Provider()
		},
	})
}
