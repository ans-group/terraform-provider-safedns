package safedns

import (
	"errors"
	"os"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/ukfast/sdk-go/pkg/client"
	"github.com/ukfast/sdk-go/pkg/connection"
	safednsservice "github.com/ukfast/sdk-go/pkg/service/safedns"
)

func Provider() *schema.Provider {
	return &schema.Provider{ 
		Schema: map[string]*schema.Schema{ 
			"api_key": &schema.Schema{
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
				DefaultFunc: func() (interface{}, error) {
					key := os.Getenv("UKF_API_KEY")
					if key != "" {
						return key, nil
					}
 
					return "", errors.New("api_key required")
				},
				Description: "API token required to authenticate with UKFast APIs. See https://developers.ukfast.io for more details",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"safedns_zone":   resourceZone(),
			"safedns_record": resourceRecord(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	return getService(d.Get("api_key").(string)), nil
}

func getClient(apiKey string) client.Client {
	return client.NewClient(connection.NewAPIKeyCredentialsAPIConnection(apiKey))
}

func getService(apiKey string) safednsservice.SafeDNSService {
	return getClient(apiKey).SafeDNSService()
}
