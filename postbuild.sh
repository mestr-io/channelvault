#!/bin/bash

USERVER="0.1"
VERSION=$(git rev-parse --short HEAD)
DATE=$(date -u +%Y%m%d%H%M%S)

# Check if the first command-line argument is "macOS"
if [ "$1" = "macOS" ]; then
   echo work to be done...
fi

if [ "$1" = "win" ]; then
    mv ChannelVault.exe ChannelVault-$USERVER-$VERSION-$DATE.exe
fi

if [ "$1" = "linux" ]; then
    mv ChannelVault ChannelVault-$USERVER-$VERSION-$DATE
fi

echo Version is $VERSION
