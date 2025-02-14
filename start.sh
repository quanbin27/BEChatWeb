#!/bin/sh

if [ "$SERVICE" = "grpc-server" ]; then
    exec /root/grpc-server
elif [ "$SERVICE" = "api-gateway" ]; then
    exec /root/api-gateway
else
    echo "Error: SERVICE not specified!"
    exit 1
fi
