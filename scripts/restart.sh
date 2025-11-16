#!/bin/bash
set -e

APP_DIR=/home/ec2-user/app
APP_NAME=order-management-service-ci
LOG_FILE=/tmp/order-management-service.log

cd "$APP_DIR"

echo "Buscando proceso anterior..."
PID=$(pgrep -f "$APP_NAME" || true)

if [ -n "$PID" ]; then
  echo "Deteniendo proceso $PID"
  kill "$PID" || true
  sleep 3
fi

echo "Arrancando nueva versión..."
nohup ./"$APP_NAME" >> "$LOG_FILE" 2>&1 &

echo "Despliegue OK (logs en $LOG_FILE)"

