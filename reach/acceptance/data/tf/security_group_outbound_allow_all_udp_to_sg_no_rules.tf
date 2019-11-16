resource "aws_security_group" "outbound_allow_all_udp_to_sg_no_rules" {
  name        = "aat_outbound_allow_all_udp_to_sg_no_rules"
  description = "Allow all outbound UDP traffic to SG no rules"
  vpc_id      = aws_vpc.aat_vpc.id

  egress {
    from_port       = 0
    to_port         = 65535
    protocol        = "udp"
    security_groups = [aws_security_group.no_rules.id]
  }
}
