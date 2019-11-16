resource "aws_instance" "source" {
  ami           = data.aws_ami.ubuntu.id
  instance_type = "t2.micro"
  subnet_id = aws_subnet.subnet_1_of_1.id


  security_groups = [
    aws_security_group.outbound_allow_https_to_ip.id
  ]

  tags = {
    Name = "aat_source"
  }
}

resource "aws_instance" "destination" {
  ami           = data.aws_ami.ubuntu.id
  instance_type = "t2.micro"
  subnet_id = aws_subnet.subnet_1_of_1.id

  security_groups = [
    aws_security_group.inbound_allow_https_from_ip.id
  ]

  tags = {
    Name = "aat_destination"
  }
}
