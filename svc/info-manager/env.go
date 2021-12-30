package main

import (
	"lab.weave.nl/nid/nid-core/pkg/environment"
	"lab.weave.nl/nid/nid-core/pkg/utilities/objectstorage"
)

type objectStorageConfig struct {
	objectstorage.ClientConfig
	Bucket string
}

// InfoManagerConfig config struct for info-manager service
type InfoManagerConfig struct {
	environment.BaseConfig
	ObjectStorage objectStorageConfig
}
