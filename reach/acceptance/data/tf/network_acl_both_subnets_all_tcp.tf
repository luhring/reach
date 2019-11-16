resource "aws_network_acl" "both_subnets_all_tcp" {
  vpc_id = "${aws_vpc.aat_vpc.id}"

  subnet_ids = [
    aws_subnet.subnet_1_of_2.id,
    aws_subnet.subnet_2_of_2.id,
  ]

  egress {
    protocol   = "tcp"
    rule_no    = 100
    action     = "allow"
    cidr_block = "0.0.0.0/0"
    from_port  = 0
    to_port    = 65535
  }

  ingress {
    protocol   = "tcp"
    rule_no    = 100
    action     = "allow"
    cidr_block = "0.0.0.0/0"
    from_port  = 0
    to_port    = 65535
  }

  tags = {
    Name = "aat_both_subnets_all_tcp"
  }
}
