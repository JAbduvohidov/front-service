package app

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"front-service/pkg/core/auth"
	"front-service/pkg/core/models"
	"github.com/JAbduvohidov/jwt"
	"github.com/JAbduvohidov/mux"
	jwt2 "github.com/JAbduvohidov/mux/middleware/jwt"
	"html/template"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
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

		user, ok := r.Context().Value(jwt2.ContextKey("jwt")).(*Payload)

		data := &struct {
			Title        string
			Registration bool
			Login        bool
			Authorized   bool
			User         *Payload
			Movies       []Movie
		}{
			Title:      "MOSEP",
			Authorized: ok,
			User:       user,
			Movies:     movies,
		}

		err = tpl.Execute(w, data)
		if err != nil {
			log.Printf("error while executing template %s %v", tpl.Name(), err)
		}
	}
}

func (s *Server) handleDeleteMovie() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, _ := context.WithTimeout(context.Background(), time.Second)

		request, err := http.NewRequestWithContext(
			ctx,
			http.MethodDelete,
			fmt.Sprintf("%s/api/movies/%s", s.authClient.Murl, r.Context().Value("id")),
			bytes.NewBuffer([]byte("")),
		)
		if err != nil {
			log.Print(err)
		}
		token, err := r.Cookie("token")
		if err != nil {
			log.Print("error getting token", err)
		}
		request.Header.Set("Authorization", "Bearer "+token.Value)

		_, err = http.DefaultClient.Do(request)
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
				return
			}
			log.Print(err)
		}
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
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

		user, ok := r.Context().Value(jwt2.ContextKey("jwt")).(*Payload)

		data := &struct {
			Title        string
			Registration bool
			Login        bool
			Authorized   bool
			User         *Payload
			Movies       []Movie
		}{
			Title:      "MOSEP",
			Authorized: ok,
			User:       user,
			Movies:     movies,
		}

		err = tpl.Execute(w, data)
		if err != nil {
			log.Printf("error while executing template %s %v", tpl.Name(), err)
		}
	}
}

func (s *Server) handleGetMoviePage() http.HandlerFunc {
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
		user, ok := r.Context().Value(jwt2.ContextKey("jwt")).(*Payload)

		data := &struct {
			Title        string
			Registration bool
			Login        bool
			Movie        Movie
			Authorized   bool
			User         *Payload
		}{
			Title:      movie.Title + " | MOSEP",
			Movie:      movie,
			Authorized: ok,
			User:       user,
		}

		err = tpl.Execute(w, data)
		if err != nil {
			log.Printf("error while executing template %s %v", tpl.Name(), err)
		}
	}
}

