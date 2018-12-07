package protocol

import (
	"encoding/json"

	// TODO remove engine dependencies
	"github.com/reinerRubin/froppyshima/back/internal/engine"
)

type Method string

const (
	MethodNew  Method = "new"
	MethodLoad Method = "load"
	MethodHit  Method = "hit"
)

type Call struct {
	ID     string
	Method Method
	Body   json.RawMessage
}

type ServerResponse struct {
	ID    *string
	Error string
	Body  json.RawMessage
}

type (
	HitArgs struct {
		X, Y int
	}

	HitResponse struct {
		Result engine.HitResult
	}
)

type (
	NewResponse struct {
		ID engine.GameID
	}
)

type (
	LoadArgs struct {
		ID engine.GameID
	}

	LoadResponse struct {
		Maxx         int
		Maxy         int
		FailHits     []*engine.Coord
		SuccessHits  []*engine.Coord
		FatalHits    []*engine.Coord
		AnyMoreShips bool
	}
)

func NewCallResponse(c *Call) *ServerResponse {
	return &ServerResponse{
		ID: &c.ID,
	}
}

func NewEvent() *ServerResponse {
	return &ServerResponse{}
}

func ParseCall(rawCall []byte) (*Call, error) {
	call := &Call{}
	if err := json.Unmarshal(rawCall, call); err != nil {
		return nil, err
	}

	return call, nil
}

func MarshallServerResponse(r *ServerResponse) ([]byte, error) {
	return json.Marshal(r)
}
