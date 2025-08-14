# WIP

FROM ubuntu:24.04 AS builder

# Prevent interactive prompts during installation
ENV DEBIAN_FRONTEND=noninteractive

# 1. Install System Dependencies
# Install C/C++ toolchain (build-essential), GTK & WebKit development
# libraries, and utilities like curl and wget.
RUN apt-get update && apt-get install -y \
    build-essential \
    libgtk-3-dev \
    libwebkit2gtk-4.0-dev \
    pkg-config \
    ca-certificates \
    curl \
    wget \
    && rm -rf /var/lib/apt/lists/*

# 2. Install Go
ENV GO_VERSION=1.24.5
RUN wget "https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz" && \
    tar -C /usr/local -xzf "go${GO_VERSION}.linux-amd64.tar.gz" && \
    rm "go${GO_VERSION}.linux-amd64.tar.gz"
ENV PATH="/usr/local/go/bin:${PATH}"

# 3. Install Node.js & npm
# Wails uses this for frontend dependency management and building.
ENV NODE_VERSION=24.5
RUN curl -fsSL "https://deb.nodesource.com/setup_${NODE_VERSION}.x" | bash - && \
    apt-get install -y nodejs && \
    rm -rf /var/lib/apt/lists/*

# 4. Install Wails CLI
# This command fetches and installs the Wails command-line tool.
RUN go install github.com/wailsapp/wails/v2/cmd/wails@latest
# The 'wails' executable is placed in /root/go/bin by default. Add it to the PATH.
ENV PATH="/root/go/bin:${PATH}"

# 5. Build Your Application
WORKDIR /app

# Argument for your application's binary name
ARG APP_NAME=myapp
ENV APP_NAME=${APP_NAME}

# Copy dependency files first to leverage Docker's cache
COPY go.mod go.sum ./
COPY wails.json ./

# Download Go dependencies
RUN go mod download

# Copy the rest of your application's source code
COPY . .

# Run Go unit tests before building
RUN echo "--- Running Go Tests ---" && go test ./...

# Finally, build the Wails application binary
RUN echo "--- Building Wails App for linux/amd64 ---" && \
    wails build -platform linux/amd64 -o $APP_NAME
