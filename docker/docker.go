package docker

import (
	dockerapi "github.com/fsouza/go-dockerclient"
	"strings"
)

func GetVirtualHosts() ([]string, error) {
	client, err := dockerapi.NewClientFromEnv()
	if err != nil {
		return nil, err
	}
	apiContainers, err := client.ListContainers(dockerapi.ListContainersOptions{All: false})
	if err != nil {
		return nil, err
	}
	var virtualHosts []string
	for _, apiContainer := range apiContainers {
		container, err := client.InspectContainerWithOptions(dockerapi.InspectContainerOptions{ID: apiContainer.ID})
		if err != nil {
			return nil, err
		}
		exists, virtualHost := getVirtualHostForContainer(container)
		if exists {
			virtualHosts = append(virtualHosts, virtualHost)
		}
	}
	return virtualHosts, nil
}

func getVirtualHostForContainer(container *dockerapi.Container) (bool, string) {
	env := container.Config.Env
	for _, envString := range env {
		if strings.Contains(envString, "VIRTUAL_HOST=") {
			virtualHost := strings.Replace(envString, "VIRTUAL_HOST=", "", -1)
			return true, virtualHost
		}
	}
	return false, ""
}
