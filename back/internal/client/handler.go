package client

import (
	"encoding/json"
	"fmt"

	"github.com/reinerRubin/froppyshima/back/internal/protocol"
)

type (
	MethodHandler func(params json.RawMessage) (response json.RawMessage, err error)
	MethodRouter  map[protocol.Method]MethodHandler
)

func (m MethodRouter) Exec(method protocol.Method, params json.RawMessage) (json.RawMessage, error) {
	handler, found := m[method]
	if !found {
		return nil, fmt.Errorf("cant find handler")
	}

	return handler(params)
}
