// Package aws will help one in creating/deleting/updating/fetching informaation of various resources in aws.
// But this is tailor made for this application if one needs
// customized result, he/she has to write one similar to this for them by calling the (master) interface.
// This package is capable of returining both custom response
// and raw from cloud depending on what you pass.
package aws

import (
	"sort"
	"strconv"

	"github.com/aws/aws-sdk-go/service/ec2"
	aws "github.com/nikhilsbhat/neuron-cloudy/cloud/aws/interface"
)

// Tag holds the info for tagging cloud resource which was or will be created.
type Tag struct {
	// Resource refers to the ID of the aws resource for which the tag has to be creted.
	Resource string
	// Name refers to the name that the tag will refered as.
	Name string
	// Value refers to the value assigned to the tag to be created.
	Value string
}

// CommonInput Implements GetAvailabilityZones, GetRegions, GetRegionFromAvail and GetUniqueNumberFromTags
type CommonInput struct {
	// AvailabilityZone is the zone name of cloud aws of which information has to be retrived.
	AvailabilityZone string
	// SortInput takes an array sorts it and returns sorted value.
	SortInput []string
	// GetRaw returns unfiltered response from the cloud if it is set to true.
	GetRaw bool
}

// CommonResponse holds the responses form the methods implemented by above structure
type CommonResponse struct {
	// Regions holds an array of the region names retrived.
	Regions       []string `json:"Regions,omitempty"`
	GetRegionsRaw *ec2.DescribeRegionsOutput
}

// GetAvailabilityZones gets the list of availability-zones in the selected region.
func (r *CommonInput) GetAvailabilityZones(con aws.EstablishConnectionInput) ([]string, error) {

	ec2, sesserr := con.EstablishConnection()
	if sesserr != nil {
		return nil, sesserr
	}

	result, err := ec2.DescribeAllAvailabilityZones(
		&aws.AwsCommonInput{},
	)
	if err != nil {
		return nil, err
	} else {
		availabilityzones := result.AvailabilityZones
		zones := make([]string, 0)
		for _, zone := range availabilityzones {
			zones = append(zones, *zone.ZoneName)
		}
		return zones, nil
	}
}

// CreateTags will create the tags to the selected resource and sends back the response.
func (t *Tag) CreateTags(con aws.EstablishConnectionInput) (string, error) {

	ec2, sesserr := con.EstablishConnection()
	if sesserr != nil {
		return "", sesserr
	}

	err := ec2.CreateTags(
		&aws.CreateTagsInput{
			Resource: t.Resource,
			Name:     t.Name,
			Value:    t.Value,
		})
	if err != nil {
		return "", err
	}
	return t.Value, nil
}

// GetRegions get the list of regions available in the selected cloud provider.
func (r *CommonInput) GetRegions(con aws.EstablishConnectionInput) (CommonResponse, error) {

	ec2, sesserr := con.EstablishConnection()
	if sesserr != nil {
		return CommonResponse{}, sesserr
	}

	result, err := ec2.GetRegions()
	if err != nil {
		return CommonResponse{}, err
	}

	if r.GetRaw == true {
		return CommonResponse{GetRegionsRaw: result}, nil
	}

	regions := make([]string, 0)
	for _, region := range result.Regions {
		regions = append(regions, *region.RegionName)
	}
	return CommonResponse{Regions: regions}, nil
}

// GetRegionFromAvail will fetch the region from the availability-zone selected.
func (r *CommonInput) GetRegionFromAvail(con aws.EstablishConnectionInput) (string, error) {

	ec2, sesserr := con.EstablishConnection()
	if sesserr != nil {
		return "", sesserr
	}

	result, err := ec2.DescribeAvailabilityZones(
		&aws.AwsCommonInput{
			AvailabilityZone: r.AvailabilityZone,
		},
	)

	if err != nil {
		return "", err
	}
	return *result.AvailabilityZones[0].RegionName, nil
}

// GetUniqueNumberFromTags will return an unique number generated from the latest created resource type.
// Newly created number will be appended to the name of the resource that will be created further.
func (r *CommonInput) GetUniqueNumberFromTags() (int, error) {

	// Sort by name, preserving original order
	sort.SliceStable(r.SortInput, func(i, j int) bool { return r.SortInput[i] < r.SortInput[j] })
	if len(r.SortInput) == 0 {
		return 0, nil
	}
	lastchr := r.SortInput[len(r.SortInput)-1]
	uniq, err := strconv.Atoi(string(lastchr[len(lastchr)-1]))
	if err != nil {
		return 0, err
	}
	return (uniq + 1), nil
}
