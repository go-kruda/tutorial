package main

import (
	"fmt"
	"log"
	"math"
	"strings"

	"github.com/go-kruda/kruda"
)

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

func main() {
	app := kruda.New()

	// TODO: Register POST /calculate
	// TODO: Register POST /greet
	// TODO: Register GET /health

	_ = math.Round
	log.Println("Server starting on :3000 ...")
	log.Fatal(app.Listen(":3000"))
}
