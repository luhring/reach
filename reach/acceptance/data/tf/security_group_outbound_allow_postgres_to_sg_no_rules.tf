resource "aws_security_group" "outbound_allow_postgres_to_sg_no_rules" {
  name        = "aat_outbound_allow_postgres_to_sg_no_rules"
  description = "Allow all outbound Postgres traffic to SG no rules"
  vpc_id      = aws_vpc.aat_vpc.id

  egress {
    from_port       = 5432
    to_port         = 5432
    protocol        = "tcp"
    security_groups = [aws_security_group.no_rules.id]
  }
}
