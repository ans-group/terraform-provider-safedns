package safedns

import (
	"context"

	safednsservice "github.com/ans-group/sdk-go/pkg/service/safedns"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceZone() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceZoneRead,

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

func dataSourceZoneRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(safednsservice.SafeDNSService)

	name := d.Get("name").(string)

	zone, err := service.GetZone(name)
	if err != nil {
		return diag.Errorf("Error retrieving zone: %s", err)
	}

	d.SetId(zone.Name)
	d.Set("name", zone.Name)
	d.Set("description", zone.Description)

	return nil
}
