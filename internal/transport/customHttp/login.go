package customHttp

import (
	"SimpleShop/internal/domain"
	"SimpleShop/pkg/logger"
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"net/http"

	"SimpleShop/internal/service/session"
)

/*
How does CSRF process work
1) When the client gets the webpage where needs to be inserted data into it, the server must provide the webpage where its csrf token value and the embedded csrf unique string
must be same.
2) When the client sends (with Post method) data upon that fields on the webpage, the server must check whether the embedded data upon html is same to the data inside csrf token
If it is the server must provide the appropriate webpage resource to that request.
Otherwise the server must send back to the client error page 403 Forbidden

1) Generate CSRF unique number
2) I have to set it into the cookie as csrf token
3) Then it has to be embedded  into html, as hidden mark.
*/
func (handler *HandlerHttp) logIn(w http.ResponseWriter, r *http.Request) {
	customLogger.DebugLogger.Println("The logIn handler is activated")

	if r.URL.Path != "/auth/login" {
		handler.InfoLog.Println(errors.New("incorrect request's endpoint"))
		clientError(w, nil, http.StatusNotFound, nil)
		return
	}

	if !(r.Method == http.MethodPost || r.Method == http.MethodGet) {
		handler.InfoLog.Println(errors.New("incorrect request's method"))
		clientError(w, nil, http.StatusMethodNotAllowed, nil)
		return
	}

	role := r.Context().Value("Role").(string)
	if role != "Guest" {
		handler.InfoLog.Println(errors.New("incorrect request's role"))
		clientError(w, nil, http.StatusForbidden, nil)
		return
	}

	if r.Method == http.MethodGet {
		customLogger.DebugLogger.Println("logIn's handler of GET request is activated")
		files := []string{
			"../ui/html/login.tmpl.html",
		}

		tmpl, err := template.ParseFiles(files...)
		if err != nil {
			customLogger.ErrorLogger.Print(logger.ErrorWrapper("Transport", "logIn", "There is a problem in the process of parsing the html files with template", err))
			serverError(w)
			return
		}
		var buf bytes.Buffer

		err = tmpl.ExecuteTemplate(&buf, "login", nil)
		if err != nil {
			customLogger.ErrorLogger.Print(logger.ErrorWrapper("Transport", "logIn", "There is a problem in the process of execution of parsed the html files", err))
			serverError(w)
			return
		}
		w.WriteHeader(http.StatusOK)
		_, err = buf.WriteTo(w)
		if err != nil {
			customLogger.ErrorLogger.Print(logger.ErrorWrapper("Transport", "logIn", "There is a problem in the process of converting data from buffer to the http.ResponseWriter", err))
			serverError(w)
			return
		}

	}

	if r.Method == http.MethodPost {
		customLogger.DebugLogger.Println("logIn's handler of POST request is activated")

		err := r.ParseForm()
		if err != nil {
			customLogger.ErrorLogger.Print(logger.ErrorWrapper("Transport", "logIn", "Failed the parsing the Form of html", err))
			serverError(w)
			return
		}

		email := r.FormValue("email")
		password := r.FormValue("password")

		//flag := r.FormValue("flag")
		//
		//if flag {
		//	authentication()
		//}

		// Think about session based on token here
		tokenSignature, err := handler.Service.LogIn(email, password)
		if err != nil {
			if errors.Is(err, domain.ErrUserNotFound) {
				handler.DebugLog.Println(fmt.Errorf("Function \"logIn\": %w", err))
				clientError(w, []string{"../ui/html/login.tmpl.html"}, http.StatusBadRequest, domain.ErrUserNotFound)
			} else {
				customLogger.ErrorLogger.Print(logger.ErrorWrapper("Transport", "logIn", "There is a problem in the process of giving the tokenSignature", err))
				serverError(w)
			}
			return
		}

		//Cookies
		session.SetTokenToCookie(w, "auth_token", tokenSignature)
		http.Redirect(w, r, "/", http.StatusSeeOther)

	}
}
