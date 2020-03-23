package app

import (
	"github.com/JAbduvohidov/mux/middleware/jwt"
	"github.com/JAbduvohidov/mux/middleware/logger"
	"reflect"
)

var (
	Root     = "/"
	Search   = "/search"
	Movieh   = "/{id}"
	Register = "/register"
)

func (s *Server) InitRoutes() {
	jwtMW := jwt.JWT(jwt.SourceCookie, reflect.TypeOf((*Payload)(nil)).Elem(), *s.secret)
	//authMW := authenticated.Authenticated(jwt.IsContextNonEmpty, true, Root)

	s.router.GET(Root, s.handleFrontPage(), jwtMW, logger.Logger("HTTP"))
	s.router.GET(Search, s.handleSearchPage(), jwtMW, logger.Logger("HTTP"))
	s.router.GET(Movieh, s.handleMoviePage(), jwtMW, logger.Logger("HTTP"))
	s.router.GET(Register, s.handleRegisterPage(), jwtMW, logger.Logger("HTTP"))
}
