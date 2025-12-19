#!/bin/bash
# Bulk set Cloud Run environment variables from .env file (production only)
# This script ONLY reads from .env - it NEVER reads from .env.local
# 
# Usage:
#   ./scripts/set-cloud-run-env.sh                    # Reads from .env file and merges
#   ./scripts/set-cloud-run-env.sh --replace           # Replaces all vars from .env
#   ./scripts/set-cloud-run-env.sh VAR1=val1 VAR2=val2 # Sets specific vars (overrides .env)

set -e

SERVICE_NAME=${SERVICE_NAME:-barracuda-api}
GCP_REGION=${GCP_REGION:-us-central1}
ENV_FILE=".env"  # Always use .env (production), never .env.local
USE_MERGE=true   # Use --update-env-vars by default (merges)

# Secret Manager mappings (env var -> secret name)
declare -A SECRET_VARS=(
  ["SUPABASE_SERVICE_ROLE_KEY"]="supabase_service_role_key"
  ["SUPABASE_JWT_SECRET"]="supabase_jwt_secret"
  ["STRIPE_SECRET_KEY"]="stripe_secret_key"
)

# Parse arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        --replace)
            USE_MERGE=false  # Use --set-env-vars (replaces all)
            shift
            ;;
        --service)
            SERVICE_NAME="$2"
            shift 2
            ;;
        --region)
            GCP_REGION="$2"
            shift 2
            ;;
        --help)
            echo "Usage: $0 [--replace] [--service NAME] [--region REGION] [VAR=value ...]"
            echo ""
            echo "This script reads from .env file (production variables only)."
            echo "It NEVER reads from .env.local to prevent accidental local config deployment."
            echo ""
            echo "Options:"
            echo "  --replace      Replace all env vars instead of merging (use --set-env-vars)"
            echo "  --service NAME Cloud Run service name (default: barracuda-api)"
            echo "  --region REGION GCP region (default: us-central1)"
            echo ""
            echo "Examples:"
            echo "  $0                                    # Read from .env and merge"
            echo "  $0 --replace                          # Replace all vars from .env"
            echo "  $0 VAR1=val1 VAR2=val2               # Set specific vars (overrides .env)"
            exit 0
            ;;
        *=*)
            # Direct variable assignment (e.g., VAR=value)
            if [ -z "$ENV_VARS" ]; then
                ENV_VARS="$1"
            else
                ENV_VARS="$ENV_VARS,$1"
            fi
            shift
            ;;
        *)
            echo "Unknown option: $1"
            echo "Use --help for usage information"
            exit 1
            ;;
    esac
done

# If no direct vars provided, read from .env file (production only)
if [ -z "$ENV_VARS" ]; then
    if [ ! -f "$ENV_FILE" ]; then
        echo "Error: .env file not found!"
        echo "Create a .env file with your production environment variables."
        exit 1
    fi
    
    echo "Reading production environment variables from $ENV_FILE..."
    echo "(Note: .env.local is never read - production only)"
    echo ""
    
    # Filter out comments and empty lines, then format as KEY=VALUE pairs
    while IFS= read -r line || [ -n "$line" ]; do
        # Skip comments and empty lines
        [[ "$line" =~ ^[[:space:]]*# ]] && continue
        [[ -z "${line// }" ]] && continue
        
        # Remove inline comments
        line="${line%%#*}"
        line="${line%"${line##*[![:space:]]}"}"  # trim trailing whitespace
        
        # Skip if no = sign
        [[ ! "$line" =~ = ]] && continue
        
        # Add to ENV_VARS or SECRET_ENV_VARS
        VAR_NAME="${line%%=*}"
        VAR_VALUE="${line#*=}"
        if [[ -n "${SECRET_VARS[$VAR_NAME]}" ]]; then
            SECRET_SPEC="${VAR_NAME}=${SECRET_VARS[$VAR_NAME]}:latest"
            if [ -z "$SECRET_ENV_VARS" ]; then
                SECRET_ENV_VARS="$SECRET_SPEC"
            else
                SECRET_ENV_VARS="$SECRET_ENV_VARS,$SECRET_SPEC"
            fi
        else
            if [ -z "$ENV_VARS" ]; then
                ENV_VARS="$line"
            else
                ENV_VARS="$ENV_VARS,$line"
            fi
        fi
    done < "$ENV_FILE"
fi

if [ -z "$ENV_VARS" ] && [ -z "$SECRET_ENV_VARS" ]; then
    echo "Error: No environment variables to set."
    echo "Provide variables via command line (VAR=value) or populate $ENV_FILE"
    exit 1
fi

echo "Service: $SERVICE_NAME"
echo "Region: $GCP_REGION"
echo "Mode: $([ "$USE_MERGE" = true ] && echo "Merge (--update-env-vars)" || echo "Replace (--set-env-vars)")"
echo "Source: $ENV_FILE (production only - .env.local is never used)"
echo ""
if [ -n "$ENV_VARS" ]; then
    echo "Variables to set:"
    echo "$ENV_VARS" | tr ',' '\n' | sed 's/=.*/=***/' | sed 's/^/  - /'
    echo ""
fi

if [ -n "$SECRET_ENV_VARS" ]; then
    echo "Secrets to set:"
    echo "$SECRET_ENV_VARS" | tr ',' '\n' | sed 's/=.*/=*** (secret)/' | sed 's/^/  - /'
    echo ""
fi

# Count variables
VAR_COUNT=$(echo "$ENV_VARS" | tr ',' '\n' | wc -l | xargs)
SECRET_COUNT=$(echo "$SECRET_ENV_VARS" | tr ',' '\n' | wc -l | xargs)
echo "Total variables: $VAR_COUNT"
echo "Total secrets: $SECRET_COUNT"
echo ""

read -p "Continue? (y/N) " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "Cancelled."
    exit 1
fi

# Execute the gcloud command
GCLOUD_ARGS=(
    gcloud run services update "$SERVICE_NAME"
    --platform managed
    --region "$GCP_REGION"
    --quiet
)

if [ "$USE_MERGE" = true ]; then
    if [ -n "$ENV_VARS" ]; then
        GCLOUD_ARGS+=(--update-env-vars="$ENV_VARS")
    fi
    if [ -n "$SECRET_ENV_VARS" ]; then
        GCLOUD_ARGS+=(--update-secrets="$SECRET_ENV_VARS")
    fi
else
    if [ -n "$ENV_VARS" ]; then
        GCLOUD_ARGS+=(--set-env-vars="$ENV_VARS")
    fi
    if [ -n "$SECRET_ENV_VARS" ]; then
        GCLOUD_ARGS+=(--set-secrets="$SECRET_ENV_VARS")
    fi
fi

"${GCLOUD_ARGS[@]}"

echo ""
echo "✓ Environment variables updated!"
echo ""
echo "Service URL:"
gcloud run services describe "$SERVICE_NAME" \
    --platform managed \
    --region "$GCP_REGION" \
    --format="value(status.url)"
