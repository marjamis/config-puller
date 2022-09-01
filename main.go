package main

import (
	"context"
	"fmt"
	"os"
	"sort"
	"strings"

	"gocloud.dev/blob"
	_ "gocloud.dev/blob/s3blob"
)

const (
	envPrefix          string = "CONFIG_PULLER_"
	schemeSuffix       string = "_SCHEME"
	bucketNameSuffix   string = "_BUCKETNAME"
	objectKeySuffix    string = "_OBJECTKEY"
	saveLocationSuffix string = "_SAVELOCATION"
	permissionsSuffix  string = "_PERMISSIONS"
)

type ConfigDetails struct {
	scheme       string
	bucketName   string
	objectKey    string
	saveLocation string
	permissions  string
}

func main() {
	fmt.Println("Getting configuration files...")

	configsToGet, failures := getConfigsFromEnvs()
	for _, config := range configsToGet {
		getFile(config)
	}

	fmt.Printf("There were %d configuration files pulled successfully and %d failures.\n", len(configsToGet), failures)
}

func findAllConfigFiles() []string {
	// Use a map to filter out duplicates of filenames
	filesnames := map[string]int{}
	for _, env := range os.Environ() {
		if strings.Contains(env, envPrefix) {
			filename := strings.Split(env, "_")[2]
			filesnames[filename] = 0
		}
	}

	// Convert dictionary keys into a string array
	files := make([]string, 0, len(filesnames))
	for filename := range filesnames {
		files = append(files, filename)
	}

	sort.Strings(files)

	return files
}

func areAllEnvsAvailable(filename string) bool {
	requiredConfigurations := []string{
		schemeSuffix,
		bucketNameSuffix,
		objectKeySuffix,
		saveLocationSuffix,
		permissionsSuffix,
	}
	for _, requiredConfiguration := range requiredConfigurations {
		if os.Getenv(envPrefix+filename+requiredConfiguration) == "" {
			return false
		}
	}

	return true
}

func getConfigsFromEnvs() (configs []ConfigDetails, failures int) {
	for _, filename := range findAllConfigFiles() {
		if !areAllEnvsAvailable(filename) {
			fmt.Println("Skipping file: " + filename + " as there are missing envs.")
			failures++
			continue
		}

		configs = append(configs, ConfigDetails{
			scheme:       strings.Replace(os.Getenv(envPrefix+filename+schemeSuffix), "\"", "", -1),
			bucketName:   strings.Replace(os.Getenv(envPrefix+filename+bucketNameSuffix), "\"", "", -1),
			objectKey:    strings.Replace(os.Getenv(envPrefix+filename+objectKeySuffix), "\"", "", -1),
			saveLocation: strings.Replace(os.Getenv(envPrefix+filename+saveLocationSuffix), "\"", "", -1),
			permissions:  strings.Replace(os.Getenv(envPrefix+filename+permissionsSuffix), "\"", "", -1),
		})
	}

	return
}

func getFile(fileConfig ConfigDetails) {
	ctx := context.Background()

	bucket, err := blob.OpenBucket(ctx, fileConfig.scheme+"://"+fileConfig.bucketName)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer bucket.Close()
	blobReader, err := bucket.NewReader(ctx, fileConfig.objectKey, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(blobReader)
}
