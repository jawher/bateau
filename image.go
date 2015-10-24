package main

import (
	"fmt"

	"strings"

	"github.com/fsouza/go-dockerclient"
	"github.com/jawher/bateau/query"
)

var (
	imgFields = map[string][]query.Operator{
		"id":             {query.EQ, query.LIKE},
		"cmd":            {query.EQ, query.LIKE},
		"entrypoint":     {query.EQ, query.LIKE},
		"comment":        {query.EQ, query.LIKE},
		"author":         {query.EQ, query.LIKE},
		"arch":           {query.EQ, query.LIKE},
		"docker_version": {query.EQ, query.LIKE},

		"label.*": {query.IS, query.EQ, query.LIKE},

		"size":    {query.EQ, query.GT},
		"created": {query.EQ, query.GT},
	}
)

type DockerImage struct {
	client    *docker.Client
	apiImage  docker.APIImages
	fullImage *docker.Image
}

func wrapImage(client *docker.Client, apiImage docker.APIImages) *DockerImage {
	return &DockerImage{
		client:   client,
		apiImage: apiImage,
	}
}

var _ query.Queryable = &DockerImage{}

func (c *DockerImage) Is(field string, operator query.Operator, value string) bool {
	switch {
	case strings.HasPrefix(field, "label."):
		label := strings.TrimPrefix(field, "label.")
		labelValue, found := c.full().Config.Labels[label]
		if operator == query.IS {
			return found
		}
		return strCompare(labelValue, operator, value)
	case field == "id":
		return strCompare(c.apiImage.ID, operator, value)
	case field == "docker_version":
		return strCompare(c.full().DockerVersion, operator, value)
	case field == "comment":
		return strCompare(c.full().Comment, operator, value)
	case field == "author":
		return strCompare(c.full().Author, operator, value)
	case field == "arch":
		return strCompare(c.full().Architecture, operator, value)
	case field == "cmd":
		return sliceCompare(c.full().Config.Cmd, operator, value)
	case field == "entrypoint":
		return sliceCompare(c.full().Config.Entrypoint, operator, value)
	case field == "size":
		return sizeCompare(c.apiImage.VirtualSize, operator, value)
	case field == "created":
		return durationCompare(c.full().Created, operator, value)
	default:
		panic(fmt.Sprintf("Invalid field %s", field))
	}
}

func (c *DockerImage) full() *docker.Image {
	if c.fullImage != nil {
		return c.fullImage
	}
	daRealImage, err := c.client.InspectImage(c.apiImage.ID)
	if err != nil {
		fail("Error while retreiving image %s: %v", c.apiImage.ID, err)
	}
	c.fullImage = daRealImage
	return c.fullImage
}
