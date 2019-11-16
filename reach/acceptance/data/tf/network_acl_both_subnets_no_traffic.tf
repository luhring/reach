resource "aws_network_acl" "both_subnets_no_traffic" {
  vpc_id = "${aws_vpc.aat_vpc.id}"

  subnet_ids = [
    aws_subnet.subnet_1_of_2.id,
    aws_subnet.subnet_2_of_2.id,
  ]

  tags = {
    Name = "aat_both_subnets_no_traffic"
  }
}
