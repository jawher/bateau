package main

import (
	"os"

	"fmt"

	"path/filepath"

	"github.com/fsouza/go-dockerclient"
	"github.com/jawher/bateau/query"
	"github.com/jawher/mow.cli"
)

func main() {
	app := cli.App("bateau", "Docker ps on steroids")
	app.LongDesc = HELP

	endpoint := app.StringOpt("e endpoint", "", "The docker socket path or TCP address")
	_ = app.BoolOpt("c containers", true, "Filter on containers")
	images := app.BoolOpt("i images", false, "Filter on images")

	queryStr := app.StringArg("QUERY", "", "The containers filtering query")

	app.Spec = "[-e] [-c|-i] QUERY"
	app.Action = func() {
		switch {
		case *images:
			queryImages(*queryStr, *endpoint)
		default:
			queryContainers(*queryStr, *endpoint)
		}
	}
	app.Run(os.Args)
}

func queryImages(queryStr, endpoint string) {
	matcher, err := query.Parse(queryStr, imgFields)
	if err != nil {
		fail("Invalid query: %v", err)
	}
	client := NewDocker(endpoint)

	images, err := client.ListImages(docker.ListImagesOptions{All: false})
	if err != nil {
		fail("Error while listing containers: %v", err)
	}
	for _, image := range images {
		if matcher.Match(wrapImage(client, image)) {
			fmt.Printf("%s\n", image.ID)
		}
	}
}

func queryContainers(queryStr, endpoint string) {
	matcher, err := query.Parse(queryStr, conFields)
	if err != nil {
		fail("Invalid query: %v", err)
	}
	client := NewDocker(endpoint)

	containers, err := client.ListContainers(docker.ListContainersOptions{All: true})
	if err != nil {
		fail("Error while listing containers: %v", err)
	}
	for _, container := range containers {
		if matcher.Match(wrapContainer(client, container)) {
			fmt.Printf("%s\n", container.ID)
		}
	}
}

func fail(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	cli.Exit(1)
}

func NewDocker(endpoint string) *docker.Client {
	endpoint = resolveDockerEndpoint(endpoint)

	if len(os.Getenv("DOCKER_TLS_VERIFY")) != 0 {
		client, err := docker.NewTLSClient(endpoint,
			filepath.Join(os.Getenv("DOCKER_CERT_PATH"), "cert.pem"),
			filepath.Join(os.Getenv("DOCKER_CERT_PATH"), "key.pem"),
			filepath.Join(os.Getenv("DOCKER_CERT_PATH"), "ca.pem"))
		if err != nil {
			fail("Error while connecting to docker: %v", err)
		}

		return client
	}
	client, err := docker.NewClient(endpoint)
	if err != nil {
		fail("Error while connecting to docker: %v", err)
	}
	return client
}

func resolveDockerEndpoint(input string) string {
	if len(input) != 0 {
		return input
	}
	if len(os.Getenv("DOCKER_HOST")) != 0 {
		return os.Getenv("DOCKER_HOST")
	}
	return "unix:///var/run/docker.sock"
}

const HELP = `Docker ps on steroids.

Container fields:
* running, paused, restarting: booleans. e.g. 'running', 'running | paused'
* label.<label-name>: boolean to test for existence, e.g. 'label.arch' or string to test value, e.g. 'label.arch=amd64'
* id, name, image. cmd, entrypoint: string, e.g. 'entrypoint~bash'
* exit: int, e.g. 'exit=1' or 'exit>0'
* created, exited: duration, e.g. 'created>2w'

Image fields:
* id, tag, cmd, entrypoint, comment, author, arch, docker_version: string, e.g. 'id~a5fde33'
* label.<label-name>: boolean to test for existence, e.g. 'label.arch' or string to test value, e.g. 'label.arch=amd64'
* size: size, e.g. 'size>200MB'
* created: duration, e.g. 'created>2w'

Duration units:
ms->milliseconds, s->seconds, m->minutes, h->hours, d->days, w->weeks, M,months->months, y->years

Size units:
KB->1024, MB->1024MB, GB->1024MB
Kb->1000, Mb->1000Mb, Gb->1000Mb

Operators:
'=' : exact equality, '~' : case-insensitive contains, '!=' : exact inequality, '!~' : inverse of ~
'>', '>=', '<', '<=' : numeric comparison
'|' : logical or, '&' : logical and, '!' : logical not
'(', ')' : to control precedence
`
