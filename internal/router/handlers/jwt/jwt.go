package jwt

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

// Конфигурация JWT
type Config struct {
	SecretKey     string
	AccessExpiry  time.Duration
	RefreshExpiry time.Duration
}

// Модели (используем ваши существующие)
type User struct {
	ID       int    `db:"id_user"`
	Password string `db:"password"`
}

type RefreshToken struct {
	ID        int       `db:"id"`
	Token     string    `db:"token"`
	UserID    int       `db:"id_user"`
	ExpiresAt time.Time `db:"expires_at"`
}

// Сервис для работы с JWT
type Service struct {
	db     *sqlx.DB
	config Config
}

// Claims для JWT токена
type Claims struct {
	UserID int `json:"user_id"`
	jwt.StandardClaims
}

// Инициализация сервиса
func NewService(db *sqlx.DB) *Service {
	return &Service{
		db: db,
		config: Config{
			SecretKey:     os.Getenv("JWT_SECRET"),
			AccessExpiry:  15 * time.Minute,
			RefreshExpiry: 7 * 24 * time.Hour,
		},
	}
}

// Обработчик входа
func (s *Service) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request format")
		return
	}

	// Получаем пользователя из БД
	var user User
	err := s.db.Get(&user, "SELECT id_user, password FROM Users WHERE login = $1", request.Login)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondError(w, http.StatusUnauthorized, "Invalid credentials")
		} else {
			respondError(w, http.StatusInternalServerError, "Database error")
		}
		return
	}

	// Проверяем пароль
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		respondError(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	// Генерируем токены
	accessToken, err := s.generateAccessToken(user.ID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	refreshToken, err := s.generateRefreshToken()
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to generate refresh token")
		return
	}

	// Сохраняем refresh токен
	_, err = s.db.Exec(`
		INSERT INTO refresh_tokens (token, id_user, expires_at) 
		VALUES ($1, $2, $3)`,
		refreshToken, user.ID, time.Now().Add(s.config.RefreshExpiry))
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to save refresh token")
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"expires_in":    int(s.config.AccessExpiry.Seconds()),
	})
}

// Обработчик обновления токена
func (s *Service) RefreshHandler(w http.ResponseWriter, r *http.Request) {
	var request struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request format")
		return
	}

	// Проверяем refresh токен
	var token RefreshToken
	err := s.db.Get(&token, `
		SELECT id, id_user, expires_at 
		FROM refresh_tokens 
		WHERE token = $1 AND expires_at > NOW()`,
		request.RefreshToken)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "Invalid or expired refresh token")
		return
	}

	// Генерируем новые токены
	newAccessToken, err := s.generateAccessToken(token.UserID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	newRefreshToken, err := s.generateRefreshToken()
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to generate refresh token")
		return
	}

	// Обновляем refresh токен
	_, err = s.db.Exec(`
		UPDATE refresh_tokens 
		SET token = $1, expires_at = $2 
		WHERE id = $3`,
		newRefreshToken, time.Now().Add(s.config.RefreshExpiry), token.ID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to update refresh token")
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"access_token":  newAccessToken,
		"refresh_token": newRefreshToken,
		"expires_in":    int(s.config.AccessExpiry.Seconds()),
	})
}

// Обработчик выхода
func (s *Service) RevokeHandler(w http.ResponseWriter, r *http.Request) {
	var request struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request format")
		return
	}

	// Удаляем refresh токен
	result, err := s.db.Exec("DELETE FROM refresh_tokens WHERE token = $1", request.RefreshToken)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to revoke token")
		return
	}

	if rows, _ := result.RowsAffected(); rows == 0 {
		respondError(w, http.StatusNotFound, "Token not found")
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"status": "success"})
}

// Middleware для аутентификации
func (s *Service) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			respondError(w, http.StatusUnauthorized, "Authorization header required")
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims := &Claims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(s.config.SecretKey), nil
		})

		if err != nil || !token.Valid {
			respondError(w, http.StatusUnauthorized, "Invalid token")
			return
		}

		// Добавляем user_id в контекст
		ctx := context.WithValue(r.Context(), "user_id", claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Генерация access токена
func (s *Service) generateAccessToken(userID int) (string, error) {
	claims := &Claims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(s.config.AccessExpiry).Unix(),
			Issuer:    "your-app-name",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.config.SecretKey))
}

// Генерация refresh токена
func (s *Service) generateRefreshToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// Вспомогательные функции для ответов
func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, map[string]string{"error": message})
}