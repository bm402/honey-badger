Content-Type: multipart/mixed; boundary="//"
MIME-Version: 1.0

--//
Content-Type: text/cloud-config; charset="us-ascii"
MIME-Version: 1.0
Content-Transfer-Encoding: 7bit
Content-Disposition: attachment; filename="cloud-config.txt"

#cloud-config
cloud_final_modules:
- [scripts-user, always]

--//
Content-Type: text/x-shellscript; charset="us-ascii"
MIME-Version: 1.0
Content-Transfer-Encoding: 7bit
Content-Disposition: attachment; filename="userdata.txt"

#!/bin/bash

# change ssh port for remote access
sed -i 's/#Port 22/Port 44422/g' /etc/ssh/sshd_config
service sshd restart

# install dependencies
yum update -y
yum install -y git
yum install -y golang

# switch to ec2-user
su ec2-user
export HOME=/home/ec2-user

# download listener code
cd ~
if [ ! -d honey-badger ]; then
    git clone https://github.com/bncrypted/honey-badger.git
fi
cd honey-badger
git fetch --all
git checkout cdk-setup

# run listeners
go run listener/cmd/main.go
