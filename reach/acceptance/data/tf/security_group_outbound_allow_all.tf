resource "aws_security_group" "outbound_allow_all" {
  name        = "aat_outbound_allow_all"
  description = "Allow all outbound traffic"
  vpc_id      = aws_vpc.aat_vpc.id

  egress {
    from_port       = 0
    to_port         = 0
    protocol        = "-1"
    cidr_blocks     = ["0.0.0.0/0"]
  }
}
