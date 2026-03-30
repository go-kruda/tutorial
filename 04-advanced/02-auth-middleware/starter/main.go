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
// JWT Helpers (provided -- focus on middleware, not JWT)
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
	claims := JWTClaims{Username: username, Role: role, Exp: time.Now().Add(1 * time.Hour).Unix()}
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

// TODO: implement authMiddleware -- verify JWT token
//
// Hint: middleware is kruda.HandlerFunc = func(c *kruda.Ctx) error
//   1. Read the header: c.Header("Authorization")
//   2. Check the format "Bearer <token>"
//   3. validateToken(token)
//   4. Store claims in context: c.Set("username", claims.Username)
//   5. Pass to the next handler: return c.Next()
//   6. If error: return kruda.Unauthorized("message")
func authMiddleware() kruda.HandlerFunc {
	return func(c *kruda.Ctx) error {
		// TODO: implement
		return kruda.Unauthorized("not implemented")
	}
}

// ============================================================
// Application Entry Point
// ============================================================

func main() {
	app := kruda.New()

	// TODO: Register public routes
	//
	// Hint:
	//   kruda.Get[struct{}, MessageResponse](app, "/health", ...)
	//   kruda.Post[LoginInput, TokenResponse](app, "/login", ...)

	// TODO: Create a protected group with auth middleware
	//
	// Hint:
	//   api := app.Group("/api")
	//   api.Guard(authMiddleware())
	//   kruda.GroupGet[struct{}, ProfileResponse](api, "/profile", ...)
	//   -- In the handler, get the username: c.Get("username").(string)

	log.Println("Server starting on :3000 ...")
	log.Fatal(app.Listen(":3000"))
}
