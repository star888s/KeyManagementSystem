resource "aws_api_gateway_rest_api" "kms_api" {
  name        = "kms"
  description = "for KMS web"
}

#########################################################################
# upsert_schedule
#########################################################################
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

resource "aws_api_gateway_method" "method_options_upsert_schedule" {
  rest_api_id          = aws_api_gateway_rest_api.kms_api.id
  resource_id          = aws_api_gateway_resource.resource_upsert_schedule.id
  api_key_required     = false
  authorization        = "NONE"
  authorization_scopes = []
  http_method          = "OPTIONS"
  request_models       = {}
  request_parameters   = {}
}

resource "aws_api_gateway_method_response" "res_options_upsert_schedule" {
  http_method = aws_api_gateway_method.method_options_upsert_schedule.http_method
  rest_api_id = aws_api_gateway_rest_api.kms_api.id
  resource_id = aws_api_gateway_resource.resource_upsert_schedule.id
  response_models = {
    "application/json" = "Empty"
  }
  response_parameters = {
    "method.response.header.Access-Control-Allow-Headers" = false
    "method.response.header.Access-Control-Allow-Methods" = false
    "method.response.header.Access-Control-Allow-Origin"  = false
  }
  status_code = "200"
}

resource "aws_api_gateway_integration" "int_options_upsert_schedule" {
  rest_api_id          = aws_api_gateway_rest_api.kms_api.id
  resource_id          = aws_api_gateway_resource.resource_upsert_schedule.id
  cache_key_parameters = []
  connection_type      = "INTERNET"
  http_method          = aws_api_gateway_method.method_options_upsert_schedule.http_method
  request_parameters   = {}
  request_templates = {
    "application/json" = jsonencode(
      {
        statusCode = 200
      }
    )
  }
  timeout_milliseconds = 29000
  type                 = "MOCK"
}

resource "aws_api_gateway_integration_response" "int_res_options_upsert_schedule" {
  rest_api_id = aws_api_gateway_rest_api.kms_api.id
  resource_id = aws_api_gateway_resource.resource_upsert_schedule.id
  http_method = "OPTIONS"
  response_parameters = {
    "method.response.header.Access-Control-Allow-Headers" = "'Content-Type,X-Amz-Date,Authorization,X-Api-Key,X-Amz-Security-Token'"
    "method.response.header.Access-Control-Allow-Methods" = "'OPTIONS'"
    "method.response.header.Access-Control-Allow-Origin"  = "'*'"
  }
  response_templates = {}
  status_code        = aws_api_gateway_method_response.res_options_upsert_schedule.status_code
}


#########################################################################
# get_schedule
#########################################################################
resource "aws_api_gateway_resource" "resource_get_schedule" {
  rest_api_id = aws_api_gateway_rest_api.kms_api.id
  parent_id   = aws_api_gateway_rest_api.kms_api.root_resource_id
  path_part   = "get_schedule"
}

resource "aws_api_gateway_method" "method_get_schedule" {
  rest_api_id      = aws_api_gateway_rest_api.kms_api.id
  resource_id      = aws_api_gateway_resource.resource_get_schedule.id
  http_method      = "GET"
  authorization    = "NONE"
  api_key_required = true
}

resource "aws_api_gateway_integration" "int_get_schedule" {
  rest_api_id             = aws_api_gateway_rest_api.kms_api.id
  resource_id             = aws_api_gateway_resource.resource_get_schedule.id
  http_method             = aws_api_gateway_method.method_get_schedule.http_method
  integration_http_method = "POST"
  content_handling        = "CONVERT_TO_TEXT"
  type                    = "AWS_PROXY"
  uri                     = aws_lambda_function.get_schedule_api.invoke_arn
}

resource "aws_api_gateway_method" "method_options_get_schedule" {
  rest_api_id          = aws_api_gateway_rest_api.kms_api.id
  resource_id          = aws_api_gateway_resource.resource_get_schedule.id
  api_key_required     = false
  authorization        = "NONE"
  authorization_scopes = []
  http_method          = "OPTIONS"
  request_models       = {}
  request_parameters   = {}
}

resource "aws_api_gateway_method_response" "res_options_get_schedule" {
  http_method = aws_api_gateway_method.method_options_get_schedule.http_method
  rest_api_id = aws_api_gateway_rest_api.kms_api.id
  resource_id = aws_api_gateway_resource.resource_get_schedule.id
  response_models = {
    "application/json" = "Empty"
  }
  response_parameters = {
    "method.response.header.Access-Control-Allow-Headers" = false
    "method.response.header.Access-Control-Allow-Methods" = false
    "method.response.header.Access-Control-Allow-Origin"  = false
  }
  status_code = "200"
}

