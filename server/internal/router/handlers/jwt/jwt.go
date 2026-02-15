package jwt

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"main/internal/repositories"
	"main/internal/repositories/models"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
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
	config Config
	logger *logrus.Logger
	repo   repositories.Repository
}

// Claims для JWT токена
type Claims struct {
	UserID int `json:"user_id"`
	jwt.StandardClaims
}

// Инициализация сервиса
func NewService(log *logrus.Logger, r repositories.Repository) *Service {
	return &Service{
		config: Config{
			SecretKey:     os.Getenv("JWT_SECRET"),
			AccessExpiry:  15 * time.Minute,
			RefreshExpiry: 7 * 24 * time.Hour,
		},
		logger: log,
		repo:   r,
	}
}

// Обработчик входа
func (s *Service) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}
	contentType := r.Header.Get("Content-Type")

	if contentType == "application/json" {
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			respondError(w, http.StatusBadRequest, "Failed decode json")
			return
		}
	} else if strings.Contains(contentType, "multipart/form-data") {
		if err := r.ParseMultipartForm(10 << 20); err != nil {
			respondError(w, http.StatusBadRequest, "Failed form-parsing")
			return
		}
		if request.Login = r.PostForm.Get("login"); request.Login == "" {
			respondError(w, http.StatusBadRequest, "Empety data")
			return
		}
		if request.Password = r.PostForm.Get("password"); request.Password == "" {
			respondError(w, http.StatusBadRequest, "Empety data")
			return
		}
	} else {
		respondError(w, http.StatusBadRequest, "Invalid request format")
		return
	}

	user, err := s.repo.Users().GetByLogin(r.Context(), request.Login)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondError(w, http.StatusForbidden, "Invalid credentials")
		} else {
			respondError(w, http.StatusInternalServerError, "Database error")
		}
		return
	}

	// Проверяем пароль
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		respondError(w, http.StatusForbidden, "Invalid credentials")
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
	tx, err := s.repo.BeginTx(r.Context())
	if err != nil {
		s.logger.Errorf("Failed begin tx(log login) by user:%d:err:%v", user.ID, err)
	}
	defer tx.Rollback()
	// Сохраняем refresh токен
	err = tx.Refresh().Create(r.Context(), &models.RefreshToken{Token: refreshToken, UserID: user.ID, ExpiresAt: time.Now().Add(s.config.RefreshExpiry)})
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to save refresh token")
		return
	}

	_, err = tx.Log().Create(r.Context(), &models.Log{IdUser: user.ID, IdEntity: 5, IdType: 4}).Exists(r.Context())
	if err != nil {
		s.logger.Errorf("Failed log login by user:%d:err:%v", user.ID, err)
	}
	err = tx.Commit()
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed commit")
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
		MaxAge:   7 * 24 * 60 * 60, // 7 дней
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
		MaxAge:   15 * 60,
	})
	w.WriteHeader(204)
}

func (s *Service) RefreshHandler(w http.ResponseWriter, r *http.Request) {

	var refreshToken string

	for _, c := range r.Cookies() {
		if c.Name == "refresh_token" {
			refreshToken = c.Value
		}
	}
	if refreshToken == "" {
		respondError(w, http.StatusUnauthorized, "Coocie unknown")
		return
	}

	token, err := s.repo.Refresh().GetByToken(r.Context(), refreshToken)
	if err != nil {
		respondError(w, http.StatusUnauthorized, "Invalid or expired refresh token")
		http.SetCookie(w, &http.Cookie{
			Name:     "refresh_token",
			Value:    "",
			HttpOnly: true,
			SameSite: http.SameSiteStrictMode,
			Path:     "/",
			MaxAge:   -1,
		})
		return
	}

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
	tx, err := s.repo.BeginTx(r.Context())
	if err != nil {
		s.logger.Errorf("Failed begin tx(log refresh) by user:%d:err:%v", token.UserID, err)
	}
	defer tx.Rollback()
	err = tx.Refresh().Update(r.Context(), &models.RefreshToken{ID: token.ID, Token: newRefreshToken, ExpiresAt: time.Now().Add(s.config.RefreshExpiry)})
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to update refresh token")
		return
	}
	_, err = tx.Log().Create(r.Context(), &models.Log{IdUser: token.UserID, IdEntity: 5, IdType: 6}).Exists(r.Context())
	if err != nil {
		s.logger.Errorf("Failed log refresh by user:%d:err:%v", token.UserID, err)
	}
	err = tx.Commit()
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed commit")
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    newRefreshToken,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
		MaxAge:   7 * 24 * 60 * 60, // 7 дней
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    newAccessToken,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
		MaxAge:   15 * 60,
	})
	w.WriteHeader(204)
}

func (s *Service) RevokeHandler(w http.ResponseWriter, r *http.Request) {
	var refreshToken string

	for _, c := range r.Cookies() {
		if c.Name == "refresh_token" {
			refreshToken = c.Value
		}
	}
	token, err := s.repo.Refresh().GetByToken(r.Context(), refreshToken)
	if err != nil {
		respondError(w, http.StatusNotFound, "Token not found")
		return
	}
	tx, err := s.repo.BeginTx(r.Context())
	if err != nil {
		s.logger.Errorf("Failed begin tx(log revoke) by user:%d:err:%v", token.UserID, err)
	}
	defer tx.Rollback()

	err = tx.Refresh().DeleteByToken(r.Context(), refreshToken)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to revoke token")
		return
	}
	_, err = tx.Log().Create(r.Context(), &models.Log{IdUser: token.UserID, IdEntity: 5, IdType: 5}).Exists(r.Context())
	if err != nil {
		s.logger.Errorf("Failed log revoke by user:%d:err:%v", token.UserID, err)
	}
	err = tx.Commit()
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed commit")
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
		MaxAge:   -1,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    "",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
		MaxAge:   -1,
	})
	respondJSON(w, http.StatusOK, map[string]string{"status": "success"})
}

// Middleware для аутентификации
func (s *Service) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var access string
		var refresh string
		for _, c := range r.Cookies() {
			if c.Name == "access_token" {
				access = c.Value
			}
			if c.Name == "refresh_token" {
				refresh = c.Value
			}
		}
		if access == "" && refresh == "" {
			respondError(w, http.StatusUnauthorized, "Invalid token")
			return
		}
		claims := &Claims{}

		token, err := jwt.ParseWithClaims(access, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(s.config.SecretKey), nil
		})

		if err != nil || !token.Valid {
			respondError(w, http.StatusForbidden, "Invalid token")
			return
		}
		tx, err := s.repo.BeginTx(r.Context())
		if err != nil {
			s.logger.Errorf("Failed begin tx(log auth) by user:%d:err:%v", claims.UserID, err)
		}
		defer tx.Rollback()

		_, err = tx.Log().Create(r.Context(), &models.Log{IdUser: claims.UserID, IdEntity: 5, IdType: 7}).Exists(r.Context())
		if err != nil {
			s.logger.Errorf("Failed log auth by user:%d:err:%v", claims.UserID, err)
		}
		err = tx.Commit()
		if err != nil {
			respondError(w, http.StatusInternalServerError, "Failed commit")
			return
		}
		ctx := context.WithValue(r.Context(), "id_user", claims.UserID)
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
