package app

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"front-service/pkg/core/auth"
	"github.com/JAbduvohidov/jwt"
	"github.com/JAbduvohidov/mux"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"time"
)

type Movie struct {
	Id          int64    `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Image       string   `json:"image"`
	Year        string   `json:"year"`
	Country     string   `json:"country"`
	Actors      []string `json:"actors"`
	Genres      []string `json:"genres"`
	Creators    []string `json:"creators"`
	Studio      string   `json:"studio"`
	ExtLink     string   `json:"ext_link"`
}

type Server struct {
	router     *mux.ExactMux
	secret     *jwt.Secret
	authClient *auth.Client
}

func NewServer(router *mux.ExactMux, secret *jwt.Secret, authClient *auth.Client) *Server {
	return &Server{router: router, secret: secret, authClient: authClient}
}

func (s *Server) Start() {
	s.InitRoutes()
}

func (s *Server) Stop() {
	// TODO: make server stop
}

type ErrorDTO struct {
	Errors []string `json:"errors"`
}

func (s *Server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	s.router.ServeHTTP(writer, request)
}

func (s *Server) handleFrontPage() http.HandlerFunc {
	var (
		tpl *template.Template
		err error
	)
	return func(writer http.ResponseWriter, request *http.Request) {
		tpl, err = template.ParseFiles(
			filepath.Join("web/templates", "index.gohtml"),
			filepath.Join("web/templates", "base.gohtml"),
		)
		if err != nil {
			panic(err)
		}

		ctx, _ := context.WithTimeout(context.Background(), time.Second)

		request, err := http.NewRequestWithContext(
			ctx,
			http.MethodGet,
			fmt.Sprintf("%s/api/movies", s.authClient.Murl),
			bytes.NewBuffer([]byte("")),
		)
		if err != nil {
			log.Print(err)
		}
		request.Header.Set("Content-Type", "application/json")

		response, err := http.DefaultClient.Do(request)
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				err = tpl.Execute(writer, struct{}{})
				if err != nil {
					log.Printf("error while executing template %s %v", tpl.Name(), err)
				}
				return
			}
			log.Print(err)
		}
		responseBody, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Print(err)
		}

		var movies []Movie

		err = json.Unmarshal(responseBody, &movies)
		if err != nil {
			log.Print(err)
		}

		for i := 0; i < len(movies); i++ {
			path := string(s.authClient.Furl) + "/" + movies[i].Image
			movies[i].Image = path
		}

		data := &struct {
			Title  string
			Movies []Movie
		}{
			Title:  "MOSEP",
			Movies: movies,
		}

		err = tpl.Execute(writer, data)
		if err != nil {
			log.Printf("error while executing template %s %v", tpl.Name(), err)
		}
	}
}

func (s *Server) handleSearchPage() http.HandlerFunc {
	var (
		tpl *template.Template
		err error
	)
	return func(w http.ResponseWriter, r *http.Request) {
		tpl, err = template.ParseFiles(
			filepath.Join("web/templates", "index.gohtml"),
			filepath.Join("web/templates", "base.gohtml"),
		)
		if err != nil {
			panic(err)
		}

		ctx, _ := context.WithTimeout(context.Background(), time.Second)

		request, err := http.NewRequestWithContext(
			ctx,
			http.MethodGet,
			fmt.Sprintf("%s/api/movies/search?q=%s", s.authClient.Murl, r.FormValue("q")),
			bytes.NewBuffer([]byte("")),
		)
		if err != nil {
			log.Print(err)
		}
		request.Header.Set("Content-Type", "application/json")

		response, err := http.DefaultClient.Do(request)
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				err = tpl.Execute(w, struct{}{})
				if err != nil {
					log.Printf("error while executing template %s %v", tpl.Name(), err)
				}
				return
			}
			log.Print(err)
		}
		responseBody, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Print(err)
		}

		var movies []Movie

		err = json.Unmarshal(responseBody, &movies)
		if err != nil {
			log.Print(err)
		}

		for i := 0; i < len(movies); i++ {
			path := string(s.authClient.Furl) + "/" + movies[i].Image
			movies[i].Image = path
		}
		data := &struct {
			Title  string
			Movies []Movie
		}{
			Title:  "MOSEP",
			Movies: movies,
		}

		err = tpl.Execute(w, data)
		if err != nil {
			log.Printf("error while executing template %s %v", tpl.Name(), err)
		}
	}
}

func (s *Server) handleMoviePage() http.HandlerFunc {
	var (
		tpl *template.Template
		err error
	)
	return func(w http.ResponseWriter, r *http.Request) {
		tpl, err = template.ParseFiles(
			filepath.Join("web/templates", "movie.gohtml"),
			filepath.Join("web/templates", "base.gohtml"),
		)
		if err != nil {
			panic(err)
		}

		ctx, _ := context.WithTimeout(context.Background(), time.Second)

		request, err := http.NewRequestWithContext(
			ctx,
			http.MethodGet,
			fmt.Sprintf("%s/api/movies/%s", s.authClient.Murl, r.Context().Value("id")),
			bytes.NewBuffer([]byte("")),
		)
		if err != nil {
			log.Print(err)
		}
		request.Header.Set("Content-Type", "application/json")

		response, err := http.DefaultClient.Do(request)
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				err = tpl.Execute(w, struct{}{})
				if err != nil {
					log.Printf("error while executing template %s %v", tpl.Name(), err)
				}
				return
			}
			log.Print(err)
		}
		responseBody, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Print(err)
		}

		var movie Movie

		err = json.Unmarshal(responseBody, &movie)
		if err != nil {
			log.Print(err)
		}

		path := string(s.authClient.Furl) + "/" + movie.Image
		movie.Image = path

		data := &struct {
			Title string
			Movie Movie
		}{
			Title: movie.Title + " | MOSEP",
			Movie: movie,
		}

		err = tpl.Execute(w, data)
		if err != nil {
			log.Printf("error while executing template %s %v", tpl.Name(), err)
		}
	}
}

func (s *Server) handleRegisterPage() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

	}
}
