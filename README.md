# Bus Ticketing

A mini-project that lets the passenger create an account, book a bus ticket, and schedule the trip. This project uses some AWS services such as API Gateway, DynamoDB, Lambda, SQS, and EventBridge. It also documents how this project is created. This serves as a familiarization for the AWS services.

## API Usage and Specification

* [User Schema API](docs/api_usage/user.md)
* [Bus Schema API](docs/api_usage/bus.md)
* [Bus Unit Schema API](docs/api_usage/bus_unit.md)
* [Bus Route Schema API](docs/api_usage/bus_trip.md)

## Using `Makefile` to install, bootstrap, and deploy the project

1. Install all the dependencies and bootstrap your project
    ```bash
    dev@dev:~:bus-ticketing$ make init
    ```

    To initialize the project with specific AWS profile, you can pass a parameter called `profile`.
    ```bash
    dev@dev:~:bus-ticketing$ make init profile=profile_name
    ```

2. Deploy the project.
    ```bash
    dev@dev:~:bus-ticketing$ make deploy
    # Deploying with specific AWS profile
    dev@dev:~:bus-ticketing$ make deploy profile=profile_name
    ```

## File Structure
In the root directory we have some configuration files, most of which are language specific.
* **`tsconfig.json`**: TypeScript configuration
* **`jest.config.js`**: Configuration for testing
* **`package.json`**: Manages our node packages and scripts
* **`cdk.json`**: Tells the CDK Toolkit how to execute your app

### `cdk.json`
The **`app`** key tells the CDK CLI how to run our code. The command points to the location of our CDK App.
```bash
npx ts-node --prefer-ts-exts bin/bus-ticketing.ts
```

The *feature flags* in the **`context`** object give us the option to enable or disable some breaking changes that have been made by the AWS CDK team outside of majore version releases. It allow the AWS CDK team to push new features that cause breaking changes without having to wait for a major version release. They can just enable the new functionality for new projects, whereas old projects without the flags will continue to work.

### `bin/bus-ticketing.ts`
Every CDK App can consist of one or more Stacks. You can think of a stack as a unit of deployment. All AWS resources defined within the scope of a stack, either directly or indirectly, are provisioned as a single unit. Because AWS CDK stacks are implemented through AWS CloudFormation stacks, they have the same limitations as in AWS CloudFormation.

If you don't specify `env`, the stack will be environment-agnostic. Account/Region-dependent features and context lookups will not work, but a single synthesized template can be deployed anywhere.

To specialize the stack for the AWS Account and Region that are implied by the current CLI configuration:
```typescript
env: {
   account: process.env.CDK_DEFAULT_ACCOUNT,
   region: process.env.CDK_DEFAULT_REGION
}
```

If you know exactly what `account` and `region` you want to deploy the stack to:
```typescript
env: {
   account: '012345678912',
   region: 'us-east-1'
}
```

For more information, see [Environments](https://docs.aws.amazon.com/cdk/latest/guide/environments.html).

## Useful commands

* `npm install`     install projects dependencies
* `npm run build`   compile typescript to js
* `npm run watch`   watch for changes and compile
* `npm run test`    perform the jest unit tests
* `cdk deploy`      deploy this stack to your default AWS account/region
* `cdk diff`        compare deployed stack with current state
* `cdk synth`       emits the synthesized CloudFormation template
* `cdk bootstrap`   deployment of AWS CloudFormation template to a specific AWS environment (account and region)
* `cdk destroy`     destroy this stack from your default AWS account/region