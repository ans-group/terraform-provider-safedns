# safedns_record Resource

This resource is for managing SafeDNS records

## Example Usage

```hcl
resource "safedns_record" "example-record-1" {
    name = "something.example.com"
    zone_name = "example.com"
    type = "A"
    content = "10.1.2.3"
}
```

## Argument Reference

* `name`: (Required) Name of record e.g. `something.example.com`
* `zone_name`: (Required) Name of zone for record e.g. `example.com`
* `type`: (Required) Type of record
* `content`: (Required) Content for record
* `priority`: Priority of record