# terraform-provider-demo

A minimal, educational [Terraform](https://www.terraform.io/) provider written in Go.
It shows how to wire a real REST API into the Terraform lifecycle
(`Create` / `Read` / `Update` / `Delete`) using a small resource called `demo_post`.

The example API used by default is the free
[JSONPlaceholder](https://jsonplaceholder.typicode.com/) — so you can run
everything without any credentials.

> ⚠️ **This is a learning project, not production-ready.**

---

## What it does

The provider exposes one resource:

| Resource    | Action on API                  |
|-------------|--------------------------------|
| `demo_post` | `POST /posts`, `GET /posts/{id}`, `DELETE /posts/{id}` |

Fields:

| Attribute | Type   | Required | Description            |
|-----------|--------|----------|------------------------|
| `title`   | string | yes      | Post title             |
| `body`    | string | no       | Post body              |

The `id` attribute is populated automatically from the API response and stored
in Terraform state.

---

## Requirements

- [Go](https://go.dev/dl/) 1.21+
- [Terraform](https://developer.hashicorp.com/terraform/install) 1.5+

---

## Project layout

```
.
├── main.go                   # plugin entry point
├── provider/
│   └── provider.go           # provider schema + demo_post CRUD
├── internal/client/
│   └── client.go             # plain HTTP client for the API
├── examples/
│   └── main.tf               # usage example
├── Makefile
└── README.md
```

---

## Build

```bash
make build
```

This produces a `terraform-provider-demo` binary in the repo root.

---

## Install (dev override)

Terraform normally fetches providers from a registry. For local development we
tell Terraform to use our locally built binary via a *dev override*.

1. Create a `.terraformrc` in the repo root:

   ```hcl
   provider_installation {
     dev_overrides {
       "hamdiBouhani/demo" = "/absolute/path/to/this/repo"
     }
     direct {}
   }
   ```

   `make install` prints this snippet with the correct absolute path filled in.

2. Tell Terraform to use it (either export for the session, or put it in
   `~/.terraformrc`):

   ```bash
   export TF_CLI_CONFIG_FILE=/absolute/path/to/repo/.terraformrc
   ```

---

## Usage

`examples/main.tf`:

```hcl
terraform {
  required_providers {
    demo = {
      source = "hamdiBouhani/demo"
    }
  }
}

provider "demo" {
  base_url = "https://jsonplaceholder.typicode.com"
  # api_key is optional; can also be set via DEMO_API_KEY env var
}

resource "demo_post" "hello" {
  title = "Hello from Terraform"
  body  = "Created by a Go provider calling the API"
}

output "post_id" {
  value = demo_post.hello.id
}
```

Run:

```bash
export TF_CLI_CONFIG_FILE=$(pwd)/.terraformrc
cd examples
terraform plan
terraform apply
terraform destroy
```

Or via the Makefile from the repo root:

```bash
make plan
make apply
make destroy
```

Expected output on apply:

```
demo_post.hello: Creating...
demo_post.hello: Creation complete after 0s [id=101]

Apply complete! Resources: 1 added, 0 changed, 0 destroyed.

Outputs:

post_id = "101"
```

---

## Configuration

The provider accepts the following arguments. Each can also be provided through
an environment variable.

| Argument   | Env var         | Required | Default                               |
|------------|-----------------|----------|---------------------------------------|
| `base_url` | `DEMO_BASE_URL` | no       | `https://jsonplaceholder.typicode.com`|
| `api_key`  | `DEMO_API_KEY`  | no       | *(empty)*                             |

`api_key` is marked `sensitive` and sent as `Authorization: Bearer <key>` when
non-empty.

---

## Development

```bash
make help     # list all targets
make fmt      # gofmt
make vet      # go vet
make build    # compile
make check    # fmt + vet + build
make clean    # remove binary and terraform artifacts
```

To use your own API, edit `internal/client/client.go`:

- change the request/response structs to match your API,
- adjust paths (`/posts` etc.),
- implement a real `Update` (currently `demo_post` does delete-and-recreate
  when `title` or `body` changes).

---

## How it works (short version)

1. `main.go` starts the plugin server via `plugin.Serve`.
2. `provider/provider.go` defines the schema for the provider and for
   `demo_post`. Each CRUD hook receives a `*client.Client` via the `meta`
   argument.
3. `internal/client/client.go` contains all HTTP logic, kept separate so it
   can be unit-tested without Terraform.
4. On `terraform apply`, Terraform calls `CreateContext`, which performs the
   HTTP request, writes the returned `id` into `d.SetId(...)`, and re-reads the
   resource to keep state consistent.

---

## License

MIT — see [LICENSE](./LICENSE).