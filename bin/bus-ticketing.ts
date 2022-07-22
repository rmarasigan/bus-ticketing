#!/usr/bin/env node
import 'source-map-support/register';
import * as cdk from 'aws-cdk-lib';
import { BusTicketingStack } from '../lib/stacks/bus-ticketing-stack';

const app = new cdk.App();

new BusTicketingStack(app, 'BusTicketingStack', {
  env: {
    account: process.env.CDK_DEFAULT_ACCOUNT,
    region: process.env.CDK_DEFAULT_REGION
  }
});