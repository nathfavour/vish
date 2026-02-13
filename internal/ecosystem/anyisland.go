package ecosystem

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"time"
)

type HandshakeResponse struct {
	Status           string `json:"status"`
	ToolID           string `json:"tool_id"`
	Version          string `json:"version"`
	AnyislandVersion string `json:"anyisland_version"`
}

var (
	isManagedCached bool
	handshakeCached *HandshakeResponse
	lastCheck       time.Time
)

// IsManaged checks if vish is being managed by Anyisland via the Pulse handshake
func IsManaged() (bool, *HandshakeResponse) {
	if time.Since(lastCheck) < 30*time.Second && lastCheck.IsZero() == false {
		return isManagedCached, handshakeCached
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return false, nil
	}

	socketPath := filepath.Join(home, ".anyisland", "anyisland.sock")
	conn, err := net.DialTimeout("unix", socketPath, 100*time.Millisecond) // Lower timeout
	if err != nil {
		lastCheck = time.Now()
		isManagedCached = false
		return false, nil
	}
	defer conn.Close()

	handshake := map[string]string{"op": "HANDSHAKE"}
	data, _ := json.Marshal(handshake)
	_, err = conn.Write(data)
	if err != nil {
		return false, nil
	}

	decoder := json.NewDecoder(conn)
	var resp HandshakeResponse
	err = decoder.Decode(&resp)
	if err != nil {
		return false, nil
	}

	isManagedCached = resp.Status == "MANAGED"
	handshakeCached = &resp
	lastCheck = time.Now()

	return isManagedCached, handshakeCached
}

// RegisterWithDaemon registers vish with the local Anyisland daemon via UDP
func RegisterWithDaemon(version string) error {
	packet := map[string]string{
		"op":      "REGISTER",
		"name":    "vish",
		"source":  "github.com/nathfavour/vish",
		"version": version,
		"type":    "binary",
	}
	data, _ := json.Marshal(packet)

	conn, err := net.Dial("udp", "localhost:1995")
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.Write(data)
	return err
}

func CheckPulse() string {
	managed, resp := IsManaged()
	if managed {
		return fmt.Sprintf("managed by anyisland %s", resp.AnyislandVersion)
	}
	return "unmanaged"
}
