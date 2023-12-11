# tfnew

tfnew is a small cli tool for generating new terraform modules in a fairly opinionated way

## installation

### bash

```shell
curl -fsSL https://raw.githubusercontent.com/fitz7/tfnew/main/install.sh | bash
```

### go

```shell
go install gitub.com/fitz7/tfnew@latest
```

## usage

### init

This step is optional in using `tfnew` as all the args are also supported by `tfnew module`.

`tfnew init` will create a `.tfnew.yaml` file in the root of your project directory (wherever the .git folder is)

#### default behaviour

```shell
tfnew init
```

By default, `tfnew init` specifies a local backend and a local `terraform.tfstate` file.

#### specifying a backend

currently the only other supported backend is `gcs`

```shell
tfnew init --backend=gcs --backend_gcs_bucket=my-state-bucket --backend_gcs_prefix=my-repo
```

This will configure subsequent runs in this project to create a terraform `backend` block configured to use gcs with the bucket and prefix

> [!IMPORTANT]
> the prefix defined is not the final prefix. when you create your first module the prefix will be of the form `prefix/module_path`

### module

#### basic usage

```shell
tfnew module modules/new-module
```

Will generate a new module with a very basic terraform block in the `versions.tf` file

```hcl
terraform {
  required_version = ">= 1.0"
}
```

#### create a new root module

A root module will be generated with a backend block to store the root modules state it will also set the required terraform version to the latest minor version

```shell
tfnew module root-module --root
```

```hcl
terraform {
  required_version = "~> 1.6"
  backend "gcs" {
    bucket = "my-state-bucket"
    prefix = "root-module"
  }
}
```

#### creating modules with required_providers

required_providers must be referenced by their source and are also generated with their latest minor version

```shell
tfnew module root-module-providers --root --required_providers=hashicorp/google,hashicorp/google-beta
```

```hcl
terraform {
  required_version = "~> 1.6"
  backend "gcs" {
    bucket = "my-state-bucket"
    prefix = "root-module-providers"
  }
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 5.8"
    }
    google-beta = {
      source  = "hashicorp/google-beta"
      version = "~> 5.8"
    }
  }
}
```
