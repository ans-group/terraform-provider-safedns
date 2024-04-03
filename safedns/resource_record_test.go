package safedns

import (
	"errors"
	"fmt"
	"strconv"
	"testing"

	safednsservice "github.com/ans-group/sdk-go/pkg/service/safedns"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccRecord_basic(t *testing.T) {
	var record safednsservice.Record
	resourceName := "safedns_record.test-record"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckRecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckRecordConfig_basic, UKF_TEST_ZONE_NAME, UKF_TEST_RECORD_NAME),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRecordExists(resourceName, &record),
					resource.TestCheckResourceAttr(resourceName, "zone_name", UKF_TEST_ZONE_NAME),
					resource.TestCheckResourceAttr(resourceName, "type", "A"),
					resource.TestCheckResourceAttr(resourceName, "content", "10.0.0.1"),
				),
			},
		},
	})
}

func testAccCheckRecordExists(n string, record *safednsservice.Record) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No record ID is set")
		}

		service := testAccProvider.Meta().(safednsservice.SafeDNSService)

		recordID, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return err
		}

		zoneName := rs.Primary.Attributes["zone_name"]

		getRecord, err := service.GetZoneRecord(zoneName, recordID)
		if err != nil {
			var zoneRecordNotFoundError *safednsservice.ZoneRecordNotFoundError
			if errors.As(err, &zoneRecordNotFoundError) {
				return nil
			}
			return err
		}

		*record = getRecord

		return nil
	}
}

func testAccCheckRecordDestroy(s *terraform.State) error {
	service := testAccProvider.Meta().(safednsservice.SafeDNSService)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "safedns_record" {
			continue
		}

		recordID, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return err
		}

		zoneName := rs.Primary.Attributes["zone_name"]

		_, err = service.GetZoneRecord(zoneName, recordID)
		if err == nil {
			return fmt.Errorf("Record with id [%d] still exists", recordID)
		}

		var zoneRecordNotFoundError *safednsservice.ZoneRecordNotFoundError
		if errors.As(err, &zoneRecordNotFoundError) {
			return nil
		}

		return err
	}

	return nil
}

var testAccCheckRecordConfig_basic = `
resource "safedns_zone" "test-zone" {
  name = "%s"
}
resource "safedns_record" "test-record" {
	zone_name = "${safedns_zone.test-zone.name}"
	name = "%s"
	type = "A"
	content = "10.0.0.1"
  }
`
