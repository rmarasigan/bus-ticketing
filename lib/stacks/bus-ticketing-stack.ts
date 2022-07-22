import { Stack, StackProps } from 'aws-cdk-lib';
import { Construct } from 'constructs';
import * as iam from 'aws-cdk-lib/aws-iam';
import * as apigw from 'aws-cdk-lib/aws-apigateway';
import { ManagedPolicy } from 'aws-cdk-lib/aws-iam';

export class BusTicketingStack extends Stack
{
  constructor(scope: Construct, id: string, props?: StackProps)
  {
    super(scope, id, props);

    // Custom IAM Role
    const Role = new iam.Role(this, 'BusTicketing_CustomRole', {
      roleName: 'BusTicketing_CustomRole',
      assumedBy: new iam.ServicePrincipal('lambda.amazonaws.com'),
      managedPolicies: [
        ManagedPolicy.fromAwsManagedPolicyName('service-role/AWSLambdaBasicExecutionRole'),
        ManagedPolicy.fromManagedPolicyArn(this, 'AmazonSQSReadOnlyAccess', 'arn:aws:iam::aws:policy/AmazonSQSReadOnlyAccess'),
      ]
    });
  }
}
