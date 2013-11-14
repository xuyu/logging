package logging

import (
	"net"
	"time"
)

const (
	DefaultReconnectDuration = 3
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
	bh, err := NewBaseHandler(conn, DEBUG, DefaultTimeLayout, DefaultFormat)
	if err != nil {
		return nil, err
	}
	h.BaseHandler = bh
	go h.check_state()
	return h, nil
}

func (h *NetHandler) DialTimeout() (net.Conn, error) {
	return net.DialTimeout(h.Network, h.Address, h.Timeout)
}

func (h *NetHandler) check_state() {
	for {
		time.Sleep(DefaultReconnectDuration * time.Second)
		state, err := h.get_state()
		if state {
			continue
		}
		_, ok := err.(net.Error)
		if !ok {
			continue
		}
		conn, er := h.DialTimeout()
		if er == nil {
			h.Writer = conn
			h.set_state(true, nil)
		}
	}
}
