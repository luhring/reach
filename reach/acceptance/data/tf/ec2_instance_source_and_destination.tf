resource "aws_instance" "source" {
  ami           = "${data.aws_ami.ubuntu.id}"
  instance_type = "t2.micro"

  tags = {
    Name = "source"
  }
}

resource "aws_instance" "destination" {
  ami           = "${data.aws_ami.ubuntu.id}"
  instance_type = "t2.micro"

  tags = {
    Name = "destination"
  }
}

output "source_id" {
  value = aws_instance.source.id
}

output "destination_id" {
  value = aws_instance.destination.id
}
