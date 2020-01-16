package aws

import (
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/service/ec2"
	neuronaws "github.com/nikhilsbhat/neuron-cloudy/cloud/aws/interface"
	err "github.com/nikhilsbhat/neuron-cloudy/errors"
)

// ImageCreateInput implements CreateImage for creation of image
type ImageCreateInput struct {
	// InstanceId refers to the ID of the aws instance of which the image has to be captured.
	InstanceId string
	// GetRaw returns unfiltered response from the cloud if it is set to true.
	GetRaw bool
}

// ImageResponse contains filtered/unfiltered response received from aws.
type ImageResponse struct {
	// Name refers to the name of the image captured or retrived.
	Name string `json:"Name,omitempty"`
	// ImageId refers to the ID of the image captured or retrived.
	ImageId string `json:"ImageId,omitempty"`
	// ImageIds refers to an array of IDs of the image captured or retrived.
	ImageIds []string `json:"ImageIds,omitempty"`
	// State defines the state of image pending/deleted etc.
	State string `json:"State,omitempty"`
	// IsPublic defines whether the image is publicly avilable.
	IsPublic bool `json:"IsPublic,omitempty"`
	// CreationDate holds the date of image creation.
	CreationDate string `json:"CreationDate,omitempty"`
	// Description describes the image captured/retrived.
	Description string `json:"Description,omitempty"`
	// DefaultResponse would be returened if function encounteres unknown circumstances.
	DefaultResponse string `json:"DefaultResponse,omitempty"`
	// DeleteResponse defines the image deletion status.
	DeleteResponse string `json:"DeleteResponse,omitempty"`
	// SnapShot return the information gathered as part of image dealt with.
	SnapShot SnapshotDetails `json:"SnapShot,omitempty"`
	// CreateImageRaw holds the unfiltered response from aws for image creation.
	CreateImageRaw *ec2.CreateImageOutput `json:"CreateImageRaw,omitempty"`
	// GetImagesRaw holds the unfiltered response from aws for retriving details of images.
	GetImagesRaw *ec2.DescribeImagesOutput `json:"GetImagesRaw,omitempty"`
	// GetImageRaw holds the unfiltered response from aws for retriving image details.
	GetImageRaw *ec2.Image `json:"GetImageRaw,omitempty"`
}

// SnapshotDetails holds the details of snapshot captured such as type disk, size of it and etc.
type SnapshotDetails struct {
	// SnapshotId refers to the ID of the snapshot created/associated to the image of which information is retrived.
	SnapshotId string `json:"SnapshotId,omitempty"`
	// VolumeType defines the volume type of the snapshot created.
	VolumeType string `json:"VolumeType,omitempty"`
	// VolumeSize is the volume size of the snapshot associated to an image
	VolumeSize int64 `json:"VolumeSize,omitempty"`
}

// CreateImage will capture the image of the server/vm based on the input received from ImageCreateInput.
func (img *ImageCreateInput) CreateImage(con neuronaws.EstablishConnectionInput) (ImageResponse, error) {

	ec2, seserr := con.EstablishConnection()
	if seserr != nil {
		return ImageResponse{}, seserr
	}

	// fetching instance details as I need to pass this while taking server backup
	searchInstance := new(CommonComputeInput)
	searchInstance.InstanceIds = []string{img.InstanceId}
	instanceResult, insterr := searchInstance.SearchInstance(con)
	if insterr != nil {
		return ImageResponse{}, insterr
	}
	if instanceResult == false {
		return ImageResponse{}, err.ServerNotFound()
	}

	getInstname := new(DescribeInstanceInput)
	getInstname.InstanceIds = []string{img.InstanceId}
	instanceName, instgeterr := getInstname.GetServersDetails(con)
	if instgeterr != nil {
		return ImageResponse{}, instgeterr
	}

	// Here where do stuff to take server backup
	nowTime := time.Now().Local().Format("2006-01-02 09:10:31")

	// fetching names from images so that we can name the new image uniquely
	result, deserr := ec2.DescribeAllImages(
		&neuronaws.DescribeComputeInput{},
	)

	if deserr != nil {
		return ImageResponse{}, deserr
	}

	imagenames := make([]string, 0)
	for _, imgs := range result.Images {
		imagenames = append(imagenames, *imgs.Name)
	}

	// Getting Unique number to name image uniquely
	uqnin := CommonInput{SortInput: imagenames}
	uqnchr, unerr := uqnin.GetUniqueNumberFromTags()
	if unerr != nil {
		return ImageResponse{}, unerr
	}

	imageCreateResult, imgerr := ec2.CreateImage(
		&neuronaws.ImageCreateInput{
			Description: "This image is captured by neuron api for " + instanceName[0].InstanceName + " @ " + nowTime,
			InstanceId:  img.InstanceId,
			ServerName:  instanceName[0].InstanceName + "-snapshot-" + strconv.Itoa(uqnchr),
		},
	)

	// handling the error if it throws while subnet is under creation process
	if imgerr != nil {
		return ImageResponse{}, imgerr
	}

	// This will take care of creation of primary tags to the image
	tags := new(Tag)
	tags.Resource = *imageCreateResult.ImageId
	tags.Name = "Name"
	tags.Value = instanceName[0].InstanceName + "-snapshot" + strconv.Itoa(uqnchr)
	_, tagErr := tags.CreateTags(con)
	if tagErr != nil {
		return ImageResponse{}, tagErr
	}

	if img.GetRaw == true {
		return ImageResponse{CreateImageRaw: imageCreateResult}, nil
	}

	return ImageResponse{Name: instanceName[0].InstanceName + "-snapshot", ImageId: *imageCreateResult.ImageId, Description: "This image is captured by Neuron api for " + instanceName[0].InstanceName + " @ " + nowTime}, nil
}
