resource "bigip_vserver" "terraform-test-virtualserver" {
    name = "tf-test-virtualserver"
    description = "some custom description goes here"
    partition = "Website"
    dest_ip = "2.2.2.2"
    dest_port = 80
    // pool = "/${bigip_pool.terraform-test-pool.partition}/${bigip_pool.terraform-test-pool.name}"

    // enabled = true (optional)
    // protocol = "tcp" (optional, can be tcp udp or any)
    // snat_automap = true (optional)
}
