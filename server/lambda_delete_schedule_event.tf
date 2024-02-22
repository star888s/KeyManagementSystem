resource "aws_lambda_function" "delete_schedule" {
  function_name    = "delete_schedule"
  filename         = "./lambda/delete_schedule/archive/bootstrap.zip"
  role             = aws_iam_role.delete_schedule_role.arn
  handler          = "go"
  runtime          = "provided.al2"
  memory_size      = 512
  timeout          = 30
  source_code_hash = data.archive_file.delete_schedule_archive.output_base64sha256
  environment {
    variables = {
      "TZ" = "Asia/Tokyo"
    }
  }

}

resource "null_resource" "delete_schedule" {
  triggers = {
    # always_run = timestamp()
    "codechange" = sha1(file("./lambda/delete_schedule/cmd/main.go"))
  }
  provisioner "local-exec" {
    command = "cd ./lambda/delete_schedule/cmd/ && GOOS=linux GOARCH=amd64 go build -o ../build/bootstrap main.go"
  }
}

data "archive_file" "delete_schedule_archive" {
  type        = "zip"
  source_file = "./lambda/delete_schedule/build/bootstrap"
  output_path = "./lambda/delete_schedule/archive/bootstrap.zip"

  depends_on = [null_resource.delete_schedule]
}


resource "aws_iam_role" "delete_schedule_role" {
  name               = "delete_schedule_role"
  assume_role_policy = file("policies/lambda-assume-role.json")
}

resource "aws_iam_policy" "delete_schedule_policy" {
  name   = "delete_schedule_policy"
  policy = file("policies/lambda_delete_schedule_policy.json")
}

resource "aws_iam_role_policy_attachment" "delete_schedule_policy_attachment" {
  role       = aws_iam_role.delete_schedule_role.name
  policy_arn = aws_iam_policy.delete_schedule_policy.arn
}

# cloudwatch
resource "aws_cloudwatch_log_group" "delete_schedule_log_group" {
  name              = "/aws/lambda/delete_schedule"
  retention_in_days = 1
}


resource "aws_lambda_permission" "delete_schedule_event_bridge_permission" {
  statement_id  = "AllowExecutionFromEventBridge"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.delete_schedule.function_name
  principal     = "events.amazonaws.com"
  source_arn    = jsondecode(file(".local/arn.json")).key_event
}


resource "aws_cloudwatch_event_rule" "invoke_delete_schedule" {
  name                = "invoke_delete_schedule"
  schedule_expression = "cron(0 0 * * ? *)"
}

resource "aws_lambda_permission" "allow_cloudwatch_delete_schedule" {
  statement_id  = "AllowExecutionFromCloudWatch"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.delete_schedule.function_name
  principal     = "events.amazonaws.com"
  source_arn    = aws_cloudwatch_event_rule.invoke_delete_schedule.arn
}

resource "aws_cloudwatch_event_target" "target_delete_schedule" {
  rule      = aws_cloudwatch_event_rule.invoke_delete_schedule.name
  target_id = "delete_schedule"
  arn       = aws_lambda_function.delete_schedule.arn
}
