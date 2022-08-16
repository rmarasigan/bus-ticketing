# Bus Ticketing

A mini-project that lets the passenger create an account, book a bus ticket, and schedule the trip. This project uses some AWS services such as API Gateway, DynamoDB, Lambda, SQS, and EventBridge. It also documents how this project is created. This serves as a familiarization for the AWS services.

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

## Create new CDK App
To list the available options we can use the `--list` parameter flag to the command.
```bash
dev@dev:~$ cdk init --list
Available templates:
* app: Template for a CDK Application
   └─ cdk init app --language=[csharp|fsharp|go|java|javascript|python|typescript]
* lib: Template for a CDK Construct Library
   └─ cdk init lib --language=typescript
* sample-app: Example CDK Application with some constructs
   └─ cdk init sample-app --language=[csharp|fsharp|go|java|javascript|python|typescript]
```

There are 3 templates we can start from:
* **`app`**: A basic starter template
* **`lib`**: A template for writing a CDK construct library
* **`sample-app`**: A starter with some constructs included

In this project we are going to use the `app` template with the `TypeScript` language. To initializie a project, we pass `--language` parameter.
```bash
dev@dev:~$ cdk init app --language typescript
```

## File Structure
In the root directory we have some configuration files, most of which are language specific.
* **`tsconfig.json`**: TypeScript configuration
* **`jest.config.js`**: Configuration for testing
* **`package.json`**: Manages our node packages and scripts
* **`cdk.json`**: Tells the CDK Toolkit how to execute your app

### `cdk.json`
The **`app`** key tells the CDK CLI how to run our code. The command points to the location of our CDK App.
```json
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

If you know exactly what Account and Region you want to deploy the stack to:
```typescript
env: {
   account: '012345678912',
   region: 'us-east-1'
}
```

For more information, see [Environments](https://docs.aws.amazon.com/cdk/latest/guide/environments.html).

## `npm`
**npm** (node package manager) is the dependency or package manager you get out of the box when you install Node.js. The npm manages all the packages and modules for Node.js and consists of command-line client `npm`. A *package* contains all the files needed for a module and *modules* are the JavaScript libraries that can be included in the Node project according to the requirements of the project.

When executables are installed via npm packages, npm creates links to them:
* **local** installs have links create at the `./node_modules/.bin/` directory
* **global** installs have links created from the global `bin/` directory

## `npx`
**npx** stands for **Node Package Execute** and it comes with the npm, when you installed npm above 5.2.0 version then automatically npx will be installed. It is an npm package runner that can execute any package that you want from the npm registry without even installing the package.

If npx is not installed you can install that by using `npm install` command.
```bash
dev@dev:~$ npm install -g npx
```

## Godoc
Godoc extracts and generates documentation for Go programs. It runs a web server and presents the documentation as a web page.

Install `godoc` with the following command:
```bash
dev@dev:~$ go get golang.org/x/tools/cmd/godoc
go get: installing executables with 'go get' in module mode is deprecated.
	Use 'go install pkg@version' instead.
	For more information, see https://golang.org/doc/go-get-install-deprecation
	or run 'go help get' or 'go help install'.
```

If it throws an issue of `module mode is deprecated`. Try this command to install the latest godoc:
```bash
dev@dev:~$ go install golang.org/x/tools/cmd/godoc@latest
```

Check whether the godoc is successfully installed:
```bash
dev@dev:~$ which godoc
/home/dev/go/bin/godoc
```

To start a `godoc` server:
```bash
dev@dev:~/go/src/github.com/username/repo-name$ godoc -http:=6060
using module mode; GOMOD=/home/dev/go/src/github.com/username/repo-name/go.mod
```

In your browser, visit: [godoc-server](http://localhost:6060)

## Reference
* [Godoc Command](https://stackoverflow.com/a/61300854)
* [Regex Numbers](https://stackoverflow.com/a/39399503)
* [Flag Features](https://github.com/aws/aws-cdk/blob/v1-main/packages/@aws-cdk/cx-api/lib/features.ts)
* [0055 Feature Flags](https://github.com/aws/aws-cdk-rfcs/blob/master/text/0055-feature-flags.md)
* [`npm` CLI Commands](https://docs.npmjs.com/cli/v6/commands)
* [Godoc Documentation](https://pkg.go.dev/golang.org/x/tools/cmd/godoc)
* [How to use Context in AWS CDK](https://bobbyhadz.com/blog/how-to-use-context-aws-cdk)
* [What does CDK.JSON do in AWS CDK](https://bobbyhadz.com/blog/cdk-json-aws-cdk)
* [npm vs npx — What’s the Difference?](https://www.freecodecamp.org/news/npm-vs-npx-whats-the-difference/)
* [Epoch & Unix Timestamp Conversion Tools](https://www.epochconverter.com/)
* [What are the differences between npm and npx?](https://www.geeksforgeeks.org/what-are-the-differences-between-npm-and-npx/)
* [AWS CDK Tutorial for Beginners - Step-by-Step Guide](https://bobbyhadz.com/blog/aws-cdk-tutorial-typescript)

### AWS Documentation
* [Stacks](https://docs.aws.amazon.com/cdk/v2/guide/stacks.html)
* [x-amazon-apigateway-integration.responseTemplates object](https://docs.aws.amazon.com/apigateway/latest/developerguide/api-gateway-swagger-extensions-integration-responseTemplates.html)
* [x-amazon-apigateway-integration.requestParameters object](https://docs.aws.amazon.com/apigateway/latest/developerguide/api-gateway-swagger-extensions-integration-requestParameters.html)
* [API Gateway mapping template and access logging variable reference](https://docs.aws.amazon.com/apigateway/latest/developerguide/api-gateway-mapping-template-reference.html)
* [Amazon API Gateway API request and response data mapping reference](https://docs.aws.amazon.com/apigateway/latest/developerguide/request-response-data-mappings.html)