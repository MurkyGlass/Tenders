package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"main/internal/repositories"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

const (
	LoginForm = `<div class="login-modal" id="loginModal">
        <div class="login-overlay" id="loginOverlay">
        </div>
        <div class="login-container">
            <div class="login-header">
                <h2>Авторизация</h2>
                <button class="login-close" id="loginClose">&times;</button>
            </div>
        
            <form class="login-form" id="loginForm">
                <div class="form-group">
                    <label for="loginInput" class="form-label">
                        <i class="fas fa-user"></i> Логин
                    </label>
                    <input 
                        type="text" 
                        id="loginInput" 
                        name="login"
                        class="form-input"
                        placeholder="Введите ваш логин"
                        required
                        pattern="[A-Za-z0-9_]{3,30}"
                        autocomplete="username">
                </div>
            
                <div class="form-group">
                    <label for="passwordInput" class="form-label">
                        <i class="fas fa-lock"></i> Пароль
                    </label>
                    <div class="password-wrapper">
                        <input 
                            type="password" 
                            id="passwordInput" 
                            name="password"
                            class="form-input"
                            placeholder="Введите ваш пароль"
                            required
                            minlength="8"
                            autocomplete="current-password">
                        <button type="button" class="password-toggle" id="passwordToggle">
                            <i class="fas fa-eye"></i>
                        </button>
                    </div>
                </div>
        
                <button type="submit" class="login-submit" id="loginSubmit">
                    <i class="fas fa-sign-in-alt"></i> Войти
                </button>
            
                <div class="form-footer">
                    <button type="button" class="form-btn register-btn" id="regisBtn">Зарегистрироваться</button>
                    <a href="#" class="form-link">Забыли пароль?</a>
                </div>

                
            </form>
        </div>
    </div>`
	RegistrationForm = `<div class="register-modal" id="registerModal">
        <div class="register-overlay" id="registerOverlay"></div>
        <div class="register-container">
            <div class="register-header">
                <h2>Регистрация</h2>
                <button class="register-close" id="registerClose">&times;</button>
            </div>
        
            <form class="register-form" id="registerForm">
                <!-- Секция компании -->
                <div class="form-section">
                    <h3 class="section-title">
                        <i class="fas fa-building"></i> Данные компании
                    </h3>
                
                    <div class="form-group">
                        <label for="companyName" class="form-label">
                            <i class="fas fa-signature"></i> Название компании *
                        </label>
                        <input 
                            type="text" 
                            id="companyName" 
                            name="companyName"
                            class="form-input"
                            placeholder="Введите название компании"
                            required
                            maxlength="100">
                        <div class="form-hint">Максимум 100 символов</div>
                    </div>
                
                    <div class="form-group">
                        <label for="companyEmail" class="form-label">
                            <i class="fas fa-envelope"></i> Email компании *
                        </label>
                        <input 
                            type="email" 
                            id="companyEmail" 
                            name="companyEmail"
                            class="form-input"
                            placeholder="example@company.com"
                            required
                            maxlength="50">
                    </div>
                
                    <div class="form-group">
                        <label for="companyAddress" class="form-label">
                            <i class="fas fa-map-marker-alt"></i> Юридический адрес *
                        </label>
                        <input 
                            type="text" 
                            id="companyAddress" 
                            name="companyAddress"
                            class="form-input"
                            placeholder="Введите полный адрес компании"
                            required
                            maxlength="100">
                    </div>
                
                    <div class="form-row">
                        <div class="form-group half">
                            <label for="companyINN" class="form-label">
                                <i class="fas fa-file-contract"></i> ИНН *
                            </label>
                            <input 
                                type="text" 
                                id="companyINN" 
                                name="companyINN"
                                class="form-input"
                                placeholder="10 или 12 цифр"
                                required
                                pattern="\d{10}|\d{12}"
                                maxlength="12">
                        </div>
                    
                        <div class="form-group half">
                            <label for="companyEGRUL" class="form-label">
                                <i class="fas fa-file-alt"></i> ЕГРЮЛ *
                            </label>
                            <input 
                                type="text" 
                                id="companyEGRUL" 
                                name="companyEGRUL"
                                class="form-input"
                                placeholder="Номер записи"
                                pattern="\d{13}"
                                required
                                maxlength="13">
                        </div>
                    </div>
                
                    <div class="form-group">
                        <label for="companyDescription" class="form-label">
                            <i class="fas fa-file-signature"></i> Описание компании
                        </label>
                        <textarea 
                            id="companyDescription" 
                            name="companyDescription"
                            class="form-textarea"
                            placeholder="Краткое описание деятельности компании (до 500 символов)"
                            rows="3"
                            maxlength="500"></textarea>
                    </div>
                </div>
            
                <!-- Секция пользователя -->
                <div class="form-section">
                    <h3 class="section-title">
                        <i class="fas fa-user"></i> Данные администратора
                    </h3>
                
                    <div class="form-group">
                        <label for="userLogin" class="form-label">
                            <i class="fas fa-user-circle"></i> Логин *
                        </label>
                        <input 
                            type="text" 
                            id="userLogin" 
                            name="userLogin"
                            class="form-input"
                            placeholder="Придумайте уникальный логин"
                            required
                            pattern="[A-Za-z0-9_]{3,30}"
                            maxlength="30">
                        <div class="form-hint">Только латинские буквы, цифры и нижнее подчёркивание</div>
                    </div>
                
                    <div class="form-group">
                        <label for="userName" class="form-label">
                            <i class="fas fa-user"></i> ФИО пользователя *
                        </label>
                        <input 
                            type="text" 
                            id="userName" 
                            name="userName"
                            class="form-input"
                            placeholder="Введите ваше имя"
                            pattern="[A-Za-zА-Яа-я_]{3,30}"
                            required
                            maxlength="30">
                    </div>
                
                    <div class="form-group">
                        <label for="userEmail" class="form-label">
                            <i class="fas fa-envelope"></i> Личный email *
                        </label>
                        <input 
                            type="email" 
                            id="userEmail" 
                            name="userEmail"
                            class="form-input"
                            placeholder="your.email@example.com"
                            required
                            maxlength="50">
                    </div>
                
                    <div class="form-group">
                        <label for="userPassword" class="form-label">
                            <i class="fas fa-lock"></i> Пароль *
                        </label>
                        <div class="password-wrapper">
                            <input 
                                type="password" 
                                id="userPassword" 
                                name="userPassword"
                                class="form-input"
                                placeholder="Придумайте надёжный пароль"
                                required
                                minlength="8"
                                autocomplete="new-password">
                            <button type="button" class="password-toggle" id="userPasswordToggle">
                                <i class="fas fa-eye"></i>
                            </button>
                        </div>
                    </div>
                
                    <div class="form-group">
                        <label for="confirmPassword" class="form-label">
                            <i class="fas fa-lock"></i> Подтверждение пароля *
                        </label>
                        <div class="password-wrapper">
                            <input 
                                type="password" 
                                id="confirmPassword" 
                                name="confirmPassword"
                                class="form-input"
                                placeholder="Повторите пароль"
                                required
                                autocomplete="new-password">
                            <button type="button" class="password-toggle" id="confirmPasswordToggle">
                                <i class="fas fa-eye"></i>
                            </button>
                        </div>
                    </div>
                </div>
            
                <div class="form-agreement">
                    <input type="checkbox" id="agreement" name="agreement" required>
                    <label for="agreement">
                        Я согласен с <a href="#" class="link">политикой конфиденциальности</a> и <a href="#" class="link">правилами использования</a> *
                    </label>
                </div>
            
                <button type="submit" class="register-submit" id="registerSubmit">
                    <i class="fas fa-user-plus"></i> Зарегистрироваться
                </button>
            
                <div class="form-footer">
                    <span>Уже есть аккаунт?</span>
                    <a href="#" class="login-link">Войти</a>
                </div>
            </form>
        
            <div class="register-status" id="registerStatus"></div>
        </div>
    </div>`
)

