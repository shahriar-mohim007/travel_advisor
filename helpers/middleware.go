package helpers

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
	"travel_advisor/pkg/config"
)

func JWTAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			resp := &Response{
				Status:  http.StatusUnauthorized,
				Message: "Authorization header required",
			}
			resp.Render(w)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			resp := &Response{
				Status:  http.StatusUnauthorized,
				Message: "Invalid authorization header format",
			}
			resp.Render(w)
			return
		}

		token := parts[1]
		userID, err := validateJWTToken(token)
		if err != nil {
			resp := &Response{
				Status:  http.StatusUnauthorized,
				Message: "Invalid or expired token",
				Error:   err.Error(),
			}
			resp.Render(w)
			return
		}
		ctx := r.Context()
		ctx = context.WithValue(ctx, "user_id", userID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func validateJWTToken(token string) (uint, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return 0, fmt.Errorf("invalid token format")
	}

	message := parts[0] + "." + parts[1]
	expectedSignature := createHMACSignature(message, config.App().JwtSecret)
	expectedSignatureEncoded := base64URLEncode(expectedSignature)

	if parts[2] != expectedSignatureEncoded {
		return 0, fmt.Errorf("invalid signature")
	}

	payloadBytes, err := base64URLDecode(parts[1])
	if err != nil {
		return 0, fmt.Errorf("failed to decode payload: %v", err)
	}

	var payload map[string]interface{}
	if err := json.Unmarshal(payloadBytes, &payload); err != nil {
		return 0, fmt.Errorf("failed to parse payload: %v", err)
	}

	if exp, ok := payload["exp"].(float64); ok {
		if time.Now().Unix() > int64(exp) {
			return 0, fmt.Errorf("token expired")
		}
	}

	if userID, ok := payload["user_id"].(float64); ok {
		return uint(userID), nil
	}

	return 0, fmt.Errorf("invalid user ID in token")
}

func GenerateJWTToken(userID uint) (string, error) {

	header := `{"alg":"HS256","typ":"JWT"}`
	headerEncoded := base64URLEncode([]byte(header))

	payload := fmt.Sprintf(`{"user_id":%d,"exp":%d}`, userID, time.Now().Add(24*time.Hour).Unix())
	payloadEncoded := base64URLEncode([]byte(payload))

	message := headerEncoded + "." + payloadEncoded
	signature := createHMACSignature(message, config.App().JwtSecret)
	signatureEncoded := base64URLEncode(signature)

	return message + "." + signatureEncoded, nil
}

func base64URLEncode(data []byte) string {
	encoded := hex.EncodeToString(data)
	return strings.TrimRight(strings.ReplaceAll(strings.ReplaceAll(encoded, "+", "-"), "/", "_"), "=")
}

func createHMACSignature(message, secret string) []byte {
	hash := sha256.New()
	hash.Write([]byte(message + secret))
	return hash.Sum(nil)
}

func base64URLDecode(data string) ([]byte, error) {
	padding := 4 - len(data)%4
	if padding != 4 {
		data += strings.Repeat("=", padding)
	}

	data = strings.ReplaceAll(data, "-", "+")
	data = strings.ReplaceAll(data, "_", "/")

	return hex.DecodeString(data)
}
