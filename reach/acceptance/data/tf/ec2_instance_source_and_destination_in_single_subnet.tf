resource "aws_instance" "source" {
  ami           = "${data.aws_ami.ubuntu.id}"
  instance_type = "t2.micro"
  subnet_id = "${aws_subnet.subnet_single.id}"

  security_groups = ["${aws_security_group.outbound_allow_all.id}"]

  tags = {
    Name = "aat_source"
  }
}

resource "aws_instance" "destination" {
  ami           = "${data.aws_ami.ubuntu.id}"
  instance_type = "t2.micro"
  subnet_id = "${aws_subnet.subnet_single.id}"

  security_groups = ["${aws_security_group.inbound_allow_all.id}"]

  tags = {
    Name = "aat_destination"
  }
}

output "source_id" {
  value = aws_instance.source.id
}

output "destination_id" {
  value = aws_instance.destination.id
}
