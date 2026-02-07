package models

import (
    "crypto/rand"
    "crypto/sha256"
    "encoding/base64"
    "time"
)

type User struct {
    ID         string    `json:"id"`
    Login      string    `json:"login"`
    Email      string    `json:"email"`
    Password   string    `json:"-"` // хешированный
    PubKey     string    `json:"pubkey"`
    PrivKey    string    `json:"-"`
    IPAddress  string    `json:"ip_address"`
    Active     bool      `json:"active"`
    CreatedAt  time.Time `json:"created_at"`
    LastLogin  time.Time `json:"last_login"`
}

type Session struct {
    UserID    string
    ExpiresAt time.Time
}

func GenerateKeys() (pubKey, privKey string, err error) {
    priv := make([]byte, 32)
    if _, err = rand.Read(priv); err != nil {
        return
    }
    privKey = base64.StdEncoding.EncodeToString(priv)
    
    h := sha256.New()
    h.Write(priv)
    pubKey = base64.StdEncoding.EncodeToString(h.Sum(nil))
    return
}

func HashPassword(password string) string {
    h := sha256.New()
    h.Write([]byte(password))
    return base64.StdEncoding.EncodeToString(h.Sum(nil))
}
