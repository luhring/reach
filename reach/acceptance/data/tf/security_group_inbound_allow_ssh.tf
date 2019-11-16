resource "aws_security_group" "inbound_allow_ssh_from_all_ip_addresses" {
  name        = "aat_inbound_allow_ssh_from_all_ip_addresses"
  description = "Allow all SSH traffic"
  vpc_id      = aws_vpc.aat_vpc.id

  ingress {
    from_port       = 22
    to_port         = 22
    protocol        = "tcp"
    cidr_blocks     = ["0.0.0.0/0"]
  }
}
