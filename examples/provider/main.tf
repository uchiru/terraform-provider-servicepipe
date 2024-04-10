terraform {
  required_providers {
    servicepipe = {
      source = "hashicorp.com/edu/servicepipe"
    }
  }
}

provider "servicepipe" {
  endpoint = "https://api.servicepipe.ru/api/v1"
  token    = var.SERVICEPIPE_TOKEN
}

resource "servicepipe_l7resource" "test" {
  l7_resource_name = "test.example.com"
  www_redir        = 1
  http_2_https     = 1
  force_ssl        = 0

  use_custom_ssl = 1
  custom_ssl_key = base64encode(var.test_example_com_key)
  custom_ssl_crt = base64encode(var.test_example_com_crt)

  origins = [
    {
      ip     = "190.90.160.33"
      weight = 50
      mode   = "primary"
    },
    # {
    #   ip     = "190.90.160.34"
    #   weight = 50
    #   mode   = "primary"
    # },
  ]
}

# resource "servicepipe_l7origin" "test" {
#   l7_resource_id = servicepipe_l7resource.test.l7_resource_id
#   ip             = "190.90.160.33"
#   weight         = 50
#   mode           = "primary"
#   # mode = "backup"
# }
#
# resource "servicepipe_l7origin" "test2" {
#   l7_resource_id = servicepipe_l7resource.test.l7_resource_id
#   ip             = "190.90.160.34"
#   weight         = 50
#   mode           = "primary"
#   # mode = "backup"
# }

# resource "servicepipe_l7origin" "test2" {
#   l7_resource_id = servicepipe_l7resource.test.l7_resource_id
#   ip             = "190.90.160.30"
#   weight = 50
#   mode = "primary"
#   # mode = "backup"
# }
