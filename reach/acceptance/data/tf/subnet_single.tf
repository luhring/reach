resource "aws_subnet" "subnet_1_of_1" {
  vpc_id     = aws_vpc.aat_vpc.id
  cidr_block = "10.0.1.0/24"
  map_public_ip_on_launch = false

  tags = {
    Name = "aat_subnet_1_of_1"
  }
}
