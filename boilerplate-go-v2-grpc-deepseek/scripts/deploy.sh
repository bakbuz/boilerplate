#!/bin/bash

set -e

echo "🚀 Starting deployment..."

# Variables
EC2_INSTANCE="ec2-user@your-ec2-instance"
APP_NAME="grpc-highperf-backend"
INSTALL_DIR="/opt/$APP_NAME"
SERVICE_NAME="grpc-backend"

# Build binary
echo "📦 Building binary..."
make build-prod

# Create deployment package
echo "📁 Creating deployment package..."
mkdir -p deploy
cp bin/server deploy/
cp deployment/systemd/$SERVICE_NAME.service deploy/
cp .env.example deploy/.env

# Copy to EC2
echo "🖥️  Copying files to EC2..."
scp -r deploy/* $EC2_INSTANCE:/tmp/$APP_NAME/

# Setup on EC2
echo "🔧 Setting up on EC2..."
ssh $EC2_INSTANCE << EOF
    set -e
    
    # Create user and directories
    sudo groupadd -f appgroup
    sudo useradd -r -s /bin/false -g appgroup appuser || true
    
    # Stop service
    sudo systemctl stop $SERVICE_NAME || true
    
    # Backup existing
    sudo mv $INSTALL_DIR $INSTALL_DIR.backup.\$(date +%Y%m%d%H%M%S) || true
    
    # Copy new files
    sudo mkdir -p $INSTALL_DIR
    sudo cp -r /tmp/$APP_NAME/* $INSTALL_DIR/
    
    # Set permissions
    sudo chown -R appuser:appgroup $INSTALL_DIR
    sudo chmod 750 $INSTALL_DIR
    sudo chmod 550 $INSTALL_DIR/server
    
    # Setup systemd
    sudo cp $INSTALL_DIR/$SERVICE_NAME.service /etc/systemd/system/
    sudo systemctl daemon-reload
    
    # Setup environment
    if [ ! -f /etc/$APP_NAME/.env ]; then
        sudo mkdir -p /etc/$APP_NAME
        sudo cp $INSTALL_DIR/.env /etc/$APP_NAME/
    fi
    
    # Start service
    sudo systemctl enable $SERVICE_NAME
    sudo systemctl start $SERVICE_NAME
    
    # Cleanup
    rm -rf /tmp/$APP_NAME
    
    echo "✅ Deployment completed!"
    echo "📊 Service status:"
    sudo systemctl status $SERVICE_NAME --no-pager
EOF

echo "🎉 Deployment completed successfully!"