# 🤖 Section 04-05 — MCP CLI Integration: Using Kruda MCP for AI-Assisted Development

⏱️ Estimated time: **30 minutes**

Welcome to the MCP CLI Integration lesson! In this lesson you'll learn how to use the **Kruda MCP CLI** (`kruda mcp`) to connect your Kruda project with AI agents such as Claude Code, Cursor, or Windsurf -- letting AI help create and manage code in your project.

---

## Lesson Objectives

- Understand what **Model Context Protocol (MCP)** is and why it matters
- Learn that MCP in Kruda is a **CLI tool**, not a library API
- Build a normal Kruda app that AI can understand and work with
- Set up `kruda mcp` to connect with AI development tools
- Learn the MCP tools that Kruda CLI provides: `kruda_new`, `kruda_add_handler`, `kruda_list_routes`, etc.

---

## What You'll Learn

By the end of this lesson you'll be able to:

- Build a normal Kruda app with typed handlers (`kruda.Post[I, O]`, `kruda.Get[I, O]`)
- Use `kruda.WithOpenAPIInfo()` so the MCP CLI can read metadata
- Use `kruda.WithDescription()` and `kruda.WithTags()` to add information for AI to understand
- Set up `kruda mcp init` to connect with AI tools
- Understand the MCP tools that the CLI provides

---

## Prerequisites

| Tool | Version |
|---|---|
| Go | 1.25 or higher |
| Git | Latest version |
| Text Editor / IDE | VS Code, GoLand, or your preferred editor |
| Terminal | For running `go run` |
| Kruda CLI | `go install github.com/go-kruda/kruda/cmd/kruda@latest` |

> If you haven't completed Section 04-04 (SSE) yet, it's recommended to go back and do it first -- [Section 04-04 -- SSE](../04-sse/)

---

## File Structure

```
04-advanced/05-mcp-server/
├── README.md          <- You are here
├── starter/           <- Starter code (with TODOs to fill in)
│   ├── go.mod
│   └── main.go
└── complete/          <- Fully working reference implementation
    ├── go.mod
    └── main.go
```

- **[starter/](./starter/)** -- A compilable skeleton with `// TODO:` placeholders for you to fill in
- **[complete/](./complete/)** -- A fully working reference implementation you can compare against

---

## What is MCP?

**Model Context Protocol (MCP)** is an open standard that allows AI development tools to automatically discover and invoke tools from your project.

### MCP in Kruda is a CLI Tool

An important thing to understand: MCP in Kruda is **not** a library API that you import in your code. It's a **CLI tool** that helps AI agents work with your Kruda project.

```
┌──────────────────┐     MCP Protocol     ┌──────────────────┐
│   AI Agent       │ <------------------> │   kruda mcp      │
│ (Claude Code,    │                      │   (CLI tool)     │
│  Cursor, etc.)   │                      │                  │
│                  │  1. List tools       │  kruda_new       │
│                  │  2. Call tool        │  kruda_add_handler│
│                  │  3. Get result       │  kruda_list_routes│
│                  │                      │  kruda_docs      │
└──────────────────┘                      └──────────────────┘
                                                   |
                                                   v
                                          ┌──────────────────┐
                                          │  Your Kruda App  │
                                          │  (normal Go code)│
                                          └──────────────────┘
```

| Concept | Description |
|---|---|
| `kruda mcp` | CLI command that starts an MCP server for AI tools |
| `kruda mcp init` | Sets up MCP configuration in the project |
| `kruda_new` | MCP tool for AI to create a new Kruda project |
| `kruda_add_handler` | MCP tool for AI to add a new handler to the project |
| `kruda_add_resource` | MCP tool for AI to add a new resource (CRUD) |
| `kruda_list_routes` | MCP tool for AI to view all routes in the project |
| `kruda_docs` | MCP tool for AI to search Kruda documentation |

> The AI agent doesn't call your app's API directly -- it uses Kruda CLI tools to **create and modify code** in your project instead

---

## Step-by-Step Guide

### Step 1: Open the starter project

```bash
cd 04-advanced/05-mcp-server/starter
```

Open the `main.go` file -- you'll see a structure with `// TODO:` comments indicating where you need to add code.

### Step 2: Understand the Types

The starter includes pre-defined structs for request/response:

```go
type CalculatorInput struct {
    Operation string  `json:"operation" validate:"required"`
    A         float64 `json:"a"`
    B         float64 `json:"b"`
}

type CalculatorResponse struct {
    Result    float64 `json:"result"`
    Operation string  `json:"operation"`
}

type GreetingInput struct {
    Name     string `json:"name" validate:"required"`
    Language string `json:"language"`
}

type GreetingResponse struct {
    Message string `json:"message"`
}
```

> Notice that it uses `validate:"required"` instead of `mcp:"..."` -- because this is a normal REST API, not an MCP server

### Step 3: Fill in the Calculator Logic

Replace the `// TODO:` in `calculate()`:

```go
func calculate(op string, a, b float64) (float64, error) {
    switch strings.ToLower(op) {
    case "add":
        return a + b, nil
    case "subtract":
        return a - b, nil
    case "multiply":
        return a * b, nil
    case "divide":
        if b == 0 {
            return 0, fmt.Errorf("division by zero")
        }
        return a / b, nil
    default:
        return 0, fmt.Errorf("unknown operation: %s", op)
    }
}
```

