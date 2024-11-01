package customHttp

import (
	"SimpleShop/pkg/logger"
	"errors"
	"net/http"
	"strconv"
)

func (handler *HandlerHttp) reaction(w http.ResponseWriter, r *http.Request) {
	customLogger.DebugLogger.Println("The reaction handler is started to operate")

	if r.URL.Path != "/reaction" {
		handler.InfoLog.Println(errors.New("incorrect request's endpoint"))
		clientError(w, nil, http.StatusNotFound, nil)
		return
	}

	if !(r.Method == http.MethodPost) {
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
	err := r.ParseForm()
	if err != nil {
		customLogger.ErrorLogger.Print(logger.ErrorWrapper("Transport", "reaction", "Failed the parsing the Form of html", err))
		serverError(w)
		return
	}

	reaction := r.FormValue("reaction")
	productIdstr := r.FormValue("product_id")
	productId, err := strconv.Atoi(productIdstr)
	if err != nil {
		customLogger.ErrorLogger.Print(logger.ErrorWrapper("Transport", "reaction", "Failed the conversion the string to the integer value", err))
		serverError(w)
		return
	}
	userId := r.Context().Value("UserId").(int)
	switch reaction {
	case "purchase":
		customLogger.DebugLogger.Println("The purchase operation is clicked")
		if err := handler.Service.Purchase(userId, productId); err != nil {
			customLogger.ErrorLogger.Print(logger.ErrorWrapper("Transport", "reaction", "Failed to implement purchase operation", err))
			serverError(w)
			return
		}
	case "like":
		if err := handler.Service.Like(userId, productId); err != nil {
			customLogger.ErrorLogger.Print(logger.ErrorWrapper("Transport", "reaction", "Failed to implement purchase operation", err))
			serverError(w)
			return
		}
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
