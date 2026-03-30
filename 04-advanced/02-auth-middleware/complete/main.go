package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/go-kruda/kruda"
)

// ============================================================
// Request / Response Types
// ============================================================

type LoginInput struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type TokenResponse struct {
	Token     string `json:"token"`
	ExpiresIn int    `json:"expires_in"`
}

type ProfileResponse struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	Message  string `json:"message"`
}

type MessageResponse struct {
	Message string `json:"message"`
}

// ============================================================
// JWT Helpers (simplified for tutorial purposes)
// ============================================================

var jwtSecret = []byte("kruda-tutorial-secret-key")

type JWTClaims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	Exp      int64  `json:"exp"`
}

func generateToken(username, role string) (string, error) {
	header := map[string]string{"alg": "HS256", "typ": "JWT"}
	headerJSON, _ := json.Marshal(header)
	headerB64 := base64.RawURLEncoding.EncodeToString(headerJSON)

	claims := JWTClaims{
		Username: username,
		Role:     role,
		Exp:      time.Now().Add(1 * time.Hour).Unix(),
	}
	claimsJSON, _ := json.Marshal(claims)
	claimsB64 := base64.RawURLEncoding.EncodeToString(claimsJSON)

	signingInput := headerB64 + "." + claimsB64
	mac := hmac.New(sha256.New, jwtSecret)
	mac.Write([]byte(signingInput))
	signature := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))

	return signingInput + "." + signature, nil
}

func validateToken(tokenStr string) (*JWTClaims, error) {
	parts := strings.Split(tokenStr, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid token format")
	}

	signingInput := parts[0] + "." + parts[1]
	mac := hmac.New(sha256.New, jwtSecret)
	mac.Write([]byte(signingInput))
	expectedSig := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))

	if !hmac.Equal([]byte(parts[2]), []byte(expectedSig)) {
		return nil, fmt.Errorf("invalid token signature")
	}

	claimsJSON, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, fmt.Errorf("decode claims: %w", err)
	}

	var claims JWTClaims
	if err := json.Unmarshal(claimsJSON, &claims); err != nil {
		return nil, fmt.Errorf("unmarshal claims: %w", err)
	}

	if time.Now().Unix() > claims.Exp {
		return nil, fmt.Errorf("token expired")
	}

	return &claims, nil
}

// ============================================================
// Auth Middleware
// ============================================================
//
// A Kruda middleware is a HandlerFunc:
//
//   func(c *kruda.Ctx) error
//
// Call c.Next() to pass control to the next middleware or handler.
// Return an error or respond directly to stop the chain.

func authMiddleware() kruda.HandlerFunc {
	return func(c *kruda.Ctx) error {
		// Read the Authorization header.
		authHeader := c.Header("Authorization")
		if authHeader == "" {
			return kruda.Unauthorized("missing Authorization header")
		}

		// Expect "Bearer <token>" format.
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			return kruda.Unauthorized("invalid Authorization format (expected: Bearer <token>)")
		}

		// Validate the token.
		claims, err := validateToken(parts[1])
		if err != nil {
			return kruda.Unauthorized(fmt.Sprintf("authentication failed: %v", err))
		}

		// Store claims in context locals so downstream handlers
		// can access the authenticated user's information.
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)

		// Pass control to the next middleware or handler.
		return c.Next()
	}
}

// ============================================================
// Application Entry Point
// ============================================================

func main() {
	app := kruda.New()

	// ── Public Routes ────────────────────────────────────────
	kruda.Get[struct{}, MessageResponse](app, "/health", func(c *kruda.C[struct{}]) (*MessageResponse, error) {
		return &MessageResponse{Message: "OK"}, nil
	})

	kruda.Post[LoginInput, TokenResponse](app, "/login", func(c *kruda.C[LoginInput]) (*TokenResponse, error) {
		if c.In.Username == "admin" && c.In.Password == "secret" {
			token, err := generateToken(c.In.Username, "admin")
			if err != nil {
				return nil, kruda.InternalError("generate token failed")
			}
			return &TokenResponse{Token: token, ExpiresIn: 3600}, nil
		}
		return nil, kruda.Unauthorized("invalid username or password")
	})

	// ── Protected Routes (with Auth Middleware) ──────────────
	api := app.Group("/api")
	api.Guard(authMiddleware())

	// Routes under /api require a valid JWT token.
	kruda.GroupGet[struct{}, ProfileResponse](api, "/profile", func(c *kruda.C[struct{}]) (*ProfileResponse, error) {
		username := c.Get("username").(string)
		role := c.Get("role").(string)
		return &ProfileResponse{
			Username: username,
			Role:     role,
			Message:  fmt.Sprintf("Welcome back, %s! You have %s access.", username, role),
		}, nil
	})

	log.Println("Server starting on :3000 ...")
	log.Println("  Public:    GET /health, POST /login")
	log.Println("  Protected: GET /api/profile (requires Bearer token)")
	log.Fatal(app.Listen(":3000"))
}
