package main

import (
    "embed"
    "log"
    "net/http"
    
    "github.com/gorilla/mux"
    "vpn-service/internal/auth"
    "vpn-service/internal/database"
    "vpn-service/internal/tun"
    "vpn-service/web"
)

var staticFS embed.FS

func main() {
    db, err := database.NewDB("vpn.db")
    if err != nil {
        log.Fatal("Database:", err)
    }
    defer db.Close()

    wg, err := tun.NewWireGuard()
    if err != nil {
        log.Fatal("WireGuard:", err)
    }
    defer wg.Shutdown()
    
    r := mux.NewRouter()

    r.PathPrefix("/static/").Handler(http.FileServer(http.FS(staticFS)))

    r.HandleFunc("/login", auth.LoginHandler(db)).Methods("GET", "POST")
    r.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
        http.Redirect(w, r, "/login", http.StatusSeeOther)
    })

    api := r.PathPrefix("/api").Subrouter()
    api.Use(auth.AuthMiddleware(db))
    
    api.HandleFunc("/users", getUsers(db)).Methods("GET")
    api.HandleFunc("/users", createUser(db, wg)).Methods("POST")
    api.HandleFunc("/users/{id}/config", getClientConfig(db, wg)).Methods("GET")
    api.HandleFunc("/users/{id}", deleteUser(db, wg)).Methods("DELETE")
    
    r.PathPrefix("/").HandlerFunc(dashboardHandler(staticFS))
    
    log.Println("VPN Service running on :8080")
    log.Fatal(http.ListenAndServe(":8080", r))
}

func getUsers(db *database.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        users, err := db.GetAllUsers()
        if err != nil {
            http.Error(w, err.Error(), 500)
            return
        }
        json.NewEncoder(w).Encode(users)
    }
}

func createUser(db *database.DB, wg *tun.WireGuard) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var req struct {
            Login string `json:"login"`
            Email string `json:"email"`
        }
        json.NewDecoder(r.Body).Decode(&req)
        
        user := &models.User{
            ID:      generateID(),
            Login:   req.Login,
            Email:   req.Email,
            Password: models.HashPassword("defaultpass123"),
        }
        user.PubKey, user.PrivKey, _ = models.GenerateKeys()
        user.IPAddress = "10.0.0." + user.ID[:3]
        
        if err := db.CreateUser(user); err != nil {
            http.Error(w, err.Error(), 500)
            return
        }
        
        wg.AddClient(user)
        w.WriteHeader(201)
    }
}

func getClientConfig(db *database.DB, wg *tun.WireGuard) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        vars := mux.Vars(r)
        id := vars["id"]

        w.Header().Set("Content-Disposition", "attachment; filename=wg.conf")
        fmt.Fprint(w, "[Interface]\n...")
    }
}
