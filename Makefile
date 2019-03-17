# Makefile for Operator

CURRENT_DIR=$(shell \pwd)
DOCKER_BIN = $(shell \which docker)
KUBE_BIN = $(shell \which kubectl)
SDK_BIN = $(shell \which operator-sdk)
OPERATOR_IMAGE_NAME = 'quay.io/drlovewizard/example-operator'
OPERATOR_IMAGE_VER = 'v0.0.1'
OPERATOR_IMAGE = $(OPERATOR_IMAGE_NAME):$(OPERATOR_IMAGE_VER)

.PHONY: all clean build deploy

all: clean build push deploy

env-check:
ifndef QUAY_USERNAME
        $(error QUAY_USERNAME is not defined)
endif
ifndef QUAY_PASSWORD
        $(error QUAY_USERNAME is not defined)
endif

build:
	@$(SDK_BIN) build $(OPERATOR_IMAGE)

clean:
	-@$(KUBE_BIN) delete -f deploy/service_account.yaml
	-@$(KUBE_BIN) delete -f deploy/role.yaml
	-@$(KUBE_BIN) delete -f deploy/role_binding.yaml
	-@$(KUBE_BIN) delete -f deploy/crds/pgcluster_v1alpha1_pgcluster_crd.yaml
	-@$(KUBE_BIN) delete -f deploy/operator.yaml
	-@$(KUBE_BIN) delete -f deploy/crds/pgcluster_v1alpha1_pgcluster_cr.yaml

deploy: clean
	@$(KUBE_BIN) create -f deploy/service_account.yaml
	@$(KUBE_BIN) create -f deploy/role.yaml
	@$(KUBE_BIN) create -f deploy/role_binding.yaml
	@$(KUBE_BIN) create -f deploy/crds/pgcluster_v1alpha1_pgcluster_crd.yaml
	@$(KUBE_BIN) create -f deploy/operator.yaml
	@$(KUBE_BIN) create -f deploy/crds/pgcluster_v1alpha1_pgcluster_cr.yaml

quay-login: env-check
	@$(DOCKER_BIN) login quay.io --username=$(QUAY_USERNAME) --password=$(QUAY_PASSWORD)

push: quay-login
	@$(DOCKER_BIN) push $(OPERATOR_IMAGE)
