resource "aws_security_group" "inbound_allow_esp" {
  name        = "aat_inbound_allow_esp"
  description = "Allow all inbound ESP traffic"
  vpc_id      = aws_vpc.aat_vpc.id

  ingress {
    from_port       = 0
    to_port         = 0
    protocol        = "50"
    cidr_blocks     = ["0.0.0.0/0"]
  }
}
