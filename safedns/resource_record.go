package safedns

import (
	"context"
	"errors"
	"strconv"

	"github.com/ans-group/sdk-go/pkg/ptr"
	safednsservice "github.com/ans-group/sdk-go/pkg/service/safedns"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceRecord() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRecordCreate,
		ReadContext:   resourceRecordRead,
		UpdateContext: resourceRecordUpdate,
		DeleteContext: resourceRecordDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
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

func resourceRecordCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
	tflog.Debug(ctx, "Created CreateRecordRequest", map[string]interface{}{
		"request": createReq,
	})

	tflog.Info(ctx, "Creating record")
	recordID, err := service.CreateZoneRecord(d.Get("zone_name").(string), createReq)
	if err != nil {
		return diag.Errorf("Error creating record: %s", err)
	}

	d.SetId(strconv.Itoa(recordID))

	return resourceRecordRead(ctx, d, meta)
}

func resourceRecordRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(safednsservice.SafeDNSService)

	zoneName := d.Get("zone_name").(string)
	recordID, _ := strconv.Atoi(d.Id())

	tflog.Info(ctx, "Retrieving record", map[string]interface{}{
		"zone_name": zoneName,
		"record_id": recordID,
	})

	record, err := service.GetZoneRecord(zoneName, recordID)
	if err != nil {
		var zoneRecordNotFoundError *safednsservice.ZoneRecordNotFoundError
		switch {
		case errors.As(err, &zoneRecordNotFoundError):
			d.SetId("")
			return nil
		default:
			return diag.FromErr(err)
		}
	}

	d.Set("name", record.Name)
	d.Set("type", record.Type)
	d.Set("content", record.Content)
	d.Set("priority", record.Priority)

	return nil
}

func resourceRecordUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	tflog.Info(ctx, "Updating record", map[string]interface{}{
		"record_id": recordID,
	})

	_, err := service.PatchZoneRecord(zoneName, recordID, patchRequest)
	if err != nil {
		return diag.Errorf("Error updating record with id [%d]", recordID)
	}

	return resourceRecordRead(ctx, d, meta)
}

func resourceRecordDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	service := meta.(safednsservice.SafeDNSService)

	recordID, _ := strconv.Atoi(d.Id())

	tflog.Info(ctx, "Deleting record", map[string]interface{}{
		"record_id": recordID,
	})

	err := service.DeleteZoneRecord(d.Get("zone_name").(string), recordID)
	if err != nil {
		return diag.Errorf("Error removing record with id [%d]: %s", recordID, err)
	}

	return nil
}
