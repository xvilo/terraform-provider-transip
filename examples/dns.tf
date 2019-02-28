variable "private_key" {}

provider "transip" {
  account_name = "aequitas"
  private_key  = "${var.private_key}"
}

# data "transip_domain" "example_com" {
#   name = "ijohan.nl"
# }

# resource "transip_dns_entry" "test" {
#   domain  = "${data.transip_domain.example_com.id}"
#   name    = "test"
#   type    = "CNAME"
#   content = "@"
# }

# output "domain_id" {
#   value = "${data.transip_domain.example_com.id}"
# }

# output "locked" {
#   value = "${data.transip_domain.example_com.is_locked}"
# }

# output "domain_name" {
#   value = "${data.transip_domain.example_com.name}"
# }

# output "ns" {
#   value = "${lookup(data.transip_domain.example_com.nameservers[0], "hostname")}"
# }

resource "transip_domain" "test" {
  name = "locohost.nl"
}

resource "transip_dns_entry" "www" {
  domain  = "${transip_domain.test.id}"
  name    = "www"
  type    = "CNAME"
  content = "@"
}

resource "transip_dns_entry" "test1" {
  domain  = "${transip_domain.test.id}"
  name    = "test"
  type    = "A"
  content = "1.2.3.4"
}

# resource "transip_dns_entry" "test2" {
#   domain  = "${transip_domain.test.id}"
#   name    = "test"
#   expire  = "600"
#   type    = "A"
#   content = "1.2.3.4"
# }
