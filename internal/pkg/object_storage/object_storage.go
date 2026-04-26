package objectstorage

import (
	"context"
	stdErrors "errors"
	"io"
	"manga-go/internal/pkg/config"
	"manga-go/internal/pkg/logger"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	s3Config "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go"
	pkgerrors "github.com/pkg/errors"
)

type ObjectStorage struct {
	s3            *s3.Client
	presignClient *s3.PresignClient
	bucketName    string
	endpoint      string
}

func NewObjectStorage(logger *logger.Logger, config *config.Config) *ObjectStorage {
	cfg := config.ObjectStorage
	logger.Infof("Initializing object storage with endpoint: %s, bucket: %s", cfg.Endpoint, cfg.BucketName)
	cfgS3, err := s3Config.LoadDefaultConfig(
		context.TODO(),
		s3Config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			cfg.AccessKeyID,
			cfg.SecretAccessKey,
			"",
		)),
		s3Config.WithRegion(cfg.Region),
	)

	clientS3 := s3.NewFromConfig(cfgS3, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(cfg.Endpoint)
		o.UsePathStyle = true
	})

	presignClient := s3.NewPresignClient(clientS3)
	if err != nil {
		logger.Fatalf("Failed to load AWS config: %v", err)
		panic(err)
	}

	return &ObjectStorage{
		s3:            clientS3,
		presignClient: presignClient,
		bucketName:    cfg.BucketName,
		endpoint:      cfg.Endpoint,
	}
}

func (o *ObjectStorage) GetS3Client() *s3.Client {
	return o.s3
}

func (o *ObjectStorage) CreatePresignedURL(ctx context.Context, key string) (string, error) {
	presignParams := &s3.GetObjectInput{
		Bucket: aws.String(o.bucketName),
		Key:    aws.String(key),
	}
	presignResult, err := o.presignClient.PresignGetObject(
		ctx,
		presignParams,
		s3.WithPresignExpires(15*time.Minute),
	)

	if err != nil {
		return "", pkgerrors.WithMessage(err, "create presigned url")
	}

	return presignResult.URL, nil
}

func (o *ObjectStorage) GetFile(ctx context.Context, fileName string) ([]byte, error) {
	getObjectInput := &s3.GetObjectInput{
		Bucket: aws.String(o.bucketName),
		Key:    aws.String(fileName),
	}

	resp, err := o.s3.GetObject(ctx, getObjectInput)
	if err != nil {
		return nil, pkgerrors.WithMessage(err, "get file from object storage")
	}
	defer resp.Body.Close()

	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, pkgerrors.WithMessage(err, "read file content")
	}

	return buf, nil
}

func (o *ObjectStorage) UploadFile(ctx context.Context, fileName string, body io.Reader, contentLength int64, contentType string) error {
	putObjectInput := &s3.PutObjectInput{
		Bucket:      aws.String(o.bucketName),
		Key:         aws.String(fileName),
		Body:        body,
		ContentType: aws.String(contentType),
	}

	if contentLength > 0 {
		putObjectInput.ContentLength = aws.Int64(contentLength)
	}

	_, err := o.s3.PutObject(ctx, putObjectInput)
	if err != nil {
		return pkgerrors.WithMessage(err, "upload file to object storage")
	}

	return nil
}

func (o *ObjectStorage) DeleteFile(ctx context.Context, fileName string) error {
	deleteObjectInput := &s3.DeleteObjectInput{
		Bucket: aws.String(o.bucketName),
		Key:    aws.String(fileName),
	}

	_, err := o.s3.DeleteObject(ctx, deleteObjectInput)
	if err != nil {
		return pkgerrors.WithMessage(err, "delete file from object storage")
	}

	return nil
}

func (o *ObjectStorage) IsNotFoundError(err error) bool {
	if err == nil {
		return false
	}

	var noSuchKey *types.NoSuchKey
	if stdErrors.As(err, &noSuchKey) {
		return true
	}

	var apiErr smithy.APIError
	if stdErrors.As(err, &apiErr) {
		code := strings.TrimSpace(apiErr.ErrorCode())
		switch code {
		case "NoSuchKey", "NotFound":
			return true
		}
	}

	lower := strings.ToLower(err.Error())
	return strings.Contains(lower, "nosuchkey") || strings.Contains(lower, "not found")
}
