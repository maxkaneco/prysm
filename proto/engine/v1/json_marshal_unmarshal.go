package enginev1

import (
	"google.golang.org/protobuf/encoding/protojson"
)

// MarshalJSON defines a custom json.Marshaler interface implementation
// that uses protojson underneath the hood, as protojson will respect
// proper struct tag naming conventions required for the JSON-RPC engine API to work.
func (e *ExecutionPayload) MarshalJSON() ([]byte, error) {
	return protojson.Marshal(e)
}

// UnmarshalJSON defines a custom json.Unmarshaler interface implementation
// that uses protojson underneath the hood, as protojson will respect
// proper struct tag naming conventions required for the JSON-RPC engine API to work.
func (e *ExecutionPayload) UnmarshalJSON(enc []byte) error {
	return protojson.Unmarshal(enc, e)
}

// MarshalJSON --
func (p *PayloadAttributes) MarshalJSON() ([]byte, error) {
	return protojson.Marshal(p)
}

// UnmarshalJSON --
func (p *PayloadAttributes) UnmarshalJSON(enc []byte) error {
	return protojson.Unmarshal(enc, p)
}

// MarshalJSON --
func (p *PayloadStatus) MarshalJSON() ([]byte, error) {
	//type Alias PayloadStatus
	//return json.Marshal(&struct {
	//Status string `json:"status"`
	//*Alias
	//}{
	//Status: p.Status.String(),
	//Alias:  (*Alias)(p),
	//})
	return protojson.Marshal(p)
}

// UnmarshalJSON --
func (p *PayloadStatus) UnmarshalJSON(enc []byte) error {
	//type Alias PayloadStatus
	//aux := &struct {
	//Status string `json:"status"`
	//*Alias
	//}{
	//Alias: (*Alias)(p),
	//}
	//if err := json.Unmarshal(enc, &aux); err != nil {
	//return err
	//}
	//p.Status = PayloadStatus_Status(PayloadStatus_Status_value[aux.Status])
	//return nil
	return protojson.Unmarshal(enc, p)
}

// MarshalJSON --
func (f *ForkchoiceState) MarshalJSON() ([]byte, error) {
	return protojson.Marshal(f)
}

// UnmarshalJSON --
func (f *ForkchoiceState) UnmarshalJSON(enc []byte) error {
	return protojson.Unmarshal(enc, f)
}