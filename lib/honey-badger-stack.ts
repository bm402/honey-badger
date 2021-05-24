import * as cdk from '@aws-cdk/core';
import * as apigw from "@aws-cdk/aws-apigateway";
import * as ec2 from '@aws-cdk/aws-ec2';
import * as dynamodb from '@aws-cdk/aws-dynamodb';
import * as iam from '@aws-cdk/aws-iam';
import * as lambda from '@aws-cdk/aws-lambda';
import * as lambdaEventSources from '@aws-cdk/aws-lambda-event-sources';
import * as ssm from '@aws-cdk/aws-ssm';
import * as sqs from '@aws-cdk/aws-sqs';
import * as fs from 'fs';
import * as path from 'path'

export class HoneyBadgerStack extends cdk.Stack {
    constructor(scope: cdk.Construct, id: string, props?: cdk.StackProps) {
        super(scope, id, props);

        // roles
        const listenerInstanceRole = new iam.Role(this, 'InstanceRole', {
            assumedBy: new iam.ServicePrincipal('ec2.amazonaws.com'),
        });
        const aggregatorLambdaExecutionRole = new iam.Role(this, 'LambdaExecutionRole', {
            assumedBy: new iam.ServicePrincipal('lambda.amazonaws.com'),
            managedPolicies: [
                iam.ManagedPolicy.fromAwsManagedPolicyName("service-role/AWSLambdaBasicExecutionRole"),
            ],
        });
        const heatmapDataLambdaExecutionRole = new iam.Role(this, 'LambdaExecutionRole', {
            assumedBy: new iam.ServicePrincipal('lambda.amazonaws.com'),
            managedPolicies: [
                iam.ManagedPolicy.fromAwsManagedPolicyName("service-role/AWSLambdaBasicExecutionRole"),
            ],
        });

        // dynamodb table for raw logs
        const rawLogsTable = new dynamodb.Table(this, 'RawLogsTable', {
            partitionKey: { 
                name: 'ingress_port',
                type: dynamodb.AttributeType.STRING,
            },
            sortKey: {
                name: 'timestamp',
                type: dynamodb.AttributeType.NUMBER,
            },
            readCapacity: 2,
            writeCapacity: 2,
            stream: dynamodb.StreamViewType.NEW_IMAGE,
        });
        rawLogsTable.grantWriteData(listenerInstanceRole);

        // put dynamodb table name in parameter store
        const rawLogsTableNameParam = new ssm.StringParameter(this, 'StringParameterRawLogsTableName', {
            parameterName: 'RawLogsTableName',
            stringValue: rawLogsTable.tableName,
        });
        rawLogsTableNameParam.grantRead(listenerInstanceRole);

        // dynamodb table for aggregated logs
        const aggregatedLogsTable = new dynamodb.Table(this, 'AggregatedLogsTable', {
            partitionKey: { 
                name: 'lat_lon',
                type: dynamodb.AttributeType.STRING,
            },
            readCapacity: 2,
            writeCapacity: 2,
        });
        aggregatedLogsTable.grantReadWriteData(aggregatorLambdaExecutionRole);
        aggregatedLogsTable.grantReadData(heatmapDataLambdaExecutionRole);

        // lambda function for aggregating log data
        const aggregatorLambda = new lambda.Function(this, 'AggregatorLambda', {
            code: lambda.Code.fromAsset(path.join(__dirname, '../aggregator'), {
                bundling: {
                    image: lambda.Runtime.GO_1_X.bundlingImage,
                    user: "root",
                    command: [
                        'bash', '-c', [
                            'GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o /asset-output/main *.go',
                        ].join(' && ')
                    ]
                },
            }),
            handler: 'main',
            runtime: lambda.Runtime.GO_1_X,
            role: aggregatorLambdaExecutionRole,
            environment: {
                'AGGREGATED_LOGS_TABLE_NAME': aggregatedLogsTable.tableName,
            },
        });

        // dead letter queue for raw logs table stream
        const rawLogsTableStreamDLQ = new sqs.Queue(this, 'DeadLetterQueue');

        // triggers the aggregator lambda when new data is written to the raw logs table
        aggregatorLambda.addEventSource(new lambdaEventSources.DynamoEventSource(rawLogsTable, {
            startingPosition: lambda.StartingPosition.TRIM_HORIZON,
            batchSize: 5,
            bisectBatchOnError: true,
            onFailure: new lambdaEventSources.SqsDlq(rawLogsTableStreamDLQ),
            retryAttempts: 10
        }));

        // api gateway for http data retrieval
        const api = new apigw.RestApi(this, "DataApi", {
            restApiName: "HoneyBadgerDataApi",
        });

        // lambda function for retrieving heatmap data from aggregated logs table
        const heatmapDataLambda = new lambda.Function(this, 'HeatmapDataLambda', {
            code: lambda.Code.fromAsset(path.join(__dirname, '../api/heatmap-data'), {
                bundling: {
                    image: lambda.Runtime.GO_1_X.bundlingImage,
                    user: "root",
                    command: [
                        'bash', '-c', [
                            'GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o /asset-output/main *.go',
                        ].join(' && ')
                    ]
                },
            }),
            handler: 'main',
            runtime: lambda.Runtime.GO_1_X,
            role: heatmapDataLambdaExecutionRole,
            environment: {
                'AGGREGATED_LOGS_TABLE_NAME': aggregatedLogsTable.tableName,
            },
        });
        const heatmapDataIntegration = new apigw.LambdaIntegration(heatmapDataLambda);
        api.root.addMethod('GET', heatmapDataIntegration)

        // set up listener on ec2
        const defaultVpc = ec2.Vpc.fromLookup(this, 'VPC', { 
            isDefault: true,
        });

        const listenerSecurityGroup = new ec2.SecurityGroup(this, 'SecurityGroup', {
            vpc: defaultVpc,
            allowAllOutbound: true,
        });

        // SSH access
        listenerSecurityGroup.addIngressRule(ec2.Peer.anyIpv4(), ec2.Port.tcp(44422));

        // honeypot port access
        const commonPorts: number[] = [21,22,23,53,80,110,135,139,143,443,445,993,995,1723,3306,3389,5900,8080];
        for (let port of commonPorts) {
            listenerSecurityGroup.addIngressRule(ec2.Peer.anyIpv4(), ec2.Port.tcp(port));
            listenerSecurityGroup.addIngressRule(ec2.Peer.anyIpv6(), ec2.Port.tcp(port));
        }

        const listenerInstanceType = ec2.InstanceType.of(ec2.InstanceClass.T2, ec2.InstanceSize.MICRO);

        const listenerMachineImage = ec2.MachineImage.latestAmazonLinux({
            generation: ec2.AmazonLinuxGeneration.AMAZON_LINUX_2,
        });

        const listenerUserData = ec2.UserData.custom(fs.readFileSync('scripts/listener-instance-user-data.sh', 'utf8'));

        const listenerInstance = new ec2.Instance(this, 'EC2', {
            vpc: defaultVpc,
            role: listenerInstanceRole,
            securityGroup: listenerSecurityGroup,
            instanceType: listenerInstanceType,
            machineImage: listenerMachineImage,
            userData: listenerUserData,

            // key pair created manually before deployment
            keyName: 'honey-badger-ec2-key',
        });

        new cdk.CfnOutput(this, 'ListenerInstancePublicIP', {
            value: listenerInstance.instancePublicIp
        });
    }
}
