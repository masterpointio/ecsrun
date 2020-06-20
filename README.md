# ecsrun

[![Go Report Card](https://goreportcard.com/badge/github.com/masterpointio/ecsrun)](https://goreportcard.com/report/github.com/masterpointio/ecsrun)
[![GitHub Workflow Status](https://img.shields.io/github/workflow/status/masterpointio/ecsrun/Go%20Build%20%26%20Test)](https://github.com/masterpointio/ecsrun/actions?query=workflow%3A%22Go+Build+%26+Test%22)
[![Powered By: GoReleaser](https://img.shields.io/badge/powered%20by-goreleaser-green.svg?style=flat-square)](https://github.com/goreleaser)
[![Release](https://img.shields.io/github/release/masterpointio/ecsrun.svg)](https://github.com/masterpointio/ecsrun/releases/latest)

Easily run one-off tasks against an ECS Task Definition 🐳

## Purpose

`ecsrun` is a small go CLI app to provide a config file based approach to executing one-off ECS Tasks. The ECS `RunTask` command is a pain to write out on the command line, so this tool provides an easy way to wrap any common `RunTask` executions you do in a simple yaml file.

## Install

#### From Homebrew

```
brew install masterpointio/tap/ecsrun
```

#### From Go Get

```
go get -u github.com/masterpointio/ecsrun
```

## Usage

#### Invoking with `ecsrun.yaml` (easiest)

Given you have an `ecsrun.yaml` like so:

```yaml
default: &default
  cluster: mp-test-cluster
  task: mp-test-alpine
  security-group: sg-06c65c3206401917e
  subnet: subnet-0c97e16b8a52b4b86
  public: false
  cmd:
    - bash
    - -c
    - echo
    - "hello world"

migrate:
  <<: *default
  task: mp-test-django
  cmd:
    - python
    - ./manage.py
    - migrate
```

You can invoke two easy commands to spin up a one-off task:

```bash
# Invoke the 'mp-test-alpine' task definition with the 'hello world' `CMD`
ecsrun

# Invoke the 'mp-test-django' task definition with the `manage.py migrate` `CMD`
ecsrun --config migrate
```

#### From Command Line

`ecsrun` supports all of the config options via CLI arguments as well:

```bash
ecsrun --cluster mp-example-task-runner \
       --subnet subnet-0c97e16b8a52b4b86 \
       --security-group sg-06c65c3206401917e \
       --cmd "bash,-c,echo,\"Hello world\"" \
       --region us-west-2 \
       --public \
       --verbose
```

You can use this in combination with a configuration file to only override certain arguments:

```bash
ecsrun --config migrate
       --subnet ${DIFFERENT_SUBNET} \
       --public
```

#### From Environment Variables

You can also pass configuration to `ecsrun` via environment variables:

```bash
export AWS_PROFILE="mp-gowiem"
export AWS_ACCESS_KEY_ID="123"
export AWS_SECRET_ACCESS_KEY="SECRET123"
export ECSRUN_CMD="bash,-c,echo,\"Hello world\""
export ECSRUN_CLUSTER="mp-testing-cluster"
export ECSRUN_TASK="mp-testing-task"
export ECSRUN_SECURITY_GROUP="sg-06c65c3206401917e"
export ECSRUN_SUBNET="subnet-0c97e16b8a52b4b86"
export ECSRUN_VERBOSE="true"

# Invoke a dry run to check the resulting `RunTaskInput` configuration
ecsrun --dry-run
```

#### Initialize an empty `ecsrun.yaml`

Don't have an `ecsrun.yaml` file yet? Initialize the scaffold of one in your current directory:

```
ecsrun init
```

#### More

Be sure to check out `ecsrun help` for more info and full configuration options.

## Inspiration

I wrote the included `run_command` bash script for a client project as admin tasks were quite common on the project (migrations, django `manage.py` jobs, debugging, etc). This script was pretty ugly (what bash script isn't honestly), but it got the job done. I wanted to build a new project in golang to try out the langauge, and converting `run_command` to something with a bit more grace seemed like a fun project. `ecsrun` is the result!

## Roadmap

- [x] Support basic CLI usage
- [x] Support local config file
- [x] Support `--dryrun` Flag
- [x] Add more tests
- [x] Add a `ecsrun init` command to generate the ecsrun.yml config file.
- [ ] Support log group / stream tailing of initiated task
- [ ] Support selection of resources similar to gossm (cluster, task def, task def revision, etc etc)
- [ ] Support validation of given params: cluster, definition name, revision, subnet ID, SG ID, ect.
- [ ] Support EC2 usage.