func (s *Server) handlePostMoviePage() http.HandlerFunc {
	var (
		tpl *template.Template
	)
	return func(w http.ResponseWriter, r *http.Request) {
		type FileURL struct {
			Id  string
			URL string
		}
		var fileInfo []FileURL
		nofile := false
		file, head, err := r.FormFile("image")
		if err != nil {
			if !errors.Is(err, http.ErrMissingFile) {
				log.Println(err)
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}
			nofile = true
		}
		if !nofile {
			defer file.Close()

			fileData, err := ioutil.ReadAll(file)
			if err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}

			ctx, _ := context.WithTimeout(r.Context(), 5*time.Second)

			var b bytes.Buffer
			writer := multipart.NewWriter(&b)
			fw, err := writer.CreateFormFile("file", head.Filename)
			if err != nil {
				fmt.Println(err)
				return
			}
			_, err = fw.Write(fileData)
			if err != nil {
				fmt.Println(err)
				return
			}
			writer.Close()

			req, err := http.NewRequestWithContext(
				ctx,
				http.MethodPost,
				fmt.Sprintf("%s/files", s.authClient.Furl),
				&b)
			if err != nil {
				fmt.Printf("%v\n", err)
			}

			req.Header.Set("Content-Type", writer.FormDataContentType())

			response, err := http.DefaultClient.Do(req)
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

			err = json.Unmarshal(responseBody, &fileInfo)
			if err != nil {
				log.Print(err)
			}
		}

		fileUrl := ""
		if len(fileInfo) != 0 {
			fileUrl = fileInfo[0].URL
		}

		movieId, _ := strconv.Atoi(strings.TrimSpace(r.Context().Value("id").(string)))

		movie := Movie{
			Id:          int64(movieId),
			Title:       r.FormValue("title"),
			Description: r.FormValue("description"),
			Image:       fileUrl,
			Year:        r.FormValue("year"),
			Country:     r.FormValue("country"),
			Actors:      strings.Split(r.FormValue("actors"), ", "),
			Genres:      strings.Split(r.FormValue("genres"), ", "),
			Creators:    strings.Split(r.FormValue("creators"), ", "),
			Studio:      r.FormValue("creators"),
			ExtLink:     r.FormValue("link"),
		}

		ctx, _ := context.WithTimeout(context.Background(), time.Second)

		userB, _ := json.Marshal(movie)
		log.Println("*** movie json value", string(userB))

		request, err := http.NewRequestWithContext(
			ctx,
			http.MethodPost,
			fmt.Sprintf("%s/api/movies", s.authClient.Murl),
			bytes.NewBuffer(userB),
		)
		if err != nil {
			log.Print("error creating request:", err)
		}
		token, err := r.Cookie("token")
		if err != nil {
			log.Print("error getting token", err)
		}
		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("Authorization", "Bearer "+token.Value)

		response2, err := http.DefaultClient.Do(request)
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				return
			}
			log.Print("error sending request:", err)
		}
		responseBody2, err := ioutil.ReadAll(response2.Body)
		if err != nil {
			log.Print("error reading response body", err)
		}
		log.Println("*** response body", string(responseBody2))

		tpl, err = template.ParseFiles(
			filepath.Join("web/templates", "movie.gohtml"),
			filepath.Join("web/templates", "base.gohtml"),
		)
		if err != nil {
			panic(err)
		}

		ctx, _ = context.WithTimeout(context.Background(), time.Second)

		request, err = http.NewRequestWithContext(
			ctx,
			http.MethodGet,
			fmt.Sprintf("%s/api/movies/%d", s.authClient.Murl, movieId),
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

		err = json.Unmarshal(responseBody, &movie)
		if err != nil {
			log.Print(err)
		}

		path := string(s.authClient.Furl) + "/" + movie.Image
		movie.Image = path
		user, ok := r.Context().Value(jwt2.ContextKey("jwt")).(*Payload)

		data := &struct {
			Title        string
			Registration bool
			Login        bool
			Movie        Movie
			Authorized   bool
			User         *Payload
		}{
			Title:      movie.Title + " | MOSEP",
			Movie:      movie,
			Authorized: ok,
			User:       user,
		}

		err = tpl.Execute(w, data)
		if err != nil {
			log.Printf("error while executing template %s %v", tpl.Name(), err)
		}
	}
}

func (s *Server) handleGetNewMoviePage() http.HandlerFunc {
	var (
		tpl *template.Template
		err error
	)
	return func(w http.ResponseWriter, r *http.Request) {
		tpl, err = template.ParseFiles(
			filepath.Join("web/templates", "new_movie.gohtml"),
			filepath.Join("web/templates", "base.gohtml"),
		)
		if err != nil {
			panic(err)
		}

		user, ok := r.Context().Value(jwt2.ContextKey("jwt")).(*Payload)

		data := &struct {
			Title        string
			Registration bool
			Login        bool
			Authorized   bool
			User         *Payload
		}{
			Title:      "Add movie | MOSEP",
			Authorized: ok,
			User:       user,
		}

		err = tpl.Execute(w, data)
		if err != nil {
			log.Printf("error while executing template %s %v", tpl.Name(), err)
		}
	}
}

