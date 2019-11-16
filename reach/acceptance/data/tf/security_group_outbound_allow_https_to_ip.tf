resource "aws_security_group" "outbound_allow_https_to_ip" {
  name        = "aat_outbound_allow_https_to_ip"
  description = "Allow outbound HTTPS traffic to IP CIDR"
  vpc_id      = aws_vpc.aat_vpc.id

  egress {
    from_port       = 443
    to_port         = 443
    protocol        = "tcp"
    cidr_blocks     = [aws_subnet.subnet_1_of_1.cidr_block]
  }
}
