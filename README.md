# reach

[![CircleCI](https://circleci.com/gh/luhring/reach.svg?style=svg)](https://circleci.com/gh/luhring/reach)
[![GitHub license](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/luhring/reach/blob/master/LICENSE)

Reach is a tool for discovering the impact your AWS configuration has on the flow of network traffic.

## Getting Started

To perform an analysis, specify a **source** EC2 instance and a **destination** EC2 instance:

```Text
$ reach <source> <destination>
```

![Image](.data/reach-demo.gif)

Reach uses your AWS configuration to analyze the potential for network connectivity between two EC2 instances in your AWS account. This means **you don't need to install Reach on any servers** — you just need access to the AWS API.

The key benefits of Reach are:

- **Solve problems faster:** Find missing links in a network path in _seconds_, not hours.

- **Don't compromise on security:** Tighten security without worrying about impacting any required network flows.

- **Learn about your network:** Gain better insight into currently allowed network flows, and discover new consequences of your network design.

- **Build better pipelines:** Discover network-level problems before running application integration or end-to-end tests by adding Reach to your CI/CD pipelines.

## Basic Usage

The values for `source` and `destination` should each uniquely identify an EC2 instance in your AWS account. You can use an **instance ID** or a **name tag**, and you can enter just the first few characters instead of the entire value, as long as what you've entered matches exactly one EC2 instance.

Some examples:

```Text
$ reach i-0452993c7efa3a314 i-02b8dfb5537e80860
```

```Text
$ reach i-04 i-02
```

```Text
$ reach web-instance database-instance
```

```Text
$ reach web data
```

## More Features

### Assertions

If you deploy infrastructure via CI/CD pipelines, it can be helpful to validate the network design itself before running any tests that rely on a correct network configuration.

You can use assertion flags to ensure that your source **can** or **cannot** reach your destination.

If an assertion succeeds, Reach exits  `0`. If an assertion fails, Reach exits `2`.

To confirm that the source **can** reach the destination:

```Text
$ reach web-server database-server --assert-reachable
```

To confirm that the source **cannot** reach the destination:

```Text
$ reach some-server super-sensitive-server --assert-not-reachable
```

### Explanations

Normally, Reach's output is very basic. It displays a simple list of zero or more kinds of network traffic that are allowed to flow from the source to the destination. However, the process Reach uses to perform its analysis is more complex.

If you're troubleshooting a network problem in AWS, it's probably more helpful to see _"why"_ the analysis result is what it is.

You can tell Reach to expose the reasoning behind the displayed result by using the `--explain` flag:

```Text
$ reach web-instance db-instance --explain
```

In this case, Reach will provide significantly more detail about the analysis. Specificially, the output will also show you:

- Exactly which "network points" were used in the analysis (not just the EC2 instance, but the EC2 instance's specific network interface, and the specific IP address attached to the network interface)
- All of the "factors" (relevant aspects of your configuration) Reach used to figure out what traffic is being allowed by specific properties of your resources (e.g. security group rules, instance state, etc.)

## Feature Ideas

- ~~**Same-subnet analysis:** Between two EC2 instances within the same subnet~~ (done!)
- **Same-VPC analysis:** Between two EC2 instances within the same VPC, including for EC2 instances in separate subnets
- **IP address analysis:** Between an EC2 instance and a specified IP address that may be outside of AWS entirely (enhancement idea: provide shortcuts for things like the user's own IP address, a specified hostname's resolved IP address, etc.)
- **Filtered analysis:** Specify a particular kind of network traffic to analyze (e.g. a single TCP port) and return results only for that filter
- **Other AWS resources:** Analyze other kinds of AWS resources than just EC2 instances (e.g. ELB, Lambda, VPC endpoints, etc.)
- **Peered VPC analysis**: Between resources from separate but peered VPCs
- Other things! Your ideas are welcome!

## Disclaimers

- This tool is a work in progress! Use at your own risk, and please submit issues as you encounter bugs or have feature requests.

- Because Reach gets all of its information from the AWS API, Reach makes no guarantees about network service accessibility with respect to the operating system or applications running on a host within the cloud environment. In other words, Reach can tell you if your VPC resources and EC2 instances are configured correctly, but Reach _cannot_ tell you if an OS firewall is blocking network traffic, or if an application listening on a port has crashed.
