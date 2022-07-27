import * as cdk from 'aws-cdk-lib';
import { Construct } from 'constructs';
import * as iam from 'aws-cdk-lib/aws-iam';
import * as lambda from 'aws-cdk-lib/aws-lambda';
import * as dynamodb from 'aws-cdk-lib/aws-dynamodb';
import * as apigw from 'aws-cdk-lib/aws-apigateway';
import { ManagedPolicy } from 'aws-cdk-lib/aws-iam';

export class BusTicketingStack extends cdk.Stack
{
  constructor(scope: Construct, id: string, props?: cdk.StackProps)
  {
    super(scope, id, props);

    //  When the resource is removed from the app,
    // it will be physically destroyed.
    const REMOVAL_POLICY = cdk.RemovalPolicy.DESTROY;

    // Custom IAM Role
    const BusTicketing_CustomRole = new iam.Role(this, 'BusTicketing_CustomRole', {
      roleName: 'BusTicketing_CustomRole',
      assumedBy: new iam.ServicePrincipal('lambda.amazonaws.com'),
      managedPolicies: [
        ManagedPolicy.fromAwsManagedPolicyName('service-role/AWSLambdaBasicExecutionRole'),
        ManagedPolicy.fromManagedPolicyArn(this, 'AmazonSQSReadOnlyAccess', 'arn:aws:iam::aws:policy/AmazonSQSReadOnlyAccess'),
      ]
    });

    // ---------- DynamoDB ---------- //
    const UsersTable = new dynamodb.Table(this, 'BusTicketing_UsersTable', {
      tableName: 'BusTicketing_UsersTable',
      partitionKey: {
        name: "id",
        type: dynamodb.AttributeType.STRING
      },
      sortKey: {
        name: "type",
        type: dynamodb.AttributeType.STRING
      },
      billingMode: dynamodb.BillingMode.PAY_PER_REQUEST,
      removalPolicy: REMOVAL_POLICY
    });

    const BusTable = new dynamodb.Table(this, 'BusTicketing_BusTable', {
      tableName: 'BusTicketing_BusTable',
      partitionKey: {
        name: "id",
        type: dynamodb.AttributeType.STRING
      },
      sortKey: {
        name: "company",
        type: dynamodb.AttributeType.STRING
      },
      billingMode: dynamodb.BillingMode.PAY_PER_REQUEST,
      removalPolicy: REMOVAL_POLICY
    });

    const BusUnitTable = new dynamodb.Table(this, 'BusTicketing_BusUnitTable', {
      tableName: 'BusTicketing_BusUnitTable',
      partitionKey: {
        name: "id",
        type: dynamodb.AttributeType.STRING
      },
      sortKey: {
        name: "bus",
        type: dynamodb.AttributeType.STRING
      },
      billingMode: dynamodb.BillingMode.PAY_PER_REQUEST,
      removalPolicy: REMOVAL_POLICY
    });

    // ---------- Lambda Functions ---------- //
    const users = new lambda.Function(this, 'BusTicketing_Users', {
      functionName: 'BusTicketing_Users',
      description: 'An API Lambda Integration that will process API requests coming from Users API.',
      handler: 'users',
      memorySize: 1024,
      role: BusTicketing_CustomRole,
      runtime: lambda.Runtime.GO_1_X,
      timeout: cdk.Duration.seconds(60),
      code: lambda.Code.fromAsset('cmd/users'),
      environment: {
        "USERS_TABLE": UsersTable.tableName
      }
    });
    UsersTable.grantReadWriteData(users);
    users.applyRemovalPolicy(REMOVAL_POLICY);

    const bus = new lambda.Function(this, 'BusTicketing_Bus', {
      functionName: 'BusTicketing_Bus',
      description: 'An API Lambda Integration that will process API requests coming from Bus API',
      handler: 'bus',
      memorySize: 1024,
      role: BusTicketing_CustomRole,
      runtime: lambda.Runtime.GO_1_X,
      timeout: cdk.Duration.seconds(60),
      code: lambda.Code.fromAsset('cmd/bus'),
      environment: {
        "BUS_TABLE": BusTable.tableName
      }
    });
    BusTable.grantReadWriteData(bus);
    bus.applyRemovalPolicy(REMOVAL_POLICY);

    const busUnit = new lambda.Function(this, 'BusTicketing_BusUnit', {
      functionName: 'BusTicketing_BusUnit',
      description: 'An API Lambda Integration that will process API requests coming from Bus Unit API',
      handler: 'busUnit',
      memorySize: 1024,
      role: BusTicketing_CustomRole,
      runtime: lambda.Runtime.GO_1_X,
      timeout: cdk.Duration.seconds(60),
      code: lambda.Code.fromAsset('cmd/busUnit'),
      environment: {
        "BUS_TABLE": BusTable.tableName,
        "BUS_UNIT_TABLE": BusUnitTable.tableName
      }
    });
    BusTable.grantReadWriteData(busUnit);
    BusUnitTable.grantReadWriteData(busUnit);
    busUnit.applyRemovalPolicy(REMOVAL_POLICY);

    const busRoute = new lambda.Function(this, 'BusTicketing_BusRoute', {
      functionName: 'BusTicketing_BusRoute',
      description: 'An API Lambda Integration that will process API requests coming from Bus Route API',
      handler: 'busRoute',
      memorySize: 1024,
      role: BusTicketing_CustomRole,
      runtime: lambda.Runtime.GO_1_X,
      timeout: cdk.Duration.seconds(60),
      code: lambda.Code.fromAsset('cmd/busRoute')
    });
    busRoute.applyRemovalPolicy(REMOVAL_POLICY);

    const busTrip = new lambda.Function(this, 'BusTicketing_BusTrip', {
      functionName: 'BusTicketing_BusTrip',
      description: 'An API Lambda Integration that will process API requests coming from Bus Trip API',
      handler: 'busTrip',
      memorySize: 1024,
      role: BusTicketing_CustomRole,
      runtime: lambda.Runtime.GO_1_X,
      timeout: cdk.Duration.seconds(60),
      code: lambda.Code.fromAsset('cmd/busTrip')
    });
    busTrip.applyRemovalPolicy(REMOVAL_POLICY);

    // ---------- API Gateway ---------- //
    const api = new apigw.RestApi(this, 'bus-ticketing-api', {
      deploy: true,
      restApiName: 'bus-ticketing-api',
      description: 'An API for Users, Bus, and Transaction.'
    });
    api.applyRemovalPolicy(REMOVAL_POLICY);

    // ---------- API Gateway Integration ---------- //
    // Integrates AWS Lambda function to an API Gateway Method.
    const usersApiIntegration = new apigw.LambdaIntegration(users);

    const usersApi = api.root.addResource('users');
    usersApi.addMethod('GET', usersApiIntegration, {
      requestParameters: {
        'method.request.querystring.id': true
      }
    });

    usersApi.addMethod('POST', usersApiIntegration, {
      requestParameters: {
        'method.request.querystring.type': true
      }
    });

    const busApiIntegration = new apigw.LambdaIntegration(bus);
    const busUnitApiIntegration = new apigw.LambdaIntegration(busUnit);
    const busTripApiIntegration = new apigw.LambdaIntegration(busTrip);
    const busRouteApiIntegration = new apigw.LambdaIntegration(busRoute);

    const busApi = api.root.addResource('bus');
    busApi.addMethod('GET', busApiIntegration);
    busApi.addMethod('POST', busApiIntegration, {
      requestParameters: {
        'method.request.querystring.type': true
      }
    });

    const busUnitApi = busApi.addResource('unit');
    busUnitApi.addMethod('GET', busUnitApiIntegration);
    busUnitApi.addMethod('POST', busUnitApiIntegration, {
      requestParameters: {
        'method.request.querystring.bus': true,
        'method.request.querystring.type': true
      }
    });

    const busRouteApi = busApi.addResource('route');
    busRouteApi.addMethod('GET', busRouteApiIntegration);
    busRouteApi.addMethod('POST', busRouteApiIntegration);

    const busTripApi = busApi.addResource('trip');
    busTripApi.addMethod('GET', busTripApiIntegration);
    busTripApi.addMethod('POST', busTripApiIntegration);
  }
}
