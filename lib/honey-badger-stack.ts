import * as cdk from '@aws-cdk/core';
import * as ec2 from '@aws-cdk/aws-ec2';
import * as dynamodb from '@aws-cdk/aws-dynamodb';
import * as iam from '@aws-cdk/aws-iam';
import * as ssm from '@aws-cdk/aws-ssm';
import * as fs from 'fs';

export class HoneyBadgerStack extends cdk.Stack {
    constructor(scope: cdk.Construct, id: string, props?: cdk.StackProps) {
        super(scope, id, props);

        // roles
        const instanceRole = new iam.Role(this, 'InstanceRole', {
            assumedBy: new iam.ServicePrincipal('ec2.amazonaws.com'),
        });

        // dynamodb table for raw logs
        const table = new dynamodb.Table(this, 'RawLogsTable', {
            partitionKey: { 
                name: 'ingress_port',
                type: dynamodb.AttributeType.STRING,
            },
            sortKey: {
                name: 'timestamp',
                type: dynamodb.AttributeType.STRING,
            },
            readCapacity: 2,
            writeCapacity: 2,
        });
        table.grantWriteData(instanceRole)

        // put dynamodb table name in parameter store
        const param = new ssm.StringParameter(this, 'StringParameterRawLogsTableName', {
            parameterName: 'RawLogsTableName',
            stringValue: table.tableName,
        });
        param.grantRead(instanceRole);

        // set up listener on ec2
        const defaultVpc = ec2.Vpc.fromLookup(this, 'VPC', { 
            isDefault: true,
        });

        const securityGroup = new ec2.SecurityGroup(this, 'SecurityGroup', {
            vpc: defaultVpc,
            allowAllOutbound: true,
        });

        // SSH access
        securityGroup.addIngressRule(ec2.Peer.anyIpv4(), ec2.Port.tcp(44422));

        // honeypot port access
        const commonPorts: number[] = [21,22,23,53,80,110,135,139,143,443,445,993,995,1723,3306,3389,5900,8080];
        for (let port of commonPorts) {
            securityGroup.addIngressRule(ec2.Peer.anyIpv4(), ec2.Port.tcp(port));
            securityGroup.addIngressRule(ec2.Peer.anyIpv6(), ec2.Port.tcp(port));
        }

        const instanceType = ec2.InstanceType.of(ec2.InstanceClass.T2, ec2.InstanceSize.MICRO);

        const machineImage = ec2.MachineImage.latestAmazonLinux({
            generation: ec2.AmazonLinuxGeneration.AMAZON_LINUX_2,
        });

        const userData = ec2.UserData.custom(fs.readFileSync('scripts/ec2-user-data.sh', 'utf8'));

        const instance = new ec2.Instance(this, 'EC2', {
            vpc: defaultVpc,
            role: instanceRole,
            securityGroup: securityGroup,
            instanceType: instanceType,
            machineImage: machineImage,
            userData: userData,

        // key pair created manually before deployment
        keyName: 'honey-badger-ec2-key',
        });

        new cdk.CfnOutput(this, 'EC2PublicIP', {
            value: instance.instancePublicIp
        });
    }
}
