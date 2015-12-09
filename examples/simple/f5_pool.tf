
resource "bigip_pool" "terraform-test-pool" {
    name = "terraform-test-pool"
    partition = "Website"
    description = "terraform-test-pool!"
}
