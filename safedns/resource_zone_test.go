package safedns

import (
	"errors"
	"fmt"
	"testing"

	safednsservice "github.com/ans-group/sdk-go/pkg/service/safedns"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccZone_basic(t *testing.T) {
	var zone safednsservice.Zone
	resourceName := "safedns_zone.test-zone"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckZoneDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckZoneConfig_basic, UKF_TEST_ZONE_NAME),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckZoneExists(resourceName, &zone),
					resource.TestCheckResourceAttr(resourceName, "name", UKF_TEST_ZONE_NAME),
					resource.TestCheckResourceAttr(resourceName, "description", "test description"),
				),
			},
		},
	})
}

func testAccCheckZoneExists(n string, zone *safednsservice.Zone) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No zone ID is set")
		}

		service := testAccProvider.Meta().(safednsservice.SafeDNSService)

		zoneName := rs.Primary.ID

		getZone, err := service.GetZone(zoneName)
		if err != nil {
			var zoneNotFoundError *safednsservice.ZoneNotFoundError
			if errors.As(err, &zoneNotFoundError) {
				return nil
			}
			return err
		}

		*zone = getZone

		return nil
	}
}

func testAccCheckZoneDestroy(s *terraform.State) error {
	service := testAccProvider.Meta().(safednsservice.SafeDNSService)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "safedns_zone" {
			continue
		}

		zoneName := rs.Primary.ID

		_, err := service.GetZone(zoneName)
		if err == nil {
			return fmt.Errorf("Zone with name [%s] still exists", zoneName)
		}

		var zoneNotFoundError *safednsservice.ZoneNotFoundError
		if errors.As(err, &zoneNotFoundError) {
			return nil
		}

		return err
	}

	return nil
}

var testAccCheckZoneConfig_basic = `
resource "safedns_zone" "test-zone" {
  name = "%s"
  description = "test description"
}`
