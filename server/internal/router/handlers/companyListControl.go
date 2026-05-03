package handler

import (
	"html/template"
	"main/internal/repositories/models"
	"net/http"
	"sort"
	"strings"

	_ "github.com/lib/pq"
)

func CompareStrings(a string, b string) int {
	runesA := []rune(strings.ToLower(a))
	runesB := []rune(strings.ToLower(b))
	if len(runesA) < len(runesB) {
		for i := range runesA {
			if runesA[i] > runesB[i] {
				return 1
			} else if runesA[i] < runesB[i] {
				return -1
			}
		}
		if len(runesA) == len(runesB) {
			return 0
		} else {
			return 1
		}
	} else {
		for i := range runesB {
			if runesA[i] > runesB[i] {
				return 1
			} else if runesA[i] < runesB[i] {
				return -1
			}
		}
		if len(runesA) == len(runesB) {
			return 0
		} else {
			return -1
		}
	}
}
func (h *Handlers) GetCompaniesListwindow() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		companies, err := h.Repo.Company().GetAll(r.Context())
		if err != nil {
			h.handleError(w, "Failed get companies:", err, 500)
			return
		}

		sortParam := r.URL.Query().Get("sort")

		switch sortParam {
		case "name_asc":
			// сортировка по названию (А-Я)
			sort.Slice(companies, func(i int, j int) bool {
				return CompareStrings(companies[i].Name, companies[j].Name) == -1
			})
		case "name_desc":
			// сортировка по названию (Я-А)
			sort.Slice(companies, func(i int, j int) bool {
				return CompareStrings(companies[i].Name, companies[j].Name) == 1
			})
		case "adress_asc":
			// сортировка по адресс а-я
			sort.Slice(companies, func(i int, j int) bool {
				return CompareStrings(companies[i].Address, companies[j].Address) == -1
			})
		case "adress_desc":
			// сортировка по адресс я-а
			sort.Slice(companies, func(i int, j int) bool {
				return CompareStrings(companies[i].Address, companies[j].Address) == 1				
			})
		case "inn_asc":
			// сортировка по ИНН 1-9
			sort.Slice(companies, func(i int, j int) bool {
				return companies[i].INN < companies[j].INN
			})
		case "inn_desc":
			// сортировка по ИНН 9-1
			sort.Slice(companies, func(i int, j int) bool {
				return companies[i].INN > companies[j].INN
			})
		case "egrul_asc":
			// сортировка по ЕГРЮЛ 1-9
			sort.Slice(companies, func(i int, j int) bool {
				return companies[i].EGRUL < companies[j].EGRUL
			})
		case "egrul_desc":
			// сортировка по ЕГРЮЛ 9-1
			sort.Slice(companies, func(i int, j int) bool {
				return companies[i].EGRUL > companies[j].EGRUL
			})
		default:
			// сортировка по умолчанию
			sort.Slice(companies, func(i int, j int) bool {
				return CompareStrings(companies[i].Name, companies[j].Name) == -1
			})
			sortParam = "name_asc"
		}

		// сортировка фильтрации ограничения итд
		search := r.URL.Query().Get("search")

		var ended []models.Company
		for _, company := range companies {

			if search != "" {
				lsearch := strings.ToLower(search)
				if !strings.Contains(strings.ToLower(company.Name), lsearch) && !strings.Contains(strings.ToLower(company.Description), lsearch) && !strings.Contains(strings.ToLower(company.Address), lsearch) && !strings.Contains(strings.ToLower(company.Email), lsearch) && !strings.Contains(strings.ToLower(company.INN), lsearch) && !strings.Contains(strings.ToLower(company.EGRUL), lsearch) {
					continue
				}
			}
			ended = append(ended, company)
		}

		tmpl, err := template.ParseFiles("./client/pages/company_list.html")
		if err != nil {
			h.handleError(w, "Failed companies load:", err, 500)
			return
		}

		type data struct {
			LoginForm        template.HTML
			RegistrationForm template.HTML
			Companies        []models.Company
			Sort             string
		}
		err = tmpl.Execute(w, &data{Companies: ended, LoginForm: LoginForm, RegistrationForm: RegistrationForm, Sort: sortParam})
		if err != nil {
			h.handleError(w, "Failed companylist render:", err, 500)
			return
		}
	}
}
