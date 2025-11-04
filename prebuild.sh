#!/bin/bash

VERSION=$(git rev-parse --short HEAD)

# Check if the first command-line argument is "macOS"
if [ "$1" = "macOS" ]; then
    #mv ChannelVault.app wails.app
    echo macOS build...
fi

if [ "$1" = "win" ]; then
    #mv ChannelVault.app wails.app
    echo win build...
    rm -f ChannelVault*exe
fi

if [ "$1" = "linux" ]; then
    #mv ChannelVault.app wails.app
    echo linux build...
    rm -f ChannelVault 
    rm -f ChannelVault-?.?-???????-??????????????
fi

echo Version is $VERSION
