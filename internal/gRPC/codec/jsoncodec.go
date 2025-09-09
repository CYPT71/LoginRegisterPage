package codec

import (
	"encoding/json"

	"google.golang.org/grpc/encoding"
)

// JSONCodec implements gRPC encoding.Codec using JSON for marshaling.
type JSONCodec struct{}

// Marshal encodes v into JSON.
func (JSONCodec) Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

// Unmarshal decodes JSON data into v.
func (JSONCodec) Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

// Name returns the codec name.
func (JSONCodec) Name() string {
	return "json"
}

func init() {
	encoding.RegisterCodec(JSONCodec{})
}
