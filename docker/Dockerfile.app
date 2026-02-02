# Build stage
FROM golang:1.24-bookworm AS builder

# 通过构建参数接收敏感信息
ARG GOPRIVATE_ARG
ARG GOPROXY_ARG=https://goproxy.cn,direct
ARG GOSUMDB_ARG=off
ARG APK_MIRROR_ARG

# 设置Go环境变量
ENV GOPRIVATE=${GOPRIVATE_ARG}
ENV GOPROXY=${GOPROXY_ARG}
ENV GOSUMDB=${GOSUMDB_ARG}

WORKDIR /app

# Install dependencies - 合并 RUN 减少层数
RUN --mount=type=cache,target=/var/cache/apt,sharing=locked \
    --mount=type=cache,target=/var/lib/apt,sharing=locked \
    if [ -n "$APK_MIRROR_ARG" ]; then \
        sed -i "s@deb.debian.org@${APK_MIRROR_ARG}@g" /etc/apt/sources.list.d/debian.sources; \
    fi && \
    apt-get update && \
    apt-get install -y --no-install-recommends git build-essential

# Copy go mod and sum files first for better caching
COPY go.mod go.sum ./

# Download dependencies and install migrate tool in parallel
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go mod download && \
    go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

COPY . .

# Get version and commit info for build injection
ARG VERSION_ARG
ARG COMMIT_ID_ARG
ARG BUILD_TIME_ARG
ARG GO_VERSION_ARG

# Set build-time variables
ENV VERSION=${VERSION_ARG}
ENV COMMIT_ID=${COMMIT_ID_ARG}
ENV BUILD_TIME=${BUILD_TIME_ARG}
ENV GO_VERSION=${GO_VERSION_ARG}

# Build the application with version info - 合并构建步骤
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    make download_spatial && \
    make build-prod && \
    cp -r /go/pkg/mod/github.com/yanyiwu/ /app/yanyiwu/

# Final stage
FROM debian:12.12-slim

ARG APK_MIRROR_ARG

WORKDIR /app

# Create a non-root user and install all dependencies in one layer
RUN useradd -m -s /bin/bash appuser && \
    if [ -n "$APK_MIRROR_ARG" ]; then \
        sed -i "s@deb.debian.org@${APK_MIRROR_ARG}@g" /etc/apt/sources.list.d/debian.sources; \
    fi && \
    apt-get update && \
    apt-get install -y --no-install-recommends \
        build-essential \
        postgresql-client \
        default-mysql-client \
        ca-certificates \
        tzdata \
        curl \
        bash \
        python3 \
        python3-pip \
        python3-venv \
        libffi-dev \
        libssl-dev \
        nodejs \
        npm && \
    # Install uv directly via pip (faster than curl script)
    python3 -m pip install --break-system-packages uv && \
    ln -sf /usr/local/bin/uvx /usr/local/bin/uvx 2>/dev/null || true && \
    # Cleanup
    apt-get clean && \
    rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/* && \
    mkdir -p /home/appuser/.local/bin && \
    chown -R appuser:appuser /home/appuser

# Create data directories and set permissions
RUN mkdir -p /data/files && \
    chown -R appuser:appuser /app /data/files

# Copy migrate tool from builder stage
COPY --from=builder /go/bin/migrate /usr/local/bin/
COPY --from=builder /app/yanyiwu/ /go/pkg/mod/github.com/yanyiwu/

# Copy the binary from the builder stage
COPY --from=builder /app/config ./config
COPY --from=builder /app/scripts ./scripts
COPY --from=builder /app/migrations ./migrations
COPY --from=builder /app/dataset/samples ./dataset/samples
COPY --from=builder /root/.duckdb /home/appuser/.duckdb
COPY --from=builder /app/WeKnora .

# Make scripts executable
RUN chmod +x ./scripts/*.sh

# Expose ports
EXPOSE 8080

# Switch to non-root user and run the application directly
USER appuser

CMD ["./WeKnora"]
