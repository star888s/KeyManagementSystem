resource "aws_lambda_function" "key_operation" {
  function_name    = "key_operation"
  filename         = "./lambda/key_operation/archive/bootstrap.zip"
  role             = aws_iam_role.key_operation_role.arn
  handler          = "go"
  runtime          = "provided.al2"
  memory_size      = 128
  source_code_hash = data.archive_file.lambda.output_base64sha256
  environment {
    variables = {
        "APIKEY" = file(".local/api.txt")
            }
        }

}

resource "null_resource" "key_operation" {
  triggers = {
    always_run = timestamp()
  }
  provisioner "local-exec" {
    command = "cd ./lambda/key_operation/cmd/ && GOOS=linux GOARCH=amd64 go build -o ../build/bootstrap main.go"
  }
}

data "archive_file" "lambda" {
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
  policy = file("policies/lambda-policy.json")
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
