resource "aws_api_gateway_rest_api" "kms_api" {
  name        = "kms"
  description = "for KMS web"
}

resource "aws_api_gateway_resource" "resource_upsert_schedule" {
  rest_api_id = aws_api_gateway_rest_api.kms_api.id
  parent_id   = aws_api_gateway_rest_api.kms_api.root_resource_id
  path_part   = "upsert_schedule"
}

resource "aws_api_gateway_method" "method_upsert_schedule" {
  rest_api_id      = aws_api_gateway_rest_api.kms_api.id
  resource_id      = aws_api_gateway_resource.resource_upsert_schedule.id
  http_method      = "POST"
  authorization    = "NONE"
  api_key_required = true
}

resource "aws_api_gateway_integration" "int_upsert_schedule" {
  rest_api_id             = aws_api_gateway_rest_api.kms_api.id
  resource_id             = aws_api_gateway_resource.resource_upsert_schedule.id
  http_method             = aws_api_gateway_method.method_upsert_schedule.http_method
  integration_http_method = "POST"
  content_handling        = "CONVERT_TO_TEXT"
  type                    = "AWS_PROXY"
  uri                     = aws_lambda_function.upsert_schedule_api.invoke_arn
}

resource "aws_api_gateway_deployment" "deploy_kms_api" {
  depends_on  = [aws_api_gateway_integration.int_upsert_schedule]
  rest_api_id = aws_api_gateway_rest_api.kms_api.id
  description = "kms_api"
  stage_name  = "prod"
}

resource "aws_api_gateway_stage" "stage_kms_api" {
  stage_name    = aws_api_gateway_deployment.deploy_kms_api.stage_name
  rest_api_id   = aws_api_gateway_rest_api.kms_api.id
  deployment_id = aws_api_gateway_deployment.deploy_kms_api.id
}


resource "aws_api_gateway_api_key" "kms_web" {
  name = "kms_web"
}

resource "aws_api_gateway_usage_plan" "plan_kms_web" {
  name        = "kms_web"
  description = "kms usage plan"

  throttle_settings {
    burst_limit = 500
    rate_limit  = 100
  }

  api_stages {
    api_id = aws_api_gateway_rest_api.kms_api.id
    stage  = aws_api_gateway_stage.stage_kms_api.stage_name
  }
}

resource "aws_api_gateway_usage_plan_key" "plan_key_kms_web" {
  key_id        = aws_api_gateway_api_key.kms_web.id
  key_type      = "API_KEY"
  usage_plan_id = aws_api_gateway_usage_plan.plan_kms_web.id
}
