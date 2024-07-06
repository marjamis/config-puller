package main

import (
	"io/fs"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
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
				scheme:       "s3",
				bucketName:   "my-bucket",
				objectKey:    "objectkey",
				saveLocation: "/save/location/with-filename.config",
				// Conversion of 777 (octet) to base 10
				permissions: fs.FileMode(511),
			},
			{
				scheme:       "s3",
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

func TestFindAllFiles(t *testing.T) {
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
