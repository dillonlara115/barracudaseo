.PHONY: build test install clean release frontend-build frontend-dev serve docker-build docker-push deploy-backend

# Build the binary (requires frontend to be built first)
build: frontend-build
	go build -o bin/barracuda .

# Run tests
test:
	go test ./... -v

# Run tests with coverage
test-coverage:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out

# Install to $GOPATH/bin (requires frontend to be built first)
install: frontend-build
	go install .

# Install alias to ~/.zshrc
install-alias:
	@echo "Adding alias to ~/.zshrc..."
	@if ! grep -q "alias barracuda=" ~/.zshrc 2>/dev/null; then \
		echo "alias barracuda=\"$(shell pwd)/bin/barracuda\"" >> ~/.zshrc; \
		echo "✓ Alias added to ~/.zshrc"; \
		echo "Run 'source ~/.zshrc' or restart your terminal to use it"; \
	else \
		echo "⚠️  Alias already exists in ~/.zshrc"; \
	fi

# Add bin directory to PATH in ~/.zshrc
install-path:
	@echo "Adding bin directory to PATH in ~/.zshrc..."
	@if ! grep -q "$(shell pwd)/bin" ~/.zshrc 2>/dev/null; then \
		echo 'export PATH="$(shell pwd)/bin:$$PATH"' >> ~/.zshrc; \
		echo "✓ PATH updated in ~/.zshrc"; \
		echo "Run 'source ~/.zshrc' or restart your terminal to use it"; \
	else \
		echo "⚠️  PATH already includes bin directory"; \
	fi

# Clean build artifacts
clean:
	rm -rf bin/
	rm -f coverage.out

# Build for multiple platforms (requires frontend to be built first)
release: frontend-build
	GOOS=linux GOARCH=amd64 go build -o bin/barracuda-linux-amd64 .
	GOOS=linux GOARCH=arm64 go build -o bin/barracuda-linux-arm64 .
	GOOS=darwin GOARCH=amd64 go build -o bin/barracuda-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 go build -o bin/barracuda-darwin-arm64 .
	GOOS=windows GOARCH=amd64 go build -o bin/barracuda-windows-amd64.exe .

# Format code
fmt:
	go fmt ./...

# Lint code
lint:
	golangci-lint run

# Run benchmarks
bench:
	go test -bench=. -benchmem ./...

# Build frontend
frontend-build:
	cd web && npm install && npm run build

# Run frontend in dev mode
frontend-dev:
	cd web && npm install && npm run dev

# Serve results (requires built frontend)
serve:
	go run . serve --results results.json --graph graph.json

# Docker and Cloud Run deployment targets
# Set these environment variables:
#   GCP_PROJECT_ID - Your Google Cloud project ID
#   GCP_REGION - Region for deployment (default: us-central1)
#   IMAGE_NAME - Docker image name (default: barracuda-api)

GCP_PROJECT_ID ?= $(shell gcloud config get-value project 2>/dev/null)
GCP_REGION ?= us-central1
IMAGE_NAME ?= barracuda-api
IMAGE_TAG ?= latest
REPOSITORY ?= barracuda
ENV_FILE ?= $(CURDIR)/.env
IMAGE_URI = $(GCP_REGION)-docker.pkg.dev/$(GCP_PROJECT_ID)/$(REPOSITORY)/$(IMAGE_NAME):$(IMAGE_TAG)

# Build Docker image
docker-build:
	@if [ -z "$(GCP_PROJECT_ID)" ]; then \
		echo "Error: GCP_PROJECT_ID not set. Set it or run: gcloud config set project YOUR_PROJECT_ID"; \
		exit 1; \
	fi
	@echo "Building Docker image for linux/amd64 (Cloud Run requirement): $(IMAGE_URI)"
	docker build --platform linux/amd64 -t $(IMAGE_URI) .
	@echo "✓ Docker image built successfully"

# Push Docker image to Artifact Registry
docker-push: docker-build
	@if [ -z "$(GCP_PROJECT_ID)" ]; then \
		echo "Error: GCP_PROJECT_ID not set"; \
		exit 1; \
	fi
	@echo "Pushing Docker image to Artifact Registry..."
	gcloud auth configure-docker $(GCP_REGION)-docker.pkg.dev --quiet
	docker push $(IMAGE_URI)
	@echo "✓ Docker image pushed successfully"

