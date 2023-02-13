# terraform-provider-safedns

## Getting Started

This provider is available via the [Terraform Registry](https://registry.terraform.io/providers/ukfast/safedns/latest) with Terraform v0.13+

## Getting Started (manual)

To get started, the `terraform-provider-safedns` binary (`.exe` extension if Windows) should be downloaded from [Releases](https://github.com/ukfast/terraform-provider-safedns/releases) and placed in the plugins directory (see [here](https://www.terraform.io/docs/configuration/providers.html#third-party-plugins) for more information). For this example, we'll place it in `~/.terraform.d/plugins/`

We'll then need to initialise terraform with our provider:

```console
terraform init
```

Finally, we can invoke `terraform apply` to apply our terraform configuration:

```console
terraform apply
```

## Documentation

Documentation is located within this repository at `/docs`, and is published in the [Terraform Registry](https://registry.terraform.io/providers/ans-group/safedns/latest/docs)

## Upgrading

:warning: This provider was originally created under the `ukfast` organisation, and was later moved to the `ans-group` organisation. Upgrading the `ukfast` provider will result in the error `checksum list has unexpected`. Updating provider config to use the `ans-group` organisation will resolve this