func (s *Server) handlePostNewMoviePage() http.HandlerFunc {
	var (
		tpl *template.Template
		err error
	)
	return func(w http.ResponseWriter, r *http.Request) {
		tpl, err = template.ParseFiles(
			filepath.Join("web/templates", "new_movie.gohtml"),
			filepath.Join("web/templates", "base.gohtml"),
		)
		if err != nil {
			panic(err)
		}

		file, head, err := r.FormFile("image")
		if err != nil {
			log.Println(err)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		defer file.Close()

		fileData, err := ioutil.ReadAll(file)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		ctx, _ := context.WithTimeout(r.Context(), 5*time.Second)

		var b bytes.Buffer
		writer := multipart.NewWriter(&b)
		fw, err := writer.CreateFormFile("file", head.Filename)
		if err != nil {
			fmt.Println(err)
			return
		}
		_, err = fw.Write(fileData)
		if err != nil {
			fmt.Println(err)
			return
		}
		writer.Close()

		req, err := http.NewRequestWithContext(
			ctx,
			http.MethodPost,
			fmt.Sprintf("%s/files", s.authClient.Furl),
			&b)
		if err != nil {
			fmt.Printf("%v\n", err)
		}

		req.Header.Set("Content-Type", writer.FormDataContentType())

		response, err := http.DefaultClient.Do(req)
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

		type FileURL struct {
			Id  string
			URL string
		}
		var fileInfo []FileURL
		err = json.Unmarshal(responseBody, &fileInfo)
		if err != nil {
			log.Print(err)
		}
		movie := Movie{
			Id:          0,
			Title:       r.FormValue("title"),
			Description: r.FormValue("description"),
			Image:       fileInfo[0].URL,
			Year:        r.FormValue("year"),
			Country:     r.FormValue("country"),
			Actors:      strings.Split(r.FormValue("actors"), ", "),
			Genres:      strings.Split(r.FormValue("genres"), ", "),
			Creators:    strings.Split(r.FormValue("creators"), ", "),
			Studio:      r.FormValue("creators"),
			ExtLink:     r.FormValue("link"),
		}

		ctx, _ = context.WithTimeout(context.Background(), time.Second)

		userB, _ := json.Marshal(movie)
		log.Println("*** movie json value", string(userB))

		request, err := http.NewRequestWithContext(
			ctx,
			http.MethodPost,
			fmt.Sprintf("%s/api/movies", s.authClient.Murl),
			bytes.NewBuffer(userB),
		)
		if err != nil {
			log.Print("error creating request:", err)
		}
		token, err := r.Cookie("token")
		if err != nil {
			log.Print("error getting token", err)
		}
		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("Authorization", "Bearer "+token.Value)

		response2, err := http.DefaultClient.Do(request)
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				return
			}
			log.Print("error sending request:", err)
		}
		responseBody2, err := ioutil.ReadAll(response2.Body)
		if err != nil {
			log.Print("error reading response body", err)
		}
		log.Println("*** response body", string(responseBody2))
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	}
}

func (s *Server) handleGetRegisterPage() http.HandlerFunc {
	var (
		tpl *template.Template
		err error
	)
	return func(w http.ResponseWriter, r *http.Request) {
		tpl, err = template.ParseFiles(
			filepath.Join("web/templates", "register.gohtml"),
			filepath.Join("web/templates", "base.gohtml"),
		)
		if err != nil {
			panic(err)
		}

		data := &struct {
			Title        string
			Registration bool
			Login        bool
			Authorized   bool
			IsError      bool
		}{
			Title:        "Registration | MOSEP",
			Registration: true,
		}

		err = tpl.Execute(w, data)
		if err != nil {
			log.Printf("error while executing template %s %v", tpl.Name(), err)
		}
	}
}