resource "aws_api_gateway_integration" "int_options_get_schedule" {
  rest_api_id          = aws_api_gateway_rest_api.kms_api.id
  resource_id          = aws_api_gateway_resource.resource_get_schedule.id
  cache_key_parameters = []
  connection_type      = "INTERNET"
  http_method          = "OPTIONS"
  request_parameters   = {}
  request_templates = {
    "application/json" = jsonencode(
      {
        statusCode = 200
      }
    )
  }
  timeout_milliseconds = 29000
  type                 = "MOCK"
}

resource "aws_api_gateway_integration_response" "int_res_options_get_schedule" {
  rest_api_id = aws_api_gateway_rest_api.kms_api.id
  resource_id = aws_api_gateway_resource.resource_get_schedule.id
  http_method = "OPTIONS"
  response_parameters = {
    "method.response.header.Access-Control-Allow-Headers" = "'Content-Type,X-Amz-Date,Authorization,X-Api-Key,X-Amz-Security-Token'"
    "method.response.header.Access-Control-Allow-Methods" = "'OPTIONS'"
    "method.response.header.Access-Control-Allow-Origin"  = "'*'"
  }
  response_templates = {}
  status_code        = aws_api_gateway_method_response.res_options_get_schedule.status_code
}

#########################################################################
# delete_schedule
#########################################################################
resource "aws_api_gateway_resource" "resource_delete_schedule" {
  rest_api_id = aws_api_gateway_rest_api.kms_api.id
  parent_id   = aws_api_gateway_rest_api.kms_api.root_resource_id
  path_part   = "delete_schedule"
}

resource "aws_api_gateway_method" "method_delete_schedule" {
  rest_api_id      = aws_api_gateway_rest_api.kms_api.id
  resource_id      = aws_api_gateway_resource.resource_delete_schedule.id
  http_method      = "POST"
  authorization    = "NONE"
  api_key_required = true
}

resource "aws_api_gateway_integration" "int_delete_schedule" {
  rest_api_id             = aws_api_gateway_rest_api.kms_api.id
  resource_id             = aws_api_gateway_resource.resource_delete_schedule.id
  http_method             = aws_api_gateway_method.method_delete_schedule.http_method
  integration_http_method = "POST"
  content_handling        = "CONVERT_TO_TEXT"
  type                    = "AWS_PROXY"
  uri                     = aws_lambda_function.delete_schedule_api.invoke_arn
}

resource "aws_api_gateway_method" "method_options_delete_schedule" {
  rest_api_id          = aws_api_gateway_rest_api.kms_api.id
  resource_id          = aws_api_gateway_resource.resource_delete_schedule.id
  api_key_required     = false
  authorization        = "NONE"
  authorization_scopes = []
  http_method          = "OPTIONS"
  request_models       = {}
  request_parameters   = {}
}

resource "aws_api_gateway_method_response" "res_options_delete_schedule" {
  http_method = aws_api_gateway_method.method_options_delete_schedule.http_method
  rest_api_id = aws_api_gateway_rest_api.kms_api.id
  resource_id = aws_api_gateway_resource.resource_delete_schedule.id
  response_models = {
    "application/json" = "Empty"
  }
  response_parameters = {
    "method.response.header.Access-Control-Allow-Headers" = false
    "method.response.header.Access-Control-Allow-Methods" = false
    "method.response.header.Access-Control-Allow-Origin"  = false
  }
  status_code = "200"
}

resource "aws_api_gateway_integration" "int_options_delete_schedule" {
  rest_api_id          = aws_api_gateway_rest_api.kms_api.id
  resource_id          = aws_api_gateway_resource.resource_delete_schedule.id
  cache_key_parameters = []
  connection_type      = "INTERNET"
  http_method          = aws_api_gateway_method.method_options_delete_schedule.http_method
  request_parameters   = {}
  request_templates = {
    "application/json" = jsonencode(
      {
        statusCode = 200
      }
    )
  }
  timeout_milliseconds = 29000
  type                 = "MOCK"
}

resource "aws_api_gateway_integration_response" "int_res_options_delete_schedule" {
  rest_api_id = aws_api_gateway_rest_api.kms_api.id
  resource_id = aws_api_gateway_resource.resource_delete_schedule.id
  http_method = "OPTIONS"
  response_parameters = {
    "method.response.header.Access-Control-Allow-Headers" = "'Content-Type,X-Amz-Date,Authorization,X-Api-Key,X-Amz-Security-Token'"
    "method.response.header.Access-Control-Allow-Methods" = "'OPTIONS'"
    "method.response.header.Access-Control-Allow-Origin"  = "'*'"
  }
  response_templates = {}
  status_code        = aws_api_gateway_method_response.res_options_delete_schedule.status_code
}


