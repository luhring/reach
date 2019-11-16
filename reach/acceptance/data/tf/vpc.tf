resource "aws_vpc" "aat_vpc" {
  cidr_block       = "10.0.0.0/16"

  tags = {
    Name = "aat_vpc"
  }
}
