package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "image/png"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
)

type ImageGenerationConfig struct {
	Style      string `json:"style"`
	StyleIsURL bool   `json:"styleIsUrl"`
	Scene      string `json:"scene"`
	SceneIsURL bool   `json:"sceneIsUrl"`
	Size       string `json:"size"`
	NumImages  int    `json:"numImages"`
}

type Image struct {
	ID     string `json:"id"`
	Prompt string `json:"prompt"`
	URL    string `json:"url"`
	Data   []byte `json:"data"`
}

var db *sql.DB

func GenerateImageOptions(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Accept")
}

func GenerateImage(w http.ResponseWriter, req *http.Request) {
	var config ImageGenerationConfig

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	err := json.NewDecoder(req.Body).Decode(&config)
	if err != nil {
		http.Error(w, "Bad request!", http.StatusBadRequest)
		return
	}

	// TODO: Re-use this
	var stylePrompt string
	if config.StyleIsURL {
		stylePrompt = PromptFromURL(config.Style, true)
	} else {
		stylePrompt = config.Style
	}

	var scenePrompt string
	if config.SceneIsURL {
		scenePrompt = PromptFromURL(config.Scene, false)
	} else {
		scenePrompt = config.Scene
	}

	var ids []string
	for i := 0; i < config.NumImages; i++ {
		imageUrl, prompt := GenerateDallEImage(scenePrompt, stylePrompt, config.Size)
		id := insertImageFromUrl(imageUrl, prompt)
		ids = append(ids, id)
	}

	err = json.NewEncoder(w).Encode(ids)
	if err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
	}
}

func GetImageIds(w http.ResponseWriter, req *http.Request) {
	ids := selectImageIds()

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	err := json.NewEncoder(w).Encode(ids)
	if err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
	}
}

func GetImageById(w http.ResponseWriter, req *http.Request) {
	id := req.PathValue("id")
	req.ParseForm()
	_, isHighResolution := req.Form["isHighResolution"]
	image := selectImageById(id, isHighResolution)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	err := json.NewEncoder(w).Encode(image)
	if err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
	}
}

func HealthCheck(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprintf(w, "ok\n")
}

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, req)
		fmt.Println(req.Method, req.URL.Path, time.Since(start))
	})
}

func main() {
	log.SetOutput(os.Stdout)

	// connStr := "postgresql://myuser:mypassword@0.0.0.0/mydb?sslmode=disable"
	connStr := "postgresql://myuser:mypassword@db/mydb?sslmode=disable"
	// Connect to database
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	router := http.NewServeMux()

	router.HandleFunc("/health-check", HealthCheck)
	router.HandleFunc("GET /images", GetImageIds)
	router.HandleFunc("GET /images/{id}", GetImageById)
	router.HandleFunc("POST /images", GenerateImage)
	router.HandleFunc("OPTIONS /images", GenerateImageOptions)

	server := http.Server{
		Addr:    ":8090",
		Handler: Logging(router),
	}

	server.ListenAndServe()
}
