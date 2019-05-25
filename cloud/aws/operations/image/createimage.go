package awsimage

import (
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/service/ec2"
	aws "github.com/nikhilsbhat/neuron-cloudy/cloud/aws/interface"
	common "github.com/nikhilsbhat/neuron-cloudy/cloud/aws/operations/common"
	server "github.com/nikhilsbhat/neuron-cloudy/cloud/aws/operations/server"
	err "github.com/nikhilsbhat/neuron-cloudy/error"
	//err "github.com/nikhilsbhat/neuron-cloudy/error"
)

// ImageCreateInput implements CreateImage for creation of image
type ImageCreateInput struct {
	InstanceId string
	GetRaw     bool
}

// ImageResponse contains filtered/unfiltered response received from aws.
type ImageResponse struct {
	Name            string                    `json:"Name,omitempty"`
	ImageId         string                    `json:"ImageId,omitempty"`
	ImageIds        []string                  `json:"ImageIds,omitempty"`
	State           string                    `json:"State,omitempty"`
	IsPublic        bool                      `json:"IsPublic,omitempty"`
	CreationDate    string                    `json:"CreationDate,omitempty"`
	Description     string                    `json:"Description,omitempty"`
	DefaultResponse string                    `json:"DefaultResponse,omitempty"`
	DeleteResponse  string                    `json:"ImageResponse,omitempty"`
	SnapShot        SnapshotDetails           `json:"SnapShot,omitempty"`
	CreateImageRaw  *ec2.CreateImageOutput    `json:"CreateImageRaw,omitempty"`
	GetImagesRaw    *ec2.DescribeImagesOutput `json:"GetImagesRaw,omitempty"`
	GetImageRaw     *ec2.Image                `json:"GetImageRaw,omitempty"`
}

// SnapshotDetails holds the details of snapshot captured such as type disk, size of it and etc.
type SnapshotDetails struct {
	SnapshotId string `json:"SnapshotId,omitempty"`
	VolumeType string `json:"VolumeType,omitempty"`
	VolumeSize int64  `json:"VolumeSize,omitempty"`
}

// CreateImage will capture the image of the server/vm based on the input received from ImageCreateInput.
func (img *ImageCreateInput) CreateImage(con aws.EstablishConnectionInput) (ImageResponse, error) {

	ec2, seserr := con.EstablishConnection()
	if seserr != nil {
		return ImageResponse{}, seserr
	}

	// fetching instance details as I need to pass this while taking server backup
	searchInstance := server.CommonComputeInput{InstanceIds: []string{img.InstanceId}}
	instanceResult, insterr := searchInstance.SearchInstance(con)
	if insterr != nil {
		return ImageResponse{}, insterr
	}
	if instanceResult == false {
		return ImageResponse{}, err.ServerNotFound()
	}

	getInstname := server.DescribeInstanceInput{InstanceIds: []string{img.InstanceId}}
	instanceName, instgeterr := getInstname.GetServersDetails(con)
	if instgeterr != nil {
		return ImageResponse{}, instgeterr
	}

	// Here where do stuff to take server backup
	nowTime := time.Now().Local().Format("2006-01-02 09:10:31")

	// fetching names from images so that we can name the new image uniquely
	result, deserr := ec2.DescribeAllImages(
		&aws.DescribeComputeInput{},
	)

	if deserr != nil {
		return ImageResponse{}, deserr
	}

	imagenames := make([]string, 0)
	for _, imgs := range result.Images {
		imagenames = append(imagenames, *imgs.Name)
	}

	// Getting Unique number to name image uniquely
	uqnin := common.CommonInput{SortInput: imagenames}
	uqnchr, unerr := uqnin.GetUniqueNumberFromTags()
	if unerr != nil {
		return ImageResponse{}, unerr
	}

	imageCreateResult, imgerr := ec2.CreateImage(
		&aws.ImageCreateInput{
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
	tags := common.Tag{*imageCreateResult.ImageId, "Name", instanceName[0].InstanceName + "-snapshot" + strconv.Itoa(uqnchr)}
	_, tagErr := tags.CreateTags(con)
	if tagErr != nil {
		return ImageResponse{}, tagErr
	}

	/* This will be versioning the images, now this has no much impact but once neuron is built completely this will be helpful
	tags2 := common.Tag{*imageCreateResult.ImageId, "Version", "1"}
	_, tag2_err := tags2.CreateTags()
	if tag2_err != nil {
		return nil, tag2_err
	}*/

	if img.GetRaw == true {
		return ImageResponse{CreateImageRaw: imageCreateResult}, nil
	}

	return ImageResponse{Name: instanceName[0].InstanceName + "-snapshot", ImageId: *imageCreateResult.ImageId, Description: "This image is captured by Neuron api for " + instanceName[0].InstanceName + " @ " + nowTime}, nil
}
