package storage

import (
	"bytes"
	"errors"
	"github.com/astaxie/beego"
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack"
	"github.com/rackspace/gophercloud/openstack/objectstorage/v1/objects"
	"github.com/rackspace/gophercloud/pagination"
	"mime"
	"path/filepath"
	"webseite/cache"
)

type openStackStorage struct {
}

var (
	client      *gophercloud.ServiceClient
	url         string
	existsCache *cache.TimeoutCache
)

func init() {
	if v, err := beego.AppConfig.Bool("OpenStackOn"); err == nil && v == true {
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

		// Get the public cdn URL
		url = beego.AppConfig.String("OpenStackCDNUrl")

		// Build up cache
		tempCache, err := cache.NewTimeoutCache(1600)
		if err != nil {
			panic(err)
		}

		existsCache = tempCache
	}
}

func (s *openStackStorage) Store(storeBytes []byte, filename string) (bool, error) {
	mimeType := mime.TypeByExtension(filepath.Ext(filename))
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	opts := objects.CreateOpts{
		ContentType: mimeType,
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

func (s *openStackStorage) Exists(filename string) bool {
	value, ok := existsCache.Get(filename)
	if !ok {
		result := make(chan bool, 1)

		go func() {
			opts := objects.ListOpts{
				Prefix: filename,
			}

			pager := objects.List(client, "testing", opts)
			pager.EachPage(func(page pagination.Page) (bool, error) {
				// Get a slice of containers.Container structs
				objectNames, err := objects.ExtractNames(page)

				if err != nil {
					result <- false
				}

				for _, c := range objectNames {
					if c == filename {
						result <- true
						return false, nil
					}
				}

				return true, nil
			})

			result <- false
		}()

		rval := <-result
		existsCache.Add(filename, rval)
		return rval
	} else {
		return value.(bool)
	}
}

func (s *openStackStorage) Delete(filename string) (bool, error) {
	res := objects.Delete(client, "testing", filename, nil)
	if res.Err == nil {
		return true, nil
	}

	return false, res.Err
}

func (s *openStackStorage) GetUrl(filename string) (string, error) {
	if s.Exists(filename) {
		return url + "/" + filename, nil
	}

	return "", errors.New("Not found")
}
