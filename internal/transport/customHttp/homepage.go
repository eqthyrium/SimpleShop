package customHttp

import (
	"SimpleShop/internal/service/session"
	"SimpleShop/pkg/logger"
	"bytes"
	"html/template"
	"net/http"
)

func (handler *HandlerHttp) homePage(w http.ResponseWriter, r *http.Request) {
	handler.DebugLog.Println("homePage handler is activated")

	if r.URL.Path != "/" {
		handler.InfoLog.Println("incorrect request's endpoint")
		clientError(w, nil, http.StatusNotFound, nil)
		return
	}
	if r.Method != http.MethodPost && r.Method != http.MethodGet {
		handler.InfoLog.Println("incorrect request's method")
		clientError(w, nil, http.StatusMethodNotAllowed, nil)
		return
	}

	if r.Method == http.MethodGet {
		var userId int
		role := r.Context().Value("Role").(string)
		if role != "Guest" {
			userId = r.Context().Value("UserId").(int)
		}

		switch role {
		case "User":
			homePageGet(w, userId, []string{"../ui/html/shoppage.tmpl.html", "../ui/html/shoppageUser.tmpl.html"}, handler)
		case "Guest":
			homePageGet(w, -1, []string{"../ui/html/shoppage.tmpl.html", "../ui/html/shoppageGuest.tmpl.html"}, handler)
		}
	}

	if r.Method == http.MethodPost {
		//
	}

}

func homePageGet(w http.ResponseWriter, userId int, files []string, handler *HandlerHttp) {

	tmpl, err := template.ParseFiles(files...)
	if err != nil {
		customLogger.ErrorLogger.Print(logger.ErrorWrapper("Transport", "homePageGet", "There is a problem in the process of parsing the html files with template", err))
		serverError(w)
		return
	}

	csrfText, err := session.GenerateRandomCSRFText()
	if err != nil {
		customLogger.ErrorLogger.Print(logger.ErrorWrapper("Transport", "homePageGet", "There is a problem in the process of generating random CSRF text", err))
		serverError(w)
		return
	}

	if userId >= 0 {
		CSRFMap[session.MapUUID[userId]] = csrfText
	}

	var buf bytes.Buffer
	data, err := handler.Service.Homepage(userId)

	if err != nil {
		customLogger.ErrorLogger.Print(logger.ErrorWrapper("Transport", "homePageGet", "There is a problem in the process of getting Product entity from data base", err))
		serverError(w)
		return
	}
	err = tmpl.ExecuteTemplate(&buf, "shoppage", map[string]interface{}{
		"CSRFText": csrfText,
		"Product":  data,
	})

	if err != nil {
		customLogger.ErrorLogger.Print(logger.ErrorWrapper("Transport", "homePageGet", "There is a problem in the process of rendering template to the buffer", err))
		serverError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = buf.WriteTo(w)
	if err != nil {
		customLogger.ErrorLogger.Print(logger.ErrorWrapper("Transport", "homePageGet", "There is a problem in the process of converting data from buffer to the http.ResponseWriter", err))
		serverError(w)
		return
	}
}
