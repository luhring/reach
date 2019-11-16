resource "aws_security_group" "inbound_allow_postgres_from_sg_no_rules" {
  name        = "aat_inbound_allow_postgres_from_sg_no_rules"
  description = "Allow all inbound Postgres traffic from SG no rules"
  vpc_id      = aws_vpc.aat_vpc.id

  ingress {
    from_port       = 5432
    to_port         = 5432
    protocol        = "tcp"
    security_groups = [aws_security_group.no_rules.id]
  }
}
