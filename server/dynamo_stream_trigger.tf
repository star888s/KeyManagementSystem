resource "aws_lambda_event_source_mapping" "schedule_watcher" {
  event_source_arn       = aws_dynamodb_table.ScheduleInfo.stream_arn
  function_name          = aws_lambda_function.create_schedule.arn
  starting_position      = "LATEST"
  batch_size             = 100
  maximum_retry_attempts = 0

}
