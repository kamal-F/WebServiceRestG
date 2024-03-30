package others

import (
	"encoding/json"
	"net/http"
	"rongsokapi/database"
)

type News struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

func NewsHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "GET":
		getNews(w, r)
	case "POST":
		createNews(w, r)
	case "PUT":
		updateNews(w, r)
	case "DELETE":
		deleteNews(w, r)
	default:
		http.Error(w, "Unsupported HTTP Method", http.StatusBadRequest)
	}
}
func getNews(w http.ResponseWriter, r *http.Request) {
	// implementation of getNews function
	db := database.DB

	news := []News{}
	query := "SELECT * FROM news"
	rows, err := db.Query(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var n News
		if err := rows.Scan(&n.ID, &n.Title, &n.Content); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		news = append(news, n)
	}

	json.NewEncoder(w).Encode(news)

	/* w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "getNews"}) */
}

func createNews(w http.ResponseWriter, r *http.Request) {
	// implementation of createNews function
}

func updateNews(w http.ResponseWriter, r *http.Request) {
	// implementation of createNews function
}

func deleteNews(w http.ResponseWriter, r *http.Request) {
	// implementation of createNews function
}
