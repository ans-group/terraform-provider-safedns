# SafeDNS Provider

Official UKFast SafeDNS Terraform provider, allowing for manipulation of SafeDNS

## Example Usage

```hcl
provider "safedns" {
  api_key = "abc"
}

resource "safedns_zone" "zone-1" {
    name = "example.com"
    description = "example zone"
}
```

## Argument Reference

* `api_key`: UKFast API key - read/write permissions for `safedns` service required. If omitted, will use `UKF_API_KEY` environment variable value