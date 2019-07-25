package safedns

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	safednsservice "github.com/ukfast/sdk-go/pkg/service/safedns"
)

func dataSourceZone() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceZoneRead,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceZoneRead(d *schema.ResourceData, meta interface{}) error {
	service := meta.(safednsservice.SafeDNSService)

	name := d.Get("name").(string)

	zone, err := service.GetZone(name)
	if err != nil {
		return fmt.Errorf("Error retrieving zone: %s", err)
	}

	d.SetId(zone.Name)
	d.Set("name", zone.Name)
	d.Set("description", zone.Description)

	return nil
}
