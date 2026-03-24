package protocol

import (
	"encoding/json"
	"fmt"
)

// MarshalRequest serializes a Request with the given params to JSON.
func MarshalRequest(id, method string, params any) ([]byte, error) {
	p, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("marshal params: %w", err)
	}
	return json.Marshal(Request{
		Type:   FrameTypeRequest,
		ID:     id,
		Method: method,
		Params: p,
	})
}

// MarshalResponse serializes a success Response.
func MarshalResponse(id string, payload any) ([]byte, error) {
	p, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal payload: %w", err)
	}
	return json.Marshal(Response{
		Type:    FrameTypeResponse,
		ID:      id,
		OK:      true,
		Payload: p,
	})
}

// MarshalErrorResponse serializes an error Response.
func MarshalErrorResponse(id string, errPayload ErrorPayload) ([]byte, error) {
	return json.Marshal(Response{
		Type:  FrameTypeResponse,
		ID:    id,
		OK:    false,
		Error: &errPayload,
	})
}

// MarshalEvent serializes an Event.
func MarshalEvent(eventName EventName, payload any) ([]byte, error) {
	p, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal payload: %w", err)
	}
	return json.Marshal(Event{
		Type:      FrameTypeEvent,
		EventName: eventName,
		Payload:   p,
	})
}

// ParseFrame performs initial deserialization of a raw JSON frame to determine
// its type. Callers should then unmarshal into the appropriate concrete type.
func ParseFrame(data []byte) (*RawFrame, error) {
	var f RawFrame
	if err := json.Unmarshal(data, &f); err != nil {
		return nil, fmt.Errorf("parse frame: %w", err)
	}
	return &f, nil
}

// UnmarshalRequest parses a full Request from raw JSON.
func UnmarshalRequest(data []byte) (*Request, error) {
	var r Request
	if err := json.Unmarshal(data, &r); err != nil {
		return nil, fmt.Errorf("unmarshal request: %w", err)
	}
	return &r, nil
}

// UnmarshalResponse parses a full Response from raw JSON.
func UnmarshalResponse(data []byte) (*Response, error) {
	var r Response
	if err := json.Unmarshal(data, &r); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}
	return &r, nil
}

// UnmarshalEvent parses a full Event from raw JSON.
func UnmarshalEvent(data []byte) (*Event, error) {
	var e Event
	if err := json.Unmarshal(data, &e); err != nil {
		return nil, fmt.Errorf("unmarshal event: %w", err)
	}
	return &e, nil
}
