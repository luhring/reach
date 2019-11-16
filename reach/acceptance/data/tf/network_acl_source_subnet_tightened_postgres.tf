resource "aws_network_acl" "source_subnet_tightened_postgres" {
  vpc_id = "${aws_vpc.aat_vpc.id}"

  subnet_ids = [
    aws_subnet.subnet_1_of_2.id,
  ]

  egress {
    protocol   = "tcp"
    rule_no    = 100
    action     = "allow"
    cidr_block = aws_subnet.subnet_2_of_2.cidr_block
    from_port  = 5432
    to_port    = 5432
  }

  ingress {
    protocol   = "tcp"
    rule_no    = 100
    action     = "allow"
    cidr_block = aws_subnet.subnet_2_of_2.cidr_block
    from_port  = 0
    to_port    = 65535
  }

  tags = {
    Name = "aat_source_subnet_tightened_postgres"
  }
}
