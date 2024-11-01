package customHttp

import (
	"SimpleShop/internal/domain"
	"SimpleShop/pkg/logger"
	"bytes"
	"errors"
	"html/template"
	"net/http"
)

func (handler *HandlerHttp) recommendation(w http.ResponseWriter, r *http.Request) {
	customLogger.DebugLogger.Println("The recommendation handler is activated")

	if r.URL.Path != "/recommendation" {
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

	userId := r.Context().Value("UserId").(int)

	behaviour, collaborative, err := handler.Service.Recommendation(userId)
	if err != nil {
		handler.ErrorLog.Println(logger.ErrorWrapper("Transport", "recommendation", "There is a problem with process of getting the behaviour and user orientated products", err))
		serverError(w)
		return
	}

	data := struct {
		Behaviour     []domain.Product
		Collaborative []domain.Product
	}{
		Behaviour:     behaviour,
		Collaborative: collaborative,
	}

	var buf bytes.Buffer

	tmpl, err := template.ParseFiles("../ui/html/recommendation.tmpl.html")
	if err != nil {
		customLogger.ErrorLogger.Print(logger.ErrorWrapper("Transport", "recommendation", "There is a problem in the process of parsing the html files with template", err))
		serverError(w)
		return
	}

	err = tmpl.ExecuteTemplate(&buf, "recommendation", map[string]interface{}{
		"Product": data,
	})

	if err != nil {
		customLogger.ErrorLogger.Print(logger.ErrorWrapper("Transport", "recommendation", "There is a problem in the process of rendering template to the buffer", err))
		serverError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = buf.WriteTo(w)
	if err != nil {
		customLogger.ErrorLogger.Print(logger.ErrorWrapper("Transport", "recommendation", "There is a problem in the process of converting data from buffer to the http.ResponseWriter", err))
		serverError(w)
		return
	}
}
