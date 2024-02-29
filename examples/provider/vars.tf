variable "SERVICEPIPE_TOKEN" {}

variable "test_runit_cc_key" {
  type    = string
  default = ""

  validation {
    condition     = var.test_runit_cc_key != ""
    error_message = "The test_runit_cc_key value must be not empty, ex: test_runit_cc_key"
  }
}

variable "test_runit_cc_crt" {
  type    = string
  default = ""

  validation {
    condition     = var.test_runit_cc_crt != ""
    error_message = "The test_runit_cc_crt value must be not empty, ex: test_runit_cc_crt"
  }
}
