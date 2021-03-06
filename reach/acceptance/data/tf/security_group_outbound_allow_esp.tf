resource "aws_security_group" "outbound_allow_esp" {
  name        = "aat_outbound_allow_esp"
  description = "Allow all outbound ESP traffic"
  vpc_id      = aws_vpc.aat_vpc.id

  egress {
    from_port       = 0
    to_port         = 0
    protocol        = "50"
    cidr_blocks     = ["0.0.0.0/0"]
  }
}
