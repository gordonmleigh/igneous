package main

import (

  "github.com/docker/go-plugins-helpers/volume"
  "github.com/gophercloud/gophercloud"
  "github.com/gophercloud/gophercloud/openstack"
  "github.com/gophercloud/gophercloud/openstack/blockstorage/v1/volumes"
  "github.com/gophercloud/gophercloud/pagination"
  log "github.com/Sirupsen/logrus"
)


type igneousDriver struct {
  fsRoot      string
  client      *gophercloud.ProviderClient
}


func newIgneousDriver(fsRoot, authOptions) {
  log.WithFields(getAuthFields(authOptions)).Info("Creating new driver instance")
  client, err := openstack.AuthenticatedClient(authOptions)

  if err != nil {
    log.Fatalln("Can't create openstack client: ", err)
  }

  d := igneousDriver{
    fsRoot:   fsRoot,
    client:   client,
  }

  return d
}


func (d igneousDriver) Create(r volume.Request) volume.Response {
  log.WithFields(log.Fields{ "Name": r.Name, "Options": r.Options }).Info("REQUEST: Create volume")

  // check that size option has been provided
  size, exist := r.Options["size"]

  if !exist {
    log.Error("must specify size to create")
    return volume.Response{ Err: "must specify size to create"}
  }

  opts := volumes.CreateOpts{
    Name: r.Name,
    Size: size,
  }
  // create the volume on OpenStack Cinder
  _, err := volumes.Create(d.client, opts).Extract()

  if err != nil {
    log.Error("Error creating volume: ", err)
    return volume.Response{ Err: err.Error() }
  }

  return volume.Response{}
}


func (d igneousDriver) Get(r volume.Request) volume.Response {
  log.WithFields(log.Fields{ "Name": r.Name }).Info("REQUEST: Get volume")

}


func (d igneousDriver) mountpoint(id string) string {
	return filepath.Join(d.fsRoot, id)
}


func (d igneousDriver) getVolumeId(name string) (string, error) {
  log.WithFields(log.Fields{ "Name": name }).Debug("Looking up volume on cinder");
  opts := volumes.ListOpts{ Name: name }
  pager := volumes.List(d.client, opts)
  var volume volume.Volume

  err := pager.EachPage(func(page pagination.Page) (bool, error) {
    vols, err := volumes.ExtractVolumes(page)
  });
}


func getAuthFields(options *gophercloud.AuthOptions) *log.Fields {
  fields := log.Fields{
    "IdentityEndpoint": options.IdentityEndpoint
  }

  if len(options.Username) > 0 {
    fields["Username"] = options.Username
  }
  if len(options.UserID) > 0 {
    fields["UserID"] = options.Username
  }
  if len(options.Password) > 0 {
    fields["Password"] = "****"
  }
  if len(options.TenantID) > 0 {
    fields["TenantID"] = options.TenantID
  }
  if len(options.TenantName) > 0 {
    fields["TenantName"] = options.TenantName
  }

  return &fields
}
