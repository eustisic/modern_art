package s3_client

import (
	"bytes"
	"fmt"
	"image"
	"os"

	"modern_art/utils"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// gets an image encodes it into a jpeg and uploads it to the bucket
func PostImage(artist string, image image.Image) error {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("REGION"))}))

	uploader := s3manager.NewUploader(sess)

	var buffer []byte
	buffer, err := utils.EncodeToJPEG(image)
	if err != nil {
		fmt.Println("failed to encode to jpeg")
		return err
	}

	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket:      aws.String(os.Getenv("S3_BUCKET")),
		Key:         aws.String(artist + "/1"),
		Body:        bytes.NewReader(buffer),
		ContentType: aws.String("image/jpeg"),
	})
	if err != nil {
		fmt.Println("failed to upload file")
		return err
	}
	fmt.Printf("file uploaded to, %s\n", aws.StringValue(&result.Location))
	return nil
}
