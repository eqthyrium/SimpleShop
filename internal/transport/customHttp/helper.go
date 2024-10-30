package customHttp

import (
	"SimpleShop/internal/domain"
	"SimpleShop/pkg/logger"
	"errors"
	"html/template"
	"net/http"
)

// Remind: After calling this kind of function, you have to return in handler immediately
var customLoggerError *logger.CustomLogger = logger.NewLogger().GetLoggerObject("../logging/info.log", "../logging/error.log", "../logging/debug.log", "ErrorHandling")

func errorWebpage(w http.ResponseWriter, paths []string, status int, err error) {

	customLogger.DebugLogger.Println("errorWebpage handler is activated")
	var base string
	if status == http.StatusInternalServerError || status == http.StatusNotFound {
		paths = append(paths, "../ui/html/error/standard.html")
	}

	if status == http.StatusBadRequest {
		customLogger.DebugLogger.Println("Bad Requested error handling part is opened")

		if errors.Is(err, domain.ErrUserNotFound) {
			customLogger.DebugLogger.Println("ErrUserNotFound part of the code in the errorWebpage")
			base = "login"
			paths = append(paths, "../ui/html/error/login.tmpl.html")

		} else if errors.Is(err, domain.ErrInvalidCredential) {
			customLogger.DebugLogger.Println("ErrInvalidCredential in the errorWebpage")
			base = "signup"
			paths = append(paths, "../ui/html/error/signup.tmpl.html")

		}

	}

	tmpl, err := template.ParseFiles(paths...)
	if err != nil {
		customLoggerError.ErrorLogger.Print(logger.ErrorWrapper("Transport", "ErrorWebpage", "There is a problem in parsing the html files with template function", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	data := struct {
		ErrorMessage string
		Status       int
	}{
		ErrorMessage: http.StatusText(status),
		Status:       status,
	}

	if base == "" {
		err = tmpl.Execute(w, data)
		if err != nil {
			customLoggerError.ErrorLogger.Print(logger.ErrorWrapper("Transport", "ErrorWebpage", "There is a problem in execution the html files with template function", err))
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	err = tmpl.ExecuteTemplate(w, base, nil)
	if err != nil {
		customLoggerError.ErrorLogger.Print(logger.ErrorWrapper("Transport", "ErrorWebpage", "There is a problem in execution the html files with template function", err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

	}
}

func serverError(w http.ResponseWriter) {

	w.WriteHeader(http.StatusInternalServerError)
	errorWebpage(w, []string{}, http.StatusInternalServerError, nil)

}

func clientError(w http.ResponseWriter, paths []string, statusCode int, err error) {

	w.WriteHeader(statusCode)

	if http.StatusNotFound == statusCode {
		errorWebpage(w, nil, http.StatusNotFound, nil)
		return
	}

	if http.StatusBadRequest == statusCode {
		errorWebpage(w, paths, http.StatusBadRequest, err)
		return
	}
}
