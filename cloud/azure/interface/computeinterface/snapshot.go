package azurecompute

import (
	"fmt"
	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2018-04-01/compute"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/to"
)

func getSnapshotsClient() compute.SnapshotsClient {
	snapshotsClient := compute.NewSnapshotsClient(subscription)
	snapshotsClient.Authorizer = autorest.NewBearerAuthorizer(token)
	return snapshotsClient
}

type SnapshotIn struct {
	ResourceGroup string
	SnapshotName  string `json:"snapshotname,omitempty"`
	SourceImageID string `json:"sourceimageid,omitempty"`
	Location      string `json:"location,omitempty"`
}

func (s SnapshotIn) CreateSnapshot() (ss compute.Snapshot, err error) {

	snapshotsClient := getSnapshotsClient()
	future, err := snapshotsClient.CreateOrUpdate(
		ctx,
		s.ResourceGroup,
		s.SnapshotName,
		compute.Snapshot{
			Location: to.StringPtr(s.Location),
			DiskProperties: &compute.DiskProperties{
				OsType: compute.Linux,
				CreationData: &compute.CreationData{
					CreateOption:     compute.Copy,
					SourceResourceID: to.StringPtr(s.SourceImageID),
				},
			},
		},
	)

	if err != nil {
		return ss, fmt.Errorf("cannot create snapshot: %v", err)
	}

	err = future.WaitForCompletion(ctx, snapshotsClient.Client)
	if err != nil {
		return ss, fmt.Errorf("cannot get the snapshot create or update future response: %v", err)
	}

	return future.Result(snapshotsClient)
}

func (s SnapshotIn) DeleteSnapshot() (ar autorest.Response, err error) {

	snapshotsClient := getSnapshotsClient()
	future, err := snapshotsClient.Delete(
		ctx,
		s.ResourceGroup,
		s.SnapshotName,
	)

	if err != nil {
		return ar, fmt.Errorf("cannot delete snapshot: %v", err)
	}

	err = future.WaitForCompletion(ctx, snapshotsClient.Client)
	if err != nil {
		return ar, fmt.Errorf("cannot get the snapshot delete future response: %v", err)
	}

	return future.Result(snapshotsClient)
}

func (s SnapshotIn) GetSnapshot() (ss compute.Snapshot, err error) {

	snapshotsClient := getSnapshotsClient()
	future, err := snapshotsClient.Get(
		ctx,
		s.ResourceGroup,
		s.SnapshotName,
	)

	if err != nil {
		return ss, fmt.Errorf("cannot get snapshot: %v", err)
	}

	return future, err
}
