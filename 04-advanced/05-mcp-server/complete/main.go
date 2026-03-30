package main

import (
	"fmt"
	"log"
	"math"
	"strings"

	"github.com/go-kruda/kruda"
)

// ============================================================
// Request / Response Types
// ============================================================
//
// This is a normal Kruda REST API. The MCP integration in Kruda
// is a CLI tool (`kruda mcp`) that helps AI agents discover and
// work with your project -- it is NOT a library API embedded in
// your application code.

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

type MessageResponse struct {
	Message string `json:"message"`
}

// ============================================================
// Business Logic
// ============================================================

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
		return 0, fmt.Errorf("unknown operation: %s (use add, subtract, multiply, divide)", op)
	}
}

func greet(name, lang string) string {
	if lang == "" {
		lang = "en"
	}
	switch strings.ToLower(lang) {
	case "th":
		return fmt.Sprintf("สวัสดีครับ %s! ยินดีต้อนรับสู่ Kruda", name)
	case "ja":
		return fmt.Sprintf("こんにちは %s さん！Kruda へようこそ。", name)
	default:
		return fmt.Sprintf("Hello, %s! Welcome to Kruda.", name)
	}
}

// ============================================================
// Application Entry Point
// ============================================================

func main() {
	app := kruda.New(
		kruda.WithOpenAPIInfo("Calculator & Greeting API", "1.0.0",
			"A simple API for the Kruda MCP CLI tutorial"),
	)

	kruda.Post[CalculatorInput, CalculatorResponse](app, "/calculate",
		func(c *kruda.C[CalculatorInput]) (*CalculatorResponse, error) {
			result, err := calculate(c.In.Operation, c.In.A, c.In.B)
			if err != nil {
				return nil, kruda.BadRequest(err.Error())
			}
			result = math.Round(result*1e6) / 1e6
			return &CalculatorResponse{
				Result:    result,
				Operation: c.In.Operation,
			}, nil
		},
		kruda.WithDescription("Perform arithmetic operations"),
		kruda.WithTags("Calculator"),
	)

	kruda.Post[GreetingInput, GreetingResponse](app, "/greet",
		func(c *kruda.C[GreetingInput]) (*GreetingResponse, error) {
			return &GreetingResponse{
				Message: greet(c.In.Name, c.In.Language),
			}, nil
		},
		kruda.WithDescription("Generate a personalised greeting"),
		kruda.WithTags("Greeting"),
	)

	kruda.Get[struct{}, MessageResponse](app, "/health",
		func(c *kruda.C[struct{}]) (*MessageResponse, error) {
			return &MessageResponse{Message: "OK"}, nil
		},
	)

	log.Println("Server starting on :3000 ...")
	log.Fatal(app.Listen(":3000"))
}
