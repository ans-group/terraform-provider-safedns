package safedns

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	safednsservice "github.com/ukfast/sdk-go/pkg/service/safedns"
	"log"
)

func resourceZone() *schema.Resource {
	return &schema.Resource{
		Create: resourceZoneCreate,
		Read:   resourceZoneRead,
		Update: resourceZoneUpdate,
		Delete: resourceZoneDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
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

func resourceZoneCreate(d *schema.ResourceData, meta interface{}) error {
	service := meta.(safednsservice.SafeDNSService)

	zoneName := d.Get("name").(string)

	createReq := safednsservice.CreateZoneRequest{
		Name:        zoneName,
		Description: d.Get("description").(string),
	}
	log.Printf("Created CreateZoneRequest: %+v", createReq)

	log.Print("Creating zone")
	err := service.CreateZone(createReq)
	if err != nil {
		return fmt.Errorf("Error creating zone: %s", err)
	}

	d.SetId(zoneName)

	return resourceZoneRead(d, meta)
}

func resourceZoneRead(d *schema.ResourceData, meta interface{}) error {
	service := meta.(safednsservice.SafeDNSService)

	zoneName := d.Id()

	log.Printf("Retrieving zone with name [%s]", zoneName)
	zone, err := service.GetZone(zoneName)
	if err != nil {
		switch err.(type) {
		case *safednsservice.ZoneNotFoundError:
			d.SetId("")
			return nil
		default:
			return err
		}
	}

	d.Set("name", zone.Name)
	d.Set("description", zone.Description)

	return nil
}

func resourceZoneUpdate(d *schema.ResourceData, meta interface{}) error {
	service := meta.(safednsservice.SafeDNSService)

	patchRequest := safednsservice.PatchZoneRequest{}

	zoneName := d.Id()

	if d.HasChange("description") {
		patchRequest.Description = d.Get("description").(string)
	}

	log.Printf("Updating zone with name [%s]", zoneName)
	err := service.PatchZone(zoneName, patchRequest)
	if err != nil {
		return fmt.Errorf("Error updating zone with name [%s]", zoneName)
	}

	return resourceZoneRead(d, meta)
}

func resourceZoneDelete(d *schema.ResourceData, meta interface{}) error {
	service := meta.(safednsservice.SafeDNSService)

	zoneName := d.Id()

	log.Printf("Removing zone with name [%s]", zoneName)
	err := service.DeleteZone(zoneName)
	if err != nil {
		return fmt.Errorf("Error removing zone with name [%s]: %s", zoneName, err)
	}

	return nil
}
