package logging

import (
	"net"
	"time"
)

const (
	DefaultReconnectDuration = 10
)

type NetHandler struct {
	*BaseHandler
	Network string
	Address string
	Timeout time.Duration
}

func NewNetHandler(network, address string, timeout time.Duration) (*NetHandler, error) {
	h := &NetHandler{
		Network: network,
		Address: address,
		Timeout: timeout,
	}
	conn, err := h.DialTimeout()
	if err != nil {
		return nil, err
	}
	bh, err := NewBaseHandler(conn, INFO, DefaultTimeLayout, DefaultFormat)
	if err != nil {
		return nil, err
	}
	h.BaseHandler = bh
	h.GotError = h.GotNetError
	return h, nil
}

func (h *NetHandler) DialTimeout() (net.Conn, error) {
	return net.DialTimeout(h.Network, h.Address, h.Timeout)
}

func (h *NetHandler) GotNetError(err error) {
	if _, ok := err.(net.Error); !ok {
		h.BaseHandler.GotError(err)
		return
	}
	for {
		conn, err := h.DialTimeout()
		if err == nil {
			h.Writer = conn
			break
		}
		time.Sleep(DefaultReconnectDuration * time.Second)
	}
}
