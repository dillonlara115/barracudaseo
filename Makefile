.PHONY: build test install clean release frontend-build frontend-dev serve docker-build docker-push deploy-backend setup-crawl-job

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
# Cloud Run Job that runs crawls (same image as the API, different command)
CRAWL_JOB_NAME ?= barracuda-crawl
CRAWL_JOB_MEMORY ?= 2Gi
CRAWL_JOB_CPU ?= 2
CRAWL_JOB_TIMEOUT ?= 6h
# Pages per crawl rendered in headless Chrome (0 turns rendering off)
CRAWL_JOB_RENDER_PAGES ?= 5

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
	if [ -z "$$PUBLIC_SUPABASE_URL" ] && [ -z "$$STRIPE_SECRET_KEY" ]; then \
		echo "Warning: No environment variables loaded. Check .env"; \
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
	# STRIPE_SECRET_KEY is managed via Secret Manager (see --set-secrets) \
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
		--set-secrets="SUPABASE_SERVICE_ROLE_KEY=supabase_service_role_key:latest,SUPABASE_JWT_SECRET=supabase_jwt_secret:latest,STRIPE_SECRET_KEY=stripe_secret_key:latest" \
		--memory=512Mi \
		--cpu=1 \
		--timeout=300 \
		--max-instances=10 \
		--port=8080 \
		--quiet; \
	if gcloud run jobs describe $(CRAWL_JOB_NAME) --region $$GCP_REG --format="value(name)" >/dev/null 2>&1; then \
		echo "Updating crawl job image..."; \
		gcloud run jobs update $(CRAWL_JOB_NAME) \
			--image $$GCP_REG-docker.pkg.dev/$$GCP_PROJ/$(REPOSITORY)/$(IMAGE_NAME):$(IMAGE_TAG) \
			--region $$GCP_REG \
			--quiet; \
	else \
		echo "Note: crawl job $(CRAWL_JOB_NAME) does not exist yet. Run 'make setup-crawl-job' to make crawls durable."; \
	fi; \
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
	if gcloud run jobs describe $(CRAWL_JOB_NAME) --region $$GCP_REG --format="value(name)" >/dev/null 2>&1; then \
		echo "Updating crawl job image..."; \
		gcloud run jobs update $(CRAWL_JOB_NAME) \
			--image $$GCP_REG-docker.pkg.dev/$$GCP_PROJ/$(REPOSITORY)/$(IMAGE_NAME):$(IMAGE_TAG) \
			--region $$GCP_REG \
			--quiet; \
	else \
		echo "Note: crawl job $(CRAWL_JOB_NAME) does not exist yet. Run 'make setup-crawl-job' to make crawls durable."; \
	fi; \
	echo "✓ Image deployment complete!"; \
	echo "Service URL:"; \
	gcloud run services describe $(IMAGE_NAME) \
		--platform managed \
		--region $$GCP_REG \
		--format="value(status.url)"

# One-time setup of the Cloud Run Job that runs crawls. Safe to run again: it
# updates the job if it already exists.
#   1. creates/updates the job from the API's image, with the Supabase config
#   2. lets the API's service account start executions of it
#   3. tells the API which job to use (CRAWL_JOB)
# PSI_API_KEY is copied to the job when it is set in $(ENV_FILE); without it
# crawls skip Core Web Vitals. The job renders CRAWL_JOB_RENDER_PAGES pages per
# crawl in headless Chrome.
setup-crawl-job:
	@set -a; \
	if [ -f "$(ENV_FILE)" ]; then . "$(ENV_FILE)"; fi; \
	set +a; \
	GCP_PROJ="$${GCP_PROJECT_ID:-$(GCP_PROJECT_ID)}"; \
	GCP_REG="$${GCP_REGION:-$(GCP_REGION)}"; \
	SUPABASE_URL="$${PUBLIC_SUPABASE_URL:-$(PUBLIC_SUPABASE_URL)}"; \
	SUPABASE_ANON="$${PUBLIC_SUPABASE_ANON_KEY:-$(PUBLIC_SUPABASE_ANON_KEY)}"; \
	if [ -z "$$GCP_PROJ" ] || [ -z "$$SUPABASE_URL" ] || [ -z "$$SUPABASE_ANON" ]; then \
		echo "Error: GCP_PROJECT_ID, PUBLIC_SUPABASE_URL and PUBLIC_SUPABASE_ANON_KEY must be set (in $(ENV_FILE) or exported)."; \
		exit 1; \
	fi; \
	ENV_VARS="PUBLIC_SUPABASE_URL=$$SUPABASE_URL,PUBLIC_SUPABASE_ANON_KEY=$$SUPABASE_ANON,RENDER_MAX_PAGES=$(CRAWL_JOB_RENDER_PAGES)"; \
	if [ -n "$$PSI_API_KEY" ]; then \
		ENV_VARS="$$ENV_VARS,PSI_API_KEY=$$PSI_API_KEY"; \
	else \
		echo "Note: PSI_API_KEY is not set in $(ENV_FILE), so crawls will not measure Core Web Vitals."; \
	fi; \
	if gcloud run jobs describe $(CRAWL_JOB_NAME) --region $$GCP_REG --format="value(name)" >/dev/null 2>&1; then \
		VERB=update; \
	else \
		VERB=create; \
	fi; \
	echo "Running: gcloud run jobs $$VERB $(CRAWL_JOB_NAME)"; \
	gcloud run jobs $$VERB $(CRAWL_JOB_NAME) \
		--image $$GCP_REG-docker.pkg.dev/$$GCP_PROJ/$(REPOSITORY)/$(IMAGE_NAME):$(IMAGE_TAG) \
		--region $$GCP_REG \
		--args=crawl-job \
		--memory=$(CRAWL_JOB_MEMORY) \
		--cpu=$(CRAWL_JOB_CPU) \
		--task-timeout=$(CRAWL_JOB_TIMEOUT) \
		--max-retries=1 \
		--update-env-vars="$$ENV_VARS" \
		--set-secrets="SUPABASE_SERVICE_ROLE_KEY=supabase_service_role_key:latest,SUPABASE_JWT_SECRET=supabase_jwt_secret:latest" \
		--quiet || exit 1; \
	SERVICE_ACCOUNT=$$(gcloud run services describe $(IMAGE_NAME) --region $$GCP_REG --format="value(spec.template.spec.serviceAccountName)"); \
	if [ -z "$$SERVICE_ACCOUNT" ]; then \
		echo "Error: could not read the service account of $(IMAGE_NAME)."; \
		exit 1; \
	fi; \
	echo "Allowing $$SERVICE_ACCOUNT to start the crawl job..."; \
	gcloud run jobs add-iam-policy-binding $(CRAWL_JOB_NAME) \
		--region $$GCP_REG \
		--member="serviceAccount:$$SERVICE_ACCOUNT" \
		--role=roles/run.jobsExecutorWithOverrides \
		--quiet >/dev/null || exit 1; \
	echo "Pointing the API at the crawl job..."; \
	gcloud run services update $(IMAGE_NAME) \
		--region $$GCP_REG \
		--update-env-vars="CRAWL_JOB=projects/$$GCP_PROJ/locations/$$GCP_REG/jobs/$(CRAWL_JOB_NAME)" \
		--quiet || exit 1; \
	echo "✓ Crawl job ready: new crawls run as executions of $(CRAWL_JOB_NAME)"

# Quick deploy (rebuild and deploy in one command)
# Use deploy-image for faster deployments that don't change env vars
deploy: deploy-backend
