resource "aws_dynamodb_table" "KeyInfo" {

  name         = "KeyInfo"
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "id"
  attribute {
    name = "id"
    type = "S"
  }
  # attribute {
  #   name = "name"
  #   type = "S"
  # }
  # attribute {
  #   name = "uuid"
  #   type = "S"
  # }
  # attribute {
  #   name = "secretKey"
  #   type = "S"
  # }
  # attribute {
  #   name = "apiKey"
  #   type = "S"
  # }

  point_in_time_recovery {
    enabled = true
  }

  lifecycle {
    prevent_destroy = true
  }
}
