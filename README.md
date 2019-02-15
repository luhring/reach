# cnct

[![Go Report Card](https://goreportcard.com/badge/github.com/luhring/cnct)](https://goreportcard.com/report/github.com/luhring/cnct)
[![CircleCI](https://circleci.com/gh/luhring/cnct.svg?style=svg)](https://circleci.com/gh/luhring/cnct)
[![GitHub license](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/luhring/cnct/blob/master/LICENSE)

"connect" / Cloud Network Connectivity Tool

A tool for diagnosing network connectivity issues in AWS, for better understanding your cloud networks, and for asserting expectations about network security in the cloud.

## Overview

cnct evaluates the potential for network connectivity between EC2 instances by querying the AWS API for network configuration data.

cnct is able to determine the ports an EC2 instance can access on another EC2 instance, taking into consideration security group rules, instance subnet placements, instance running states, network ACL rules, and route tables.

cnct doesn't need to run on any of the EC2 instances you're evaluating, it just needs to run on a system that has access to the AWS API.

**Disclaimer:** Because cnct gets all of its information from the AWS API, cnct makes no guarantees about network accessibility with respect to the operating system or applications running on an EC2 instance. In other words, cnct can tell you if your VPC resources and EC2 instances are configured correctly, but cnct _can't_ tell you if an OS firewall is blocking network traffic, or if an application listening on a port has crashed.

## Uses

You can ask what ports are reachable on an EC2 instance from the perspective of another instance.

```
$ cnct "client instance" "server instance"
Is "client instance" able to connect to "server instance"?
Yes, on these ports:
- TCP 80
- TCP 443
```

You can ask about connectivity to just one specific port.

```
$ cnct "web-server" "db-server" --port 1433
Is "web-server" able to connect to "db-server" on port TCP 1433?
No.
Why not?
- The instance "db-server" doesn't have any security groups with an inbound rule that allows access on port TCP 1433.
- The subnet "database-private-subnet" in which the instance "db-server" resides doesn't have any network ACL rules that allow inbound traffic on port TCP 1433.
```

You can ask about the outbound access a specific instance has to all other instances.

```
$ cnct "first-instance" --outbound
To which instances can "first-instance" connect?
- "second-instance" (i-0028edd2bac34109d) on ports TCP 80, 443
- "another-instance" (i-004cf8d8f0921234c) on ports TCP 80, 443, UDP 53
```

You can ask about inbound access to a specific instance from all other instances.

```
$ cnct "lonely-instance" --inbound
Which instances can connect to "lonely-instance"?
No instances.
```

## CLI Syntax

`cnct "first-instance" ["second-instance"] [OPTIONS]`

### Options

`--port`, `-p` Filter the connectivity assessment down to a single port. Useful when you need to check only a single network service, such as HTTP (i.e. `--port 80`) or MySQL (i.e. `--port 3306`). Can only be used to specify **TCP** ports, for now.

`--inbound` (Can't be used if more than one instance is specified.) Evaluate network access **to** specified instance from all other instances.

`--outbound` (Can't be used if more than one instance is specified.) Evaluate network access **from** specified instance to all other instances.

`--error-if-none`, `-n` Return an error (exit code `2`) if the result set contains **no** items.

`--error-if-any`, `-a` Return an error (exit code `2`) if the result set contains **any** items.

### Specifying an instance

cnct is able to handle various methods of specifying an EC2 instance.

- **By Name tag.** Most people assign a descriptive name to each of their EC2 instances via a "Name" [tag](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/Using_Tags.html). Example: `"my-database-server"`.
- **By Instance ID.** This is a reliably unique identifier within the EC2 service, and it's assigned by AWS. Example: `"i-08a43985c56df54e4"`.

#### Notes about instance specification strings

1. **Quotes:** Use quotes to surround instance specification strings as necessary within your shell environment. Quotes never hurt, but sometimes they can be left off -- for example, when using an instance name tag that contains no spaces, only letters and hyphens.

1. **Shortened strings:** To make CLI text entry less tedious, cnct only requires enough of an instance specification string to be unique within the current AWS account and region. For example, let's say you have three EC2 instances in a particular region in your AWS account, and these instances are named "web-01", "web-02", and "db-master". You could refer to the "db-master" instance by typing just `"db"`, since no other instance begins with the same text. This would let you type the command `cnct db --inbound`, and cnct would understand that you were asking about inbound access to the "db-master" instance. However, it would not be sufficient to use the string `"web"`, since multiple instances have name tags that begin with the text "web". The rules for shortened strings apply to both name tags and instance IDs.

## Road map

- ~~Connectivity between EC2 instances within a single subnet~~
- Connectivity between EC2 instances in separate subnets in a single VPC
- Connectivity betweeen an EC2 instance and an IP address
