package factory

import (
	"sync"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/retailcloud"
	log "github.com/sirupsen/logrus"
)

var client *retailcloud.Client
var onceCli sync.Once

// GetRetailCloudClient is return a retailcloud client.
func GetRetailCloudClient() *retailcloud.Client {
	onceCli.Do(func() {
		config := GlobalConfig()
		var err error
		client, err = retailcloud.NewClientWithAccessKey(config.Aksk.RegionID, config.Aksk.AccessKeyID, config.Aksk.AccessKeySecret)
		if err != nil {
			log.Error(err)
		}
	})
	return client
}
