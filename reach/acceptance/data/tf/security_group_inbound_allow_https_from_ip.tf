resource "aws_security_group" "inbound_allow_https_from_ip" {
  name        = "aat_inbound_allow_https_from_ip"
  description = "Allow inbound HTTPS traffic from IP CIDR"
  vpc_id      = aws_vpc.aat_vpc.id

  ingress {
    from_port       = 443
    to_port         = 443
    protocol        = "tcp"
    cidr_blocks     = [aws_subnet.subnet_1_of_1.cidr_block]
  }
}
