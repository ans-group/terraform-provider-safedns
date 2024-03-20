package safedns

import (
	"fmt"
	"log"
	"strconv"

	"github.com/ans-group/sdk-go/pkg/ptr"
	safednsservice "github.com/ans-group/sdk-go/pkg/service/safedns"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceRecord() *schema.Resource {
	return &schema.Resource{
		Create: resourceRecordCreate,
		Read:   resourceRecordRead,
		Update: resourceRecordUpdate,
		Delete: resourceRecordDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"zone_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"content": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"priority": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceRecordCreate(d *schema.ResourceData, meta interface{}) error {
	service := meta.(safednsservice.SafeDNSService)

	createReq := safednsservice.CreateRecordRequest{
		Name:    d.Get("name").(string),
		Type:    d.Get("type").(string),
		Content: d.Get("content").(string),
	}

	priority := d.Get("priority").(int)
	if priority > 0 {
		createReq.Priority = ptr.Int(priority)
	}
	log.Printf("Created CreateRecordRequest: %+v", createReq)

	log.Print("Creating record")
	recordID, err := service.CreateZoneRecord(d.Get("zone_name").(string), createReq)
	if err != nil {
		return fmt.Errorf("Error creating record: %s", err)
	}

	d.SetId(strconv.Itoa(recordID))

	return resourceRecordRead(d, meta)
}

func resourceRecordRead(d *schema.ResourceData, meta interface{}) error {
	service := meta.(safednsservice.SafeDNSService)

	zoneName := d.Get("zone_name").(string)
	recordID, _ := strconv.Atoi(d.Id())

	log.Printf("Retrieving record with id [%d] in zone [%s]", recordID, zoneName)
	record, err := service.GetZoneRecord(zoneName, recordID)
	if err != nil {
		switch err.(type) {
		case *safednsservice.ZoneRecordNotFoundError:
			d.SetId("")
			return nil
		default:
			return err
		}
	}

	d.Set("name", record.Name)
	d.Set("type", record.Type)
	d.Set("content", record.Content)
	d.Set("priority", record.Priority)

	return nil
}

func resourceRecordUpdate(d *schema.ResourceData, meta interface{}) error {
	service := meta.(safednsservice.SafeDNSService)

	patchRequest := safednsservice.PatchRecordRequest{}

	zoneName := d.Get("zone_name").(string)
	recordID, _ := strconv.Atoi(d.Id())

	if d.HasChange("content") {
		patchRequest.Content = d.Get("content").(string)
	}
	if d.HasChange("priority") {
		patchRequest.Priority = ptr.Int(d.Get("priority").(int))
	}

	log.Printf("Updating record with id [%d]", recordID)
	_, err := service.PatchZoneRecord(zoneName, recordID, patchRequest)
	if err != nil {
		return fmt.Errorf("Error updating record with id [%d]", recordID)
	}

	return resourceRecordRead(d, meta)
}

func resourceRecordDelete(d *schema.ResourceData, meta interface{}) error {
	service := meta.(safednsservice.SafeDNSService)

	recordID, _ := strconv.Atoi(d.Id())

	log.Printf("Removing record with id [%d]", recordID)
	err := service.DeleteZoneRecord(d.Get("zone_name").(string), recordID)
	if err != nil {
		return fmt.Errorf("Error removing record with id [%d]: %s", recordID, err)
	}

	return nil
}