func (s *Server) handlePostRegisterPage() http.HandlerFunc {
	var (
		tpl *template.Template
		err error
	)
	return func(w http.ResponseWriter, r *http.Request) {
		tpl, err = template.ParseFiles(
			filepath.Join("web/templates", "register.gohtml"),
			filepath.Join("web/templates", "base.gohtml"),
		)
		if err != nil {
			panic(err)
		}

		ctx, _ := context.WithTimeout(context.Background(), time.Second)

		user := models.UserRequestDTO{
			Id:       0,
			Name:     r.FormValue("name"),
			Surname:  r.FormValue("surname"),
			Login:    r.FormValue("login"),
			Password: r.FormValue("password"),
			Avatar:   "",
		}

		userB, _ := json.Marshal(user)

		request, err := http.NewRequestWithContext(
			ctx,
			http.MethodPost,
			fmt.Sprintf("%s/api/users", s.authClient.Aurl),
			bytes.NewBuffer(userB),
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

		var errorData models.ErrorDTO

		err = json.Unmarshal(responseBody, &errorData)
		if err != nil {
			log.Print(err)
		}

		if len(errorData.Errors) == 0 {
			ctx, _ := context.WithTimeout(context.Background(), time.Second)

			token := models.TokenRequestDTO{
				Login:    user.Login,
				Password: user.Password,
			}

			tokenB, _ := json.Marshal(token)

			request, err := http.NewRequestWithContext(
				ctx,
				http.MethodPost,
				fmt.Sprintf("%s/api/tokens", s.authClient.Aurl),
				bytes.NewBuffer(tokenB),
			)
			if err != nil {
				log.Print(err)
			}
			request.Header.Set("Content-Type", "application/json")

			response, err := http.DefaultClient.Do(request)
			if err != nil {
				if errors.Is(err, context.DeadlineExceeded) {
					return
				}
				log.Print(err)
			}
			responseBody, err := ioutil.ReadAll(response.Body)
			if err != nil {
				log.Print(err)
			}

			var tokenData models.TokenResponseDTO
			var errorsData models.ErrorDTO

			err = json.Unmarshal(responseBody, &tokenData)
			if err != nil {
				log.Print(err)
			}

			err = json.Unmarshal(responseBody, &errorsData)
			if err != nil {
				log.Print(err)
			}

			if len(errorsData.Errors) != 0 {
				_, _ = w.Write([]byte("<h1>" + errorsData.Errors[0] + "</h1>"))
				return
			}

			cookie := http.Cookie{
				Name:     "token",
				Value:    tokenData.Token,
				HttpOnly: true,
			}
			http.SetCookie(w, &cookie)

			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}

		data := &struct {
			Title        string
			Registration bool
			Login        bool
			Authorized   bool
			IsError      bool
			ErrorMessage string
		}{
			Title:        "Registration | MOSEP",
			Registration: true,
			IsError:      true,
			ErrorMessage: errorData.Errors[0],
		}

		err = tpl.Execute(w, data)
		if err != nil {
			log.Printf("error while executing template %s %v", tpl.Name(), err)
		}
	}
}

func (s *Server) handleGetLoginPage() http.HandlerFunc {
	var (
		tpl *template.Template
		err error
	)
	return func(w http.ResponseWriter, r *http.Request) {
		tpl, err = template.ParseFiles(
			filepath.Join("web/templates", "login.gohtml"),
			filepath.Join("web/templates", "base.gohtml"),
		)
		if err != nil {
			panic(err)
		}

		data := &struct {
			Title        string
			Registration bool
			Login        bool
			Authorized   bool
			IsError      bool
			ErrorMessage string
		}{
			Title: "Login | MOSEP",
			Login: true,
		}

		err = tpl.Execute(w, data)
		if err != nil {
			log.Printf("error while executing template %s %v", tpl.Name(), err)
		}
	}
}

func (s *Server) handlePostLoginPage() http.HandlerFunc {
	var (
		tpl *template.Template
	)
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, _ := context.WithTimeout(context.Background(), time.Second)

		token := models.TokenRequestDTO{
			Login:    r.FormValue("login"),
			Password: r.FormValue("password"),
		}

		tokenB, _ := json.Marshal(token)

		request, err := http.NewRequestWithContext(
			ctx,
			http.MethodPost,
			fmt.Sprintf("%s/api/tokens", s.authClient.Aurl),
			bytes.NewBuffer(tokenB),
		)
		if err != nil {
			log.Print(err)
		}
		request.Header.Set("Content-Type", "application/json")

		response, err := http.DefaultClient.Do(request)
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				return
			}
			log.Print(err)
		}
		responseBody, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Print(err)
		}

		var tokenData models.TokenResponseDTO
		var errorsData models.ErrorDTO

		err = json.Unmarshal(responseBody, &tokenData)
		if err != nil {
			log.Print(err)
		}

		err = json.Unmarshal(responseBody, &errorsData)
		if err != nil {
			log.Print(err)
		}

		if len(errorsData.Errors) == 0 {
			cookie := http.Cookie{
				Name:     "token",
				Value:    tokenData.Token,
				HttpOnly: true,
			}
			http.SetCookie(w, &cookie)

			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}

		tpl, err = template.ParseFiles(
			filepath.Join("web/templates", "login.gohtml"),
			filepath.Join("web/templates", "base.gohtml"),
		)
		if err != nil {
			panic(err)
		}

		data := &struct {
			Title        string
			Registration bool
			Login        bool
			IsError      bool
			ErrorMessage string
			Authorized   bool
		}{
			Title:        "Login | MOSEP",
			Login:        true,
			IsError:      true,
			ErrorMessage: errorsData.Errors[0],
		}

		err = tpl.Execute(w, data)
		if err != nil {
			log.Printf("error while executing template %s %v", tpl.Name(), err)
		}
	}
}

