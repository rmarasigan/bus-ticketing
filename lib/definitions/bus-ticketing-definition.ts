import * as cdk from 'aws-cdk-lib';
/**
 ** Reference:
 ** https://docs.aws.amazon.com/apigateway/latest/developerguide/request-response-data-mappings.html
 ** https://docs.aws.amazon.com/apigateway/latest/developerguide/api-gateway-mapping-template-reference.html
**/
const ContentTemplate = `#foreach($type in $allParams.keySet())
#set($params = $allParams.get($type))
#if($type == "querystring" || $type == "header")
"$type": {
#foreach($paramName in $params.keySet())
"$paramName" : "$util.escapeJavaScript($params.get($paramName))"
    #if($foreach.hasNext),#end
#end
#end
#end`;

const ContextTemplate = `"user-agent": "$context.identity.userAgent",
"method": "$context.httpMethod",
"protocol": "$context.protocol",
"request-time": "$context.requestTime",
"epoch": "$context.requestTimeEpoch",
"path": "$context.path",
"error": "$context.error.message"`;

export function IntegrationResponse(statusCode: number)
{
  let integrationResponses = {
    statusCode: statusCode.toString(),
    responseTemplates: {
      'application/json':`#set($allParams = $input.params())
      {
      "status": 200,
      "content": {${ContentTemplate}},
      "context": {${ContextTemplate}}
      }`
    },
    responseParameters: {
      'method.response.header.Content-Type': 'integration.response.header.Content-Type',
    }
  }

  return integrationResponses;
}

export function MethodResponse(statusCode: number, model: cdk.aws_apigateway.Model)
{
  return {
    statusCode: statusCode.toString(),
    responseModels: {
      'application/json': model
    },
    responseParameters: {
      'method.response.header.Content-Type': true
    }
  }
}