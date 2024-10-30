package customHttp

import (
	"net/http"
)

func (handler *HandlerHttp) Routering() http.Handler {

	Middleware := func(next func(w http.ResponseWriter, r *http.Request)) http.Handler {
		return LoggingMiddleware(SecurityMiddleware(PanicMiddleware(RoleAdjusterMiddleware(CSRFMiddleware((http.HandlerFunc(next)))))))
	}

	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("../ui/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))
	mux.Handle("/", Middleware(handler.homePage))
	mux.Handle("/auth/login", Middleware(handler.logIn))
	mux.Handle("/auth/signup", Middleware(handler.signUp))
	mux.Handle("/logout", Middleware(handler.logOut))
	mux.Handle("/reaction", Middleware(handler.reaction))
	// Example of serving static files

	//mux.Handle("/post", postPagePath)
	return http.HandlerFunc(mux.ServeHTTP)
}
