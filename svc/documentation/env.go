package main

import (
	"lab.weave.nl/nid/nid-core/pkg/environment"
	"lab.weave.nl/nid/nid-core/pkg/utilities/objectstorage"
)

type objectStorageConfig struct {
	objectstorage.ClientConfig
	Bucket string
}

// documentationConfig implements the used environment variables
type documentationConfig struct {
	environment.BaseConfig
	GitlabAccessToken       string `envconfig:"GITLAB_ACCESS_TOKEN"`
	GitlabProjectIdentifier string `envconfig:"GITLAB_PROJECT_IDENTIFIER"`
	GitlabBaseURL           string `envconfig:"GITLAB_BASE_URL"`
	ObjectStorage           objectStorageConfig
}
