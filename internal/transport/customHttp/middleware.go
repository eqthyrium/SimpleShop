package customHttp

import (
	"SimpleShop/internal/domain"
	"SimpleShop/internal/service/session"
	"SimpleShop/pkg/logger"
	"context"
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"
)

/*
ToDo

1)Security Middleware (done)

2) Logging Middleware (done)

3) Token Validation Middleware (done)
Check presence of token in cookie from client request
├── True
│   ├── Check for verification of token
│   │   ├── True
│   │   │   ├── Check UserId existence in the MapUUID
│   │   │   │   ├── True
│   │   │   │   │   ├── Check expiration of the token
│   │   │   │   │   │   ├── True
│   │   │   │   │   │   │   ├── Check if token's time surpasses threshold time
│   │   │   │   │   │   │   │   ├── True
│   │   │   │   │   │   │   │   │   ├── Extend token time by 45 minutes and send new token in cookie to client
│   │   │   │   │   │   │   │   └── False
│   │   │   │   │   │   │   │       ├── Send appropriate webpage
│   │   │   │   │   │   └── False
│   │   │   │   │   │       ├── Send guest homepage (token expired)
│   │   │   │   └── False
│   │   │   │       ├── Send guest homepage (UserId not in MapUUID or another UUID found)
│   │   └── False
│   │       ├── Send guest webpage (failed token verification)
└── False
├── Send guest webpage

4) I have to implement CSRF checking middleware

5) Panic middleware
*/

var customLogger *logger.CustomLogger = logger.NewLogger().GetLoggerObject("../logging/info.log", "../logging/error.log", "../logging/debug.log", "Middleware")
var CSRFMap map[string]string = make(map[string]string)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		customLogger.InfoLogger.Print(fmt.Sprintf("Method:[%v], URL_Path: %v, Remote_Address: %v\n", r.Method, r.URL.Path, r.RemoteAddr))
		next.ServeHTTP(w, r)
		duration := time.Since(startTime)
		// Here must be reponse's status code
		customLogger.InfoLogger.Print(fmt.Sprintf("The End of the client request, and its Duration:%v\n", duration))
	})
}

func SecurityMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		next.ServeHTTP(w, r)
	})
}

func RoleAdjusterMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		customLogger.DebugLogger.Println("The RoleAdjusterMiddleware is started")

		tokenString, err := session.GetTokenFromCookie(r, "auth_token")

		if err != nil {
			customLogger.DebugLogger.Println("There is an error about getting the token from the cookie!!!")
			if errors.Is(err, http.ErrNoCookie) {
				customLogger.InfoLogger.Println("There is no cookie in the request of the client")
				ctx := context.WithValue(r.Context(), "Role", "Guest")
				next.ServeHTTP(w, r.WithContext(ctx))
			} else {
				customLogger.InfoLogger.Println("There is a problem in the process of Extraction token from cookie")
				customLogger.ErrorLogger.Print(logger.ErrorWrapper("Transport", "RoleAdjusterMiddleware", "There is a problem in the process of Extraction token from cookie", err))
				session.DeleteSessionCookie(w, "auth_token")
				serverError(w)
			}
			return
		}

		if tokenString == "" {
			customLogger.DebugLogger.Println("Entered into the absence of the token in the cookie")
			customLogger.InfoLogger.Println("There is cookie, but there is no token")
			session.DeleteSessionCookie(w, "auth_token")
			ctx := context.WithValue(r.Context(), "Role", "Guest")
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		err = session.VerifyToken(tokenString)
		if err != nil {
			customLogger.DebugLogger.Println("Entered into error handling of verification of the token")
			session.DeleteSessionCookie(w, "auth_token")

			if errors.Is(err, domain.ErrInvalidToken) {
				customLogger.InfoLogger.Println("There is an invalid token")
				http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
			} else {
				customLogger.ErrorLogger.Print(logger.ErrorWrapper("Transport", "RoleAdjusterMiddleware", "There is a problem in the process of verification of the token", err))
				serverError(w)
			}
			return
		}

		extractedToken, err := session.ExtractDataFromToken(tokenString)
		if err != nil {
			customLogger.DebugLogger.Println("Entered into error handling of extraction token")
			customLogger.ErrorLogger.Print(logger.ErrorWrapper("Transport", "RoleAdjusterMiddleware", "There is a problem in the process of extraction of the token", err))
			session.DeleteSessionCookie(w, "auth_token")
			serverError(w)
			return
		}

		if session.MapUUID[extractedToken.UserId] != extractedToken.UUID {
			customLogger.DebugLogger.Println("Entered into error handling of the check up of MappUUID")
			customLogger.InfoLogger.Println("There is not current token for the client")
			session.DeleteSessionCookie(w, "auth_token")
			delete(CSRFMap, session.MapUUID[extractedToken.UserId])
			http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
			return
		}
		customLogger.DebugLogger.Println("Checking part of token's time")
		switch session.CheckTokenTime(extractedToken) {
		case "Expired-Token":
			customLogger.InfoLogger.Println("There is an expired token")
			session.DeleteSessionCookie(w, "auth_token")
			delete(CSRFMap, session.MapUUID[extractedToken.UserId])
			delete(session.MapUUID, extractedToken.UserId)
			http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
			return
		case "Extend-Token":
			extendedToken, err := session.ExtendTokenExistence(extractedToken)
			if err != nil {
				customLogger.ErrorLogger.Print(logger.ErrorWrapper("Transport", "RoleAdjusterMiddleware", "There is a problem in the process of extension of time of the token", err))
				session.DeleteSessionCookie(w, "auth_token")
				serverError(w)
				return
			}
			customLogger.InfoLogger.Println("The member with userId:", extractedToken.UserId, "and its role:", extractedToken.Role, ", its expireTime is refreshed by adding 45 min to previous left time.")
			session.SetTokenToCookie(w, "auth_token", extendedToken)
		}

		ctx := context.WithValue(r.Context(), "Role", extractedToken.Role)
		ctx = context.WithValue(ctx, "UserId", extractedToken.UserId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func CSRFMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		customLogger.DebugLogger.Println("The CSRFMiddleware is started")

		role := r.Context().Value("Role").(string)

		if role != "Guest" && r.Method == http.MethodPost {
			customLogger.DebugLogger.Println("inside of the CSRFMiddleware's Post method checking part ")

			formCSRFText := r.FormValue("csrf_text")
			userId := r.Context().Value("UserId").(int)

			if formCSRFText != CSRFMap[session.MapUUID[userId]] {
				customLogger.InfoLogger.Println("The CSRF attack is detected, its IP is:", r.RemoteAddr)
				if _, ok := CSRFMap[session.MapUUID[userId]]; ok {
					delete(CSRFMap, session.MapUUID[userId])
				}
				delete(session.MapUUID, userId)
				session.DeleteSessionCookie(w, "auth_token")
				http.Redirect(w, r, "/auth/login", http.StatusSeeOther)
				return
			}
			customLogger.DebugLogger.Println("The CSRF checking part was good")

		}

		next.ServeHTTP(w, r)
	})
}

func PanicMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				customLogger.ErrorLogger.Println("Panic:\n", err, string(debug.Stack()))
				serverError(w)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
