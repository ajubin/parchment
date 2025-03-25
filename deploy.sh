#!/bin/bash

# Fail on error
set -e


SERVICE_NAME="parchment"
DEPLOY_DIR="/home/pi/printer-serial"
USER="pi"
HOST="192.168.1.97"
PI_SERVER_ADDR=${USER}@${HOST}
SERVICE_PATH="/etc/systemd/system/${SERVICE_NAME}.service"

source .env
echo "testing before deploying"
make test

echo "🚀 Deploying ${SERVICE_NAME} to Raspberry Pi (${HOST})..."

rsync -avz --exclude '.git' --exclude 'deploy.sh' --exclude 'assets' . ${PI_SERVER_ADDR}:${DEPLOY_DIR}



echo "🚀 Updating ${SERVICE_NAME} systemd service"
ssh -q -T ${PI_SERVER_ADDR} << EOF
    # Add Go to the system PATH
    export PATH=\$PATH:/home/pi/.local/share/go/bin/
		
    echo "🔧 Creating/updating systemd service..."
    echo "[Unit]
Description=Parchment Print Server in Go
After=network.target

[Service]
User=pi
WorkingDirectory=${DEPLOY_DIR}
ExecStart=${DEPLOY_DIR}/${SERVICE_NAME} --serialPort /dev/ttyS0 --apiUser ${API_USER} --apiToken ${API_PASSWORD}
Restart=always
# Not needed anymore
# Environment=\"MODE=prod\"

[Install]
WantedBy=multi-user.target
" | sudo tee ${SERVICE_PATH} > /dev/null

    # 3️⃣ Reload systemd and enable the service
    echo "🔄 Reloading systemd..."
    sudo systemctl daemon-reload
    sudo systemctl enable ${SERVICE_NAME}

    # 4️⃣ Build the Go binary
    cd ${DEPLOY_DIR}
    echo "🔨 Building ${SERVICE_NAME}..."
    go build -o ${SERVICE_NAME} main.go

    # 5️⃣ Restart the service
    echo "🔄 Restarting service..."
    sudo systemctl restart ${SERVICE_NAME}
    echo "✅ Deployment complete!"
EOF
