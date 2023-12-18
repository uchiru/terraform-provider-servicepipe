terraform {
  required_providers {
    servicepipe = {
      source = "hashicorp.com/edu/servicepipe"
    }
    hashicups = {
      source = "hashicorp.com/edu/hashicups"
    }
  }
}

provider "servicepipe" {
  endpoint = "https://api.servicepipe.ru/api/v1"
  token = var.SERVICEPIPE_TOKEN
}

resource "servicepipe_l7resource" "test" {
  l7_resource_name = "testdomain.xyz"
  origin_data = "190.90.160.30"
  www_redir = 1
  http_2_https = 1
  force_ssl = 0
}