type Handlers struct {
	Logger *logrus.Logger
	Repo   repositories.Repository
}

func NewHandlers(logger *logrus.Logger, rep repositories.Repository) *Handlers {
	return &Handlers{
		Logger: logger,
		Repo:   rep,
	}
}

func (h *Handlers) HealthCheck() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	}
}

// Логирование
func LoggingMiddleware(logger *logrus.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.Infof("%s %s", r.Method, r.URL.Path)
			next.ServeHTTP(w, r)
		})
	}
}

// CORS
func CorsMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Max-Age", "86400")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func (h *Handlers) getIDFromRequest(r *http.Request) (int, error) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	return strconv.Atoi(idStr)
}
func (h *Handlers) getIDOfferFromRequest(r *http.Request) (int, error) {
	vars := mux.Vars(r)
	idStr := vars["idoffer"]
	return strconv.Atoi(idStr)
}
func (h *Handlers) getLoginFromRequest(r *http.Request) (string, error) {
	vars := mux.Vars(r)
	Str := vars["login"]
	l := len(Str)
	if l > 3 && l < 31 {
		return Str, nil
	}
	return Str, fmt.Errorf("Validation Error;%s; len in bait no is 3<l<31", Str)
}

// Runing in TX
func (h *Handlers) executeInTransaction(r *http.Request, fn func(tx repositories.Transaction) error) error {
	tx, err := h.Repo.BeginTx(r.Context())
	if err != nil {
		h.Logger.Errorf("Failed to begin transaction: %v", err)
		return err
	}
	defer tx.Rollback()

	if err := fn(tx); err != nil {
		h.Logger.Errorf("Operation failed in TX: %v", err)
		return err
	}

	if err := tx.Commit(); err != nil {
		h.Logger.Errorf("Failed to commit transaction: %v", err)
		return err
	}
	return nil
}

