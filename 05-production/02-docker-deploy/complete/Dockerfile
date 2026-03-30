# ============================================================
# 🐳 Multi-Stage Dockerfile for Kruda Application
# ============================================================
#
# This Dockerfile uses a multi-stage build to produce a minimal
# production image. The build stage compiles the Go binary, and
# the runtime stage copies only the binary into a tiny Alpine
# image — resulting in an image that is typically under 20 MB.
#
# Why multi-stage?
# ----------------
# 1. Security: The final image has no compiler, source code, or
#    build tools — reducing the attack surface.
# 2. Size: A Go binary + Alpine is ~15-20 MB vs ~1 GB for a
#    full golang image.
# 3. Speed: Smaller images push/pull faster in CI/CD pipelines.

# ── Stage 1: Build ────────────────────────────────────────────
# Use the official Go Alpine image to compile the application.
# The "AS builder" label lets us reference this stage later.
FROM golang:1.25-alpine AS builder

# Install git and ca-certificates (needed for fetching modules
# over HTTPS and for the final image to make TLS connections).
RUN apk add --no-cache git ca-certificates

# Set the working directory inside the build container.
WORKDIR /app

# Copy go.mod and go.sum first to leverage Docker layer caching.
# If dependencies haven't changed, Docker reuses the cached module
# download layer — speeding up subsequent builds.
COPY go.mod go.sum ./
RUN go mod download

# Copy the application source code.
COPY . .

# Build the Go binary with optimisations:
#   CGO_ENABLED=0  — produce a statically linked binary (no libc dependency)
#   -ldflags="-s -w" — strip debug info and symbol table for smaller binary
#   -o /app/server  — output path for the compiled binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/server .

# ── Stage 2: Runtime ─────────────────────────────────────────
# Use a minimal Alpine image for the final production container.
# Alpine is ~5 MB and includes a shell for debugging if needed.
# For even smaller images, you could use `scratch` instead.
FROM alpine:latest

# Install ca-certificates so the application can make outbound
# HTTPS requests (e.g., to external APIs or databases with TLS).
RUN apk add --no-cache ca-certificates

# Create a non-root user for security best practices.
# Running as root inside a container is a security risk.
RUN adduser -D -g '' appuser

# Set the working directory.
WORKDIR /app

# Copy ONLY the compiled binary from the build stage.
# This is the key benefit of multi-stage builds — nothing
# else from the build stage (source code, compiler, modules)
# makes it into the final image.
COPY --from=builder /app/server .

# Switch to the non-root user.
USER appuser

# Expose the application port.
# This is documentation only — you still need -p when running.
EXPOSE 3000

# Configure Docker health check.
# The orchestrator will probe /health every 30s to verify
# the application is alive and responding.
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget -qO- http://localhost:3000/health || exit 1

# Run the application.
# Use exec form (JSON array) so the binary receives signals
# directly — important for graceful shutdown in Docker.
CMD ["./server"]
