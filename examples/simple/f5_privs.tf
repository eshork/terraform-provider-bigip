
// Currently this file must be present and correctly populated.
// Many other terraform providers support extracting this info from environment variables, but this provider isn't quite there just yet
provider "bigip" {
    username = "your_rest_enabled_username" //  <- change this to your f5 username
    password = "some_secret_password" //  <- change this to your f5 password
    management_ip = "f5_ip_address" //  <- you can alternatively use a dns hostname here
}
