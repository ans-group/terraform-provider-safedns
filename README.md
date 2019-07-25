# terraform-provider-safedns

## Getting Started

To get started, the `terraform-provider-safedns` binary (`.exe` extension if Windows) should be downloaded from [Releases](https://github.com/ukfast/terraform-provider-safedns/releases) and placed in a directory. For this example,
we'll place it at `/tmp/terraform-provider-safedns`.

Next, we'll go ahead and create a new directory to hold our `terraform` file and state:

```console
mkdir /home/user/terraform
```

We'll then create an example terraform file `/home/user/terraform/test.tf`:

```console
cat <<EOF > /home/user/terraform/test.tf
provider "safedns" {
  api_key = "abc"
}

resource "safedns_zone" "zone-1" {
    name = "example.com"
    description = "example zone"
}
EOF
```

We'll then need to initialise terraform with our provider (specifying `plugin-dir` as the path to where the provider was downloaded to earlier):

```console
terraform init -get-plugins=false -plugin-dir=/tmp/terraform-provider-safedns
```

Finally, we can invoke `terraform apply` to apply our terraform configuration:

```console
terraform apply
```

## Provider

**Parameters**

- `api_key`: UKFast API key - read/write permissions for `safedns` service required. If omitted, will use `UKF_API_KEY` environment variable value

## Resources

### safedns_zone

**Schema**

- `name`: (Required) Name of zone e.g. `example.com`
- `description`: Description for zone

### safedns_record

**Schema**

- `name`: (Required) Name of record e.g. `something.example.com`
- `zone_name`: (Required) Name of zone for record e.g. `example.com`
- `type`: (Required) Type of record
- `content`: (Required) Content for record
- `priority`: Priority of 

## Data sources

### safedns_zone

**Schema**

- `name`: (Required) Name of zone e.g. `example.com`
- `description`: Description for zone