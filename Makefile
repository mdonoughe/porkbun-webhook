.DEFAULT_GOAL:=help

# VARIABLES
GO ?= $(shell which go)
OS ?= $(shell $(GO) env GOOS)
ARCH ?= $(shell $(GO) env GOARCH)

KUBE_VERSION=1.26.0
export TEST_ASSET_ETCD=_test/kubebuilder/etcd
export TEST_ASSET_KUBE_APISERVER=_test/kubebuilder/kube-apiserver
export TEST_ASSET_KUBECTL=_test/kubebuilder/kubectl

REGISTRY = ghcr.io
IMAGE_NAME = lost-woods/cert-manager-porkbun-webhook
IMAGE_TAG  = 1.0.3


##@ Help
.PHONY: help
help: ## Display this help screen
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Test
.PHONY: test
test: ## Usage: TEST_DOMAIN_NAME=<your-domain> TEST_API_KEY=$(echo -n '<api-key>' | base64) TEST_SECRET_KEY=$(echo -n '<secret-key>' | base64) make test
	curl -fsSL https://go.kubebuilder.io/test-tools/$(KUBE_VERSION)/$(OS)/$(ARCH) -o kubebuilder-tools.tar.gz
	mkdir -p _test/kubebuilder
	tar -xvf kubebuilder-tools.tar.gz
	mv kubebuilder/bin/* _test/kubebuilder/
	rm kubebuilder-tools.tar.gz
	rm -rf kubebuilder/
	mkdir -p _test/data
	$(GO) test -v .

##@ Clean
.PHONY: clean
clean: ## Clean kubebuilder, helm and test data artifacts
	@rm -rf _test/
	@rm -rf _chart/

##@ Build
.PHONY: build
build: ## Build the docker image
	@DOCKER_BUILDKIT=1 docker build -t $(REGISTRY)/$(IMAGE_NAME):$(IMAGE_TAG) .

##@ Release
.PHONY: release
release: ## Push and release the docker image to the public registry
	@docker push $(REGISTRY)/$(IMAGE_NAME):$(IMAGE_TAG)

##@ Build Latest
.PHONY: build-latest
build-latest: ## Build the docker image with latest tag
	@DOCKER_BUILDKIT=1 docker build -t $(REGISTRY)/$(IMAGE_NAME):latest .

##@ Release Latest
.PHONY: release-latest
release-latest: ## Push and release the docker image to the public registry with latest tag
	@docker push $(REGISTRY)/$(IMAGE_NAME):latest

##@ Verify Chart
.PHONY: verify-chart
verify-chart: ## Lint the helm chart for errors
	@helm lint chart/

##@ Generate Chart
.PHONY: generate-chart
generate-chart: ## Generate the helm chart artifacts for release
	@helm package chart/
	@mkdir -p _chart/
	@mv cert-manager-porkbun-webhook-*.tgz _chart/