### Step 4: Fill in the Greeting Logic

Replace the `// TODO:` in `greet()`:

```go
func greet(name, lang string) string {
    if lang == "" {
        lang = "en"
    }
    switch strings.ToLower(lang) {
    case "th":
        return fmt.Sprintf("สวัสดีครับ %s!", name)
    case "ja":
        return fmt.Sprintf("こんにちは %s さん！", name)
    default:
        return fmt.Sprintf("Hello, %s!", name)
    }
}
```

### Step 5: Create the Kruda App with Typed Handlers

This is the heart of the lesson! Replace the `// TODO:` in `main()`:

```go
func main() {
    // Create the app with OpenAPI metadata
    // The MCP CLI will read this metadata to understand your API
    app := kruda.New(
        kruda.WithOpenAPIInfo("Calculator & Greeting API", "1.0.0", "A simple API for Kruda MCP CLI tutorial"),
    )

    // Calculator endpoint with description for MCP CLI
    kruda.Post[CalculatorInput, CalculatorResponse](app, "/calculate", func(c *kruda.C[CalculatorInput]) (*CalculatorResponse, error) {
        result, err := calculate(c.In.Operation, c.In.A, c.In.B)
        if err != nil {
            return nil, kruda.BadRequest(err.Error())
        }
        result = math.Round(result*1e6) / 1e6
        return &CalculatorResponse{Result: result, Operation: c.In.Operation}, nil
    }, kruda.WithDescription("Perform arithmetic operations"), kruda.WithTags("Calculator"))

    // Greeting endpoint
    kruda.Post[GreetingInput, GreetingResponse](app, "/greet", func(c *kruda.C[GreetingInput]) (*GreetingResponse, error) {
        return &GreetingResponse{Message: greet(c.In.Name, c.In.Language)}, nil
    }, kruda.WithDescription("Generate a greeting"), kruda.WithTags("Greeting"))

    // Health check endpoint
    kruda.Get[struct{}, MessageResponse](app, "/health", func(c *kruda.C[struct{}]) (*MessageResponse, error) {
        return &MessageResponse{Message: "OK"}, nil
    })

    log.Fatal(app.Listen(":3000"))
}
```

> `kruda.Post[I, O]` and `kruda.Get[I, O]` are typed handler functions -- Kruda reads the type parameters to parse the request body and serialize the response automatically

> `kruda.WithDescription()` and `kruda.WithTags()` add metadata that the MCP CLI uses when AI needs to understand your API

### Step 6: Run and test

```bash
go run main.go
```

Test the calculator endpoint:

```bash
curl -X POST http://localhost:3000/calculate \
  -H "Content-Type: application/json" \
  -d '{"operation":"add","a":10,"b":25}'
```

Result:

```json
{"result":35,"operation":"add"}
```

Test the greeting endpoint:

```bash
curl -X POST http://localhost:3000/greet \
  -H "Content-Type: application/json" \
  -d '{"name":"สมชาย","language":"th"}'
```

Result:

```json
{"message":"สวัสดีครับ สมชาย!"}
```

Test the health check:

```bash
curl http://localhost:3000/health
```

### Step 7: Set up MCP CLI Integration

Now that your app is working, let's set up the MCP CLI so AI tools can work with this project:

```bash
# Install Kruda CLI (if not already installed)
go install github.com/go-kruda/kruda/cmd/kruda@latest

# Set up MCP in the project
kruda mcp init
```

The `kruda mcp init` command will create a configuration file that AI tools use to discover your project.

### Step 8: Connect with AI Tools

#### Claude Code

Add to `.claude/mcp.json`:

```json
{
  "mcpServers": {
    "kruda": {
      "command": "kruda",
      "args": ["mcp"]
    }
  }
}
```

#### Cursor

Add to `.cursor/mcp.json`:

```json
{
  "mcpServers": {
    "kruda": {
      "command": "kruda",
      "args": ["mcp"]
    }
  }
}
```

Once configured, the AI agent will be able to:

- View all routes in the project with `kruda_list_routes`
- Add new handlers with `kruda_add_handler`
- Add CRUD resources with `kruda_add_resource`
- Search documentation with `kruda_docs`

---

## Compare with complete/

If you get stuck anywhere, check the reference implementation in **[complete/](./complete/)** and compare it with your code, or use the `diff` command:

```bash
diff starter/main.go complete/main.go
```

---

## Key Takeaways

| Concept | Description |
|---|---|
| MCP is a CLI tool | `kruda mcp` is a CLI command, not a library API |
| `kruda mcp init` | Sets up MCP configuration in the project |
| `kruda.WithOpenAPIInfo()` | Adds metadata for the MCP CLI to read |
| `kruda.WithDescription()` | Describes an endpoint for AI |
| `kruda.WithTags()` | Groups endpoints for AI |
| `kruda.Post[I, O]` | Typed POST handler -- Kruda reads type params automatically |
| `kruda.Get[I, O]` | Typed GET handler -- no manual parse/serialize needed |
| `kruda.C[I]` | Typed context with `c.In` as the parsed input |
| `kruda.BadRequest()` | Creates a 400 Bad Request error response |

---

## Next Lesson

Awesome! You've learned how to build a Kruda app that works with MCP CLI for AI-assisted development. In the next lesson we'll learn about **WebSocket** -- bidirectional real-time communication.

[Section 04-06 -- WebSocket](../06-websocket/)
