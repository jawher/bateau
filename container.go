package main

import (
	"fmt"

	"strings"

	"github.com/fsouza/go-dockerclient"
	"github.com/jawher/bateau/query"
)

var (
	conFields = map[string][]query.Operator{
		"running":    {query.IS},
		"paused":     {query.IS},
		"restarting": {query.IS},

		"label.*": {query.IS, query.EQ, query.LIKE},

		"id":         {query.EQ, query.LIKE},
		"name":       {query.EQ, query.LIKE},
		"image":      {query.EQ, query.LIKE},
		"cmd":        {query.EQ, query.LIKE},
		"entrypoint": {query.EQ, query.LIKE},

		"exit":    {query.EQ, query.GT},
		"created": {query.EQ, query.GT},
		"exited":  {query.EQ, query.GT}}
)

type DockerContainer struct {
	client        *docker.Client
	apiContainer  docker.APIContainers
	fullContainer *docker.Container
}

func wrapContainer(client *docker.Client, apiContainer docker.APIContainers) *DockerContainer {
	return &DockerContainer{
		client:       client,
		apiContainer: apiContainer,
	}
}

var _ query.Queryable = &DockerContainer{}

func (c *DockerContainer) Is(field string, operator query.Operator, value string) bool {
	switch {
	case field == "running":
		return c.full().State.Running
	case field == "paused":
		return c.full().State.Paused
	case field == "restarting":
		return c.full().State.Restarting
	case strings.HasPrefix(field, "label."):
		label := strings.TrimPrefix(field, "label.")
		labelValue, found := c.full().Config.Labels[label]
		if operator == query.IS {
			return found
		}
		return strCompare(labelValue, operator, value)
	case field == "id":
		return strCompare(c.apiContainer.ID, operator, value)
	case field == "name":
		return strCompare(strings.TrimPrefix(c.full().Name, "/"), operator, value)
	case field == "image":
		return strCompare(c.apiContainer.Image, operator, value)
	case field == "exit":
		code := c.full().State.ExitCode
		return code != -1 && intCompare(code, operator, value)
	case field == "cmd":
		return sliceCompare(c.full().Config.Cmd, operator, value)
	case field == "entrypoint":
		return sliceCompare(c.full().Config.Entrypoint, operator, value)
	case field == "created":
		return durationCompare(c.full().Created, operator, value)
	case field == "exited":
		finishedAt := c.full().State.FinishedAt
		return !finishedAt.IsZero() && durationCompare(finishedAt, operator, value)
	default:
		panic(fmt.Sprintf("Invalid field %s", field))
	}
}

func (c *DockerContainer) full() *docker.Container {
	if c.fullContainer != nil {
		return c.fullContainer
	}
	daRealContainer, err := c.client.InspectContainer(c.apiContainer.ID)
	if err != nil {
		fail("Error while retreiving container %s: %v", c.apiContainer.ID, err)
	}
	c.fullContainer = daRealContainer
	return c.fullContainer
}
