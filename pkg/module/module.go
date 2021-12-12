package module

import (
	msgpack "github.com/wapc/tinygo-msgpack"
	wapc "github.com/wapc/wapc-guest-tinygo"
)

type Host struct {
	binding string
}

func NewHost(binding string) *Host {
	return &Host{
		binding: binding,
	}
}

func (h *Host) SayHello(name string) (string, error) {
	inputArgs := SayHelloArgs{
		Name: name,
	}
	inputBytes, err := msgpack.ToBytes(&inputArgs)
	if err != nil {
		return "", err
	}
	payload, err := wapc.HostCall(
		h.binding,
		"tiero:test",
		"sayHello",
		inputBytes,
	)
	if err != nil {
		return "", err
	}
	decoder := msgpack.NewDecoder(payload)
	return decoder.ReadString()
}

type Handlers struct {
	SayHello func(name string) (string, error)
}

func (h Handlers) Register() {
	if h.SayHello != nil {
		sayHelloHandler = h.SayHello
		wapc.RegisterFunction("sayHello", sayHelloWrapper)
	}
}

var (
	sayHelloHandler func(name string) (string, error)
)

func sayHelloWrapper(payload []byte) ([]byte, error) {
	decoder := msgpack.NewDecoder(payload)
	var inputArgs SayHelloArgs
	inputArgs.Decode(&decoder)
	response, err := sayHelloHandler(inputArgs.Name)
	if err != nil {
		return nil, err
	}
	var sizer msgpack.Sizer
	sizer.WriteString(response)

	ua := make([]byte, sizer.Len())
	encoder := msgpack.NewEncoder(ua)
	encoder.WriteString(response)

	return ua, nil
}

type SayHelloArgs struct {
	Name string
}

func DecodeSayHelloArgsNullable(decoder *msgpack.Decoder) (*SayHelloArgs, error) {
	if isNil, err := decoder.IsNextNil(); isNil || err != nil {
		return nil, err
	}
	decoded, err := DecodeSayHelloArgs(decoder)
	return &decoded, err
}

func DecodeSayHelloArgs(decoder *msgpack.Decoder) (SayHelloArgs, error) {
	var o SayHelloArgs
	err := o.Decode(decoder)
	return o, err
}

func (o *SayHelloArgs) Decode(decoder *msgpack.Decoder) error {
	numFields, err := decoder.ReadMapSize()
	if err != nil {
		return err
	}

	for numFields > 0 {
		numFields--
		field, err := decoder.ReadString()
		if err != nil {
			return err
		}
		switch field {
		case "name":
			o.Name, err = decoder.ReadString()
		default:
			err = decoder.Skip()
		}
		if err != nil {
			return err
		}
	}

	return nil
}

func (o *SayHelloArgs) Encode(encoder msgpack.Writer) error {
	if o == nil {
		encoder.WriteNil()
		return nil
	}
	encoder.WriteMapSize(1)
	encoder.WriteString("name")
	encoder.WriteString(o.Name)

	return nil
}