func (h *Handlers) handleError(w http.ResponseWriter, message string, err error, code int) {
	h.Logger.Errorf("%s: %v", message, err)
	http.Error(w, message, code)
}

func jsonResponse(w http.ResponseWriter, data interface{}, statusCode ...int) {
	w.Header().Set("Content-Type", "application/json")
	if len(statusCode) > 0 {
		w.WriteHeader(statusCode[0])
	}
	json.NewEncoder(w).Encode(data)
}
func GetDateString(t time.Time) string {
	month := t.Month()
	day := t.Day()
	hour := t.Hour()
	minute := t.Minute()
	var monstr string
	if int(month) < 10 {
		monstr = "0" + strconv.Itoa(int(month))
	} else {
		monstr = strconv.Itoa(int(month))
	}
	var daystr string
	if day < 10 {
		daystr = "0" + strconv.Itoa(day)
	} else {
		daystr = strconv.Itoa(day)
	}
	var hstr string
	if hour < 10 {
		hstr = "0" + strconv.Itoa(hour)
	} else {
		hstr = strconv.Itoa(hour)
	}
	var mstr string
	if minute < 10 {
		mstr = "0" + strconv.Itoa(minute)
	} else {
		mstr = strconv.Itoa(minute)
	}
	return fmt.Sprintf("%s.%s.%d %s:%s", daystr, monstr, t.Year(), hstr, mstr)
}
func (h *Handlers) readRightsInContext(ctx context.Context, rights []string) (bool, error) {
	id, ok := ctx.Value("id_user").(int)
	if !ok {
		return false, fmt.Errorf("Bad user id")
	}
	user, err := h.Repo.Users().GetByID(ctx, id)
	if err != nil {
		return false, fmt.Errorf("db request failed")
	}
	if user.IdRoleInCompany != 1 {
		for _, r := range rights {
			if ctx.Value(r) == "no" {
				return false, nil
			}
		}
	}
	return true, nil
}
//tenders
var ViewOwnTendersPerms = []string{"view_own_tenders"}//view my tenders

var EditTenderPerms = []string{"view_own_tenders", "edit_own_tender"}//edit draft
var EditPublishTenderPerms = []string{"view_own_tenders", "edit_own_tender","change_tender_status"}//edit

var ChooseWinnerPerms = []string{"view_own_tenders", "choose_winner"}//await tenders

var CreateTenderPerms = []string{"create_tender"}//create draft
var PublishTenderPerms = []string{"create_tender", "change_tender_status"}//create

//offers
var ViewOwnOffersPerms = []string{"view_own_offers"}

var EditOfferPerms = []string{"view_own_offers", "edit_own_offer"}//edit draft
var EditPublishOfferPerms = []string{"view_own_offers", "edit_own_offer","change_offer_status"}//edit

var CreateOfferPerms = []string{"create_offer"}//create draft
var PublishOfferPerms = []string{"create_offer", "change_offer_status"}//create



//company

var EditCompanyDataPerms = []string{"edit_company_data"}//edit company
var UploadCompanyDocsPerms = []string{"upload_company_docs"}//docs

var ManageUsersPerms = []string{"manage_company_users"}
var ManageRolesPerms = []string{"manage_roles"}
