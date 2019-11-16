resource "aws_security_group" "inbound_allow_udp_dns_from_sg_no_rules" {
  name        = "aat_inbound_allow_udp_dns_from_sg_no_rules"
  description = "Allow DNS (UDP) from SG no rules"
  vpc_id      = aws_vpc.aat_vpc.id

  ingress {
    from_port       = 53
    to_port         = 53
    protocol        = "udp"
    security_groups = [aws_security_group.no_rules.id]
  }
}
