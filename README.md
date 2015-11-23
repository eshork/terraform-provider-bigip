# terraform-provider-bigip
We wanted to start using terraform, but sadly there was no F5 BIGIP provider, so we wrote one.

## Caveats
 * This was originally developed for some very specific needs, and F5 does a LOT of things that we didn't implemnt, so don't be terribly upset if it doesn't do what you want out-of-the-box
 * The author of this module is self-absorbed and often on "too much scotch" (is that a real thing?), so expect plenty of "fun" comments along the way. We tried to cull out all the f-bombs, but please bare with us!
 * We only have 1 level of partition depth - so if that matters at all (and we honestly never checked), this could totally break in weird/awesome ways
 * Again, there are plenty of missing features, but we feel this is a decent start for an F5 provider. If you fix bugs or add new features, pull requests are always appreciated

# Installation
You can either build from source or install the binaries. Building from source is left as an exercise for the reader, but pre-built binaries for a platform or two should be available for download as well. 

--

# Usage

### Provisioner Configuration

You currently need to use a provider config within your *.tf files. _Adding support to also accept environment variables would be a nice-to-have._

    provider "bigip" {
        username = "your_RestAPI_username"
        password = "your_RestAPI_password"
        management_ip = "your_F5_IP_or_hostname"
    }

### Create a Node

    resource "bigip_node" "myWebHostNode" {
        name = "myWebHostNode"
        partition = "MyHostingPartition"
        address = "1.1.1.1"
        description = "my host desecription" // optional
        // enabled = true (optional)
        // connection_limit = 0 (optional, defaults to 0 for unlimited)
        // connection_rate_limit = 0 (optional, defaults to 0 for unlimited)
    }



### Create a Resource Pool

    resource "bigip_pool" "myWebHostPool" {
    	name = "myWebHostPool"
    	partition = "Website"
    	description = "pool description" // (optional)
    }




### Create a Resource Pool Member
This resource element could use some work. Due to time contraints on the original coding, we had to make some dirty design decisions to get the provider up and working as quickly as possible. For now, it's typically best to just use terraform variables to ensure object dependencies and information are disseminated appropriately.

    resource "bigip_pool_member" "myWebHostPool-member" {
    	node_id = "${bigip_node.myWebHostNode.id}"
    	pool_id = "${bigip_pool.myWebHostPool.id}"
    	node_name = "${element(bigip_node.myWebHostNode.*.name, count.index)}"
    	address = "${element(bigip_node.myWebHostNode.*.address, count.index)}"
    	partition = "${bigip_pool.myWebHostPool.partition}"
    	port = 80
    	description = "myWebHostPool-member description" // (optional)
    }



### Create a Virtual Server

    resource "bigip_vserver" "myVirtualServer_http" {
    	name = "myVirtualServer_http"
    	description = "myVirtualServer_http description"
    	partition = "Website"
    	dest_ip = "2.2.2.2"
    	dest_port = 80
    	pool = "/${bigip_pool.myWebHostPool-member.partition}/${bigip_pool.myWebHostPool-member.name}"
    	// enabled = false (optional)
    	// protocol = "tcp" (optional, can be tcp udp or any)
    	// snat_automap = true (optional)
    }


