import * as apigw from 'aws-cdk-lib/aws-apigateway';
/**
 * Reference:
 * https://goessner.net/articles/JsonPath/
 * https://github.com/tdegrunt/jsonschema
 * https://docs.aws.amazon.com/cdk/api/v2/docs/aws-cdk-lib.aws_apigateway.JsonSchema.html
 * https://docs.aws.amazon.com/apigateway/latest/developerguide/request-response-data-mappings.html
 * https://docs.aws.amazon.com/apigateway/latest/developerguide/api-gateway-mapping-template-reference.html
**/

/**
 * Represents the data structure of the creation of *User* payload and
 * accepts an object with the following fields and are validated:
 * 
 * `user_type`, `first_name`, `last_name`, `username`, `password`, `address`, `email`, `mobile_number`
 * 
 * @param api REST API that this model is part of.
 */
export function UserApiModel(api: apigw.RestApi) {
  return api.addModel('BusTicketingUserModel', {
    modelName: 'BusTicketingUserModel',
    description: 'A User Schema that will be validated before it goes to the Lambda Function',
    schema: {
      type: apigw.JsonSchemaType.OBJECT,
      properties: {
        user_type: {
          pattern: '^.+',
          type: apigw.JsonSchemaType.STRING
        },
        first_name: {
          pattern: '^.+',
          type: apigw.JsonSchemaType.STRING
        },
        last_name: {
          pattern: '^.+',
          type: apigw.JsonSchemaType.STRING
        },
        username: {
          pattern: '^.+',
          type: apigw.JsonSchemaType.STRING
        },
        password: {
          minLength: 8,
          type: apigw.JsonSchemaType.STRING
        },
        address: {
          pattern: '^.+',
          type: apigw.JsonSchemaType.STRING
        },
        email: {
          pattern: '^.+',
          type: apigw.JsonSchemaType.STRING
        },
        mobile_number: {
          pattern: '^.+',
          type: apigw.JsonSchemaType.STRING
        }
      },
      required: [ 'user_type', 'first_name', 'last_name', 'username', 'password', 'address', 'email', 'mobile_number' ]
    }
  });
}

/**
 * Represents the data structure of the *User Login* payload that has the login
 * information and accepts an object with the following fields that are validated:
 * 
 * `username`, `password`
 * 
 * @param api REST API that this model is part of.
 */
export function UserLoginApiModel(api: apigw.RestApi) {
  return api.addModel('BusTicketingUserLoginModel', {
    modelName: 'BusTicketingUserLoginModel',
    description: 'A User Login Schema that will be validated before it goes to the Lambda Function',
    schema: {
      type: apigw.JsonSchemaType.OBJECT,
      properties: {
        username: {
          pattern: '^.+',
          type: apigw.JsonSchemaType.STRING
        },
        password: {
          minLength: 8,
          type: apigw.JsonSchemaType.STRING
        },
      },
      required: [ 'username', 'password' ]
    }
  });
}

/**
 * Represents the data structure of the creation of *Bus* payload and
 * accepts an array of objects with the following fields and are validated:
 * 
 * `name`, `owner`, `email`, `address`, `company`, `mobile_number`
 *
 * @param api REST API that this model is part of.
**/
export function BusLineApiModel(api: apigw.RestApi) {
  return api.addModel('BusTicketingBusLineModel', {
    modelName: 'BusTicketingBusLineModel',
    description: 'A Bus Schema that will be validated before it goes to the Lambda Function',
    schema: {
      minItems: 1,
      type: apigw.JsonSchemaType.ARRAY,
      items: {
        type: apigw.JsonSchemaType.OBJECT,
        properties: {
          name: {
            pattern: '^.+',
            type: apigw.JsonSchemaType.STRING
          },
          owner: {
            pattern: '^.+',
            type: apigw.JsonSchemaType.STRING
          },
          email: {
            pattern: '^.+',
            type: apigw.JsonSchemaType.STRING
          },
          address: {
            pattern: '^.+',
            type: apigw.JsonSchemaType.STRING
          },
          company: {
            pattern: '^.+',
            type: apigw.JsonSchemaType.STRING
          },
          mobile_number: {
            pattern: '^.+',
            type: apigw.JsonSchemaType.STRING
          }
        },
        required: [ 'name', 'owner', 'email', 'address', 'company', 'mobile_number' ]
      }
    }
  });
}

/**
 * Represents the data structure of the creation of *Bus Unit* payload and
 * accepts an array of objects with the following fields and are validated:
 * 
 * `bus_id`, `code`, `active`, `min_capacity`, `max_capacity`
 *
 * @param api REST API that this model is part of.
**/
export function BusUnitApiModel(api: apigw.RestApi) {
  return api.addModel('BusTicketingBusUnitModel', {
    modelName: 'BusTicketingBusUnitModel',
     description: 'A Bus Unit Schema that will be validated before it goes to the Lambda Function',
     schema: {
      minItems: 1,
      type: apigw.JsonSchemaType.ARRAY,
      items: {
        type: apigw.JsonSchemaType.OBJECT,
        properties: {
          bus_id: {
            pattern: '^.+',
            type: apigw.JsonSchemaType.STRING
          },
          code: {
            pattern: '^.+',
            type: apigw.JsonSchemaType.STRING
          },
          active: {
            type: apigw.JsonSchemaType.BOOLEAN
          },
          min_capacity: {
            minimum: 25,
            type: apigw.JsonSchemaType.NUMBER
          },
          max_capacity: {
            minimum: 25,
            type: apigw.JsonSchemaType.NUMBER,
          },
        }
      },
      required: [ 'bus_id', 'code', 'active', 'min_capacity', 'max_capacity' ]
     }
  });
}

/**
 * Represents the data structure of the creation of *Bus Route* payload and
 * accepts an object with the following fields and are validated:
 * 
 * `bus_id`, `bus_unit_id`, `currency_code`, `rate`, `available`, `departure_time`
 * `arrival_time`, `from_route`, `to_route`
 *
 * @param api REST API that this model is part of.
**/
export function BusRouteApiModel(api: apigw.RestApi) {
  return api.addModel('BusTickingBusRouteApiModel', {
    modelName: 'BusTickingBusRouteApiModel',
    schema: {
      type: apigw.JsonSchemaType.OBJECT,
      properties: {
        bus_id: {
          pattern: '^.+',
          type: apigw.JsonSchemaType.STRING
        },
        bus_unit_id: {
          pattern: '^.+',
          type: apigw.JsonSchemaType.STRING
        },
        currency_code: {
          pattern: '^.+',
          type: apigw.JsonSchemaType.STRING
        },
        rate: {
          minimum: 1,
          type: apigw.JsonSchemaType.NUMBER
        },
        available: {
          type: apigw.JsonSchemaType.BOOLEAN
        },
        departure_time: {
          pattern: '^.+',
          type: apigw.JsonSchemaType.STRING
        },
        arrival_time: {
          pattern: '^.+',
          type: apigw.JsonSchemaType.STRING
        },
        from_route: {
          pattern: '^.+',
          type: apigw.JsonSchemaType.STRING
        },
        to_route: {
          pattern: '^.+',
          type: apigw.JsonSchemaType.STRING
        }
      },
      required: [ 'bus_id', 'bus_unit_id', 'currency_code', 'rate', 'available', 'departure_time', 'arrival_time', 'from_route', 'to_route' ]
    }
  });
}