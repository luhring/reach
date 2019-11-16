resource "aws_security_group" "outbound_allow_all_tcp" {
  name        = "aat_outbound_allow_all_tcp"
  description = "Allow all outbound TCP traffic"
  vpc_id      = aws_vpc.aat_vpc.id

  egress {
    from_port       = 0
    to_port         = 65535
    protocol        = "tcp"
    cidr_blocks     = ["0.0.0.0/0"]
  }
}
