resource "aws_network_acl" "both_subnets_all_traffic" {
  vpc_id = "${aws_vpc.aat_vpc.id}"

  subnet_ids = [
    aws_subnet.subnet_1_of_2.id,
    aws_subnet.subnet_2_of_2.id,
  ]

  egress {
    protocol   = "-1"
    rule_no    = 100
    action     = "allow"
    cidr_block = "0.0.0.0/0"
    from_port  = 0
    to_port    = 0
  }

  ingress {
    protocol   = "-1"
    rule_no    = 100
    action     = "allow"
    cidr_block = "0.0.0.0/0"
    from_port  = 0
    to_port    = 0
  }

  tags = {
    Name = "aat_both_subnets_all_traffic"
  }
}
