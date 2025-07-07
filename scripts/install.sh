#!/bin/bash

# XFS Quota Kit Installation Script

set -e

BINARY_NAME="xfs-quota-kit"
INSTALL_DIR="/usr/local/bin"
CONFIG_DIR="/etc/xfs-quota-kit"
LOG_DIR="/var/log/xfs-quota-kit"
BACKUP_DIR="/var/backups/xfs-quota-kit"

# 检查权限
if [ "$EUID" -ne 0 ]; then
    echo "This script must be run as root" 
    exit 1
fi

echo "Installing XFS Quota Kit..."

# 创建目录
echo "Creating directories..."
mkdir -p "$CONFIG_DIR"
mkdir -p "$LOG_DIR"
mkdir -p "$BACKUP_DIR"

# 复制二进制文件
echo "Installing binary..."
if [ -f "build/$BINARY_NAME" ]; then
    cp "build/$BINARY_NAME" "$INSTALL_DIR/"
    chmod +x "$INSTALL_DIR/$BINARY_NAME"
else
    echo "Binary not found. Please run 'make build' first."
    exit 1
fi

# 复制配置文件
echo "Installing configuration..."
if [ -f "configs/config.yaml" ]; then
    cp "configs/config.yaml" "$CONFIG_DIR/"
fi

# 设置权限
echo "Setting permissions..."
chown -R root:root "$CONFIG_DIR"
chown -R root:root "$LOG_DIR"
chown -R root:root "$BACKUP_DIR"
chmod 755 "$CONFIG_DIR"
chmod 755 "$LOG_DIR"
chmod 755 "$BACKUP_DIR"

# 创建systemd服务文件（可选）
if command -v systemctl >/dev/null 2>&1; then
    echo "Creating systemd service..."
    cat > /etc/systemd/system/xfs-quota-kit.service << EOF
[Unit]
Description=XFS Quota Kit Server
After=network.target

[Service]
Type=simple
User=root
ExecStart=$INSTALL_DIR/$BINARY_NAME server --config $CONFIG_DIR/config.yaml
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
EOF

    systemctl daemon-reload
    echo "Service created. Enable with: systemctl enable xfs-quota-kit"
fi

echo "Installation completed successfully!"
echo ""
echo "Usage:"
echo "  $BINARY_NAME --help"
echo ""
echo "Configuration file: $CONFIG_DIR/config.yaml"
echo "Log directory: $LOG_DIR"
echo "Backup directory: $BACKUP_DIR" 