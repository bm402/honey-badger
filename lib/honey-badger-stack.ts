import * as cdk from '@aws-cdk/core';
import * as ec2 from '@aws-cdk/aws-ec2';
import * as iam from '@aws-cdk/aws-iam';
import * as fs from 'fs';

export class HoneyBadgerStack extends cdk.Stack {
    constructor(scope: cdk.Construct, id: string, props?: cdk.StackProps) {
        super(scope, id, props);

        const defaultVpc = ec2.Vpc.fromLookup(this, 'VPC', { 
            isDefault: true,
        });

        const role = new iam.Role(this, 'InstanceRole', {
            assumedBy: new iam.ServicePrincipal('ec2.amazonaws.com'),
        });

        const securityGroup = new ec2.SecurityGroup(this, 'SecurityGroup', {
            vpc: defaultVpc,
            allowAllOutbound: true,
        });

        // temporary SSH access
        securityGroup.addIngressRule(ec2.Peer.anyIpv4(), ec2.Port.tcp(22));
        securityGroup.addIngressRule(ec2.Peer.anyIpv4(), ec2.Port.tcp(44422));

        const instanceType = ec2.InstanceType.of(ec2.InstanceClass.T2, ec2.InstanceSize.MICRO);

        const machineImage = ec2.MachineImage.latestAmazonLinux({
            generation: ec2.AmazonLinuxGeneration.AMAZON_LINUX_2,
        });

        const userData = ec2.UserData.custom(fs.readFileSync('scripts/ec2-user-data.sh', 'utf8'));

        const instance = new ec2.Instance(this, 'EC2', {
            vpc: defaultVpc,
            role: role,
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
