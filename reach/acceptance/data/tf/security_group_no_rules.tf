resource "aws_security_group" "no_rules" {
  name        = "aat_no_rules"
  description = "No rules associated with this group"
  vpc_id      = aws_vpc.aat_vpc.id
}
