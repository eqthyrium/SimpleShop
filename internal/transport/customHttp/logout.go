package customHttp

import (
	"SimpleShop/internal/service/session"
	"net/http"
)

func (handler *HandlerHttp) logOut(w http.ResponseWriter, r *http.Request) {
	customLogger.DebugLogger.Println("the logOut handler is activated")

	if r.URL.Path != "/logout" {
		clientError(w, nil, http.StatusNotFound, nil)
		return
	}

	if r.Method != http.MethodPost {
		clientError(w, nil, http.StatusMethodNotAllowed, nil)
		return
	}

	role := r.Context().Value("Role").(string)
	if role == "Guest" {
		handler.DebugLog.Println("The guest is trying to log out")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	userId := r.Context().Value("UserId").(int)
	delete(session.MapUUID, userId)
	session.DeleteSessionCookie(w, "auth_token")
	customLogger.DebugLogger.Println("The cookie is deleted because of the logout operation")
	http.Redirect(w, r, "/", http.StatusSeeOther)

}
