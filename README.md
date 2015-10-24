# bateau

[![Build Status](https://travis-ci.org/jawher/bateau.svg?branch=master)](https://travis-ci.org/jawher/bateau)
[![GoDoc](https://godoc.org/github.com/jawher/bateau?status.svg)](https://godoc.org/github.com/jawher/bateau)

Docker ps on steroids, aka rich querying language for filtering docker containers and images  

## Examples:

Find all containers created more than 2 weeks ago which are not running, except for the container named `precious`,
and delete them:

```
$ bateau 'created > 2w & !running & name!=precious' | xargs docker rm -fv
```

Find all containers with either a `org/web-app` image, or whose name container `my-app` or with a `role` label set to `web`,
 and which exited with a non-zero code, and pipe them to docker inspect.

```
$ bateau '(image=org/web-app | name~my-app | label.role=web) & exit!=0' | xargs docker inspect

```

Find all images weighing more than 300MB which were created more than 2 months ago or which were built by docker 1.5 or 1.6:
 
```
$ bateau -i 'size>300MB & (created > 2M | docker_version~1.5 | docker_version~1.6)'
```

## Motivation

The default `docker ps` command has a couple of shortcomings:

* Inconsistent handling of invalid filtering keys. For example, `docker images` will reject unknown fields, whereas `docker ps` would not.
* No handling of invalid syntax: `docker ps` will accept  `"a & b.c=42"` as a filter without complaining
* No support for other operators than `=`
* Due to the previous point, `docker ps` and co need to expose various flags (`--before`, `--since`, ...) to implement features which could have been implemented as regular filter if more operators were supported, e.g. `>`, `<`, or `before` ...
* Filters can only be combined using the `and` boolean operation. there is no way to combine them using `or`. It is not possible to retrieve in one-shot containers having one of two different labels for example.
* Filters cannot be negated. It is not possible to retrieve the containers without a specific label for example.

## Usage

```
Usage: bateau [-e] [-c|-i] QUERY

Docker ps on steroids

Arguments:
  QUERY=""     The containers filtering query

Options:
  -e, --endpoint=""       The docker socket path or TCP address
  -c, --containers=true   Filter on containers
  -i, --images=false      Filter on images
```

## Query syntax

### Conditions

Conditions can be written as just a field name, e.g.:

```
running
```

or:

```
label.arch
```

For boolean fields.
Most fields require an operator and a value.

Supported operators are:

* `=`: strict equal: `name=container1`
* `~`: like operator: case insensitive string contains: `name~tAiNe`
* `!=`: not equal: `exit!=0`
* `!~`: not like: `image!~db`
* `>`, `>=`, `<` and `<=`: greater than, greater than or equal, less than and less than or equal: `size>15MB`, `created<=1d`, ...

### Negation
You can use the `!` operator to negate an expression: `!running`, `!exit=42`

### Or/And
Expressions can be combined together using the `|` and `&` boolean operators: `running & created>1d`, `size>10MB | name~junk`
 
### Parenthesis
Expressions can be wrapped inside parenthesis to control the operator precedence: `!(running | paused)`, `image~server & (running | exit=0)` 

## Supported fields:
### Containers

|     field      |       supported operators       |                          desc                         |
| -------------- | ------------------------------- | ----------------------------------------------------- |
| `running`      | <none>                          | matches running containers                            |
| `paused`       | <none>                          | matches paused containers                             |
| `restarting`   | <none>                          | matches restarting containers                         |
| `label.<name>` | <none>                          | matches containers with a `<name>` label`             |
| `label.<name>` | `=`, `~`, `!=`, `!~`            | match against the label value                         |
| `id`           | `=`, `~`, `!=`, `!~`            | match against the container id                        |
| `name`         | `=`, `~`, `!=`, `!~`            | match against the container name                      |
| `image`        | `=`, `~`, `!=`, `!~`            | match against the container image                     |
| `cmd`          | `=`, `~`, `!=`, `!~`            | match against the container command                   |
| `entrypoint`   | `=`, `~`, `!=`, `!~`            | match against the container entrypoint                |
| `exit`         | `=`, `!=`, `>`, `>=`, `<`, `<=` | match against the container exit code                 |
| `created`      | `=`, `!=`, `>`, `>=`, `<`, `<=` | match against the container age   (since creation)    |
| `exited`       | `=`, `!=`, `>`, `>=`, `<`, `<=` | match against the duration since the container exited |

### Images

|      field       |       supported operators       |                      desc                      |
| ---------------- | ------------------------------- | ---------------------------------------------- |
| `label.<name>`   | <none>                          | matches containers with a `<name>` label`      |
| `label.<name>`   | `=`, `~`, `!=`, `!~`            | match against the label value                  |
| `id`             | `=`, `~`, `!=`, `!~`            | match against the image id                     |
| `comment`        | `=`, `~`, `!=`, `!~`            | match against the image comment                |
| `author`         | `=`, `~`, `!=`, `!~`            | match against the image author                 |
| `arch`           | `=`, `~`, `!=`, `!~`            | match against the image architecture           |
| `docker_version` | `=`, `~`, `!=`, `!~`            | match against the image docker version         |
| `cmd`            | `=`, `~`, `!=`, `!~`            | match against the image command                |
| `entrypoint`     | `=`, `~`, `!=`, `!~`            | match against the image entrypoint             |
| `size`           | `=`, `!=`, `>`, `>=`, `<`, `<=` | match against the image size                   |
| `created`        | `=`, `!=`, `>`, `>=`, `<`, `<=` | match against the image age   (since creation) |



## Value formats
### Durations
The `created` and `exited` fields accept values using the duration syntax:

```
1s
2h 40m
30w
```

A duration is a sequence of one or more number and unit pairs.
The supported units are:

* `ms` for milliseconds
* `s` for seconds
* `m` for minutes
* `h` for hours
* `d` for days (24 hours)
* `w` for weeks (7 days)
* `M` or `months` for months (30 days)
* `y` for years (365 days)

### Sizes
The `size` field accepts values using the size syntax:

```
12
400MB
1GB 250MB 14Kb
```

A size is a sequence of one or more number and unit pairs.
The supported units are:

* <none> for bytes
* `kb` or `Kb`: 1000 bytes
* `KB`: 1024 bytes
* `Mb`: 1000Kb
* `MB`: 1024KB
* `Gb`: 1000Mb
* `GB`: 1024MB

## License

This work is published under the MIT license.

Please see the `LICENSE` file for details.