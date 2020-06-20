# ecsrun [![Powered By: GoReleaser](https://img.shields.io/badge/powered%20by-goreleaser-green.svg?style=flat-square)](https://github.com/goreleaser) [![Release](https://img.shields.io/github/release/masterpointio/ecsrun.svg)](https://github.com/masterpointio/ecsrun/releases/latest)

Easily run one-off tasks against an ECS Task Definition.

## Purpose

`ecsrun` is a small go CLI app to provide a config file based approach to executing one-off ECS Tasks. The ECS `RunTask` command is a pain to write out on the command line, so this tool provides an easy way to wrap any common `RunTask` executions you do in a simple yaml file.

## Install

#### Homebrew

```
brew install masterpointio/tap/ecsrun
```

#### Go Get

```
go get -u github.com/masterpointio/ecsrun
```

## Usage

### Invoking with `ecsrun.yaml`

Given you have an `ecsrun.yaml` like so:

```yaml
default: &default
  cluster: "mp-test-cluster"
  task: "mp-test-alpine"
  security-group: "sg-06c65c3206401917e"
  subnet: "subnet-0c97e16b8a52b4b86"
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

```
# Invoke the 'mp-test-alpine' task definition with the 'hello world' `CMD`
ecsrun --config default

# Invoke the 'mp-test-django' task definition with the `manage.py migrate` `CMD`
ecsrun --config migrate
```

TODO: Add more here

## Roadmap

- [x] Support basic CLI usage
- [x] Support local config file
- [x] Support `--dryrun` Flag
- [x] Add more tests
- [ ] Add a `ecsrun init` command to generate the ecsrun.yml config file.
- [ ] Support log group / stream tailing of initiated task
- [ ] Support selection of resources similar to gossm (cluster, task def, task def version, etc etc)
- [ ] Support validation of given params: cluster, definition name, revision, subnet ID, SG ID, ect.
- [ ] Support EC2 usage.
