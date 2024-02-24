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
  attribute {
    name = "scheduled"
    type = "S"
  }

  global_secondary_index {
    hash_key           = "id"
    name               = "id-scheduled-index"
    non_key_attributes = []
    projection_type    = "ALL"
    range_key          = "scheduled"
    read_capacity      = 0
    write_capacity     = 0
  }

  point_in_time_recovery {
    enabled = true
  }

  lifecycle {
    prevent_destroy = true
  }
}