#########################################################################
# get_info
#########################################################################
resource "aws_api_gateway_resource" "resource_get_info" {
  rest_api_id = aws_api_gateway_rest_api.kms_api.id
  parent_id   = aws_api_gateway_rest_api.kms_api.root_resource_id
  path_part   = "get_info"
}

resource "aws_api_gateway_method" "method_get_info" {
  rest_api_id      = aws_api_gateway_rest_api.kms_api.id
  resource_id      = aws_api_gateway_resource.resource_get_info.id
  http_method      = "GET"
  authorization    = "NONE"
  api_key_required = true
}

resource "aws_api_gateway_integration" "int_get_info" {
  rest_api_id             = aws_api_gateway_rest_api.kms_api.id
  resource_id             = aws_api_gateway_resource.resource_get_info.id
  http_method             = aws_api_gateway_method.method_get_info.http_method
  integration_http_method = "POST"
  content_handling        = "CONVERT_TO_TEXT"
  type                    = "AWS_PROXY"
  uri                     = aws_lambda_function.get_info_api.invoke_arn
}

resource "aws_api_gateway_method" "method_options_get_info" {
  rest_api_id          = aws_api_gateway_rest_api.kms_api.id
  resource_id          = aws_api_gateway_resource.resource_get_info.id
  api_key_required     = false
  authorization        = "NONE"
  authorization_scopes = []
  http_method          = "OPTIONS"
  request_models       = {}
  request_parameters   = {}
}

resource "aws_api_gateway_method_response" "res_options_get_info" {
  http_method = aws_api_gateway_method.method_options_get_info.http_method
  rest_api_id = aws_api_gateway_rest_api.kms_api.id
  resource_id = aws_api_gateway_resource.resource_get_info.id
  response_models = {
    "application/json" = "Empty"
  }
  response_parameters = {
    "method.response.header.Access-Control-Allow-Headers" = false
    "method.response.header.Access-Control-Allow-Methods" = false
    "method.response.header.Access-Control-Allow-Origin"  = false
  }
  status_code = "200"
}

resource "aws_api_gateway_integration" "int_options_get_info" {
  rest_api_id          = aws_api_gateway_rest_api.kms_api.id
  resource_id          = aws_api_gateway_resource.resource_get_info.id
  cache_key_parameters = []
  connection_type      = "INTERNET"
  http_method          = aws_api_gateway_method.method_options_get_info.http_method
  request_parameters   = {}
  request_templates = {
    "application/json" = jsonencode(
      {
        statusCode = 200
      }
    )
  }
  timeout_milliseconds = 29000
  type                 = "MOCK"
}

resource "aws_api_gateway_integration_response" "int_res_options_get_info" {
  rest_api_id = aws_api_gateway_rest_api.kms_api.id
  resource_id = aws_api_gateway_resource.resource_get_info.id
  http_method = "OPTIONS"
  response_parameters = {
    "method.response.header.Access-Control-Allow-Headers" = "'Content-Type,X-Amz-Date,Authorization,X-Api-Key,X-Amz-Security-Token'"
    "method.response.header.Access-Control-Allow-Methods" = "'OPTIONS'"
    "method.response.header.Access-Control-Allow-Origin"  = "'*'"
  }
  response_templates = {}
  status_code        = aws_api_gateway_method_response.res_options_get_info.status_code
}

#########################################################################
# deploy stage and usage plan
#########################################################################
resource "aws_api_gateway_deployment" "deploy_kms_api" {
  depends_on  = [aws_api_gateway_integration.int_upsert_schedule, aws_api_gateway_integration.int_get_schedule, aws_api_gateway_integration.int_delete_schedule, aws_api_gateway_integration.int_options_upsert_schedule, aws_api_gateway_integration.int_options_get_schedule, aws_api_gateway_integration.int_options_delete_schedule, aws_api_gateway_integration.int_get_info, aws_api_gateway_integration.int_options_get_info]
  rest_api_id = aws_api_gateway_rest_api.kms_api.id
  description = "kms_api"
  stage_name  = "prod"
}

resource "aws_api_gateway_stage" "stage_kms_api" {
  stage_name = aws_api_gateway_deployment.deploy_kms_api.stage_name
  # stage_name         = "prod"
  rest_api_id        = aws_api_gateway_rest_api.kms_api.id
  cache_cluster_size = "0.5"
  # ステージのみ更新がなされるとapiに繋がらなくなるので、デプロイメントを指定する
  # deployment_id = "1fs2vm"
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
  depends_on = [aws_api_gateway_stage.stage_kms_api]
}

resource "aws_api_gateway_usage_plan_key" "plan_key_kms_web" {
  key_id        = aws_api_gateway_api_key.kms_web.id
  key_type      = "API_KEY"
  usage_plan_id = aws_api_gateway_usage_plan.plan_kms_web.id
}
