package customHttp

import (
	"SimpleShop/internal/domain"
	"SimpleShop/pkg/logger"
	"bytes"
	"errors"
	"html/template"
	"net/http"
)

func (handler *HandlerHttp) signUp(w http.ResponseWriter, r *http.Request) {

	customLogger.DebugLogger.Println("The signup handler is activated")
	if r.URL.Path != "/auth/signup" {
		// Think about error handling, and logging it properly
		handler.InfoLog.Println("incorrect request's endpoint")
		clientError(w, nil, http.StatusNotFound, nil)
		return
	}
	if !(r.Method == http.MethodGet || r.Method == http.MethodPost) {
		handler.InfoLog.Println("incorrect request's method")
		clientError(w, nil, http.StatusMethodNotAllowed, nil)
		return
	}

	role := r.Context().Value("Role").(string)
	if role != "Guest" {
		handler.InfoLog.Println("The Non-guest client attempts to use the signUp resource")
		clientError(w, nil, http.StatusForbidden, nil)
		return
	}

	if r.Method == http.MethodGet {
		customLogger.DebugLogger.Println("The signup handler's GET request handler is activated")
		files := []string{
			"../ui/html/signup.tmpl.html",
		}

		tmpl, err := template.ParseFiles(files...)
		if err != nil {
			customLogger.ErrorLogger.Print(logger.ErrorWrapper("Transport", "signUp", "There is a problem in the process of parsing the html files with template", err))
			serverError(w)
			return
		}
		var buf bytes.Buffer
		err = tmpl.ExecuteTemplate(&buf, "signup", nil)
		if err != nil {
			customLogger.ErrorLogger.Print(logger.ErrorWrapper("Transport", "signUp", "There is a problem in the process of execution of the template", err))
			serverError(w)
			return
		}

		w.WriteHeader(http.StatusOK)
		_, err = buf.WriteTo(w)
		if err != nil {
			customLogger.ErrorLogger.Print(logger.ErrorWrapper("Transport", "signUp", "There is a problem in the process of converting data from buffer to the http.ResponseWriter", err))
			serverError(w)
			return
		}
	}

	if r.Method == http.MethodPost {
		customLogger.DebugLogger.Println("The signup handler's POST request handler is activated")

		err := r.ParseForm()
		if err != nil {
			customLogger.ErrorLogger.Print(logger.ErrorWrapper("Transport", "signUp", "There is a problem in the process of parsing the Form of html", err))
			serverError(w)
			return
		}

		nickname := r.FormValue("nickname")
		email := r.FormValue("email")
		password := r.FormValue("password")

		err = handler.Service.SignUp(nickname, email, password)
		if err != nil {
			if errors.Is(err, domain.ErrInvalidCredential) {
				handler.DebugLog.Println("There is invalid entered Credentials")
				clientError(w, []string{"../ui/html/signup.tmpl.html"}, http.StatusBadRequest, err)
			} else {
				customLogger.ErrorLogger.Print(logger.ErrorWrapper("Transport", "signUp", "Failed  Sign up operation", err))
				serverError(w)
			}
			return
		}

		http.Redirect(w, r, "/auth/login", http.StatusSeeOther)

	}
}
