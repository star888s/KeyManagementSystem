resource "aws_lambda_function" "delete_schedule_api" {
  function_name    = "delete_schedule_api"
  filename         = "./lambda/delete_schedule_api/archive/bootstrap.zip"
  role             = aws_iam_role.delete_schedule_api_role.arn
  handler          = "go"
  runtime          = "provided.al2"
  memory_size      = 512
  timeout          = 30
  source_code_hash = data.archive_file.delete_schedule_api_archive.output_base64sha256
  environment {
    variables = {
      "ARN" = jsondecode(file(".local/arn.json")).arn
      "TZ"  = "Asia/Tokyo"
    }
  }

}

resource "null_resource" "delete_schedule_api" {
  triggers = {
    # always_run = timestamp()
    "codechange" = sha1(file("./lambda/delete_schedule_api/cmd/main.go"))
  }
  provisioner "local-exec" {
    command = "cd ./lambda/delete_schedule_api/cmd/ && GOOS=linux GOARCH=amd64 go build -o ../build/bootstrap main.go"
  }
}

data "archive_file" "delete_schedule_api_archive" {
  type        = "zip"
  source_file = "./lambda/delete_schedule_api/build/bootstrap"
  output_path = "./lambda/delete_schedule_api/archive/bootstrap.zip"

  depends_on = [null_resource.delete_schedule_api]
}


resource "aws_iam_role" "delete_schedule_api_role" {
  name               = "delete_schedule_api-role"
  assume_role_policy = file("policies/lambda-assume-role.json")
}

resource "aws_iam_policy" "delete_schedule_api_policy" {
  name   = "delete_schedule_api-policy"
  policy = file("policies/lambda_delete_schedule_api_policy.json")
}

resource "aws_iam_role_policy_attachment" "delete_schedule_api_policy_attachment" {
  role       = aws_iam_role.delete_schedule_api_role.name
  policy_arn = aws_iam_policy.delete_schedule_api_policy.arn
}

# cloudwatch
resource "aws_cloudwatch_log_group" "delete_schedule_api_log_group" {
  name              = "/aws/lambda/delete_schedule_api"
  retention_in_days = 1
}


resource "aws_lambda_permission" "delete_schedule_api_permisson" {
  statement_id  = "AllowExecutionFromAPIGateway"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.delete_schedule_api.arn
  principal     = "apigateway.amazonaws.com"

  source_arn = "${aws_api_gateway_rest_api.kms_api.execution_arn}/*/POST/delete_schedule"
}
