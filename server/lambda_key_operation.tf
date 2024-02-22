resource "aws_lambda_function" "key_operation" {
  function_name    = "key_operation"
  filename         = "./lambda/key_operation/archive/bootstrap.zip"
  role             = aws_iam_role.key_operation_role.arn
  handler          = "go"
  runtime          = "provided.al2"
  memory_size      = 512
  timeout          = 30
  source_code_hash = data.archive_file.key_operation_archive.output_base64sha256
  environment {
    variables = {
      "APIKEY" = jsondecode(file(".local/api.json")).apiKey
      "URL"    = "https://app.candyhouse.co/api/sesame2/%s/cmd"
      "TZ"     = "Asia/Tokyo"
      "ARN"    = jsondecode(file(".local/arn.json")).arn_create_schedule
    }
  }

}

resource "null_resource" "key_operation" {
  triggers = {
    # always_run = timestamp()
    "codechange" = sha1(file("./lambda/key_operation/cmd/main.go"))
  }
  provisioner "local-exec" {
    command = "cd ./lambda/key_operation/cmd/ && GOOS=linux GOARCH=amd64 go build -o ../build/bootstrap main.go"
  }
}

data "archive_file" "key_operation_archive" {
  type        = "zip"
  source_file = "./lambda/key_operation/build/bootstrap"
  output_path = "./lambda/key_operation/archive/bootstrap.zip"

  depends_on = [null_resource.key_operation]
}


resource "aws_iam_role" "key_operation_role" {
  name               = "key_operation-role"
  assume_role_policy = file("policies/lambda-assume-role.json")
}

resource "aws_iam_policy" "key_operation_policy" {
  name   = "key_operation-policy"
  policy = file("policies/lambda-key-operation-policy.json")
}

resource "aws_iam_role_policy_attachment" "key_operation_policy_attachment" {
  role       = aws_iam_role.key_operation_role.name
  policy_arn = aws_iam_policy.key_operation_policy.arn
}

# cloudwatch
resource "aws_cloudwatch_log_group" "key_operation_log_group" {
  name              = "/aws/lambda/key_operation"
  retention_in_days = 1
}


resource "aws_lambda_permission" "key_operation_event_bridge_permission" {
  statement_id  = "AllowExecutionFromEventBridge"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.key_operation.function_name
  principal     = "events.amazonaws.com"
  source_arn    = jsondecode(file(".local/arn.json")).key_event
}
