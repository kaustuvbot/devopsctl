terraform {
  required_providers {
    aws = {
      source = "hashicorp/aws"
      # Missing version constraint
    }
  }
}

resource "aws_instance" "example" {
  ami           = "ami-12345"
  instance_type = "t2.micro"
}
