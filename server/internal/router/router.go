package router

import (
	handler "main/internal/router/handlers"
	"main/internal/router/handlers/jwt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func NewRouter(h *handler.Handlers, db *sqlx.DB) *mux.Router {
	r := mux.NewRouter()
	// Инициализация
	jwtService := jwt.NewService(h.Logger, h.Repo)
	//client static
	fs := http.FileServer(http.Dir("./client/static"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))
	//main page
	r.HandleFunc("", h.GetMainwindow()).Methods("GET")
	mainRouter := r.PathPrefix("/main").Subrouter()
	//public
	//main
	mainRouter.HandleFunc("", h.GetMainwindow()).Methods("GET")
	//tenders
	mainRouter.HandleFunc("/tenders", h.GetTendersListwindow(nil)).Methods("GET")
	mainRouter.HandleFunc("/tenders", h.FilterParams()).Methods("POST")
	//completed tenders
	mainRouter.HandleFunc("/tenders/completed", h.GetCompletedTendersListwindow(nil)).Methods("GET")
	mainRouter.HandleFunc("/tenders/completed", h.CompletedFilterParams()).Methods("POST")
	//tender
	mainRouter.HandleFunc("/tenders/{id}", h.GetTenderwindow()).Methods("GET")
	//registration
	mainRouter.HandleFunc("/registration", h.Registration()).Methods("POST")
	//companies
	mainRouter.HandleFunc("/companies", h.GetCompaniesListwindow()).Methods("GET")
	mainRouter.HandleFunc("/companies/{id}", h.GetCompanyWindow()).Methods("GET")
	//reset password
	mainRouter.HandleFunc("/password/reset", h.GetResetWindow()).Methods("GET")
	mainRouter.HandleFunc("/password/reset", h.ResetPassword()).Methods("POST")
	mainRouter.HandleFunc("/password/reset/form", h.GetResetForm()).Methods("GET")
	mainRouter.HandleFunc("/password/reset/form", h.UpdatePassword()).Methods("POST")
	//authentification
	r.HandleFunc("/auth/login", jwtService.LoginHandler).Methods("POST")
	r.HandleFunc("/auth/refresh", jwtService.RefreshHandler).Methods("GET")
	r.HandleFunc("/auth/revoke", jwtService.RevokeHandler).Methods("GET")
	//protected
	prRouter := r.PathPrefix("/protected").Subrouter()
	prRouter.Use(jwtService.Middleware)
	//tenders docs
	prRouter.HandleFunc("/tenders/documents/{id}", h.GetDocumentById()).Methods("GET")
	prRouter.HandleFunc("/tenders/{id}/documents", h.GetTenderDocuments()).Methods("GET")
	//--------------------------------------------------------------------------------------
	//company edit docs
	prRouter.HandleFunc("/companies/{id}/docs/edit", h.GetEditCompanyDocsWindow()).Methods("GET")
	prRouter.HandleFunc("/companies/{id}/docs/edit", h.EditCompanyDocs()).Methods("POST")
	//--------------------------------------------------------------------------------------
	//offer-create
	prRouter.HandleFunc("/tenders/{id}/offer/create", h.GetCreateOfferWindow()).Methods("GET")
	prRouter.HandleFunc("/tenders/{id}/offer/create", h.CreateOffer(2)).Methods("POST")
	prRouter.HandleFunc("/tenders/{id}/offer/create/draft", h.CreateOffer(1)).Methods("POST")
	//--------------------------------------------------------------------------------------
	//offers
	prRouter.HandleFunc("/offers/{id}", h.GetOfferWindow()).Methods("GET")
	//offers docs
	prRouter.HandleFunc("/offers/documents/{id}", h.GetDocumentById()).Methods("GET")
	prRouter.HandleFunc("/offers/{id}/documents", h.GetOfferDocuments()).Methods("GET")
	//companies docs
	prRouter.HandleFunc("/companies/documents/{id}", h.GetDocumentById()).Methods("GET")
	prRouter.HandleFunc("/companies/{id}/documents", h.GetCompanyDocuments()).Methods("GET")
	//--------------------------------------------------------------------------------------
	lkRouter := prRouter.PathPrefix("/lk").Subrouter()
	lkRouter.HandleFunc("", h.GetProfilwindow()).Methods("GET")
	lkRouter.HandleFunc("/edit", h.EditingLK()).Methods("POST")
	lkRouter.HandleFunc("/company/role/create", h.CreateRoleInCompany()).Methods("POST")
	lkRouter.HandleFunc("/company/user/create", h.CreateNewUser()).Methods("POST")
	//company tenders
	lkRouter.HandleFunc("/tenders", h.GetMyTendersListwindow(nil)).Methods("GET")
	lkRouter.HandleFunc("/tenders", h.FilterParamsByMyTenders()).Methods("POST")
	lkRouter.HandleFunc("/tenders/{id}", h.GetMyTenderwindow()).Methods("GET")
	lkRouter.HandleFunc("/tenders/{id}/edit", h.GetEditTenderWindow()).Methods("GET")
	lkRouter.HandleFunc("/tenders/{id}/edit", h.EditTender()).Methods("POST")
	lkRouter.HandleFunc("/tenders/{id}/edit/draft", h.EditDraftTender(1)).Methods("POST")
	lkRouter.HandleFunc("/tenders/{id}/edit/active", h.EditDraftTender(2)).Methods("POST")
	//awaiting tenders
	lkRouter.HandleFunc("/await/tenders", h.GetMyAwaitTendersListwindow()).Methods("GET")
	lkRouter.HandleFunc("/await/tenders/{id}", h.GetCompletedTenderwindow()).Methods("GET")
	lkRouter.HandleFunc("/await/tenders/{id}/offers/{idoffer}", h.CompletTender()).Methods("POST")
	//company offers
	lkRouter.HandleFunc("/offers", h.GetOfferListWindow()).Methods("GET")
	lkRouter.HandleFunc("/offers/{id}", h.GetMyOfferWindow()).Methods("GET")
	lkRouter.HandleFunc("/offers/{id}/edit", h.GetEditOfferWindow()).Methods("GET")
	lkRouter.HandleFunc("/offers/{id}/edit", h.EditOffer(2)).Methods("POST")
	lkRouter.HandleFunc("/offers/{id}/edit/draft", h.EditOffer(1)).Methods("POST")
	//--------------------------------------------------------------------------------------
	tenderRouter := prRouter.PathPrefix("/tender").Subrouter()
	tenderRouter.HandleFunc("/create", h.GetCreateTenderWindow()).Methods("GET")
	tenderRouter.HandleFunc("/create", h.CreateTender(2)).Methods("POST")
	tenderRouter.HandleFunc("/create/draft", h.CreateTender(1)).Methods("POST")
	//---------------------------------------------------------------------------------------------\\
	// Health check
	r.HandleFunc("/health", h.HealthCheck()).Methods("GET")
	//
	//protected admin
	adminRouter := r.PathPrefix("/admin").Subrouter()
	adminRouter.Use(jwtService.MiddlewareAdmin)
	//admin panel
	adminRouter.HandleFunc("/panel", h.GetAdminPanelWindow()).Methods("GET")
	//category create
	adminRouter.HandleFunc("/category/create", h.GetAdminCategoryCreateWindow()).Methods("GET")
	adminRouter.HandleFunc("/category/create", h.AdminCategoryCreate()).Methods("POST")
	//category update
	adminRouter.HandleFunc("/category/edit", h.GetAdminCategoryEditWindow()).Methods("GET")
	adminRouter.HandleFunc("/category/edit", h.AdminCategoryEdit()).Methods("POST")
	//------------------------------------------------------------------------------------------------
	return r
}
