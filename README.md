# 🚀 go-kruda/tutorial — Learn Kruda from Zero to Production

> A progressive, hands-on tutorial course for the **Kruda** Go web framework.
> Build real-world APIs step by step — from your first route to production deployment. ✅

---

## 🎯 Course Overview

Welcome to the official **Kruda Tutorial** repository! This is a structured, section-based learning resource that takes you from absolute beginner to production-ready Kruda developer.

Every section ships with two independent Go modules:

- **`starter/`** — A compiling skeleton with `// TODO:` placeholders. Your starting point.
- **`complete/`** — A fully working reference implementation with explanatory comments.

Compare them, learn from the diff, and build your skills one section at a time. Let's go! 🎉

---

## ✨ Kruda Feature Highlights

Kruda is a modern Go web framework packed with powerful, developer-friendly features:

| Feature | Description |
|---------|-------------|
| 🛩️ **Wing Transport** | Custom transport layer using epoll + eventfd, optimised for Linux — blazing fast throughput and low latency |
| 🔒 **Typed Handler** | Generic-based handlers for type-safe request/response — no more manual parsing or casting |
| ⚡ **Auto CRUD** | Automatic CRUD endpoint generation from model structs — eliminate boilerplate in seconds |
| 📦 **DI Container** | Built-in dependency injection container — no external libraries needed |
| 📡 **SSE** (Server-Sent Events) | Built-in push-notification mechanism over HTTP — real-time updates made simple |
| 🤖 **MCP Server** | Model Context Protocol server support — integrate AI tooling directly into your app |
| 📄 **OpenAPI Generator** | Automatic OpenAPI spec generation driven by Typed Handlers — always up-to-date docs |
| 🔗 **Middleware Chain** | Standard middleware composition pattern — auth, logging, CORS, and more |

---

## 📚 Learning Progression

Follow the sections in order for the best experience. Each builds on the previous one.

| # | Section | Topics | ⏱️ Est. Time |
|---|---------|--------|-------------|
| 00 | **Why Kruda** | Framework comparison, benchmarks | 10 min |
| 01 | **Beginner** | REST API basics, Typed Handler | 30 min |
| 02 | **Auto CRUD** | Model-driven CRUD generation | 30 min |
| 03 | **Intermediate** | DB integration, Docker, config, error handling | 45 min |
| 04 | **Advanced** | DI Container, Auth Middleware, OpenAPI, SSE, MCP Server, WebSocket, Testing, Architecture | 2–3 hrs |
| 05 | **Production** | Monitoring, Docker deploy, benchmarking | 1–2 hrs |

> 💡 **Total estimated time:** ~5–6 hours for the full course.

---

## 🧠 Skills Matrix

Each section maps to concrete skills you'll gain:

| Skill | 00 | 01 | 02 | 03 | 04 | 05 |
|-------|:--:|:--:|:--:|:--:|:--:|:--:|
| Framework evaluation | ✅ | | | | | |
| REST routing | | ✅ | | | | |
| Typed Handler (generics) | | ✅ | | | ✅ | |
| JSON request/response | | ✅ | ✅ | | | |
| Auto CRUD generation | | | ✅ | | | |
| Database integration | | | | ✅ | | |
| Docker Compose | | | | ✅ | | ✅ |
| Config management | | | | ✅ | | |
| Error handling | | | | ✅ | | |
| DI Container | | | | | ✅ | |
| Auth Middleware Chain | | | | | ✅ | |
| OpenAPI generation | | | | | ✅ | |
| SSE (Server-Sent Events) | | | | | ✅ | |
| MCP Server | | | | | ✅ | |
| WebSocket | | | | | ✅ | |
| Unit testing | | | | | ✅ | |
| Clean architecture | | | | | ✅ | |
| Prometheus monitoring | | | | | | ✅ |
| Docker deployment | | | | | | ✅ |
| Benchmarking | | | | | | ✅ |

---

## 📋 Prerequisites

Before you start, make sure you have these installed:

- ✅ **Go 1.25+** — [Download Go](https://go.dev/dl/)
- ✅ **Docker** — [Install Docker](https://docs.docker.com/get-docker/) (needed for sections 03 and 05)
- ✅ **Git** — [Install Git](https://git-scm.com/downloads)

---

## ⚡ Quick Start

Get up and running in seconds:

```bash
# Clone the repository
git clone https://github.com/go-kruda/tutorial.git
cd tutorial

# Run the first beginner example
cd 01-beginner/complete
go run main.go
```

🎉 That's it! Your first Kruda server is running. Now open the section README and start learning!

---

## 📁 Directory Structure

```
go-kruda/tutorial/
├── README.md
├── go.work
├── .gitignore
├── 00-why-kruda/
│   └── README.md
├── 01-beginner/
│   ├── README.md
│   ├── starter/
│   │   ├── go.mod
│   │   └── main.go
│   └── complete/
│       ├── go.mod
│       └── main.go
├── 02-auto-crud/
│   ├── README.md
│   ├── starter/
│   │   ├── go.mod
│   │   └── main.go
│   └── complete/
│       ├── go.mod
│       └── main.go
├── 03-intermediate/
│   ├── README.md
│   ├── docker-compose.yml
│   ├── starter/
│   │   ├── go.mod
│   │   └── main.go
│   └── complete/
│       ├── go.mod
│       └── main.go
├── 04-advanced/
│   ├── 01-di-container/
│   │   ├── README.md
│   │   ├── starter/ ...
│   │   └── complete/ ...
│   ├── 02-auth-middleware/ ...
│   ├── 03-openapi/ ...
│   ├── 04-sse/ ...
│   ├── 05-mcp-server/ ...
│   ├── 06-websocket/ ...
│   ├── 07-testing/
│   │   ├── README.md
│   │   ├── starter/ ...
│   │   └── complete/
│   │       ├── go.mod
│   │       ├── main.go
│   │       └── handler_test.go
│   └── 08-architecture/
│       ├── README.md
│       ├── starter/ ...
│       └── complete/
│           ├── go.mod
│           ├── main.go
│           ├── handler/
│           ├── service/
│           └── repository/
└── 05-production/
    ├── 01-monitoring/
    │   ├── README.md
    │   ├── starter/ ...
    │   └── complete/
    │       ├── go.mod
    │       ├── main.go
    │       └── dashboard.json
    ├── 02-docker-deploy/
    │   ├── README.md
    │   ├── starter/ ...
    │   └── complete/
    │       ├── go.mod
    │       ├── main.go
    │       └── Dockerfile
    └── 03-benchmark/
        ├── README.md
        ├── starter/ ...
        └── complete/
            ├── go.mod
            ├── main.go
            └── benchmark_test.go
```

---

## 🤝 Contributing

Found a typo? Have a suggestion? PRs and issues are welcome! 🙏

---

## 📜 License

This project is licensed under the terms specified in the [LICENSE](./LICENSE) file.

---

> 🚀 **Ready to start?** Head to [00-why-kruda/](./00-why-kruda/) to see why Kruda stands out, then jump into [01-beginner/](./01-beginner/) to write your first API!

Happy coding! 🎯✨
