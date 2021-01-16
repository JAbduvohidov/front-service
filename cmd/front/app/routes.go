package app

import (
	"context"
	"github.com/JAbduvohidov/mux/middleware/authenticated"
	"github.com/JAbduvohidov/mux/middleware/jwt"
	"github.com/JAbduvohidov/mux/middleware/logger"
	"reflect"
)

var (
	Root          = "/"
	Search        = "/search"
	Register      = "/register"
	Login         = "/login"
	Logout        = "/logout"
	Profile       = "/profile"
	Movies        = "/movies/{id}"
	DeleteMovie   = "/movies/{id}/delete"
	MoviesNew     = "/movies/new"
	Users         = "/users"
	DeleteUser    = "/users/{id}/delete"
	UpgradeUser   = "/users/{id}/upgrade/{role}"
)

func (s *Server) InitRoutes() {
	jwtMW := jwt.JWT(jwt.SourceCookie, reflect.TypeOf((*Payload)(nil)).Elem(), *s.secret)

	s.router.GET(Root,
		s.handleFrontPage(),
		jwtMW,
		logger.Logger("FRONT_PAGE"))

	s.router.POST(Root,
		s.handleFrontPage(),
		jwtMW,
		logger.Logger("FRONT_PAGE"))

	s.router.GET(Search,
		s.handleSearchPage(),
		jwtMW,
		logger.Logger("SEARCH_PAGE"))

	s.router.GET(Register,
		s.handleGetRegisterPage(),
		authenticated.Authenticated(jwt.IsContextNonEmpty, true, Root),
		jwtMW,
		logger.Logger("REGISTER_PAGE"))

	s.router.POST(Register,
		s.handlePostRegisterPage(),
		authenticated.Authenticated(jwt.IsContextNonEmpty, true, Root),
		jwtMW,
		logger.Logger("REGISTER_PAGE"))

	s.router.GET(Login,
		s.handleGetLoginPage(),
		authenticated.Authenticated(jwt.IsContextNonEmpty, true, Root),
		jwtMW,
		logger.Logger("LOGIN_PAGE"))

	s.router.POST(Login,
		s.handlePostLoginPage(),
		authenticated.Authenticated(jwt.IsContextNonEmpty, true, Root),
		jwtMW,
		logger.Logger("LOGIN_PAGE"))

	s.router.GET(Logout,
		s.handlePostLogoutPage(),
		authenticated.Authenticated(func(ctx context.Context) bool {
			return !jwt.IsContextNonEmpty(ctx)
		}, true, Root),
		jwtMW,
		logger.Logger("LOGOUT_PAGE"))

	s.router.GET(Profile,
		s.handleGetProfilePage(),
		authenticated.Authenticated(func(ctx context.Context) bool {
			return !jwt.IsContextNonEmpty(ctx)
		}, true, Root),
		jwtMW,
		logger.Logger("PROFILE_PAGE"))

	s.router.POST(Profile,
		s.handlePostProfilePage(),
		authenticated.Authenticated(func(ctx context.Context) bool {
			return !jwt.IsContextNonEmpty(ctx)
		}, true, Root),
		jwtMW,
		logger.Logger("PROFILE_PAGE"))

	s.router.GET(Movies,
		s.handleGetMoviePage(),
		jwtMW,
		logger.Logger("MOVIE_PAGE"))

	s.router.POST(Movies,
		s.handlePostMoviePage(),
		jwtMW,
		logger.Logger("MOVIE_PAGE"))

	s.router.POST(DeleteMovie,
		s.handleDeleteMovie(),
		jwtMW,
		logger.Logger("MOVIE_PAGE"))

	s.router.GET(MoviesNew,
		s.handleGetNewMoviePage(),
		authenticated.Authenticated(func(ctx context.Context) bool {
			return !jwt.IsContextNonEmpty(ctx)
		}, true, Root),
		jwtMW,
		logger.Logger("NEW_MOVIE_PAGE"))

	s.router.POST(MoviesNew,
		s.handlePostNewMoviePage(),
		authenticated.Authenticated(func(ctx context.Context) bool {
			return !jwt.IsContextNonEmpty(ctx)
		}, true, Root),
		jwtMW,
		logger.Logger("NEW_MOVIE_PAGE"))

	s.router.GET(Users,
		s.handleUsersPage(),
		authenticated.Authenticated(func(ctx context.Context) bool {
			return !jwt.IsContextNonEmpty(ctx)
		}, true, Root),
		jwtMW,
		logger.Logger("USERS_PAGE"))

	s.router.POST(Users,
		s.handleUsersPage(),
		authenticated.Authenticated(func(ctx context.Context) bool {
			return !jwt.IsContextNonEmpty(ctx)
		}, true, Root),
		jwtMW,
		logger.Logger("USERS_PAGE"))

	s.router.POST(DeleteUser,
		s.handleDeleteUser(),
		authenticated.Authenticated(func(ctx context.Context) bool {
			return !jwt.IsContextNonEmpty(ctx)
		}, true, Root),
		jwtMW,
		logger.Logger("USERS_PAGE"))

	s.router.POST(UpgradeUser,
		s.handleUpgradeUser(),
		authenticated.Authenticated(func(ctx context.Context) bool {
			return !jwt.IsContextNonEmpty(ctx)
		}, true, Root),
		jwtMW,
		logger.Logger("USERS_PAGE"))

}
