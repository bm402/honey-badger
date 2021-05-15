#!/usr/bin/env node
import 'source-map-support/register';
import * as cdk from '@aws-cdk/core';
import { HoneyBadgerStack } from '../lib/honey-badger-stack';

const app = new cdk.App();
new HoneyBadgerStack(app, 'HoneyBadgerStack', {
  env: { 
      account: process.env.CDK_DEFAULT_ACCOUNT,
      region: process.env.CDK_DEFAULT_REGION
    },
});
