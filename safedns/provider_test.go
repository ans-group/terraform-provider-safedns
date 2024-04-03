package safedns

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var testAccProviders map[string]func() (*schema.Provider, error)
var testAccProvider *schema.Provider

var (
	UKF_TEST_ZONE_NAME   = os.Getenv("UKF_TEST_ZONE_NAME")
	UKF_TEST_RECORD_NAME = os.Getenv("UKF_TEST_RECORD_NAME")
)

func init() {
	testAccProvider = Provider()
	testAccProviders = map[string]func() (*schema.Provider, error){
		"safedns": func() (*schema.Provider, error) { return testAccProvider, nil },
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ *schema.Provider = Provider()
}

func testAccPreCheck(t *testing.T) {
	testAccPreCheckRequiredEnvVars(t)
}

func testAccPreCheckRequiredEnvVars(t *testing.T) {
	if UKF_TEST_ZONE_NAME == "" {
		t.Fatal("UKF_TEST_ZONE_NAME must be set for acceptance tests")
	}
	if UKF_TEST_RECORD_NAME == "" {
		t.Fatal("UKF_TEST_RECORD_NAME must be set for acceptance tests")
	}
}
