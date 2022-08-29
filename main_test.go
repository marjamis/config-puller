package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAreEnvsAvailable(t *testing.T) {
	os.Clearenv()

	os.Setenv("CONFIG_PULLER_ADMIN_SCHEME", "s3")
	os.Setenv("CONFIG_PULLER_ADMIN_BUCKETNAME", "bucket")
	assert.False(t, areEnvsAvailable("ADMIN"))

	os.Setenv("CONFIG_PULLER_ADMIN_OBJECTKEY", "object/key")
	os.Setenv("CONFIG_PULLER_ADMIN_SAVELOCATION", "/save/location")
	os.Setenv("CONFIG_PULLER_ADMIN_PERMISSIONS", "0777")
	assert.True(t, areEnvsAvailable("ADMIN"))
}

func TestGetConfigFromEnvs(t *testing.T) {
	os.Clearenv()

	os.Setenv("CONFIG_PULLER_ADMIN_SCHEME", "s3")
	os.Setenv("CONFIG_PULLER_ADMIN_BUCKETNAME", "bucket")
	os.Setenv("CONFIG_PULLER_ADMIN_OBJECTKEY", "object/key")
	os.Setenv("CONFIG_PULLER_ADMIN_SAVELOCATION", "/save/location")
	os.Setenv("CONFIG_PULLER_ADMIN_PERMISSIONS", "0777")

	assert.Equal(t, ConfigDetails{
		scheme:       "s3",
		bucketName:   "bucket",
		objectKey:    "object/key",
		saveLocation: "/save/location",
		permissions:  "0777",
	}, getConfigsFromEnvs()[0])
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
		os.Setenv("CONFIG_PULLER_"+file+"_SCHEME", "testValue")
	}
	expected := []string{
		"123ADMIN",
		"123admin",
		"ADMIN1",
		"Admin1",
	}

	assert.Equal(t, expected, findAllConfigFiles())
}
