package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "image/png"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	_ "github.com/lib/pq"
)

type ImageGenerationRequest struct {
	Style      string `json:"style"`
	StyleIsURL bool   `json:"styleIsUrl"`
	Scene      string `json:"scene"`
	SceneIsURL bool   `json:"sceneIsUrl"`
	Size       string `json:"size"`
	NumImages  int    `json:"numImages"`
}

func (req ImageGenerationRequest) validate() bool {
	if req.Style == "" {
		return false
	}

	if req.Scene == "" {
		return false
	}

	return true
}

type Image struct {
	ID     string `json:"id"`
	Prompt string `json:"prompt"`
	URL    string `json:"url"`
	Data   []byte `json:"data"`
}

var db *sql.DB

type ImageGenerationResults struct {
	Id    string
	Error error
}

func generateImageTask(scenePrompt string, stylePrompt string, size string, res chan<- ImageGenerationResults, wg *sync.WaitGroup) {
	defer wg.Done()

	imageUrl, prompt, err := GenerateDallEImage(scenePrompt, stylePrompt, size)
	if err != nil {
		res <- ImageGenerationResults{"", err}
		return
	}

	id, err := insertImageFromUrl(imageUrl, prompt)
	if err != nil {
		res <- ImageGenerationResults{"", err}
		return
	}

	res <- ImageGenerationResults{id, nil}
}

func GenerateImage(w http.ResponseWriter, r *http.Request) {
	var req ImageGenerationRequest

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusUnprocessableEntity)
		return
	}

	// ERROR: Add validation
	if ok := req.validate(); !ok {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// TODO: Re-use this
	var stylePrompt string
	var err error
	if req.StyleIsURL {
		stylePrompt, err = PromptFromURL(req.Style, true)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	} else {
		stylePrompt = req.Style
	}

	var scenePrompt string
	if req.SceneIsURL {
		scenePrompt, err = PromptFromURL(req.Scene, false)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	} else {
		scenePrompt = req.Scene
	}

	wg := new(sync.WaitGroup)
	wg.Add(req.NumImages)
	chanRes := make(chan ImageGenerationResults, req.NumImages)
	for i := 0; i < req.NumImages; i++ {
		go generateImageTask(scenePrompt, stylePrompt, req.Size, chanRes, wg)
	}
	wg.Wait()

	var ids []string
	for i := 0; i < req.NumImages; i++ {
		res := <-chanRes

		if res.Error != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		ids = append(ids, res.Id)
	}

	if err := json.NewEncoder(w).Encode(ids); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func GetImageIds(w http.ResponseWriter, _ *http.Request) {
	ids, err := selectImageIds()
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(ids)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func GetImageById(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	r.ParseForm()
	_, isHighResolution := r.Form["isHighResolution"]
	image, err := selectImageById(id, isHighResolution)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(image); err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
		return
	}
}

func HealthCheck(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprintf(w, "ok\n")
}

type Middleware func(http.Handler) http.HandlerFunc

func MiddlewareChain(middlewares ...Middleware) Middleware {
	return func(next http.Handler) http.HandlerFunc {
		for i := len(middlewares) - 1; i >= 0; i-- {
			next = middlewares[i](next)
		}

		return next.ServeHTTP
	}
}

func LoggingMiddleware(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		fmt.Println(r.Method, r.URL.Path, time.Since(start))
	}
}

func CORSMiddleware(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Accept")
		next.ServeHTTP(w, r)
	}
}

func main() {
	log.SetOutput(os.Stdout)

	// connStr := "postgresql://myuser:mypassword@0.0.0.0/mydb?sslmode=disable"
	connStr := "postgresql://myuser:mypassword@db/mydb?sslmode=disable"
	// Connect to database
	var err error
	if db, err = sql.Open("postgres", connStr); err != nil {
		log.Fatal(err)
	}

	router := http.NewServeMux()

	router.HandleFunc("/health-check", HealthCheck)
	router.HandleFunc("GET /images", GetImageIds)
	router.HandleFunc("GET /images/{id}", GetImageById)
	router.HandleFunc("POST /images", GenerateImage)

	middlewareChain := MiddlewareChain(
		LoggingMiddleware,
		CORSMiddleware,
	)

	server := http.Server{
		Addr:    ":8090",
		Handler: middlewareChain(router),
	}

	server.ListenAndServe()
}
