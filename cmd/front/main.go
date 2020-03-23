package main

import (
	"flag"
	"front-service/cmd/front/app"
	"front-service/pkg/core/auth"
	"github.com/JAbduvohidov/di/pkg/di"
	"github.com/JAbduvohidov/jwt"
	"github.com/JAbduvohidov/mux"
	"log"
	"net"
	"net/http"
)

var (
	hostF     = flag.String("host", "", "Server host")
	portF     = flag.String("port", "", "Server port")
	secretF   = flag.String("secret", "", "Secret key")
	fileUrlF  = flag.String("fileUrl", "", "File Service URL")
	authUrlF  = flag.String("authUrl", "", "Auth Service URL")
	movieUrlF = flag.String("movieUrl", "", "Movie Service URL")
	rateUrlF  = flag.String("rateUrl", "", "Rate Service URL")
)

var (
	EHOST      = "HOST"
	EPORT      = "PORT"
	ESECRET    = "SECRET"
	EFILE_URL  = "FILE_URL"
	EAUTH_URL  = "AUTH_URL"
	EMOVIE_URL = "MOVIE_URL"
	ERATE_URL  = "RATE_URL"
)

func main() {
	flag.Parse()
	host, ok := FlagOrEnv(*hostF, EHOST)
	if !ok {
		log.Panic("can't get host")
	}
	port, ok := FlagOrEnv(*portF, EPORT)
	if !ok {
		log.Panic("can't get port")
	}
	secret, ok := FlagOrEnv(*secretF, ESECRET)
	if !ok {
		log.Panic("can't get secret")
	}
	fileUrl, ok := FlagOrEnv(*fileUrlF, EFILE_URL)
	if !ok {
		log.Panic("can't get file url")
	}
	authUrl, ok := FlagOrEnv(*authUrlF, EAUTH_URL)
	if !ok {
		log.Panic("can't get auth url")
	}
	movieUrl, ok := FlagOrEnv(*movieUrlF, EMOVIE_URL)
	if !ok {
		log.Panic("can't get movie url")
	}
	rateUrl, ok := FlagOrEnv(*rateUrlF, ERATE_URL)
	if !ok {
		log.Panic("can't get rate url")
	}
	addr := net.JoinHostPort(host, port)
	start(addr, jwt.Secret(secret), auth.FUrl(fileUrl), auth.AUrl(authUrl), auth.MUrl(movieUrl), auth.RUrl(rateUrl))
}

func start(addr string, secret jwt.Secret, fileUrl auth.FUrl, authUrl auth.AUrl, movieUrl auth.MUrl, rateUrl auth.RUrl) {
	container := di.NewContainer()

	err := container.Provide(
		app.NewServer,
		mux.NewExactMux,
		func() *jwt.Secret { return &secret },
		func() *auth.Client { return auth.NewClient(movieUrl, rateUrl, fileUrl, authUrl)},
	)

	if err != nil {
		log.Print("can't provide di: ", err)
		return
	}
	container.Start()
	var appServer *app.Server
	container.Component(&appServer)

	panic(http.ListenAndServe(addr, appServer))
}
