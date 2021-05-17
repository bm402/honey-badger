#!/bin/bash

# change ssh port for remote access
sed -i 's/#Port 22/Port 44422/g' /etc/ssh/sshd_config
sudo service sshd restart

# install dependencies
yum update -y
yum install -y git
yum install -y golang

# download listener code
cd ~
git clone https://github.com/bncrypted/honey-badger.git
cd honey-badger
git fetch --all
git checkout cdk-setup

# run listeners
go run listener/cmd/main.go
