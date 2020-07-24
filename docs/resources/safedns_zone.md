# safedns_zone Resource

This resource is for managing SafeDNS zones

## Example Usage

```hcl
resource "safedns_zone" "example-1" {
    name = "example.com"
}
```

## Argument Reference

* `name`: (Required) Name of zone e.g. `example.com`
* `description`: Description for zone