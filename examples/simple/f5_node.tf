resource "bigip_node" "terraform-test-node" {
    name = "terraform-test-node" // name of the node, as submitted to your F5
    partition = "Website" // which partition does this resource belong within?
    address = "10.3.130.130" // put a valid IP for your node in here
    description = "terraform-test-node!" // description stirng is optional, but often helpful later, so pick something useful
    // ratio = 1 // (not required; default is 1; valid is 1 through 65535)
}
