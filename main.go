package main

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"gocloud.dev/blob/s3blob"
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

type configDetails struct {
	scheme       string
	bucketName   string
	objectKey    string
	saveLocation string
	permissions  fs.FileMode
}

func main() {
	fmt.Println("Creating S3 client...")
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic(err)
	}

	client := s3.NewFromConfig(cfg)

	fmt.Println("Getting configuration files...")
	configsToGet, failures := getConfigsFromEnvs()
	for _, config := range configsToGet {
		err := getFile(config, client)
		if err != nil {
			fmt.Println(fmt.Errorf("Failure to get config file \"%s\" due to \"%s\"", fmt.Sprintf("%s/%s", config.bucketName, config.objectKey), err))
		}
	}

	fmt.Printf("There were %d configuration files pulled successfully and %d failures.\n", len(configsToGet), failures)
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

func getConfigsFromEnvs() (configs []configDetails, failures int) {
	for _, filename := range findAllConfigFiles() {
		if !areAllEnvsAvailable(filename) {
			fmt.Println("Skipping file: " + filename + " as there are missing envs.")
			failures++
			continue
		}

		// Converting the permissions octet into the fs.FileMode 32 bit integer. Basically a translation between the two formats but has the same resultant permissions
		perms, err := strconv.ParseUint(strings.Replace(os.Getenv(envPrefix+filename+permissionsSuffix), "\"", "", -1), 8, 32)
		if err != nil {
			fmt.Println("Skipping file: " + filename + " as couldn't convert permissions to fileMode.")
			failures++
			continue
		}

		configs = append(configs, configDetails{
			scheme:       strings.Replace(os.Getenv(envPrefix+filename+schemeSuffix), "\"", "", -1),
			bucketName:   strings.Replace(os.Getenv(envPrefix+filename+bucketNameSuffix), "\"", "", -1),
			objectKey:    strings.Replace(os.Getenv(envPrefix+filename+objectKeySuffix), "\"", "", -1),
			saveLocation: strings.Replace(os.Getenv(envPrefix+filename+saveLocationSuffix), "\"", "", -1),
			permissions:  fs.FileMode(perms),
		})
	}

	return
}

func findAllConfigFiles() []string {
	// Use a map to filter out duplicates of filenames
	filenames := map[string]int{}
	for _, env := range os.Environ() {
		if strings.Contains(env, envPrefix) {
			filename := strings.Split(env, "_")[2]
			filenames[filename] = 0
		}
	}

	// Convert dictionary keys into a string array
	files := make([]string, 0, len(filenames))
	for filename := range filenames {
		files = append(files, filename)
	}

	sort.Strings(files)

	return files
}

func getFile(fileConfig configDetails, client *s3.Client) (err error) {
	ctx := context.Background()

	bucket, err := s3blob.OpenBucketV2(ctx, client, fileConfig.bucketName, nil)
	if err != nil {
		return err
	}
	defer bucket.Close()

	blobReader, err := bucket.NewReader(ctx, fileConfig.objectKey, nil)
	if err != nil {
		return err
	}

	buf := new(strings.Builder)
	_, err = io.Copy(buf, blobReader)
	if err != nil {
		return err
	}

	fmt.Println("Content-Type:", blobReader.ContentType())
	err = os.WriteFile(fileConfig.saveLocation, []byte(buf.String()), fileConfig.permissions)
	if err != nil {
		return err
	}

	return nil
}
