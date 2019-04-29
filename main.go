package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/hashicorp/terraform/terraform"
	"gitlab.devops.ukfast.co.uk/ukfast/api.ukfast/client-libraries/terraform-provider-safedns/safedns"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() terraform.ResourceProvider {
			return safedns.Provider()
		},
	})
}
