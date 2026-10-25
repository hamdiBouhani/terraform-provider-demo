terraform {
  required_providers {
    demo = {
      source = "hamdibouhani/demo"
    }
  }
}

provider "demo" {
  base_url = "https://jsonplaceholder.typicode.com"
}

resource "demo_post" "hello" {
  title = "Hello from Terraform"
  body  = "Created by a Go provider calling the API"
}

output "post_id" {
  value = demo_post.hello.id
}