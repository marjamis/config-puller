package main

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/localstack"
)

const (
	DEFAULT_SCHEME = "s3"
)

func TestAreEnvsAvailable(t *testing.T) {
	os.Clearenv()
	godotenv.Load("./test/envfile.txt")

	assert.False(t, areAllEnvsAvailable("MISSINGSCHEME"))
	assert.False(t, areAllEnvsAvailable("MISSINGBUCKETNAME"))
	assert.False(t, areAllEnvsAvailable("MISSINGOBJECTKEY"))
	assert.False(t, areAllEnvsAvailable("MISSINGLOCATION"))
	assert.False(t, areAllEnvsAvailable("MISSINGPERMISSIONS"))

	assert.True(t, areAllEnvsAvailable("SUCCESS1"))
	assert.True(t, areAllEnvsAvailable("SUCCESS2"))
}

func TestGetConfigFromEnvs(t *testing.T) {
	os.Clearenv()
	godotenv.Load("./test/envfile.txt")
	configs, failureCount := getConfigsFromEnvs()

	t.Run("Success count and values", func(t *testing.T) {
		assert.ElementsMatch(t, []configDetails{
			{
				scheme:       DEFAULT_SCHEME,
				bucketName:   "my-bucket",
				objectKey:    "objectkey",
				saveLocation: "/save/location/with-filename.config",
				// Conversion of 777 (octet) to base 10
				permissions: fs.FileMode(511),
			},
			{
				scheme:       DEFAULT_SCHEME,
				bucketName:   "my-bucket2",
				objectKey:    "object/key",
				saveLocation: "/save/location/with-filename2.config",
				// Conversion of 400 (octet) to base 10
				permissions: fs.FileMode(256),
			},
		}, configs)
	})

	t.Run("Failure count", func(t *testing.T) {
		assert.Equal(t, 5, failureCount)
	})
}

func TestFindAllConfigFiles(t *testing.T) {
	os.Clearenv()

	files := []string{
		"Admin1",
		"Admin1",
		"ADMIN1",
		"123ADMIN",
		"123admin",
	}
	for _, file := range files {
		os.Setenv(envPrefix+file+schemeSuffix, "testValue")
	}

	expected := []string{
		"123ADMIN",
		"123admin",
		"ADMIN1",
		"Admin1",
	}

	assert.Equal(t, expected, findAllConfigFiles())
}

func s3Client(ctx context.Context, l *localstack.LocalStackContainer) (*s3.Client, error) {
	mappedPort, err := l.MappedPort(ctx, "4566/tcp")
	if err != nil {
		return nil, err
	}

	provider, err := testcontainers.NewDockerProvider()
	if err != nil {
		return nil, err
	}
	defer provider.Close()

	host, err := provider.DaemonHost(ctx)
	if err != nil {
		return nil, err
	}

	awsCfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("us-east-1"),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("AccessKey", "SecretAccessKey", "Token")),
	)
	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.UsePathStyle = true
		o.BaseEndpoint = aws.String(fmt.Sprintf("http://%s:%d", host, mappedPort.Int()))
	})

	return client, nil
}

func TestGetFile(t *testing.T) {
	ctx := context.Background()

	bucketName := "my-bucket2"
	localstackContainer, err := localstack.Run(ctx, "localstack/localstack:2.0.0")
	assert.NoError(t, err)
	defer func() {
		assert.NoError(t, localstackContainer.Terminate(ctx))
	}()

	client, err := s3Client(ctx, localstackContainer)
	assert.NoError(t, err)

	_, err = client.CreateBucket(ctx, &s3.CreateBucketInput{
		Bucket: &bucketName,
	})
	assert.NoError(t, err)

	t.Run("Getting configuration file from S3", func(t *testing.T) {
		keyName := "object/key"
		saveLocation := "/tmp/with-filename2.config"
		contentType := "application/json"

		_, err := client.PutObject(ctx, &s3.PutObjectInput{
			Bucket:      &bucketName,
			Key:         &keyName,
			Body:        strings.NewReader("content"),
			ContentType: &contentType,
		})
		assert.NoError(t, err)

		err = getFile(configDetails{
			scheme:       DEFAULT_SCHEME,
			bucketName:   bucketName,
			objectKey:    keyName,
			saveLocation: saveLocation,
			permissions:  fs.FileMode(511),
		}, client)
		assert.NoError(t, err)
	})
}
