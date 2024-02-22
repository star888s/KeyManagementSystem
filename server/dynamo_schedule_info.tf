resource "aws_dynamodb_table" "ScheduleInfo" {

  name           = "ScheduleInfo"
  billing_mode   = "PAY_PER_REQUEST"
  hash_key       = "id"
  range_key      = "startTime"
  stream_enabled = true
  attribute {
    name = "id"
    type = "S"
  }
  attribute {
    name = "startTime"
    type = "S"
  }
  # attribute {
  #   name = "endTime"
  #   type = "S"
  # }
  # attribute {
  #   name = "name"
  #   type = "S"
  # }
  # attribute {
  #   name = "memo"
  #   type = "S"
  # }

  point_in_time_recovery {
    enabled = true
  }

  lifecycle {
    prevent_destroy = true
  }
}
