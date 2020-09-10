ARCH = $(shell uname -m)
BUILDX = docker buildx build --push --platform linux/amd64,linux/arm64

include .env

# DOCKER_ARCH is used to build images with a suffix -arm if and arm environment is used
# the preferred way to build would be to use buildx to keep the image name consistent
# and offload the architecture detection to docker
ifeq ($(ARCH),aarch64)
	DOCKER_ARCH = -arm
endif

export DOCKER_CLI_EXPERIMENTAL=enabled
export GO111MODULE=on
export CAMERA_INDEX=$(DETECT_CAMERA_INDEX)
export WAIT_INTERVAL=$(DETECT_WAIT_INTERVAL)
export DOCKER_REPO=$(DOCKERHUB_REPO)
export MQTT_TOPIC=$(MQTT_PUBLISH_TOPIC)
export TF_VAR_region=$(AWS_REGION)
export TF_VAR_instance_type=$(AWS_INSTANCE_TYPE)
export TF_VAR_key_name=$(AWS_SSH_KEY_NAME)
export TF_VAR_bucket_name=$(AWS_S3_BUCKET_NAME)

virtualenv:
	python3 -mvenv virtualenv && pip install -r requirements.txt

.PHONY: build-detector
build-detector:
	cd detector && docker build -t $(DOCKER_REPO)/face-detector$(DOCKER_ARCH) .

.PHONY: buildx-detector
buildx-detector:
	cd detector && $(BUILDX) -t imander/face-detector .

.PHONY: build-processor
build-processor:
	cd processor && docker build -t $(DOCKER_REPO)/image-processor$(DOCKER_ARCH) .

.PHONY: buildx-processor
buildx-processor:
	cd processor && $(BUILDX) -t imander/image-processor .

.PHONY: build-forwarder
build-forwarder:
	cd mqtt_forwarder && docker build -t $(DOCKER_REPO)/mqtt-forwarder$(DOCKER_ARCH) .

.PHONY: buildx-forwarder
buildx-forwarder:
	cd mqtt_forwarder && $(BUILDX) -t imander/mqtt-forwarder .

.PHONY: build-broker
build-broker:
	cd mqtt_broker && docker build -t $(DOCKER_REPO)/mqtt-broker$(DOCKER_ARCH) .

.PHONY: buildx-broker
buildx-broker:
	cd mqtt_broker && $(BUILDX) -t imander/mqtt-broker .

.PHONY: push-detector
push-detector:
	docker push $(DOCKER_REPO)/face-detector$(DOCKER_ARCH)

.PHONY: push-processor
push-processor:
	docker push $(DOCKER_REPO)/image-processor$(DOCKER_ARCH)

.PHONY: push-forwarder
push-forwarder:
	docker push $(DOCKER_REPO)/mqtt-forwarder$(DOCKER_ARCH)

.PHONY: push-broker
push-broker:
	docker push $(DOCKER_REPO)/mqtt-broker$(DOCKER_ARCH)

.PHONY: build-all
build-all: build-detector build-processor build-forwarder build-broker

.PHONY: buildx-all
buildx-all: buildx-detector buildx-processor buildx-forwarder buildx-broker

.PHONY: push-all
push-all: push-detector push-processor push-forwarder push-broker

.PHONY: run
run:
	docker run \
		--net=host \
		--privileged \
		-v "$(PWD)/data:/tmp/data:rw" \
		-e CAMERA_INDEX=$(CAMERA_INDEX) \
		-e WAIT_INTERVAL=$(WAIT_INTERVAL) \
		-e DISPLAY=$(DISPLAY) \
		-it --rm $(DOCKER_REPO)/face-detector$(DOCKER_ARCH)

.PHONY: plan
plan:
	cd infrastructure/terraform && terraform plan

.PHONY: cloud-up
cloud-up:
	cd infrastructure/terraform && terraform apply -auto-approve

.PHONY: cloud-down
cloud-down:
	cd infrastructure/terraform && \
		terraform destroy -auto-approve --target=aws_instance.w251_image_server

.PHONY: cloud-down
cloud-down-all:
	cd infrastructure/terraform && terraform destroy -auto-approve

.PHONY: ansible-inventory
ansible-inventory:
	sed -e "s/IMAGE_SERVER/$$(jq -r '.outputs.image_server_public_ip.value' < infrastructure/terraform/terraform.tfstate)/" \
		-e "s/EDGE_SERVER/$(EDGE_SERVER)/" infrastructure/ansible/inventory.tmpl > infrastructure/ansible/inventory

.PHONY: config-up
config-up: virtualenv ansible-inventory
	source virtualenv/bin/activate && \
	cd infrastructure/ansible && \
	ansible-playbook deploy.yml --tags "start" -i inventory

.PHONY: config-down
config-down: ansible-inventory
	source virtualenv/bin/activate && \
	cd infrastructure/ansible && \
	ansible-playbook deploy.yml --tags "stop" -i inventory

.PHONY: deploy
deploy: cloud-up config-up

.PHONY: destroy
destroy: config-down cloud-down

.PHONY: destroy-all
destroy-all: config-down cloud-down-all
