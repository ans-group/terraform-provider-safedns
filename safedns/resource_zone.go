package safedns

import (
	"context"
	"errors"
	"log"

	safednsservice "github.com/ans-group/sdk-go/pkg/service/safedns"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceZone() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceZoneCreate,
		ReadContext:   resourceZoneRead,
		UpdateContext: resourceZoneUpdate,
		DeleteContext: resourceZoneDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceZoneCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(safednsservice.SafeDNSService)

	zoneName := d.Get("name").(string)

	createReq := safednsservice.CreateZoneRequest{
		Name:        zoneName,
		Description: d.Get("description").(string),
	}
	tflog.Debug(ctx, "Created CreateZoneRequest", map[string]interface{}{
		"createReq": createReq,
	})

	tflog.Info(ctx, "Creating zone")

	err := service.CreateZone(createReq)
	if err != nil {
		return diag.Errorf("Error creating zone: %s", err)
	}

	d.SetId(zoneName)

	return resourceZoneRead(ctx, d, meta)
}

func resourceZoneRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(safednsservice.SafeDNSService)

	zoneName := d.Id()

	log.Printf("Retrieving zone with name [%s]", zoneName)
	zone, err := service.GetZone(zoneName)
	if err != nil {
		var zoneNotFoundError *safednsservice.ZoneNotFoundError
		switch {
		case errors.As(err, &zoneNotFoundError):
			d.SetId("")
			return nil
		default:
			return diag.FromErr(err)
		}
	}

	d.Set("name", zone.Name)
	d.Set("description", zone.Description)

	return nil
}

func resourceZoneUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(safednsservice.SafeDNSService)

	patchRequest := safednsservice.PatchZoneRequest{}

	zoneName := d.Id()

	if d.HasChange("description") {
		patchRequest.Description = d.Get("description").(string)
	}

	tflog.Info(ctx, "Updating zone", map[string]interface{}{
		"zone_name": zoneName,
	})

	err := service.PatchZone(zoneName, patchRequest)
	if err != nil {
		return diag.Errorf("Error updating zone with name [%s]", zoneName)
	}

	return resourceZoneRead(ctx, d, meta)
}

func resourceZoneDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(safednsservice.SafeDNSService)

	zoneName := d.Id()

	tflog.Info(ctx, "Removing zone", map[string]interface{}{
		"zone_name": zoneName,
	})

	err := service.DeleteZone(zoneName)
	if err != nil {
		return diag.Errorf("Error removing zone with name [%s]: %s", zoneName, err)
	}

	return nil
}
