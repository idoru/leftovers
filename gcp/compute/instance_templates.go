package compute

import (
	"fmt"
	"strings"

	"github.com/genevieve/leftovers/gcp/common"
	gcpcompute "google.golang.org/api/compute/v1"
)

type instanceTemplatesClient interface {
	ListInstanceTemplates() (*gcpcompute.InstanceTemplateList, error)
	DeleteInstanceTemplate(template string) error
}

type InstanceTemplates struct {
	client instanceTemplatesClient
	logger logger
}

func NewInstanceTemplates(client instanceTemplatesClient, logger logger) InstanceTemplates {
	return InstanceTemplates{
		client: client,
		logger: logger,
	}
}

func (i InstanceTemplates) List(filter string) ([]common.Deletable, error) {
	templates, err := i.client.ListInstanceTemplates()
	if err != nil {
		return nil, fmt.Errorf("List Instance Templates: %s", err)
	}

	var resources []common.Deletable
	for _, template := range templates.Items {
		resource := NewInstanceTemplate(i.client, template.Name)

		if !strings.Contains(resource.Name(), filter) {
			continue
		}

		proceed := i.logger.Prompt(fmt.Sprintf("Are you sure you want to delete %s %s?", resource.Type(), resource.Name()))
		if !proceed {
			continue
		}

		resources = append(resources, resource)
	}

	return resources, nil
}
