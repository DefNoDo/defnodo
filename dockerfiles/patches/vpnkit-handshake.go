package libproxy

import (
	"fmt"
	"io"
)

type handshake struct {
	// In future we could add flags here for feature negotiation.
	// The length is dynamic but it must be < 64k.
	payload []byte // uninterpreted
}

const handshakeMagic = "https://github.com/moby/vpnkit multiplexer protocol\n"

func unmarshalHandshake(r io.Reader) (*handshake, error) {
	fmt.Printf("unmarshalHandshake\n")
	return &handshake{}, nil
}

func (h *handshake) Write(w io.Writer) error {
	fmt.Printf("handshake.Write\n")
	return nil
}
