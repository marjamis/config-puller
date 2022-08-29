package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAreEnvsAvailable(t *testing.T) {
	os.Clearenv()
	filename := "ADMIN"

	os.Setenv(envPrefix+filename+schemeSuffix, "s3")
	assert.False(t, areAllEnvsAvailable(filename))

	os.Setenv(envPrefix+filename+bucketNameSuffix, "bucket")
	assert.False(t, areAllEnvsAvailable(filename))

	os.Setenv(envPrefix+filename+objectKeySuffix, "object/key")
	assert.False(t, areAllEnvsAvailable(filename))

	os.Setenv(envPrefix+filename+saveLocationSuffix, "/save/location")
	assert.False(t, areAllEnvsAvailable(filename))

	os.Setenv(envPrefix+filename+permissionsSuffix, "0777")
	assert.True(t, areAllEnvsAvailable(filename))
}

func TestGetConfigFromEnvs(t *testing.T) {
	t.Run("All envs available", func(t *testing.T) {
		os.Clearenv()
		filename := "ADMIN"

		os.Setenv(envPrefix+filename+schemeSuffix, "s3")
		os.Setenv(envPrefix+filename+bucketNameSuffix, "bucket")
		os.Setenv(envPrefix+filename+objectKeySuffix, "object/key")
		os.Setenv(envPrefix+filename+saveLocationSuffix, "/save/location")
		os.Setenv(envPrefix+filename+permissionsSuffix, "0777")

		assert.Equal(t, ConfigDetails{
			scheme:       "s3",
			bucketName:   "bucket",
			objectKey:    "object/key",
			saveLocation: "/save/location",
			permissions:  "0777",
		}, getConfigsFromEnvs()[0])
	})

	t.Run("Missing env", func(t *testing.T) {
		os.Clearenv()
		filename := "ADMIN"

		os.Setenv(envPrefix+filename+schemeSuffix, "s3")
		os.Setenv(envPrefix+filename+bucketNameSuffix, "bucket")
		os.Setenv(envPrefix+filename+saveLocationSuffix, "/save/location")
		os.Setenv(envPrefix+filename+permissionsSuffix, "0777")

		var conf []ConfigDetails
		assert.Equal(t, conf, getConfigsFromEnvs())
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
