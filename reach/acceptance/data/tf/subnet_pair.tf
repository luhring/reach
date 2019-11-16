resource "aws_subnet" "subnet_1_of_2" {
  vpc_id     = aws_vpc.aat_vpc.id
  cidr_block = "10.0.1.0/24"
  map_public_ip_on_launch = false

  tags = {
    Name = "aat_subnet_1_of_2"
  }
}

resource "aws_subnet" "subnet_2_of_2" {
  vpc_id     = aws_vpc.aat_vpc.id
  cidr_block = "10.0.2.0/24"
  map_public_ip_on_launch = false

  tags = {
    Name = "aat_subnet_2_of_2"
  }
}
