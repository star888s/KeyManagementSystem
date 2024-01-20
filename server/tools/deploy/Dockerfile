# FROM golang:1.21.0
FROM centos:7
# FROM　alpine:latest

# RUN yum install -y epel-release
# RUN yum install -y wget
# RUN wget https://go.dev/dl/go1.21.6.linux-amd64.tar.gz
# RUN tar -C /usr/local -xzf go1.21.6.linux-amd64.tar.gz

# ENV export PATH=$PATH:/usr/local/go/bin
# ENV export GOPATH=$HOME/go
# ENV export PATH=$PATH:$GOPATH/bin

# RUN go version

# git入れてあげるといいかも
# dockerの初期設定系も
# sudo顕現周り調べろ
RUN yum update -y
RUN yum install -y wget
RUN wget https://go.dev/dl/go1.21.6.linux-amd64.tar.gz
RUN  rm -rf /usr/local/go && tar -C /usr/local -xzf go1.21.6.linux-amd64.tar.gz

ENV PATH $PATH:/usr/local/go/bin

# Terraform install
RUN yum install -y yum-utils \
    && yum-config-manager --add-repo https://rpm.releases.hashicorp.com/RHEL/hashicorp.repo \
    && yum -y install terraform

# aws vault install 
RUN curl -L https://github.com/99designs/aws-vault/releases/download/v7.2.0/aws-vault-linux-amd64 -o /usr/local/bin/aws-vault \
    && chmod 755 /usr/local/bin/aws-vault \
    &&  aws-vault --version

# RUN yum install python3 -y
