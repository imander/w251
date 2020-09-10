# Assignment 3: Face Detection

## Prerequisites
- A jetson Xavier device with key based SSH access as local user
- An AWS account with CLI credentials configured

**Software Dependencies:**
- terraform
- python 3
- make
- docker

## Building Images

Each image can be build separately via a make target. The `build-all` target will build all docker dependencies. Once images are build they should be pushed to the remote docker repository so that they can be pulled by edge devices and the cloud image processor.

A dockerhub account is required to push images. This step is only required if you plan on building your own images. Otherwise, keep the docker repository environment variable set to the default value to pull previously built images.

Once the images are built, push them to dockerhub with the following command:

```
$ make push-all
```

Both of these steps can be performed at once as follows:

```
$ make build-all push-all
```

## Multi-Architecture Builds

Building for multiple-architectures simultaneously requires more configuration that the previous method. The benefit is that all docker images can be built on a single machine. In order to perform a multi-architecture build the following prerequisites must be met:

- Docker with buildx and experimental features enabled
- A buildx environment that can build both x86 and arm images
- A dockerhub account since build and push happens in a single step

If all prerequisites are met run the following command to build all images with x86 and ARM architectures and push to dockerhub.

```
$ make buildx-all
```

When using the buildx option images won't be named based on the CPU architecture so the [docker-compose](infrastructure/ansible/roles/edge_device/templates/docker-compose.yml.j2) file for the edge devices must be edited in remove `-arm` from the image names. The following linux command can accomplish this without having to manually edit the file.

```
$ sed -i 's/-arm//' infrastructure/ansible/roles/edge_device/templates/docker-compose.yml.j2
```

## Configuration

The `.env` file should be used to configure various aspects of the infrastructure deployment. Refer to the comments in that file for details on each configuration option.

## Deploy

Ansible and Terraform are used to together to deploy the cloud infrastructure and configure the edge devices and cloud image processing server. Before running the make target to deploy you will likely need to initialize terraform.

```
$ cd infrastructure/terraform && terraform init
```

Once terraform is initialized you can deploy the facial detection infrastructure

```
$ make deploy
```

## Teardown

To stop the image detector and cloud image processor run `make destroy`. This will leave all cloud infrastructure up with the exception of the EC2 instance used as the image processor. To destroy all cloud infrastructure (including the image capture S3 bucket) run `make destroy-all`.


## MQTT Details

The MQTT QOS chosen for this service is level 0. There is no guarantee of delivery for messages with this QOS and it is chosen because the frequency of images taken would likely result in duplicate faces being detected. The MQTT topic can be configured by the user or set to a default value of the edge device host name. The image detector converts images to PNG format and appends `_png` to the MQTT topic. The cloud image processor splits the topic name on the underscore to obtain the file extension of the published message. The rest of the topic name is used to create the key of where the file is stored in S3.

For example, given the following configuration:

```
S3_BUCKET = my-bucket
MQTT_TOPIC = face-detection
```
A detected face would be published to a topic called `face-detection_png` and result in a file written to S3 at `s3://my-bucket/face-detection/<UNIX_NANO_TIMESTAMP>.png`


## Results

The results of the facial detection pipeline can be viewed in the `imander-image-capture.w251.edu` S3 bucket or at the following S3 static website.

http://imander-image-capture.w251.edu.s3-website-us-west-2.amazonaws.com/
