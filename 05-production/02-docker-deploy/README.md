# 🐳 Section 05-02 — Docker Deploy: Containerise Kruda Application

⏱️ Estimated time: **30 minutes**

Welcome to the Docker Deploy lesson! In this lesson you will learn how to containerise a Kruda application with a **multi-stage Dockerfile** -- a technique that produces small, secure Docker images ready for production.

---

## Target Learning Outcomes

- Understand the concept of **multi-stage Docker build** and why it matters for Go applications
- Write a Dockerfile that separates build stage from runtime stage
- Create a Docker image under 20 MB for production
- Add a **health check endpoint** for container orchestration
- Use environment variables for configuration following 12-factor app methodology
- Run a Kruda application in a Docker container

---

## What You Will Learn

By the end of this lesson you will be able to:

- Write a multi-stage Dockerfile for a Go application
- Use Docker layer caching to speed up builds
- Create a statically linked binary with `CGO_ENABLED=0`
- Set up `HEALTHCHECK` in the Dockerfile
- Run a container with a non-root user for security
- Pass configuration via environment variables (`-e PORT=8080`)
- Map ports between host and container (`-p 8080:3000`)

---

## Prerequisites

| Tool | Version |
|---|---|
| Go | 1.25+ |
| Docker | Latest |
| Git | Latest |
| Text Editor / IDE | VS Code, GoLand, or your preferred editor |

> If you haven't completed Section 05-01, consider going back first -- see [Section 05-01 -- Monitoring](../01-monitoring/)

---

## File Structure

```
05-production/02-docker-deploy/
|-- README.md              <-- You are here
|-- starter/               <-- Starter code (with TODOs to fill in)
|   |-- go.mod
|   +-- main.go
+-- complete/              <-- Complete solution
    |-- go.mod
    |-- main.go
    +-- Dockerfile         <-- Multi-stage Dockerfile
```

- **[starter/](./starter/)** -- Skeleton code that compiles but has `// TODO:` markers for you to complete
- **[complete/](./complete/)** -- Full working solution with Dockerfile for production

---

## Why Multi-Stage Build?

| Single Stage | Multi-Stage Build |
|---|---|
| Image ~1 GB (includes Go compiler) | Image ~15-20 MB (binary only) |
| Source code in image | No source code in final image |
| Unnecessary build tools included | Only runtime essentials |
| Wide attack surface | Narrow attack surface, more secure |

> Multi-stage build separates **build** (compile Go binary) from **runtime** (run binary) -- the final image contains only the compiled binary

---

## Step-by-Step Guide

### Step 1: Open the starter project

```bash
cd 05-production/02-docker-deploy/starter
```

Open `main.go` -- you will see the skeleton with `// TODO:` comments.

### Step 2: Add Health Check Endpoint

Health checks are essential for container orchestration -- Docker, Kubernetes, and load balancers use this endpoint to verify the application is alive:

```go
kruda.Get[struct{}, HealthResponse](app, "/health", func(c *kruda.C[struct{}]) (*HealthResponse, error) {
    return &HealthResponse{
        Status:    "ok",
        Timestamp: time.Now().UTC().Format(time.RFC3339),
        Version:   getEnv("APP_VERSION", "1.0.0"),
    }, nil
})
```

> Health endpoints should return quickly with no side effects -- the orchestrator calls them every 30 seconds

### Step 3: Read Configuration from Environment Variables

In Docker containers, configuration is passed via environment variables following the 12-factor app methodology:

```go
func getEnv(key, fallback string) string {
    if val, ok := os.LookupEnv(key); ok {
        return val
    }
    return fallback
}

port := getEnv("PORT", "3000")
```

> Use `docker run -e PORT=8080` to change the port without modifying code

### Step 4: Create the App and Register Routes

