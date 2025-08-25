FROM ubuntu:latest

# Set non-interactive mode for apt
ENV DEBIAN_FRONTEND=noninteractive

# Update and install basic dependencies
RUN apt-get update
RUN apt-get install -y --no-install-recommends \
        ca-certificates \
        curl \
        wget \
        git \
        build-essential \
        xvfb \
        xauth \
        libgtk-3-0 \
        libwebkit2gtk-4.1-dev \
        gcc-aarch64-linux-gnu
RUN rm -rf /var/lib/apt/lists/*
