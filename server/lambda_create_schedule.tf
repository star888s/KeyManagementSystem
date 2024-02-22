resource "aws_lambda_function" "create_schedule" {
  function_name    = "create_schedule"
  filename         = "./lambda/create_schedule/archive/bootstrap.zip"
  role             = aws_iam_role.create_schedule_role.arn
  handler          = "go"
  runtime          = "provided.al2"
  memory_size      = 512
  timeout          = 30
  source_code_hash = data.archive_file.create_schedule_archive.output_base64sha256
  environment {
    variables = {
      "ARN" = jsondecode(file(".local/arn.json")).arn
      "TZ"  = "Asia/Tokyo"
    }
  }

}

resource "null_resource" "create_schedule" {
  triggers = {
    # always_run = timestamp()
    "codechange" = sha1(file("./lambda/create_schedule/cmd/main.go"))
  }
  provisioner "local-exec" {
    command = "cd ./lambda/create_schedule/cmd/ && GOOS=linux GOARCH=amd64 go build -o ../build/bootstrap main.go"
  }
}

data "archive_file" "create_schedule_archive" {
  type        = "zip"
  source_file = "./lambda/create_schedule/build/bootstrap"
  output_path = "./lambda/create_schedule/archive/bootstrap.zip"

  depends_on = [null_resource.create_schedule]
}


resource "aws_iam_role" "create_schedule_role" {
  name               = "create_schedule-role"
  assume_role_policy = file("policies/lambda-assume-role.json")
}

resource "aws_iam_policy" "create_schedule_policy" {
  name   = "create_schedule-policy"
  policy = file("policies/lambda-create-schedule-policy.json")
}

resource "aws_iam_role_policy_attachment" "create_schedule_policy_attachment" {
  role       = aws_iam_role.create_schedule_role.name
  policy_arn = aws_iam_policy.create_schedule_policy.arn
}

# cloudwatch
resource "aws_cloudwatch_log_group" "create_schedule_log_group" {
  name              = "/aws/lambda/create_schedule"
  retention_in_days = 1
}
