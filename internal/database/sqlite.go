package database

import (
    "database/sql"
    "log"
    "time"
    
    "vpn-service/internal/models"
    _ "github.com/mattn/go-sqlite3"
)

type DB struct {
    *sql.DB
}

func NewDB(path string) (*DB, error) {
    db, err := sql.Open("sqlite3", path)
    if err != nil {
        return nil, err
    }
    
    db.SetConnMaxLifetime(time.Minute * 3)
    db.SetMaxOpenConns(10)
    db.SetMaxIdleConns(5)
    
    err = db.Ping()
    if err != nil {
        return nil, err
    }
    
    return &DB{db}, initSchema(db)
}

func initSchema(db *sql.DB) error {
    schema := `
    CREATE TABLE IF NOT EXISTS users (
        id TEXT PRIMARY KEY,
        login TEXT UNIQUE NOT NULL,
        email TEXT,
        password TEXT NOT NULL,
        pubkey TEXT NOT NULL,
        privkey TEXT NOT NULL,
        ip_address TEXT,
        active BOOLEAN DEFAULT 1,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        last_login DATETIME
    );
    
    CREATE TABLE IF NOT EXISTS sessions (
        id TEXT PRIMARY KEY,
        user_id TEXT,
        expires_at DATETIME NOT NULL,
        FOREIGN KEY(user_id) REFERENCES users(id)
    );
    `
    _, err := db.Exec(schema)
    return err
}

func (db *DB) CreateUser(user *models.User) error {
    stmt := `
    INSERT INTO users (id, login, email, password, pubkey, privkey, ip_address, active) 
    VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
    _, err := db.Exec(stmt, user.ID, user.Login, user.Email, user.Password, 
                     user.PubKey, user.PrivKey, user.IPAddress, user.Active)
    return err
}

func (db *DB) GetUserByLogin(login string) (*models.User, error) {
    user := &models.User{}
    err := db.QueryRow("SELECT id, login, email, password, pubkey, privkey, ip_address, active, created_at, last_login FROM users WHERE login = ?", login).
        Scan(&user.ID, &user.Login, &user.Email, &user.Password, &user.PubKey, &user.PrivKey, &user.IPAddress, &user.Active, &user.CreatedAt, &user.LastLogin)
    if err != nil {
        return nil, err
    }
    return user, nil
}

func (db *DB) GetAllUsers() ([]models.User, error) {
    rows, err := db.Query("SELECT id, login, email, pubkey, ip_address, active, created_at, last_login FROM users WHERE active = 1")
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var users []models.User
    for rows.Next() {
        var u models.User
        err := rows.Scan(&u.ID, &u.Login, &u.Email, &u.PubKey, &u.IPAddress, &u.Active, &u.CreatedAt, &u.LastLogin)
        if err != nil {
            return nil, err
        }
        users = append(users, u)
    }
    return users, nil
}
