import * as cdk from 'aws-cdk-lib';
import { Construct } from 'constructs';
import * as lambda from 'aws-cdk-lib/aws-lambda';
import * as dynamodb from 'aws-cdk-lib/aws-dynamodb';
import * as apigw from 'aws-cdk-lib/aws-apigateway';

export class BusTicketingStack extends cdk.Stack
{
  constructor(scope: Construct, id: string, props?: cdk.StackProps)
  {
    super(scope, id, props);

    //  When the resource is removed from the app, it will be physically destroyed.
    const REMOVAL_POLICY = cdk.RemovalPolicy.DESTROY;

    // ******************** DynamoDB ******************** //
    // 1. Create a DynamoDB Table that will contain the basic user record/information
    // that has a partition and sort key.
    const UsersTable = new dynamodb.Table(this, 'BusTicketing_UsersTable', {
      tableName: 'BusTicketing_UsersTable',
      partitionKey: {
        name: "username",
        type: dynamodb.AttributeType.STRING
      },
      sortKey: {
        name: "id",
        type: dynamodb.AttributeType.STRING
      },
      billingMode: dynamodb.BillingMode.PAY_PER_REQUEST,
      removalPolicy: REMOVAL_POLICY
    });

    // 2. Create a DynamoDB Table that will contain the bus line information that has
    // a partition and sort key.
    const BusTable = new dynamodb.Table(this, 'BusTicketing_BusTable', {
      tableName: 'BusTicketing_BusTable',
      partitionKey: {
        name: "name",
        type: dynamodb.AttributeType.STRING
      },
      sortKey: {
        name: "company",
        type: dynamodb.AttributeType.STRING
      },
      billingMode: dynamodb.BillingMode.PAY_PER_REQUEST,
      removalPolicy: REMOVAL_POLICY
    });

    // 3. Create a DynamoDB Table that will contain the Bus Line Units information that
    // has a partition and sorty key.
    const BusUnitTable = new dynamodb.Table(this, 'BusTicketing_BusUnitTable', {
      tableName: 'BusTicketing_BusUnitTable',
      partitionKey: {
        name: "code",
        type: dynamodb.AttributeType.STRING
      },
      sortKey: {
        name: "bus_id",
        type: dynamodb.AttributeType.STRING
      },
      billingMode: dynamodb.BillingMode.PAY_PER_REQUEST,
      removalPolicy: REMOVAL_POLICY
    });

    // ******************** Lambda Functions ******************** //
    // ***** User Lambda Functions Specification ***** //
    const createUser = new lambda.Function(this, 'createUser', {
      memorySize: 1024,
      handler: 'createUser',
      functionName: 'createUser',
      runtime: lambda.Runtime.GO_1_X,
      timeout: cdk.Duration.seconds(60),
      code: lambda.Code.fromAsset('cmd/user/createUser'),
      description: 'A Lambda Function that will process API requests and create a new user account',
      environment: {
        "USERS_TABLE": UsersTable.tableName
      }
    });
    UsersTable.grantReadWriteData(createUser);
    createUser.applyRemovalPolicy(REMOVAL_POLICY);

    const login = new lambda.Function(this, 'login', {
      memorySize: 1024,
      handler: 'login',
      functionName: 'login',
      runtime: lambda.Runtime.GO_1_X,
      timeout: cdk.Duration.seconds(60),
      code: lambda.Code.fromAsset('cmd/user/login'),
      description: 'A Lambda Function that will process API requests and login the user account',
      environment: {
        "USERS_TABLE": UsersTable.tableName
      }
    });
    UsersTable.grantReadWriteData(login);
    login.applyRemovalPolicy(REMOVAL_POLICY);

    const getUser = new lambda.Function(this, 'getUser', {
      memorySize: 1024,
      handler: 'getUser',
      functionName: 'getUser',
      runtime: lambda.Runtime.GO_1_X,
      timeout: cdk.Duration.seconds(60),
      code: lambda.Code.fromAsset('cmd/user/getUser'),
      description: 'A Lambda Function that will process API requests and get the user account record',
      environment: {
        "USERS_TABLE": UsersTable.tableName
      }
    });
    UsersTable.grantReadData(getUser);
    getUser.applyRemovalPolicy(REMOVAL_POLICY);

    const updateUser = new lambda.Function(this, 'updateUser', {
      memorySize: 1024,
      handler: 'updateUser',
      functionName: 'updateUser',
      runtime: lambda.Runtime.GO_1_X,
      timeout: cdk.Duration.seconds(60),
      code: lambda.Code.fromAsset('cmd/user/updateUser'),
      description: 'A Lambda Function that will process API requests and update the user account record',
      environment: {
        "USERS_TABLE": UsersTable.tableName
      }
    });
    UsersTable.grantReadWriteData(updateUser);
    updateUser.applyRemovalPolicy(REMOVAL_POLICY);

    // ***** Bus Lambda Functions Specification ***** //
    const createBus = new lambda.Function(this, 'createBus', {
      memorySize: 1024,
      handler: 'createBus',
      functionName: 'createBus',
      runtime: lambda.Runtime.GO_1_X,
      timeout: cdk.Duration.seconds(60),
      code: lambda.Code.fromAsset('cmd/bus/createBus'),
      description: 'A Lambda Function that will process API requests and create a new bus line record',
      environment: {
        "BUS_TABLE": BusTable.tableName
      }
    });
    BusTable.grantReadWriteData(createBus);
    createBus.applyRemovalPolicy(REMOVAL_POLICY);

    const getBus = new lambda.Function(this, 'getBus', {
      memorySize: 1024,
      handler: 'getBus',
      functionName: 'getBus',
      runtime: lambda.Runtime.GO_1_X,
      timeout: cdk.Duration.seconds(60),
      code: lambda.Code.fromAsset('cmd/bus/getBus'),
      description: 'A Lambda Function that will process API requets and returns the bus line record',
      environment: {
        "BUS_TABLE": BusTable.tableName
      }
    });
    BusTable.grantReadData(getBus);
    getBus.applyRemovalPolicy(REMOVAL_POLICY);

    const updateBus = new lambda.Function(this, 'updateBus', {
      memorySize: 1024,
      handler: 'updateBus',
      functionName: 'updateBus',
      runtime: lambda.Runtime.GO_1_X,
      timeout:cdk.Duration.seconds(60),
      code: lambda.Code.fromAsset('cmd/bus/updateBus'),
      description: 'A Lambda Function that will process API requests and update the bus line record',
      environment: {
        "BUS_TABLE": BusTable.tableName
      }
    });
    BusTable.grantReadWriteData(updateBus);
    updateBus.applyRemovalPolicy(REMOVAL_POLICY);

    // ***** Bus Unit Lambda Functions Specification ***** //
    const createBusUnit = new lambda.Function(this, 'createBusUnit', {
      memorySize: 1024,
      handler: 'createBusUnit',
      functionName: 'createBusUnit',
      runtime: lambda.Runtime.GO_1_X,
      timeout: cdk.Duration.seconds(60),
      code: lambda.Code.fromAsset('cmd/bus_unit/createBusUnit'),
      description: 'A Lambda Function that will process API requests and create a new bus line unit record',
      environment: {
        "BUS_UNIT_TABLE": BusUnitTable.tableName
      }
    });
    BusUnitTable.grantReadWriteData(createBusUnit);
    createBusUnit.applyRemovalPolicy(REMOVAL_POLICY);

    const getBusUnit = new lambda.Function(this, 'getBusUnit', {
      memorySize: 1024,
      handler: 'getBusUnit',
      functionName: 'getBusUnit',
      runtime: lambda.Runtime.GO_1_X,
      timeout: cdk.Duration.seconds(60),
      code: lambda.Code.fromAsset('cmd/bus_unit/getBusUnit'),
      description: 'A Lambda Function that will process API requests and fetch the bus line unit record',
      environment: {
        "BUS_UNIT_TABLE": BusUnitTable.tableName
      }
    });
    BusUnitTable.grantReadData(getBusUnit);
    getBusUnit.applyRemovalPolicy(REMOVAL_POLICY);

    const updateBusUnit = new lambda.Function(this, 'updateBusUnit', {
      memorySize: 1024,
      handler: 'updateBusUnit',
      functionName: 'updateBusUnit',
      runtime: lambda.Runtime.GO_1_X,
      timeout: cdk.Duration.seconds(60),
      code: lambda.Code.fromAsset('cmd/bus_unit/updateBusUnit'),
      description: 'A Lambda Function that will process API requests and update the bus line unit record',
      environment: {
        "BUS_UNIT_TABLE": BusUnitTable.tableName
      }
    });
    BusUnitTable.grantReadWriteData(updateBusUnit);
    updateBusUnit.applyRemovalPolicy(REMOVAL_POLICY);

    // ******************** API Gateway ******************** //
    const api = new apigw.RestApi(this, 'bus-ticketing-api', {
      deploy: true,
      restApiName: 'bus-ticketing-api',
      description: 'An API for Users, Bus, and Transaction.',
      deployOptions: {
        stageName: 'prod',
        metricsEnabled: true,
        tracingEnabled: true,
        loggingLevel: apigw.MethodLoggingLevel.INFO
      }
    });
    api.applyRemovalPolicy(REMOVAL_POLICY);

    const ApiParameterValidator = new apigw.RequestValidator(this, 'BusTicketing_ApiParameterValidator', {
      restApi: api,
      validateRequestParameters: true,
      requestValidatorName: 'BusTicketing_ApiParameterValidator'
    });

    // ***** User API Specification ***** //
    const UserApiRoot = api.root.addResource('user');
    const UserAccountApiRoot = UserApiRoot.addResource('account');

    const createUserApiIntegration = new apigw.LambdaIntegration(createUser);
    const createUserApi = UserApiRoot.addResource('create');
    createUserApi.addMethod('POST', createUserApiIntegration);

    const loginUserApiIntegration = new apigw.LambdaIntegration(login);
    const loginUserApi = UserApiRoot.addResource('login');
    loginUserApi.addMethod('POST', loginUserApiIntegration);

    const getUserApiIntegration = new apigw.LambdaIntegration(getUser);
    const getUserApi = UserAccountApiRoot.addResource('get');
    getUserApi.addMethod('GET', getUserApiIntegration, {
      requestParameters: {
        'method.request.querystring.id': true,
        'method.request.querystring.username': true
      },
      requestValidator: ApiParameterValidator
    });

    const updateUserApiIntegration = new apigw.LambdaIntegration(updateUser);
    const updateUserApi = UserAccountApiRoot.addResource('update');
    updateUserApi.addMethod('POST', updateUserApiIntegration, {
      requestParameters: {
        'method.request.querystring.id': true,
        'method.request.querystring.username': true
      },
      requestValidator: ApiParameterValidator
    });

    // ***** Bus API Specification ***** //
    const BusApiRoot = api.root.addResource('bus');

    const createBusApiIntegration = new apigw.LambdaIntegration(createBus);
    const createBusApi = BusApiRoot.addResource('create');
    createBusApi.addMethod('POST', createBusApiIntegration);

    const getBusApiIntegration = new apigw.LambdaIntegration(getBus);
    const getBusApi = BusApiRoot.addResource('get');
    getBusApi.addMethod('GET', getBusApiIntegration, {
      requestParameters: {
        'method.request.querystring.id': true,
        'method.request.querystring.name': true
      },
      requestValidator: ApiParameterValidator
    });

    const updateBusApiIntegration = new apigw.LambdaIntegration(updateBus);
    const updateBusApi = BusApiRoot.addResource('update');
    updateBusApi.addMethod('POST', updateBusApiIntegration, {
      requestParameters: {
        'method.request.querystring.id': true,
        'method.request.querystring.name': true
      },
      requestValidator: ApiParameterValidator
    });

    // ***** Bus Unit API Specification ***** //
    const BusUnitApiRoot = api.root.addResource('bus_unit');

    const createBusUnitApiIntegration = new apigw.LambdaIntegration(createBusUnit);
    const createBusUnitApi = BusUnitApiRoot.addResource('create');
    createBusUnitApi.addMethod('POST', createBusUnitApiIntegration);

    const getBusUnitApiIntegration = new apigw.LambdaIntegration(getBusUnit);
    const getBusUnitApi = BusUnitApiRoot.addResource('get');
    getBusUnitApi.addMethod('GET', getBusUnitApiIntegration, {
      requestParameters: {
        'method.request.querystring.code': true,
        'method.request.querystring.bus_id': true
      },
      requestValidator: ApiParameterValidator
    });

    const updateBusUnitApiIntegration = new apigw.LambdaIntegration(updateBusUnit);
    const updateBusUnitApi = BusUnitApiRoot.addResource('update');
    updateBusUnitApi.addMethod('POST', updateBusUnitApiIntegration, {
      requestParameters: {
        'method.request.querystring.code': true,
        'method.request.querystring.bus_id': true
      },
      requestValidator: ApiParameterValidator
    });
  }
}