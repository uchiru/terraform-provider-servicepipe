variable "SERVICEPIPE_TOKEN" {}

variable "test_example_com_key" {
  type    = string
  default = ""

  validation {
    condition     = var.test_example_com_key != ""
    error_message = "The test_example_com_key value must be not empty, ex: test_example_com_key"
  }
}

variable "test_example_com_crt" {
  type    = string
  default = ""

  validation {
    condition     = var.test_example_com_crt != ""
    error_message = "The test_example_com_crt value must be not empty, ex: test_example_com_crt"
  }
}