# Deploy to Cloud Run
deploy-backend: docker-push
	@set -a; \
	if [ -f "$(ENV_FILE)" ]; then \
		echo "Loading environment variables from $(ENV_FILE)..."; \
		. "$(ENV_FILE)"; \
	fi; \
	if [ -f "$(CURDIR)/.env.local" ]; then \
		echo "Loading overrides from .env.local..."; \
		. "$(CURDIR)/.env.local"; \
	fi; \
	if [ -z "$$PUBLIC_SUPABASE_URL" ] && [ -z "$$STRIPE_SECRET_KEY" ]; then \
		echo "Warning: No environment variables loaded. Check .env or .env.local"; \
	fi; \
	set +a; \
	if [ -z "$$GCP_PROJECT_ID" ] && [ -z "$(GCP_PROJECT_ID)" ]; then \
		GCP_PROJECT_ID=$$(gcloud config get-value project 2>/dev/null); \
		if [ -z "$$GCP_PROJECT_ID" ]; then \
			echo "Error: GCP_PROJECT_ID not set"; \
			exit 1; \
		fi; \
	fi; \
	if [ -z "$$PUBLIC_SUPABASE_URL" ] && [ -z "$(PUBLIC_SUPABASE_URL)" ]; then \
		echo "Error: PUBLIC_SUPABASE_URL not set. Set it in $(ENV_FILE) or export it."; \
		exit 1; \
	fi; \
	if [ -z "$$PUBLIC_SUPABASE_ANON_KEY" ] && [ -z "$(PUBLIC_SUPABASE_ANON_KEY)" ]; then \
		echo "Error: PUBLIC_SUPABASE_ANON_KEY not set. Set it in $(ENV_FILE) or export it."; \
		exit 1; \
	fi; \
	echo "Deploying to Cloud Run..."; \
	SUPABASE_URL="$${PUBLIC_SUPABASE_URL:-$(PUBLIC_SUPABASE_URL)}"; \
	SUPABASE_ANON="$${PUBLIC_SUPABASE_ANON_KEY:-$(PUBLIC_SUPABASE_ANON_KEY)}"; \
	GCP_PROJ="$${GCP_PROJECT_ID:-$(GCP_PROJECT_ID)}"; \
	GCP_REG="$${GCP_REGION:-$(GCP_REGION)}"; \
	ENV_VARS="PUBLIC_SUPABASE_URL=$$SUPABASE_URL,PUBLIC_SUPABASE_ANON_KEY=$$SUPABASE_ANON"; \
	if [ -n "$$STRIPE_SECRET_KEY" ]; then \
		ENV_VARS="$$ENV_VARS,STRIPE_SECRET_KEY=$$STRIPE_SECRET_KEY"; \
	fi; \
	if [ -n "$$STRIPE_WEBHOOK_SECRET" ]; then \
		ENV_VARS="$$ENV_VARS,STRIPE_WEBHOOK_SECRET=$$STRIPE_WEBHOOK_SECRET"; \
	fi; \
	if [ -n "$$STRIPE_PRICE_ID_PRO" ]; then \
		ENV_VARS="$$ENV_VARS,STRIPE_PRICE_ID_PRO=$$STRIPE_PRICE_ID_PRO"; \
	fi; \
	if [ -n "$$STRIPE_PRICE_ID_PRO_ANNUAL" ]; then \
		ENV_VARS="$$ENV_VARS,STRIPE_PRICE_ID_PRO_ANNUAL=$$STRIPE_PRICE_ID_PRO_ANNUAL"; \
	fi; \
	if [ -n "$$STRIPE_PRICE_ID_TEAM_SEAT" ]; then \
		ENV_VARS="$$ENV_VARS,STRIPE_PRICE_ID_TEAM_SEAT=$$STRIPE_PRICE_ID_TEAM_SEAT"; \
	fi; \
	if [ -n "$$STRIPE_SUCCESS_URL" ]; then \
		ENV_VARS="$$ENV_VARS,STRIPE_SUCCESS_URL=$$STRIPE_SUCCESS_URL"; \
	fi; \
	if [ -n "$$STRIPE_CANCEL_URL" ]; then \
		ENV_VARS="$$ENV_VARS,STRIPE_CANCEL_URL=$$STRIPE_CANCEL_URL"; \
	fi; \
	gcloud run deploy $(IMAGE_NAME) \
		--image $$GCP_REG-docker.pkg.dev/$$GCP_PROJ/$(REPOSITORY)/$(IMAGE_NAME):$(IMAGE_TAG) \
		--platform managed \
		--region $$GCP_REG \
		--allow-unauthenticated \
		--update-env-vars="$$ENV_VARS" \
		--set-secrets="SUPABASE_SERVICE_ROLE_KEY=supabase-service-role-key:latest" \
		--memory=512Mi \
		--cpu=1 \
		--timeout=300 \
		--max-instances=10 \
		--port=8080 \
		--quiet; \
	echo "✓ Deployment complete!"; \
	echo "Service URL:"; \
	GCP_REG="$${GCP_REGION:-$(GCP_REGION)}"; \
	gcloud run services describe $(IMAGE_NAME) \
		--platform managed \
		--region $$GCP_REG \
		--format="value(status.url)"

# Deploy only the image (preserves existing env vars) - faster for regular deployments
deploy-image: docker-push
	@if [ -z "$(GCP_PROJECT_ID)" ] && [ -z "$(GCP_PROJECT_ID)" ]; then \
		GCP_PROJECT_ID=$$(gcloud config get-value project 2>/dev/null); \
		if [ -z "$$GCP_PROJECT_ID" ]; then \
			echo "Error: GCP_PROJECT_ID not set"; \
			exit 1; \
		fi; \
	fi; \
	GCP_PROJ="$${GCP_PROJECT_ID:-$(GCP_PROJECT_ID)}"; \
	GCP_REG="$${GCP_REGION:-$(GCP_REGION)}"; \
	echo "Deploying new image (preserving existing env vars)..."; \
	gcloud run services update $(IMAGE_NAME) \
		--image $$GCP_REG-docker.pkg.dev/$$GCP_PROJ/$(REPOSITORY)/$(IMAGE_NAME):$(IMAGE_TAG) \
		--platform managed \
		--region $$GCP_REG \
		--quiet; \
	echo "✓ Image deployment complete!"; \
	echo "Service URL:"; \
	gcloud run services describe $(IMAGE_NAME) \
		--platform managed \
		--region $$GCP_REG \
		--format="value(status.url)"

# Quick deploy (rebuild and deploy in one command)
# Use deploy-image for faster deployments that don't change env vars
deploy: deploy-backend
