package awsimage

import (
	aws "github.com/nikhilsbhat/neuron-cloudy/cloud/aws/interface"
	err "github.com/nikhilsbhat/neuron-cloudy/error"
)

// DeleteImageInput holds the required value terminate image and implements DeleteImage for the same.
type DeleteImageInput struct {
	ImageIds []string
}

// DeleteImage actually terminates the specified captured images.
// Make sure right image-id is passed through DeleteImageInput, because once image is deleted there is no way of bringing it back.
func (img *DeleteImageInput) DeleteImage(con aws.EstablishConnectionInput) (ImageResponse, error) {

	ec2, seserr := con.EstablishConnection()
	if seserr != nil {
		return ImageResponse{}, seserr
	}

	// desribing image to check if image exists
	searchImage := GetImageInput{ImageIds: img.ImageIds}
	imagexists, desImgErr := searchImage.IsImageAvailable(con)

	if desImgErr != nil {
		return ImageResponse{}, desImgErr
	}

	if imagexists != true {
		return ImageResponse{}, err.ImageNotFound()
	}

	imageResult, imageErr := ec2.DescribeImages(
		&aws.DescribeComputeInput{
			ImageIds: img.ImageIds,
		},
	)
	if imageErr != nil {
		return ImageResponse{}, imageErr
	}

	for _, image := range imageResult.Images {

		// Deregistering image will be done by following code
		derErr := ec2.DeregisterImage(
			&aws.DeleteComputeInput{
				ImageId: *image.ImageId,
			},
		)

		if derErr != nil {
			return ImageResponse{}, derErr
		}

		// Deletion of snapshot will addressed by following code
		snapErr := ec2.DeleteSnapshot(
			&aws.DeleteComputeInput{
				SnapshotId: *image.BlockDeviceMappings[0].Ebs.SnapshotId,
			},
		)

		if snapErr != nil {
			return ImageResponse{}, snapErr
		}
	}
	return ImageResponse{DeleteResponse: "Image is successfully deleted"}, nil
}
