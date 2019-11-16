resource "aws_instance" "source" {
  ami           = data.aws_ami.ubuntu.id
  instance_type = "t2.micro"
  subnet_id = aws_subnet.subnet_1_of_1.id

  vpc_security_group_ids = [
    aws_security_group.no_rules.id
  ]

  tags = {
    Name = "aat_source"
  }
}

resource "aws_instance" "destination" {
  ami           = data.aws_ami.ubuntu.id
  instance_type = "t2.micro"
  subnet_id = aws_subnet.subnet_1_of_1.id

  vpc_security_group_ids = [
    aws_security_group.no_rules.id
  ]

  tags = {
    Name = "aat_destination"
  }
}
