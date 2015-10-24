package main

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"net/http"
	"os"

	"fmt"

	"github.com/fsouza/go-dockerclient"
	"github.com/jawher/bateau/query"
	"github.com/jawher/mow.cli"
)

func main() {
	app := cli.App("bateau", "Docker ps on steroids")

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
	client, err := docker.NewClient(endpoint)
	if err != nil {
		fail("Error while connecting to docker: %v", err)
	}

	if len(os.Getenv("DOCKER_CERT_PATH")) != 0 {
		cert, err := tls.LoadX509KeyPair(os.Getenv("DOCKER_CERT_PATH")+"/cert.pem", os.Getenv("DOCKER_CERT_PATH")+"/key.pem")
		if err != nil {
			fail("%v", err)
		}

		caCert, err := ioutil.ReadFile(os.Getenv("DOCKER_CERT_PATH") + "/ca.pem")
		if err != nil {
			fail("%v", err)
		}
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)

		tlsConfig := &tls.Config{
			Certificates: []tls.Certificate{cert},
			RootCAs:      caCertPool,
		}
		tlsConfig.BuildNameToCertificate()
		tr := &http.Transport{
			TLSClientConfig: tlsConfig,
		}
		client.HTTPClient.Transport = tr
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
