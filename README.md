# config-puller-image

A go app that pulls configuration files from a bucket and writes them locally. Primarily for use as a sidecar container where an application needs configuration files and is able to access them via a docker data volume.

## Installation

Download the docker image.

## Usage

Run the docker image directly or via an orchestrator. The generic docker run equivalent would be:

```bash
docker run -it --rm --env-file <envfile> ghcr.io/marjamis/config-puller-image:latest
```

### Configuration

For each file you would like to download from an S3 or GCP bucket specify these 5 environment variables to allow config-puller-image to get the required details:

Environment Variable | What does it do? | Valid Options
--- | --- | ---
CONFIG_PULLER_UNIQUEID_SCHEME | Defines the scheme of the bucket | gcp, s3
CONFIG_PULLER_UNIQUEID_BUCKETNAME | Name of the buket | Any valid bucketname for gcp or s3
CONFIG_PULLER_UNIQUEID_OBJECTKEY | Key of the object in the bucket | Any valid object key name and path
CONFIG_PULLER_UNIQUEID_SAVELOCATION | Specify where you would like the downloaded object to be stored | Any valid path the application can write to
CONFIG_PULLER_UNIQUEID_PERMISSIONS | The permissions of the file | Any valid Linux permissions in the octect format, such as 755, 600, etc.

### Samples

An example of both working and non-working configurations can be found in this [test file](./test/envfile.txt).

### How to

Pending...

## Support

This is very experimental but if you happen to find this project and have any issues, thoughts, or questions please don't hesitate to create a GitHub issue.

## Project status

Experimental. No promises on the future but it's a simple app so smaller updates shouldn't be too much of an issue.
