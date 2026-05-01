# Variables
DOCKER_REGISTRY=192.168.1.27:81/ai-saas
#DOCKER_REGISTRY=registry.yingquntech.com/ai-saas
TENANT_ADMIN_API_IMAGE=$(DOCKER_REGISTRY)/tenant-admin-api
TENANT_API_IMAGE=$(DOCKER_REGISTRY)/tenant-api
TENANT_RMQ_IMAGE=$(DOCKER_REGISTRY)/tenant-rmq
PLATFORM_OPENAPI_IMAGE=$(DOCKER_REGISTRY)/platform-openapi
ENV=prod
VERSION=$(ENV)-latest

# Build Docker images
build-tenant-admin-api:
	@echo "Building web Docker image: $(TENANT_ADMIN_API_IMAGE):$(VERSION)..."
	docker build -t $(TENANT_ADMIN_API_IMAGE):$(VERSION) -f ./cmd/tenant/admin-api/Dockerfile .
	@echo "Tenant Admin Api Docker image built successfully: $(TENANT_ADMIN_API_IMAGE):$(VERSION)"

build-tenant-api:
	@echo "Building web Docker image: $(TENANT_API_IMAGE):$(VERSION)..."
	docker build -t $(TENANT_API_IMAGE):$(VERSION) -f ./cmd/tenant/api/Dockerfile .
	@echo "Tenant Api Docker image built successfully: $(TENANT_API_IMAGE):$(VERSION)"

build-tenant-rmq:
	@echo "Building web Docker image: $(TENANT_RMQ_IMAGE):$(VERSION)..."
	docker build -t $(TENANT_RMQ_IMAGE):$(VERSION) -f ./cmd/tenant/rmq/Dockerfile .
	@echo "Tenant RMQ Docker image built successfully: $(TENANT_RMQ_IMAGE):$(VERSION)"

build-platform-openapi:
	@echo "Building web Docker image: $(PLATFORM_OPENAPI_IMAGE):$(VERSION)..."
	docker build -t $(PLATFORM_OPENAPI_IMAGE):$(VERSION) -f ./cmd/platform/openapi/Dockerfile .
	@echo "Platform OpenApi Docker image built successfully: $(PLATFORM_OPENAPI_IMAGE):$(VERSION)"

push-tenant-admin-api:
	@echo "Pushing Docker image: $(TENANT_ADMIN_API_IMAGE):$(VERSION)..."
	docker push $(TENANT_ADMIN_API_IMAGE):$(VERSION)
	@echo "Tenant Admin Api Docker image pushed successfully: $(TENANT_ADMIN_API_IMAGE):$(VERSION)"

push-tenant-api:
	@echo "Pushing Docker image: $(TENANT_API_IMAGE):$(VERSION)..."
	docker push $(TENANT_API_IMAGE):$(VERSION)
	@echo "Tenant Api Docker image pushed successfully: $(TENANT_API_IMAGE):$(VERSION)"

push-tenant-rmq:
	@echo "Pushing Docker image: $(TENANT_RMQ_IMAGE):$(VERSION)..."
	docker push $(TENANT_RMQ_IMAGE):$(VERSION)
	@echo "Tenant RMQ Docker image pushed successfully: $(TENANT_RMQ_IMAGE):$(VERSION)"

push-platform-openapi:
	@echo "Pushing Docker image: $(PLATFORM_OPENAPI_IMAGE):$(VERSION)..."
	docker push $(PLATFORM_OPENAPI_IMAGE):$(VERSION)
	@echo "Platform OpenApi Docker image pushed successfully: $(PLATFORM_OPENAPI_IMAGE):$(VERSION)"

build-push-tenant-admin-api: build-tenant-admin-api push-tenant-admin-api

build-push-tenant-api: build-tenant-api push-tenant-api

build-push-tenant-rmq: build-tenant-rmq push-tenant-rmq

build-push-platform-openapi: build-platform-openapi push-platform-openapi

# Build all images
build-all: build-tenant-admin-api build-tenant-api build-tenant-rmq build-platform-openapi

# Push all images
push-all: push-tenant-admin-api push-tenant-api push-tenant-rmq push-platform-openapi

# Build and push all images
build-push-all: build-all push-all

.PHONY: build-tenant-admin-api build-tenant-api build-tenant-rmq build-platform-openapi push-tenant-admin-api push-tenant-api push-tenant-rmq push-platform-openapi build-push-tenant-admin-api build-push-tenant-api build-push-tenant-rmq build-push-platform-openapi build-all push-all build-push-all
