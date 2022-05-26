terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "4.13.0"
    }
  }
}

provider "aws" {
  profile = "default"
  region  = "us-east-2"
}

resource "aws_vpc" "main" {
  cidr_block           = "172.30.0.0/16"
  enable_dns_hostnames = true
}

# Subnet groups require subnets that span at least two availability zones.
# Externally routed traffic will go through the primary subnet.
resource "aws_subnet" "primary" {
  vpc_id            = aws_vpc.main.id
  cidr_block        = "172.30.0.0/24"
  availability_zone = "us-east-2a"
}

resource "aws_subnet" "secondary" {
  vpc_id            = aws_vpc.main.id
  cidr_block        = "172.30.1.0/24"
  availability_zone = "us-east-2b"
}

resource "aws_db_subnet_group" "postgres" {
  name       = "postgres"
  subnet_ids = [aws_subnet.primary.id, aws_subnet.secondary.id]
}

resource "aws_internet_gateway" "gateway" {
  vpc_id = aws_vpc.main.id
}

resource "aws_security_group" "postgres" {
  name        = "postgres"
  description = "Allow PostgreSQL traffic"
  vpc_id      = aws_vpc.main.id

  ingress {
    description = "TLS from VPC"
    from_port   = 5432
    to_port     = 5432
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port        = 0
    to_port          = 0
    protocol         = "-1"
    cidr_blocks      = ["0.0.0.0/0"]
    ipv6_cidr_blocks = ["::/0"]
  }

  tags = {
    Name = "allow_postgres"
  }
  depends_on = [aws_internet_gateway.gateway]
}

# Wire up the VPC to the external gateway we created above, then associate the
# gateway to the primary subnet. This is what allows us to connect to the
# database from outside the VPC. When we add support for terraform to deploy
# the application, this association should be removed so that only resources
# deployed within with the AWS VPC can reach the database.
resource "aws_route_table" "main" {
  vpc_id = aws_vpc.main.id

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.gateway.id
  }
}

resource "aws_route_table_association" "main" {
  subnet_id      = aws_subnet.primary.id
  route_table_id = aws_route_table.main.id
  depends_on     = [aws_route_table.main]
}

# Make sure we use the route table associated with the gateway. Otherwise there
# is a chance the VPC will use the route table from the secondary subnet,
# leaving the RDS instance unreachable.
resource "aws_main_route_table_association" "main" {
  vpc_id         = aws_vpc.main.id
  route_table_id = aws_route_table.main.id
}

# Generate a random password and store it in AWS Secret Manager. The user must
# fetch the secret value from AWS to see it. The application should only accept
# a secret ARN to authenticate to the database and never a plain text password.
resource "random_password" "password" {
  length           = 40
  special          = true
  override_special = "!#$%&*()-_=+[]{}<>:?"
}

# We generate a random UUID to append to the secret name so that we can
# uniquely identify the secret across multiple deployments. AWS secrets are
# soft deleted so we can't use the same name within a given timeframe without
# conflicts.
resource "random_uuid" "secret_uuid" {
}

resource "aws_secretsmanager_secret" "database_secret" {
  name        = "postgres-secret-${random_uuid.secret_uuid.result}"
  description = "Postgres database secret"
}

resource "aws_secretsmanager_secret_version" "database_secret" {
  secret_id     = aws_secretsmanager_secret.database_secret.id
  secret_string = random_password.password.result
}

resource "aws_db_instance" "postgres_database" {
  engine                 = "postgres"
  instance_class         = "db.t3.micro"
  allocated_storage      = var.storage
  username               = var.username
  password               = random_password.password.result
  vpc_security_group_ids = [aws_security_group.postgres.id]
  db_subnet_group_name   = aws_db_subnet_group.postgres.name
  publicly_accessible    = true
  skip_final_snapshot    = true
  db_name                = var.database_name
}
