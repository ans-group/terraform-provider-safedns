package safedns

import (
	"fmt"
	"testing"

	safednsservice "github.com/ans-group/sdk-go/pkg/service/safedns"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceZone(t *testing.T) {
	var zone safednsservice.Zone

	resourceName := "data.safedns_zone.test-zone"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDataSourceZoneConfig_basic, UKF_TEST_ZONE_NAME),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceZoneExists(resourceName, &zone),
					resource.TestCheckResourceAttr(resourceName, "name", UKF_TEST_ZONE_NAME),
				),
			},
		},
	})
}

func testAccCheckDataSourceZoneExists(n string, zone *safednsservice.Zone) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No zone ID is set")
		}

		service := testAccProvider.Meta().(safednsservice.SafeDNSService)

		getZone, err := service.GetZone(rs.Primary.ID)
		if err != nil {
			return err
		}

		*zone = getZone

		return nil
	}
}

var testAccCheckDataSourceZoneConfig_basic = `
data "safedns_zone" "test-zone" {
	name = "%s"
}`
