package storage

import (
	"bytes"
	"errors"
	"github.com/astaxie/beego"
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack"
	"github.com/rackspace/gophercloud/openstack/objectstorage/v1/objects"
)

type openStackStorage struct {
}

var client *gophercloud.ServiceClient

func init() {
	// Build up Auth Informations
	opts := gophercloud.AuthOptions{
		IdentityEndpoint: beego.AppConfig.String("OpenStackEndPoint"),
		Username:         beego.AppConfig.String("OpenStackUser"),
		Password:         beego.AppConfig.String("OpenStackPass"),
		TenantID:         beego.AppConfig.String("OpenStackTenantId"),
	}

	// Check if we can auth
	provider, err := openstack.AuthenticatedClient(opts)
	if err != nil {
		panic(err)
	}

	// Create new ObjectStorage Client
	tempClient, err := openstack.NewObjectStorageV1(provider, gophercloud.EndpointOpts{
		Region: beego.AppConfig.String("OpenStackRegion"),
	})

	// Check if we made an error
	if err != nil {
		panic(err)
	}

	// Store the client
	client = tempClient
}

func (s *openStackStorage) Store(storeBytes []byte, filename string) (bool, error) {
	opts := objects.CreateOpts{
		ContentType: "application/octet-stream",
	}

	res := objects.Create(client, "testing", filename, bytes.NewReader(storeBytes), opts)

	headers, err := res.ExtractHeader()
	if err != nil {
		return false, err
	}

	if headers.Get("X-Trans-Id") != "" {
		return true, nil
	}

	return false, errors.New("No Transaction Id")
}
