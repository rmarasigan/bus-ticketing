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
  }
}
