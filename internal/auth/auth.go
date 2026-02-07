package auth

import (
    "net/http"
    "time"
    
    "github.com/gorilla/sessions"
    "vpn-service/internal/database"
    "vpn-service/internal/models"
)

var store = sessions.NewCookieStore([]byte("vpn-service-secret-key-2026"))

func LoginHandler(db *database.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        if r.Method == http.MethodPost {
            login := r.FormValue("login")
            password := r.FormValue("password")
            
            user, err := db.GetUserByLogin(login)
            if err != nil || models.HashPassword(password) != user.Password {
                http.Error(w, "Invalid credentials", http.StatusUnauthorized)
                return
            }

            _, _ = db.Exec("UPDATE users SET last_login = CURRENT_TIMESTAMP WHERE id = ?", user.ID)
            
            session, _ := store.Get(r, "vpn-session")
            session.Values["user_id"] = user.ID
            session.Options.MaxAge = 3600 * 24 // 24h
            session.Save(r, w)
            
            http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
            return
        }

        w.Header().Set("Content-Type", "text/html")
        http.ServeFile(w, r, "web/static/login.html")
    }
}

func AuthMiddleware(db *database.DB) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            session, err := store.Get(r, "vpn-session")
            if err != nil || session.Values["user_id"] == nil {
                http.Redirect(w, r, "/login", http.StatusSeeOther)
                return
            }
            
            next.ServeHTTP(w, r)
        })
    }
}
