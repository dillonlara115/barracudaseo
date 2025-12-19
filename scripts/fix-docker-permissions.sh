#!/bin/bash

# Script to fix Docker permissions
# This adds your user to the docker group so you don't need sudo

echo "🔧 Fixing Docker Permissions"
echo "============================"
echo ""

# Check if user is already in docker group
if groups | grep -q docker; then
    echo "✅ You are already in the docker group"
    echo ""
    echo "If you still get permission errors, try:"
    echo "1. Log out and log back in"
    echo "2. Or restart your terminal session"
    exit 0
fi

echo "Adding user to docker group..."
echo ""

# Add user to docker group
sudo usermod -aG docker $USER

echo "✅ Added $USER to docker group"
echo ""
echo "⚠️  IMPORTANT: You need to log out and log back in for this to take effect!"
echo ""
echo "After logging back in, you should be able to run:"
echo "  docker ps"
echo "  supabase start"
echo "  make deploy-backend"
echo ""
echo "Without needing sudo!"