func (s *Server) handlePostLogoutPage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie := http.Cookie{
			Name:     "token",
			Value:    "",
			Expires:  time.Now(),
			HttpOnly: true,
		}
		http.SetCookie(w, &cookie)

		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	}
}

func (s *Server) handleGetProfilePage() http.HandlerFunc {
	var (
		tpl *template.Template
		err error
	)
	return func(w http.ResponseWriter, r *http.Request) {
		tpl, err = template.ParseFiles(
			filepath.Join("web/templates", "profile.gohtml"),
			filepath.Join("web/templates", "base.gohtml"),
		)
		if err != nil {
			panic(err)
		}

		ctx, _ := context.WithTimeout(context.Background(), time.Second)

		userPayload, _ := r.Context().Value(jwt2.ContextKey("jwt")).(*Payload)

		request, err := http.NewRequestWithContext(
			ctx,
			http.MethodGet,
			fmt.Sprintf("%s/api/users/%d", s.authClient.Aurl, userPayload.Id),
			bytes.NewBuffer([]byte("")),
		)
		if err != nil {
			log.Print(err)
		}

		cookie, err := r.Cookie("token")
		if err != nil {
			log.Println("no cookie found :(")
			return
		}

		request.Header.Set("Authorization", "Bearer "+cookie.Value)

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

		var user models.UserRequestDTO

		err = json.Unmarshal(responseBody, &user)
		if err != nil {
			log.Print(err)
			return
		}

		data := &struct {
			Title        string
			Registration bool
			Login        bool
			Authorized   bool
			IsError      bool
			User         struct {
				User2 models.UserRequestDTO
				*Payload
			}
			ErrorMessage string
		}{
			Title:      "Profile | MOSEP",
			Authorized: true,
			User: struct {
				User2 models.UserRequestDTO
				*Payload
			}{User2: user, Payload: userPayload},
		}

		err = tpl.Execute(w, data)
		if err != nil {
			log.Printf("error while executing template %s %v", tpl.Name(), err)
		}
	}
}

func (s *Server) handlePostProfilePage() http.HandlerFunc {
	var (
		tpl *template.Template
		err error
	)
	return func(w http.ResponseWriter, r *http.Request) {
		tpl, err = template.ParseFiles(
			filepath.Join("web/templates", "profile.gohtml"),
			filepath.Join("web/templates", "base.gohtml"),
		)
		if err != nil {
			panic(err)
		}

		userPayload, _ := r.Context().Value(jwt2.ContextKey("jwt")).(*Payload)

		userReq := models.UserRequestDTO{
			Name:     r.FormValue("name"),
			Surname:  r.FormValue("surname"),
			Password: r.FormValue("password"),
		}

		ctx, _ := context.WithTimeout(context.Background(), time.Second)

		userReqB, _ := json.Marshal(userReq)

		request, err := http.NewRequestWithContext(
			ctx,
			http.MethodPost,
			fmt.Sprintf("%s/api/users/%d", s.authClient.Aurl, userPayload.Id),
			bytes.NewBuffer(userReqB),
		)
		if err != nil {
			log.Print("error creating request:", err)
		}
		token, err := r.Cookie("token")
		if err != nil {
			log.Print("error getting token", err)
		}
		request.Header.Set("Content-Type", "application/json")
		request.Header.Set("Authorization", "Bearer "+token.Value)

		response2, err := http.DefaultClient.Do(request)
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				return
			}
			log.Print("error sending request:", err)
		}
		responseBody2, err := ioutil.ReadAll(response2.Body)
		if err != nil {
			log.Print("error reading response body", err)
		}
		log.Println("*** response body", string(responseBody2))

		ctx, _ = context.WithTimeout(context.Background(), time.Second)

		request, err = http.NewRequestWithContext(
			ctx,
			http.MethodGet,
			fmt.Sprintf("%s/api/users/%d", s.authClient.Aurl, userPayload.Id),
			bytes.NewBuffer([]byte("")),
		)
		if err != nil {
			log.Print(err)
		}

		cookie, err := r.Cookie("token")
		if err != nil {
			log.Println("no cookie found :(")
			return
		}

		request.Header.Set("Authorization", "Bearer "+cookie.Value)

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

		var user models.UserRequestDTO

		err = json.Unmarshal(responseBody, &user)
		if err != nil {
			log.Print(err)
			return
		}

		data := &struct {
			Title        string
			Registration bool
			Login        bool
			Authorized   bool
			IsError      bool
			User         struct {
				User2 models.UserRequestDTO
				*Payload
			}
			ErrorMessage string
		}{
			Title:      "Profile | MOSEP",
			Authorized: true,
			User: struct {
				User2 models.UserRequestDTO
				*Payload
			}{User2: user, Payload: userPayload},
		}

		err = tpl.Execute(w, data)
		if err != nil {
			log.Printf("error while executing template %s %v", tpl.Name(), err)
		}
	}
}

