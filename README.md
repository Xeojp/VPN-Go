# VPN Service на Go

# Бесплатный self-hosted VPN-сервис с веб-панелью и поддержкой.

# Технический стек
Backend: Go 1.22 + WireGuard-Go + Gorilla Mux
DB: SQLite (file-based)
Frontend: Vue 3 + HTMX
Infra: Docker + systemd
Протокол: WireGuard (UDP 51820)

# Конфигурация
.env (опционально):

VPN_PORT=51820
ADMIN_LOGIN=admin
ADMIN_PASS=admin123
DB_PATH=vpn.db
BIND_ADDR=:8080
🔍 API (после логина)

# Список юзеров
curl -H "Authorization: Bearer $TOKEN" http://localhost:8080/api/users

# Создать юзера
curl -X POST -d '{"login":"user1","email":"user1@test.com"}' \
  -H "Authorization: Bearer $TOKEN" http://localhost:8080/api/users

# Скачать config
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:8080/api/users/user123/config > wg.conf

# Лицензия
MIT