package main

import (
	"github.com/ans-group/terraform-provider-safedns/safedns"
	"github.com/hashicorp/terraform/plugin"
	"github.com/hashicorp/terraform/terraform"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() terraform.ResourceProvider {
			return safedns.Provider()
		},
	})
}
