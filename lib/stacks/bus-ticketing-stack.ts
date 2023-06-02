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
      description: 'A Lambda Function that will process API requests and get the user account information',
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
      description: 'A Lambda Function that will process API requests and update the user account information',
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
      description: 'A Lambda Function that will process API requests and create a new bus line information',
      environment: {
        "BUS_TABLE": BusTable.tableName
      }
    });
    BusTable.grantReadWriteData(createBus);
    createBus.applyRemovalPolicy(REMOVAL_POLICY);

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
      }
    });

    const updateUserApiIntegration = new apigw.LambdaIntegration(updateUser);
    const updateUserApi = UserAccountApiRoot.addResource('update');
    updateUserApi.addMethod('POST', updateUserApiIntegration, {
      requestParameters: {
        'method.request.querystring.id': true,
        'method.request.querystring.username': true
      }
    });

    // ***** Bus API Specification ***** //
    const BusApiRoot = api.root.addResource('bus');

    const createBusApiIntegration = new apigw.LambdaIntegration(createBus);
    const createBusApi = BusApiRoot.addResource('create');
    createBusApi.addMethod('POST', createBusApiIntegration);
  }
}
