module vpn-service

go 1.22

require (
    github.com/gorilla/mux v1.8.1
    github.com/gorilla/securecookie v1.1.2
    github.com/gorilla/sessions v1.3.0
    github.com/mattn/go-sqlite3 v1.14.22
    github.com/pion/transport v0.15.0
    golang.zx2c4.com/wireguard/device v0.0.0-20240520191018-3d931cd65a96
    golang.zx2c4.com/wireguard/tun v0.0.0-20240115185543-a0f7b77cf22c
    golang.zx2c4.com/wireguard/wgctrl v0.0.0-20230429144221-2e7eff7f53fd
    golang.zx2c4.com/wireguard/wgtypes v0.0.0-20230328165227-77c3af67bb11
)
