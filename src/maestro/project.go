package main

import (
	"log"

	"github.com/docker/libcompose/config"
	"github.com/docker/libcompose/docker/auth"
	"github.com/docker/libcompose/docker/client"
	"github.com/docker/libcompose/docker/ctx"
	"github.com/docker/libcompose/docker/network"
	"github.com/docker/libcompose/docker/service"
	"github.com/docker/libcompose/docker/volume"
	"github.com/docker/libcompose/project"
)

// NewProject creates a Project with the specified context.
func NewProject(context *ctx.Context, parseOptions *config.ParseOptions) (*project.Project, error) {
	if context.AuthLookup == nil {
		context.AuthLookup = auth.NewConfigLookup(context.ConfigFile)
	}

	if context.ServiceFactory == nil {
		context.ServiceFactory = service.NewFactory(context)
	}

	if context.ClientFactory == nil {
		factory, err := client.NewDefaultFactory(client.Options{})
		if err != nil {
			return nil, err
		}
		context.ClientFactory = factory
	}

	if context.NetworksFactory == nil {
		networksFactory := &network.DockerFactory{
			ClientFactory: context.ClientFactory,
		}
		context.NetworksFactory = networksFactory
	}

	if context.VolumesFactory == nil {
		volumesFactory := &volume.DockerFactory{
			ClientFactory: context.ClientFactory,
		}
		context.VolumesFactory = volumesFactory
	}

	p := project.NewProject(&context.Context, nil, parseOptions)

	err := p.Parse()
	if err != nil {
		return nil, err
	}

	if err = context.LookupConfig(); err != nil {
		log.Printf("Failed to open project %s: %v", p.Name, err)
		return nil, err
	}

	return p, err
}
