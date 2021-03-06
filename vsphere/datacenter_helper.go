package vsphere

import (
	"fmt"

	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/types"
	"golang.org/x/net/context"
)

// getDatacenter gets the higher-level datacenter object for the datacenter
// name supplied by dc.
//
// The default datacenter is denoted by using an empty string. When working
// with ESXi directly, the default datacenter is always selected.
func getDatacenter(c *govmomi.Client, dc string) (*object.Datacenter, error) {
	finder := find.NewFinder(c.Client, true)
	t := c.ServiceContent.About.ApiType
	switch t {
	case "HostAgent":
		return finder.DefaultDatacenter(context.TODO())
	case "VirtualCenter":
		if dc != "" {
			return finder.Datacenter(context.TODO(), dc)
		}
		return finder.DefaultDatacenter(context.TODO())
	}
	return nil, fmt.Errorf("unsupported ApiType: %s", t)
}

// datacenterFromID locates a Datacenter by its managed object reference ID.
func datacenterFromID(client *govmomi.Client, id string) (*object.Datacenter, error) {
	finder := find.NewFinder(client.Client, false)

	ref := types.ManagedObjectReference{
		Type:  "Datacenter",
		Value: id,
	}

	ctx, cancel := context.WithTimeout(context.Background(), defaultAPITimeout)
	defer cancel()
	ds, err := finder.ObjectReference(ctx, ref)
	if err != nil {
		return nil, fmt.Errorf("could not find datacenter with id: %s: %s", id, err)
	}
	return ds.(*object.Datacenter), nil
}
