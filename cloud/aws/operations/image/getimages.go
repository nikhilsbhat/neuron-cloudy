package awsimage

import (
	aws "github.com/nikhilsbhat/neuron-cloudy/cloud/aws/interface"
)

// GetImageInput can hold all the required details for fecthing the details of image
type capturedype GetImageInput struct {
	Kind     string   `json:"Kind"`
	ImageIds []string `json:"ImageIds"`
	GetRaw   bool     `json:"GetRaw"`
}

// GetImage will fetch the details of the images, who's Id is passed to it.
func (i *GetImageInput) GetImage(con aws.EstablishConnectionInput) ([]ImageResponse, error) {

	ec2, sesErr := con.EstablishConnection()
	if sesErr != nil {
		return nil, sesErr
	}

	// desribing image to check if image exists
	imageResult, imageErr := ec2.DescribeImages(
		&aws.DescribeComputeInput{
			ImageIds: i.ImageIds,
		},
	)
	if imageErr != nil {
		return nil, imageErr
	}

	imageResponse := make([]ImageResponse, 0)
	for _, img := range imageResult.Images {
		if i.GetRaw == true {
			imageResponse = append(imageResponse, ImageResponse{GetImageRaw: img})
		} else {
			resp := new(ImageResponse)
			resp.Name = *img.Name
			resp.ImageId = *img.ImageId
			resp.CreationDate = *img.CreationDate
			resp.State = *img.State
			resp.IsPublic = *img.Public

			snap := new(SnapshotDetails)
			snap.SnapshotId = *img.BlockDeviceMappings[0].Ebs.SnapshotId
			snap.VolumeType = *img.BlockDeviceMappings[0].Ebs.VolumeType
			snap.VolumeSize = *img.BlockDeviceMappings[0].Ebs.VolumeSize
			resp.SnapShot = *snap
			imageResponse = append(imageResponse, *resp)
		}
	}
	return imageResponse, nil
}

// This function is tuned to get the details of all images present under this account in the entered region.
func (i *GetImageInput) GetAllImage(con aws.EstablishConnectionInput) ([]ImageResponse, error) {

	ec2, seserr := con.EstablishConnection()
	if seserr != nil {
		return nil, seserr
	}

	// desribing image to check if image exists
	result, err := ec2.DescribeAllImages(
		&aws.DescribeComputeInput{},
	)

	if err != nil {
		return nil, err
	}

	imageResponse := make([]ImageResponse, 0)
	for _, img := range result.Images {
		if i.GetRaw == true {
			imageResponse = append(imageResponse, ImageResponse{GetImageRaw: img})
		} else {
			resp := new(ImageResponse)
			resp.Name = *img.Name
			resp.ImageId = *img.ImageId
			resp.CreationDate = *img.CreationDate
			resp.State = *img.State
			resp.IsPublic = *img.Public

			snap := new(SnapshotDetails)
			snap.SnapshotId = *img.BlockDeviceMappings[0].Ebs.SnapshotId
			snap.VolumeType = *img.BlockDeviceMappings[0].Ebs.VolumeType
			snap.VolumeSize = *img.BlockDeviceMappings[0].Ebs.VolumeSize
			resp.SnapShot = *snap
			imageResponse = append(imageResponse, *resp)
		}
	}
	return imageResponse, nil
}
