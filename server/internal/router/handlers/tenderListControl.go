package handler

import (
	"fmt"
	"html/template"
	"main/internal/repositories/models"
	"main/internal/router/handlers/views"
	"net/http"
	"sort"
	"strconv"
	"time"

	_ "github.com/lib/pq"
)

func BuildCategoryTree(categories []models.Category, links []models.LinkView) []views.CategoryView {
	var cm = make(map[int]models.Category)
	var lpk = make(map[int][]int) //parrent key
	var lck = make(map[int]int)   //child key
	for _, category := range categories {
		cm[category.ID] = category
	}
	for _, link := range links {
		lpk[link.FirstID] = append(lpk[link.FirstID], link.SecondID)
		lck[link.SecondID] = link.FirstID
	}

	var view []views.CategoryView
	var childs = make(map[int]models.Category)
	for _, c := range categories {
		_, f0 := lck[c.ID]
		if f0 { //если явл ребенком
			childs[c.ID] = c
			continue
		}
		_, f1 := lpk[c.ID]
		if f1 && !f0 { //если явл отцом и не ребенок
			view = append(view, views.CategoryView{Category: c, Childs: nil})
			continue
		}
		view = append(view, views.CategoryView{Category: c, Childs: nil})
	}

	var road []*views.CategoryView
	for i := range view {
		road = append(road, &view[i])
	}

	for len(childs) > 0 {

		for _, r := range road {
			for _, l := range links {
				if l.FirstID == r.Category.ID {
					r.Childs = append(r.Childs, views.CategoryView{Category: childs[l.SecondID], Childs: nil})
					delete(childs, l.SecondID)
				}
			}
		}

		newRoad := []*views.CategoryView{}
		for _, r := range road {
			for i := range r.Childs {
				newRoad = append(newRoad, &r.Childs[i])
			}
		}
		road = newRoad
	}
	return view
}

func (h *Handlers) GetTendersListwindow() func(w http.ResponseWriter, r *http.Request) {
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
				return tenders[i].DateTimeStart.Before(tenders[j].DateTimeStart)
			})
		case "date_old":
			// сортировка по дате начала (старые сначала)
			sort.Slice(tenders, func(i int, j int) bool {
				return tenders[i].DateTimeStart.After(tenders[j].DateTimeStart)
			})
		case "deadline":
			// сортировка по дате окончания (заканчиваются скоро)
			sort.Slice(tenders, func(i int, j int) bool {
				return tenders[i].DateTimeEnd.Before(tenders[j].DateTimeEnd)
			})
		case "deadline_far":
			// сортировка по дате окончания (заканчиваются позже)
			sort.Slice(tenders, func(i int, j int) bool {
				return tenders[i].DateTimeEnd.After(tenders[j].DateTimeEnd)
			})
		case "name_asc":
			// сортировка по названию (А-Я)
			sort.Slice(tenders, func(i int, j int) bool {
				return tenders[i].Name < tenders[j].Name
			})
		case "name_desc":
			// сортировка по названию (Я-А)
			sort.Slice(tenders, func(i int, j int) bool {
				return tenders[i].Name > tenders[j].Name
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
				return tenders[i].DateTimeStart.Before(tenders[j].DateTimeStart)
			})
			sortParam = "date_new"
		}
		// сортировка фильтрации ограничения итд
		var TenderViews []views.TenderView
		for _, tender := range tenders {
			// черновик, завершен(отобр. в отдельном списке)
			if tender.IdStatus == 1 || tender.IdStatus == 3 {
				continue
			}

			TenderViews = append(TenderViews, views.TenderView{ID: tender.ID, Name: tender.Name,
				Description: tender.Description, DateTimeStart: GetDateString(tender.DateTimeStart),
				DateTimeEnd: GetDateString(tender.DateTimeEnd), Company: CompMap[tender.IdCompany], Status: StMap[tender.IdStatus],
				District: DistMap[tender.IdDistrict]})
		}

		tmpl, err := template.ParseFiles("./client/pages/tender_list.html")
		if err != nil {
			h.handleError(w, "Failed profil load:", err, 500)
			return
		}
		type data struct {
			LoginForm        template.HTML
			RegistrationForm template.HTML
			Tenders          []views.TenderView
			Sort             string
		}
		err = tmpl.Execute(w, &data{Tenders: TenderViews, LoginForm: LoginForm, RegistrationForm: RegistrationForm, Sort: sortParam})
		if err != nil {
			h.handleError(w, "Failed tenderlist render:", err, 500)
			return
		}
	}
}
func (h *Handlers) GetTenderwindow() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := h.getIDFromRequest(r)
		if err != nil {
			h.handleError(w, "Invalid ID", err, 500)
			return
		}

		tender, err := h.Repo.Tenders().GetByID(r.Context(), id)
		if err != nil {
			h.handleError(w, "Failed get tender:", err, 500)
			return
		}
		tmpl, err := template.ParseFiles("./client/pages/tender_view.html")
		if err != nil {
			h.handleError(w, "Failed profil load:", err, 500)
			return
		}
		company, err := h.Repo.Company().GetByID(r.Context(), tender.IdCompany)
		if err != nil {
			h.handleError(w, "Failed get company:", err, 500)
			return
		}
		status, err := h.Repo.Status().GetByID(r.Context(), tender.IdStatus)
		if err != nil {
			h.handleError(w, "Failed get status:", err, 500)
			return
		}
		district, err := h.Repo.District().GetByID(r.Context(), tender.IdDistrict)
		if err != nil {
			h.handleError(w, "Failed get district:", err, 500)
			return
		}
		var tenview *views.TenderView
		tenview = &views.TenderView{ID: tender.ID, Name: tender.Name, Description: tender.Description,
			DateTimeStart: GetDateString(tender.DateTimeStart), DateTimeEnd: GetDateString(tender.DateTimeEnd), Company: company.Name,
			Status: status.Name, District: district.Name}
		type data struct {
			LoginForm        template.HTML
			RegistrationForm template.HTML
			Tender           *views.TenderView
		}
		err = tmpl.Execute(w, &data{Tender: tenview, LoginForm: LoginForm, RegistrationForm: RegistrationForm})
		if err != nil {
			h.handleError(w, "Failed tender render:", err, 500)
			return
		}
	}
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