```go
func main() {
    port := getEnv("PORT", "3000")
    app := kruda.New()

    kruda.Get[struct{}, HealthResponse](app, "/health", func(c *kruda.C[struct{}]) (*HealthResponse, error) {
        return &HealthResponse{
            Status:    "ok",
            Timestamp: time.Now().UTC().Format(time.RFC3339),
            Version:   getEnv("APP_VERSION", "1.0.0"),
        }, nil
    })

    kruda.Get[struct{}, MessageResponse](app, "/hello", func(c *kruda.C[struct{}]) (*MessageResponse, error) {
        hostname, _ := os.Hostname()
        return &MessageResponse{
            Message: fmt.Sprintf("Hello from Docker! (host: %s)", hostname),
        }, nil
    })

    addr := fmt.Sprintf(":%s", port)
    log.Fatal(app.Listen(addr))
}
```

### Step 5: Write the Multi-Stage Dockerfile

Create a `Dockerfile` in the same directory as `main.go`:

```dockerfile
# -- Stage 1: Build -----------------------------------------
FROM golang:1.25-alpine AS builder

RUN apk add --no-cache git ca-certificates
WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/server .

# -- Stage 2: Runtime ----------------------------------------
FROM alpine:latest

RUN apk add --no-cache ca-certificates
RUN adduser -D -g '' appuser

WORKDIR /app
COPY --from=builder /app/server .

USER appuser
EXPOSE 3000

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget -qO- http://localhost:3000/health || exit 1

CMD ["./server"]
```

> Notice two `FROM` statements -- this is the core of multi-stage build. The first stage compiles the code, the second stage copies only the resulting binary.

### Step 6: Build Docker Image

```bash
docker build -t kruda-app .
```

Check image size:

```bash
docker images kruda-app
```

You will see the image is only ~15-20 MB!

### Step 7: Run the Container

```bash
docker run -d --name kruda-app -p 8080:3000 kruda-app
```

Test:

```bash
# Health check
curl http://localhost:8080/health

# Hello endpoint
curl http://localhost:8080/hello
```

### Step 8: Pass Configuration via Environment Variables

```bash
# Change port and version
docker run -d --name kruda-custom \
  -p 9090:8080 \
  -e PORT=8080 \
  -e APP_VERSION=2.0.0 \
  kruda-app
```

### Step 9: Verify Health Check

```bash
# Check health status of the container
docker inspect --format='{{.State.Health.Status}}' kruda-app

# View logs
docker logs kruda-app
```

---

## Compare with complete/

If you get stuck, check the solution in **[complete/](./complete/)** and compare:

```bash
diff starter/main.go complete/main.go
```

---

## Key Concepts Summary

| Concept | Description |
|---|---|
| Multi-Stage Build | Separate build stage from runtime stage for small images |
| `FROM ... AS builder` | Name the build stage to reference in later stages |
| `COPY --from=builder` | Copy files from the build stage into the runtime stage |
| `CGO_ENABLED=0` | Create a statically linked binary that does not depend on libc |
| `-ldflags="-s -w"` | Strip debug info to reduce binary size |
| `HEALTHCHECK` | Define health check for container orchestration |
| `USER appuser` | Run the process as non-root user for security |
| Environment Variables | Pass configuration via `-e` flag following 12-factor app |

---

## Deployment Tips

### Production Checklist

- Use multi-stage build to reduce image size
- Run container with non-root user
- Add `HEALTHCHECK` to Dockerfile
- Use environment variables for configuration
- Set resource limits (`--memory`, `--cpus`)
- Use `.dockerignore` to exclude unnecessary files

### Docker Compose Example

```yaml
version: "3.8"
services:
  kruda-app:
    build: .
    ports:
      - "8080:3000"
    environment:
      - PORT=3000
      - APP_VERSION=1.0.0
    restart: unless-stopped
    deploy:
      resources:
        limits:
          memory: 128M
          cpus: "0.5"
```

---

## Next Lesson

Great work! You have learned how to containerise a Kruda application with a multi-stage Dockerfile. In the next lesson you will learn **Benchmark** -- how to measure Kruda application performance and compare with other frameworks.

--> [Section 05-03 -- Benchmark](../03-benchmark/)
