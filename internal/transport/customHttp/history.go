package customHttp

import (
	"SimpleShop/pkg/logger"
	"bytes"
	"errors"
	"html/template"
	"net/http"
)

func (handler *HandlerHttp) history(w http.ResponseWriter, r *http.Request) {
	customLogger.DebugLogger.Println("The history handler is activated")

	if r.URL.Path != "/history" {
		handler.InfoLog.Println(errors.New("incorrect request's endpoint"))
		clientError(w, nil, http.StatusNotFound, nil)
		return
	}

	if r.Method != http.MethodGet {
		handler.InfoLog.Println(errors.New("incorrect request's method"))
		clientError(w, nil, http.StatusMethodNotAllowed, nil)
		return
	}

	role := r.Context().Value("Role").(string)

	if role == "Guest" {
		handler.InfoLog.Println(errors.New("incorrect request's role"))
		clientError(w, nil, http.StatusForbidden, nil)
		return
	}

	var path []string = []string{"../ui/html/history.tmpl.html"}
	userId := r.Context().Value("UserId").(int)

	tmpl, err := template.ParseFiles(path...)
	if err != nil {
		customLogger.ErrorLogger.Print(logger.ErrorWrapper("Transport", "history", "There is a problem in the process of parsing the html files with template", err))
		serverError(w)
		return
	}

	var buf bytes.Buffer
	purchased, liked, err := handler.Service.History(userId)

	if err != nil {

		customLogger.ErrorLogger.Print(logger.ErrorWrapper("Transport", "history", "There is a problem in the process of getting Product entity from data base", err))
		serverError(w)
		return

	}
	err = tmpl.ExecuteTemplate(&buf, "history", map[string]interface{}{
		"Purchased": purchased,
		"Liked":     liked,
	})

	if err != nil {
		customLogger.ErrorLogger.Print(logger.ErrorWrapper("Transport", "history", "There is a problem in the process of rendering template to the buffer", err))
		serverError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = buf.WriteTo(w)
	if err != nil {
		customLogger.ErrorLogger.Print(logger.ErrorWrapper("Transport", "history", "There is a problem in the process of converting data from buffer to the http.ResponseWriter", err))
		serverError(w)
		return
	}

}
