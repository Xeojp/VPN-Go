package tun

import (
    "context"
    "fmt"
    "log"
    "net"
    "os/exec"
    "strings"
    "sync"
    
    "golang.zx2c4.com/wireguard/device"
    "golang.zx2c4.com/wireguard/tun"
    "golang.zx2c4.com/wireguard/wgctrl"
    "golang.zx2c4.com/wireguard/wgctrl/wgtypes"
    "vpn-service/internal/models"
)

type WireGuard struct {
    dev     *device.Device
    clients map[string]*wgtypes.PeerConfig
    mu      sync.RWMutex
    ctx     context.Context
    cancel  context.CancelFunc
}

func NewWireGuard() (*WireGuard, error) {
    tunIface, err := tun.CreateTUN("wg0", 1420)
    if err != nil {
        return nil, fmt.Errorf("create TUN: %w", err)
    }
    
    cfg := device.Config{
        PrivateKey: wgtypes.GeneratePrivateKey(),
        ListenPort: 51820,
    }
    
    dev, err := device.NewDevice(tunIface, rand.Reader, logger{}, cfg)
    if err != nil {
        return nil, fmt.Errorf("new device: %w", err)
    }
    
    wg := &WireGuard{
        dev:     dev,
        clients: make(map[string]*wgtypes.PeerConfig),
    }
    
    wg.dev.Up()
    go wg.watchPeers()

    exec.Command("sysctl", "-w", "net.ipv4.ip_forward=1").Run()
    exec.Command("iptables", "-t", "nat", "-A", "POSTROUTING", "-s", "10.0.0.0/24", "-o", "eth0", "-j", "MASQUERADE").Run()
    
    return wg, nil
}

func (wg *WireGuard) AddClient(user *models.User) error {
    wg.mu.Lock()
    defer wg.mu.Unlock()
    
    ip := net.ParseIP(user.IPAddress)
    peer := &wgtypes.PeerConfig{
        PublicKey:           wgtypes.Key(base64.StdEncoding.DecodeString(user.PubKey)),
        AllowedIPs:          []net.IPNet{*net.IPNet(ip.String() + "/32")},
        PersistentKeepalive: 25,
    }
    
    wg.clients[user.ID] = peer

    client := wgctrl.NewClient(wgctrl.Config{})
    defer client.Close()
    
    config := wgtypes.Config{
        Peers: []wgtypes.PeerConfig{*peer},
    }
    
    return client.ConfigureDevice("wg0", config)
}

func (wg *WireGuard) RemoveClient(userID string) error {
    wg.mu.Lock()
    delete(wg.clients, userID)
    wg.mu.Unlock()

    client := wgctrl.NewClient(wgctrl.Config{})
    defer client.Close()
    
    peers := make([]wgtypes.PeerConfig, 0, len(wg.clients))
    wg.mu.RLock()
    for _, peer := range wg.clients {
        peers = append(peers, *peer)
    }
    wg.mu.RUnlock()
    
    return client.ConfigureDevice("wg0", wgtypes.Config{Peers: peers})
}

func (wg *WireGuard) Shutdown() {
    wg.cancel()
    wg.dev.Close()
}