func (s *Server) handleUsersPage() http.HandlerFunc {
	var (
		tpl *template.Template
		err error
	)
	return func(w http.ResponseWriter, r *http.Request) {
		tpl, err = template.ParseFiles(
			filepath.Join("web/templates", "users.gohtml"),
			filepath.Join("web/templates", "base.gohtml"),
		)
		if err != nil {
			panic(err)
		}

		ctx, _ := context.WithTimeout(context.Background(), time.Second * 5)

		request, err := http.NewRequestWithContext(
			ctx,
			http.MethodGet,
			fmt.Sprintf("%s/api/users", s.authClient.Aurl),
			bytes.NewBuffer([]byte("")),
		)
		if err != nil {
			log.Print(err)
		}

		token, err := r.Cookie("token")
		if err != nil {
			log.Print("error getting token", err)
		}
		request.Header.Set("Authorization", "Bearer "+token.Value)

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

		var users []User

		err = json.Unmarshal(responseBody, &users)
		if err != nil {
			log.Println(string(responseBody))
			log.Print(err)
		}

		user, ok := r.Context().Value(jwt2.ContextKey("jwt")).(*Payload)

		data := &struct {
			Title        string
			Registration bool
			Login        bool
			Authorized   bool
			User         *Payload
			Users        []User
		}{
			Title:      "Users | MOSEP",
			Authorized: ok,
			User:       user,
			Users:      users,
		}

		err = tpl.Execute(w, data)
		if err != nil {
			log.Printf("error while executing template %s %v", tpl.Name(), err)
		}
	}
}

func (s *Server) handleDeleteUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, _ := context.WithTimeout(context.Background(), time.Second)

		request, err := http.NewRequestWithContext(
			ctx,
			http.MethodDelete,
			fmt.Sprintf("%s/api/users/%s", s.authClient.Aurl, r.Context().Value("id")),
			bytes.NewBuffer([]byte("")),
		)
		if err != nil {
			log.Print(err)
		}
		token, err := r.Cookie("token")
		if err != nil {
			log.Print("error getting token", err)
		}
		request.Header.Set("Authorization", "Bearer "+token.Value)

		response, err := http.DefaultClient.Do(request)
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				http.Redirect(w, r, "/users", http.StatusTemporaryRedirect)
				return
			}
			log.Println("res:", response)
			log.Print(err)
		}
		http.Redirect(w, r, "/users", http.StatusTemporaryRedirect)
	}
}

func (s *Server) handleUpgradeUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, _ := context.WithTimeout(context.Background(), time.Second)

		request, err := http.NewRequestWithContext(
			ctx,
			http.MethodPut,
			fmt.Sprintf("%s/api/users/%s/%s", s.authClient.Aurl, r.Context().Value("id"), r.Context().Value("role")),
			bytes.NewBuffer([]byte("")),
		)
		if err != nil {
			log.Print(err)
		}
		token, err := r.Cookie("token")
		if err != nil {
			log.Print("error getting token", err)
		}
		request.Header.Set("Authorization", "Bearer "+token.Value)

		response, err := http.DefaultClient.Do(request)
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				http.Redirect(w, r, "/users", http.StatusTemporaryRedirect)
				return
			}
			log.Println("res:", response)
			log.Print(err)
		}
		http.Redirect(w, r, "/users", http.StatusTemporaryRedirect)
	}
}
