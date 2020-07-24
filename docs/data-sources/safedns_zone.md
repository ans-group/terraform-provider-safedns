# safedns_zone Data Source

This resource represents a SafeDNS zone

## Example Usage

```hcl
data "safedns_zone" "example-1" {
    name = "example.com"
}
```

## Argument Reference

* `name`: (Required) Name of zone e.g. `example.com`
* `description`: Description for zone