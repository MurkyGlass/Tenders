package handler

import (
	"fmt"
	"html/template"
	"main/internal/router/handlers/views"
	"net/http"
	"slices"
	"sort"
	"strconv"
	"strings"

	_ "github.com/lib/pq"
)

func (h *Handlers) FilterParamsByMyTenders() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		contentType := r.Header.Get("Content-Type")

		if strings.Contains(contentType, "multipart/form-data") {

			if err := r.ParseMultipartForm(10 << 20); err != nil {
				h.handleError(w, "Failed form-parsing", err, 500)
				return
			}
			IdsStr := r.PostForm.Get("category_ids")
			if IdsStr == "" {
				a := h.GetMyTendersListwindow(nil)
				a(w, r)
				return
			}
			IDsStrArr := r.PostForm["category_ids"]
			var Ids []int
			for _, str := range IDsStrArr {
				id, err := strconv.Atoi(str)
				if err != nil {
					h.handleError(w, "Invalid Parsing category id by Filter", err, 400)
					return
				}
				Ids = append(Ids, id)
			}
			a := h.GetMyTendersListwindow(Ids)
			a(w, r)
			return
		}
		h.handleError(w, "Invalid Content-Type", fmt.Errorf("expected multipart/form-data, got %s", contentType), 400)
		return
	}
}
func (h *Handlers) GetMyTendersListwindow(FilterParams []int) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		tenders, err := h.Repo.Tenders().GetAll(r.Context())
		if err != nil {
			h.handleError(w, "Failed get tenders:", err, 500)
			return
		}
		companies, err := h.Repo.Company().GetAll(r.Context())
		if err != nil {
			h.handleError(w, "Failed get companies:", err, 500)
			return
		}
		statuses, err := h.Repo.Status().GetAll(r.Context())
		if err != nil {
			h.handleError(w, "Failed get statuses:", err, 500)
			return
		}
		districts, err := h.Repo.District().GetAll(r.Context())
		if err != nil {
			h.handleError(w, "Failed get districts:", err, 500)
			return
		}

		var CompMap map[int]string = make(map[int]string)
		for _, comp := range companies {
			CompMap[comp.ID] = comp.Name
		}
		var DistMap map[int]string = make(map[int]string)
		for _, dist := range districts {
			DistMap[dist.ID] = dist.Name
		}
		var StMap map[int]string = make(map[int]string)
		for _, st := range statuses {
			StMap[st.ID] = st.Name
		}
		sortParam := r.URL.Query().Get("sort")

		switch sortParam {
		case "date_new":
			// сортировка по дате начала (новые сначала)
			sort.Slice(tenders, func(i int, j int) bool {
				return tenders[i].DateTimeStart.After(tenders[j].DateTimeStart)
			})
		case "date_old":
			// сортировка по дате начала (старые сначала)
			sort.Slice(tenders, func(i int, j int) bool {
				return tenders[i].DateTimeStart.Before(tenders[j].DateTimeStart)
			})
		case "deadline":
			// сортировка по дате окончания (заканчиваются скоро)
			sort.Slice(tenders, func(i int, j int) bool {
				return tenders[i].DateTimeEnd.After(tenders[j].DateTimeEnd)
			})
		case "deadline_far":
			// сортировка по дате окончания (заканчиваются позже)
			sort.Slice(tenders, func(i int, j int) bool {
				return tenders[i].DateTimeEnd.Before(tenders[j].DateTimeEnd)
			})
		case "name_asc":
			// сортировка по названию (А-Я)
			sort.Slice(tenders, func(i int, j int) bool {
				return CompareStrings(tenders[i].Name, tenders[j].Name) == -1
			})
		case "name_desc":
			// сортировка по названию (Я-А)
			sort.Slice(tenders, func(i int, j int) bool {
				return CompareStrings(tenders[i].Name, tenders[j].Name) == 1
			})
		case "duration_asc":
			// сортировка по длительности (сначала короткие)
			sort.Slice(tenders, func(i int, j int) bool {
				return tenders[i].DateTimeEnd.Sub(tenders[i].DateTimeStart) < tenders[j].DateTimeEnd.Sub(tenders[j].DateTimeStart)
			})
		case "duration_desc":
			// сортировка по длительности (сначала длинные)
			sort.Slice(tenders, func(i int, j int) bool {
				return tenders[i].DateTimeEnd.Sub(tenders[i].DateTimeStart) > tenders[j].DateTimeEnd.Sub(tenders[j].DateTimeStart)
			})
		case "status":
			// сортировка по статусу
			sort.Slice(tenders, func(i int, j int) bool {
				return tenders[i].IdStatus < tenders[j].IdStatus
			})
		case "company":
			// сортировка по компании
			sort.Slice(tenders, func(i int, j int) bool {
				return tenders[i].IdCompany < tenders[j].IdCompany
			})
		case "district":
			// сортировка по району
			sort.Slice(tenders, func(i int, j int) bool {
				return tenders[i].IdDistrict < tenders[j].IdDistrict
			})
		default:
			// сортировка по умолчанию
			sort.Slice(tenders, func(i int, j int) bool {
				return tenders[i].DateTimeStart.After(tenders[j].DateTimeStart)
			})
			sortParam = "date_new"
		}
		tenderLinks, err := h.Repo.LinkerTCategory(0).GetAll(r.Context())
		if err != nil {
			h.handleError(w, "Failed get tenders-categories:", err, 500)
			return
		}
		var TCmap = make(map[int][]int) //key - tenderID
		for _, tl := range tenderLinks {
			TCmap[tl.IdTender] = append(TCmap[tl.IdTender], tl.IdCategory)
		}
		// сортировка фильтрации ограничения итд
		search := r.URL.Query().Get("search")
		//pesronal data
		id, ok := r.Context().Value("id_user").(int)
		if !ok {
			h.handleError(w, "Bad user id", fmt.Errorf("Bad User id"), 500)
			return
		}
		user, err := h.Repo.Users().GetByID(r.Context(), id)
		if err != nil {
			h.handleError(w, "db request failed", err, 500)
			return
		}
		company, err := h.Repo.Company().GetByID(r.Context(), user.IdCompany)
		if err != nil {
			h.handleError(w, "db request failed", err, 500)
			return
		}
		//
		var TenderViews []views.TenderView
		for _, tender := range tenders {
			//  не черновик и не активный
			if tender.IdStatus != 1 && tender.IdStatus != 2 {
				continue
			}
			if tender.IdCompany != company.ID {
				continue
			}
			if search != "" {
				lsearch := strings.ToLower(search)
				if !strings.Contains(strings.ToLower(tender.Name), lsearch) && !strings.Contains(strings.ToLower(tender.Description), lsearch) && !strings.Contains(strings.ToLower(CompMap[tender.IdCompany]), lsearch) {
					continue
				}
			}
			if FilterParams != nil {
				if TCmap[tender.ID] != nil {
					f := false
					for _, c := range TCmap[tender.ID] {
						if slices.Contains(FilterParams, c) {
							f = true
						}
					}
					if !f {
						continue
					}
				} else {
					continue
				}
			}

			TenderViews = append(TenderViews, views.TenderView{ID: tender.ID, Name: tender.Name,
				Description: tender.Description, DateTimeStart: GetDateString(tender.DateTimeStart),
				DateTimeEnd: GetDateString(tender.DateTimeEnd), Company: CompMap[tender.IdCompany], Status: StMap[tender.IdStatus],
				District: DistMap[tender.IdDistrict]})
		}

		tmpl, err := template.ParseFiles("./client/pages/my_tender_list.html")
		if err != nil {
			h.handleError(w, "Failed my tenders load:", err, 500)
			return
		}

		categories, err := h.Repo.Category().GetAll(r.Context())
		if err != nil {
			h.handleError(w, "Failed get categories:", err, 500)
			return
		}
		catlinks, err := h.Repo.CategoryLink().GetAll(r.Context())
		if err != nil {
			h.handleError(w, "Failed get category links:", err, 500)
			return
		}
		cv := BuildCategoryTree(categories, catlinks)

		type data struct {
			LoginForm        template.HTML
			RegistrationForm template.HTML
			Tenders          []views.TenderView
			CatView          []views.CategoryView
			Sort             string
		}
		err = tmpl.Execute(w, &data{Tenders: TenderViews, LoginForm: LoginForm, RegistrationForm: RegistrationForm, Sort: sortParam, CatView: cv})
		if err != nil {
			h.handleError(w, "Failed tenderlist render:", err, 500)
			return
		}
	}
}
