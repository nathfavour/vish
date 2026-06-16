package spine

import (
	"encoding/binary"
	"encoding/json"
	"net"
	"time"

	"github.com/nathfavour/spine/pkg/types"
)

type Client struct {
	socketPath string
}

func NewClient(socketPath string) *Client {
	return &Client{socketPath: socketPath}
}

func (c *Client) Park(namespace [16]byte, duration time.Duration, state []byte) error {
	conn, err := net.Dial("unix", c.socketPath)
	if err != nil {
		return err
	}
	defer conn.Close()

	req := types.SleepRequest{
		Namespace: namespace,
		Duration:  duration,
		StateBlob: state,
	}
	data, _ := json.Marshal(req)

	var head [6]byte
	head[0] = types.MagicByte
	head[1] = byte(types.OpPark)
	binary.BigEndian.PutUint32(head[2:6], uint32(len(data)))

	conn.Write(head[:])
	conn.Write(data)

	// Block until OpWake is received
	for {
		if _, err := conn.Read(head[:]); err != nil {
			return err
		}
		if head[0] == types.MagicByte && types.OpCode(head[1]) == types.OpWake {
			return nil
		}
	}
}
