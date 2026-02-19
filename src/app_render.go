package pok

import (
	"html/template"
	"net/http"
)

type AppBaseData struct {
	PageTitle     string
	User          string
	Active        string
	NoticeSuccess string
	NoticeError   string

	// Shared stats
	TotalCards     int
	DistinctTypes  int
	TypesChartJSON string

	// Home
	RarestName  string
	RarestCount int

	

	// Stats
	MostObtainedName  string
	MostObtainedCount int
	MostObtainedImage string
	MostObtainedTypes string
}

func renderApp(w http.ResponseWriter, contentFile string, data AppBaseData) {
	t := template.Must(template.ParseFiles(
		"static/app_layout.html",
		contentFile,
	))
	_ = t.Execute(w, data)
}
