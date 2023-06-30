import * as cdk from 'aws-cdk-lib';
import { Construct } from 'constructs';
import * as sqs from 'aws-cdk-lib/aws-sqs';
import * as lambda from 'aws-cdk-lib/aws-lambda';
import * as apigw from 'aws-cdk-lib/aws-apigateway';
import * as dynamodb from 'aws-cdk-lib/aws-dynamodb';
import * as eventsource from 'aws-cdk-lib/aws-lambda-event-sources';
import { UserApiModel, UserLoginApiModel, BusLineApiModel, BusUnitApiModel, BusRouteApiModel, BookingApiModel } from '../definitions/bus-ticketing-api-model';

export class BusTicketingStack extends cdk.Stack
{
  constructor(scope: Construct, id: string, props?: cdk.StackProps)
  {
    super(scope, id, props);

    //  When the resource is removed from the app, it will be physically destroyed.
    const REMOVAL_POLICY = cdk.RemovalPolicy.DESTROY;

    // ******************** DynamoDB ******************** //
    // 1. Create a DynamoDB Table that will contain the basic user record
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

    // 4. Create a DynamoDB Table that will contain the Bus Route information that has a
    // partition and sort key.
    const BusRouteTable = new dynamodb.Table(this, 'BusTicketing_BusRouteTable', {
      tableName: 'BusTicketing_BusRouteTable',
      partitionKey: {
        name: "id",
        type: dynamodb.AttributeType.STRING
      },
      sortKey: {
        name: "bus_id",
        type: dynamodb.AttributeType.STRING
      }
    });

    // 5. Create a DynamoDB Table that will contain the Booking information that has a
    // partition and sort key.
    const BookingTable = new dynamodb.Table(this, 'BusTicketing_BookingTable', {
      tableName: 'BusTicketing_BookingTable',
      partitionKey: {
        name: "id",
        type: dynamodb.AttributeType.STRING
      },
      sortKey: {
        name: "bus_route_id",
        type: dynamodb.AttributeType.STRING
      }
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

    const filterBus = new lambda.Function(this, 'filterBus', {
      memorySize: 1024,
      handler: 'filterBus',
      functionName: 'filterBus',
      runtime: lambda.Runtime.GO_1_X,
      timeout: cdk.Duration.seconds(60),
      code: lambda.Code.fromAsset('cmd/bus/filterBus'),
      description: 'A Lambda Function that will process API requests and filter the bus line record depending on the passed query',
      environment: {
        "BUS_TABLE": BusTable.tableName
      }
    });
    BusTable.grantReadData(filterBus);
    filterBus.applyRemovalPolicy(REMOVAL_POLICY);

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

    const filterBusUnit = new lambda.Function(this, 'filterBusUnit', {
      memorySize: 1024,
      handler: 'filterBusUnit',
      functionName: 'filterBusUnit',
      runtime: lambda.Runtime.GO_1_X,
      timeout: cdk.Duration.seconds(60),
      code: lambda.Code.fromAsset('cmd/bus_unit/filterBusUnit'),
      description: 'A Lambda Function that will process API requests and filter the bus unit record depending on the passed query',
      environment: {
        "BUS_UNIT_TABLE": BusUnitTable.tableName
      }
    });
    BusUnitTable.grantReadData(filterBusUnit);
    filterBusUnit.applyRemovalPolicy(REMOVAL_POLICY);

    // ***** Bus Route Lambda Functions Specification ***** //
    const createBusRoute = new lambda.Function(this, 'createBusRoute', {
      memorySize: 1024,
      handler: 'createBusRoute',
      functionName: 'createBusRoute',
      runtime: lambda.Runtime.GO_1_X,
      timeout: cdk.Duration.seconds(60),
      code: lambda.Code.fromAsset('cmd/bus_route/createBusRoute'),
      description: 'A Lambda Function that will process API requests and create a new bus route record',
      environment: {
        "BUS_ROUTE_TABLE": BusRouteTable.tableName
      }
    });
    BusRouteTable.grantReadWriteData(createBusRoute);
    createBusRoute.applyRemovalPolicy(REMOVAL_POLICY);

    const getBusRoute = new lambda.Function(this, 'getBusRoute', {
      memorySize: 1024,
      handler: 'getBusRoute',
      functionName: 'getBusRoute',
      runtime: lambda.Runtime.GO_1_X,
      timeout: cdk.Duration.seconds(60),
      code: lambda.Code.fromAsset('cmd/bus_route/getBusRoute'),
      description: 'A Lambda Function that will process API requests and fetch the bus unit route record',
      environment: {
        "BUS_ROUTE_TABLE": BusRouteTable.tableName
      }
    });
    BusRouteTable.grantReadData(getBusRoute);
    getBusRoute.applyRemovalPolicy(REMOVAL_POLICY);

    const filterBusRoute = new lambda.Function(this, 'filterBusRoute', {
      memorySize: 1024,
      handler: 'filterBusRoute',
      functionName: 'filterBusRoute',
      runtime: lambda.Runtime.GO_1_X,
      timeout: cdk.Duration.seconds(60),
      code: lambda.Code.fromAsset('cmd/bus_route/filterBusRoute'),
      description: 'A Lambda Function that will process API requests and filter the bus unit route record depending on the passed query',
      environment: {
        "BUS_ROUTE_TABLE": BusRouteTable.tableName
      }
    });
    BusRouteTable.grantReadData(filterBusRoute);
    filterBusRoute.applyRemovalPolicy(REMOVAL_POLICY);

    const updateBusRoute = new lambda.Function(this, 'updateBusRoute', {
      memorySize: 1024,
      handler: 'updateBusRoute',
      functionName: 'updateBusRoute',
      runtime: lambda.Runtime.GO_1_X,
      timeout: cdk.Duration.seconds(60),
      code: lambda.Code.fromAsset('cmd/bus_route/updateBusRoute'),
      description: 'A Lambda Function that will process API requests and update the bus route record',
      environment: {
        "BUS_ROUTE_TABLE": BusRouteTable.tableName
      }
    });
    BusRouteTable.grantReadWriteData(updateBusRoute);
    updateBusRoute.applyRemovalPolicy(REMOVAL_POLICY);

    // ***** Booking Lambda Functions and SQS Specification ***** //
    // SQS QUEUE
    // 1. Create a deadletter queue that will contain the unsuccessfully
    // processed and should have a ".fifo" to the queue name.
    const deadLetterQueue = new sqs.Queue(this, 'bus-ticketing-booking-dlq.fifo', {
      fifo: true,
      contentBasedDeduplication: true,
      queueName: 'bus-ticketing-booking-dlq.fifo',
      removalPolicy: REMOVAL_POLICY
    });

    // 2. Create a queue that is configured to be a FIFO queue with deadletter
    // queue. It is needed to add a ".fifo" to the queue name.
    const bookingQueue = new sqs.Queue(this, 'bus-ticketing-booking.fifo', {
      fifo: true,
      queueName: 'bus-ticketing-booking.fifo',
      deadLetterQueue: {
        queue: deadLetterQueue,
        maxReceiveCount: 5
      },
      removalPolicy: REMOVAL_POLICY,
      contentBasedDeduplication: true,
      visibilityTimeout: cdk.Duration.seconds(120)
    });
    
    const createBooking = new lambda.Function(this, 'createBooking', {
      memorySize: 1024,
      handler: 'createBooking',
      functionName: 'createBooking',
      runtime: lambda.Runtime.GO_1_X,
      timeout: cdk.Duration.seconds(90),
      code: lambda.Code.fromAsset('cmd/bookings/createBooking'),
      description: 'A Lambda Function that will process API requests and sends new booking record to the SQS Queue',
      environment: {
        "BOOKING_QUEUE": bookingQueue.queueUrl
      }
    });
    bookingQueue.grantSendMessages(createBooking);
    createBooking.applyRemovalPolicy(REMOVAL_POLICY);

    const processBooking = new lambda.Function(this, 'processBooking', {
      memorySize: 1024,
      handler: 'processBooking',
      functionName: 'processBooking',
      runtime: lambda.Runtime.GO_1_X,
      timeout: cdk.Duration.seconds(90),
      code: lambda.Code.fromAsset('cmd/bookings/processBooking'),
      description: 'A Lambda Function that will process SQS events and process booking record',
      environment: {
        "BOOKING_TABLE": BookingTable.tableName
      }
    });
    BookingTable.grantReadWriteData(processBooking);
    processBooking.applyRemovalPolicy(REMOVAL_POLICY);

    processBooking.addEventSource(new eventsource.SqsEventSource(bookingQueue, {
      enabled: true,
      batchSize: 1,
      reportBatchItemFailures: true
    }));

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
    api.addGatewayResponse('BusTicketing_BadRequestGatewayResponse', {
      statusCode: '400',
      type: apigw.ResponseType.BAD_REQUEST_BODY,
      templates: {
        'application/json': `{"message": "$context.error.validationErrorString"}`
      }
    });

    const ApiParameterValidator = new apigw.RequestValidator(this, 'BusTicketing_ApiParameterValidator', {
      restApi: api,
      validateRequestParameters: true,
      requestValidatorName: 'BusTicketing_ApiParameterValidator'
    });

    const ApiRequestBodyValidator = new apigw.RequestValidator(this, 'BusTicketing_ApiRequestBodyValidator', {
      restApi: api,
      validateRequestBody: true,
      requestValidatorName: 'BusTicketing_ApiRequestBodyValidator'
    });

    // ***** User API Specification ***** //
    const UserApiRoot = api.root.addResource('user');
    UserApiRoot.applyRemovalPolicy(REMOVAL_POLICY);

    const UserAccountApiRoot = UserApiRoot.addResource('account');
    UserAccountApiRoot.applyRemovalPolicy(REMOVAL_POLICY);

    const UserModel = UserApiModel(api);
    const createUserApiIntegration = new apigw.LambdaIntegration(createUser);
    const createUserApi = UserApiRoot.addResource('create');
    createUserApi.addMethod('POST', createUserApiIntegration, {
      requestModels: {
        'application/json': UserModel
      },
      requestValidator: ApiRequestBodyValidator
    });

    const UserLoginModel = UserLoginApiModel(api);
    const loginUserApiIntegration = new apigw.LambdaIntegration(login);
    const loginUserApi = UserApiRoot.addResource('login');
    loginUserApi.addMethod('POST', loginUserApiIntegration, {
      requestModels: {
        'application/json': UserLoginModel
      },
      requestValidator: ApiRequestBodyValidator
    });

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
    BusApiRoot.applyRemovalPolicy(REMOVAL_POLICY);

    const BusLineModel = BusLineApiModel(api);
    const createBusApiIntegration = new apigw.LambdaIntegration(createBus);
    const createBusApi = BusApiRoot.addResource('create');
    createBusApi.addMethod('POST', createBusApiIntegration, {
      requestModels: {
        'application/json': BusLineModel
      },
      requestValidator: ApiRequestBodyValidator
    });

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

    const filterBusApiIntegration = new apigw.LambdaIntegration(filterBus);
    const filterBusApi = BusApiRoot.addResource('search');
    filterBusApi.addMethod('GET', filterBusApiIntegration, {
      requestParameters: {
        'method.request.querystring.name': true,
        'method.request.querystring.company': true
      }
    });

    // ***** Bus Unit API Specification ***** //
    const BusUnitApiRoot = api.root.addResource('bus-unit');
    BusUnitApiRoot.applyRemovalPolicy(REMOVAL_POLICY);

    const BusUnitModel = BusUnitApiModel(api);
    const createBusUnitApiIntegration = new apigw.LambdaIntegration(createBusUnit);
    const createBusUnitApi = BusUnitApiRoot.addResource('create');
    createBusUnitApi.addMethod('POST', createBusUnitApiIntegration, {
      requestModels: {
        'application/json': BusUnitModel
      },
      requestValidator: ApiRequestBodyValidator
    });

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

    const filterBusUnitApiIntegration = new apigw.LambdaIntegration(filterBusUnit);
    const filterBusUnitApi = BusUnitApiRoot.addResource('search');
    filterBusUnitApi.addMethod('GET', filterBusUnitApiIntegration, {
      requestParameters: {
        'method.request.querystring.bus_id': true
      },
      requestValidator: ApiParameterValidator
    });

    // ***** Bus Route API Specification ***** //
    const BusRouteApiRoot = api.root.addResource('bus-route');
    BusRouteApiRoot.applyRemovalPolicy(REMOVAL_POLICY);

    const BusRouteModel = BusRouteApiModel(api);
    const createBusRouteApiIntegration = new apigw.LambdaIntegration(createBusRoute);
    const createBusRouteApi = BusRouteApiRoot.addResource('create');
    createBusRouteApi.addMethod('POST', createBusRouteApiIntegration, {
      requestModels: {
        'application/json': BusRouteModel
      },
      requestValidator: ApiRequestBodyValidator
    });

    const getBusRouteApiIntegration = new apigw.LambdaIntegration(getBusRoute);
    const getBusRouteApi = BusRouteApiRoot.addResource('get');
    getBusRouteApi.addMethod('GET', getBusRouteApiIntegration, {
      requestParameters: {
        'method.request.querystring.id': true,
        'method.request.querystring.bus_id': true
      },
      requestValidator: ApiParameterValidator
    });

    const filterBusRouteApiIntegration = new apigw.LambdaIntegration(filterBusRoute);
    const filterBusRouteApi = BusRouteApiRoot.addResource('search');
    filterBusRouteApi.addMethod('GET', filterBusRouteApiIntegration, {
      requestParameters: {
        'method.request.querystring.bus_id': true
      },
      requestValidator: ApiParameterValidator
    });

    const updateBusRouteApiIntegration = new apigw.LambdaIntegration(updateBusRoute);
    const updateBusRouteApi = BusRouteApiRoot.addResource('update');
    updateBusRouteApi.addMethod('POST', updateBusRouteApiIntegration, {
      requestParameters: {
        'method.request.querystring.id': true,
        'method.request.querystring.bus_id': true
      },
      requestValidator: ApiParameterValidator
    });

    // ***** Booking API Specification ***** //
    const BookingApiRoot = api.root.addResource('bookings');
    BookingApiRoot.applyRemovalPolicy(REMOVAL_POLICY);

    const BookingModel = BookingApiModel(api);
    const createBookingApiIntegration = new apigw.LambdaIntegration(createBooking);
    const createBookingApi = BookingApiRoot.addResource('create');
    createBookingApi.addMethod('POST', createBookingApiIntegration, {
      requestModels: {
        'application/json': BookingModel
      },
      requestValidator: ApiRequestBodyValidator
    });
  }
}