package web

import (
	"a21hc3NpZ25tZW50/client"
	"embed"
	"html/template"
	"log"
	"net/http"
	"path"
)

type DashboardWeb interface {
	Dashboard(w http.ResponseWriter, r *http.Request)
}

type dashboardWeb struct {
	categoryClient client.CategoryClient
	embed          embed.FS
}

func NewDashboardWeb(catClient client.CategoryClient, embed embed.FS) *dashboardWeb {
	return &dashboardWeb{catClient, embed}
}

func (d *dashboardWeb) Dashboard(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("id")

	categories, err := d.categoryClient.GetCategories(userId.(string))
	if err != nil {
		log.Println("error get cat: ", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var dataTemplate = map[string]interface{}{
		"categories": categories,
	}

	var funcMap = template.FuncMap{
		"categoryInc": func(catId int) int {
			return catId + 1
		},
		"categoryDec": func(catId int) int {
			return catId - 1
		},
	}

	// ignore this
	_ = dataTemplate
	_ = funcMap
	

	// TODO: answer here
	header := path.Join("views", "general", "header.html")
	viewDashborad := path.Join("views", "main", "dashboard.html")
	
	tmp, err := template.ParseFS(d.embed, viewDashborad, header)
	if err != nil {
		log.Println(err)
		log.Println("parshing error")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = tmp.Funcs(funcMap).Execute(w, dataTemplate)
	if err != nil {
		log.Println(err)
		log.Println("execution error")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}
