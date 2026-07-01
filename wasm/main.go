//go:build js && wasm

package main

import (
	"encoding/json"
	"syscall/js"

	"github.com/yashdalfthegray/gallery-wall/layout"
)

func main() {
	js.Global().Set("goLayout", js.FuncOf(goLayout))
	<-make(chan struct{})
}

func goLayout(_ js.Value, args []js.Value) any {
	if len(args) != 1 {
		return mustJSON(layoutResponse{OK: false, Message: "expected one JSON argument"})
	}

	var params layout.Params
	if err := json.Unmarshal([]byte(args[0].String()), &params); err != nil {
		return mustJSON(layoutResponse{OK: false, Message: err.Error()})
	}

	result, err := layout.Layout(params)
	if err != nil {
		return mustJSON(errorResponse(err))
	}

	return mustJSON(layoutResponse{OK: true, Result: &result})
}

type layoutResponse struct {
	OK      bool           `json:"ok"`
	Result  *layout.Result `json:"result,omitempty"`
	Error   *errorPayload  `json:"error,omitempty"`
	Message string         `json:"message,omitempty"`
}

type errorPayload struct {
	Code    string   `json:"code"`
	Message string   `json:"message"`
	ItemID  string   `json:"itemId,omitempty"`
	ItemIDs []string `json:"itemIds,omitempty"`
}

func errorResponse(err error) layoutResponse {
	switch e := err.(type) {
	case *layout.ValidationError:
		return layoutResponse{
			OK: false,
			Error: &errorPayload{
				Code:    string(e.Code),
				Message: e.Message,
				ItemID:  e.ItemID,
				ItemIDs: e.ItemIDs,
			},
		}
	case *layout.PlacementError:
		return layoutResponse{
			OK: false,
			Error: &errorPayload{
				Code:    string(e.Code),
				Message: e.Message,
				ItemID:  e.ItemID,
			},
		}
	default:
		return layoutResponse{OK: false, Message: err.Error()}
	}
}

func mustJSON(v layoutResponse) string {
	data, err := json.Marshal(v)
	if err != nil {
		return `{"ok":false,"message":"marshal failed"}`
	}
	return string(data)
}
