package codec

import (
	"github.com/gogo/protobuf/proto"
	"github.com/icodeface/link"
	"github.com/lunixbochs/struc"
	"io"
	"reflect"
)

type ProtoProtocol struct {
	types map[uint16]reflect.Type // msgType => Type
	names map[reflect.Type]uint16 // Type => msgType
}

func Proto() *ProtoProtocol {
	return &ProtoProtocol{
		types: make(map[uint16]reflect.Type),
		names: make(map[reflect.Type]uint16),
	}
}

func (p *ProtoProtocol) Register(msgType uint16, t interface{}) {
	rt := reflect.TypeOf(t)
	if rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
	}
	p.types[msgType] = rt
	p.names[rt] = msgType
}

func (p *ProtoProtocol) NewCodec(rw io.ReadWriter) (link.Codec, error) {
	codec := &protoCodec{
		p:      p,
		reader: rw,
		writer: rw,
	}
	codec.closer, _ = rw.(io.Closer)
	return codec, nil
}

type protoCodec struct {
	p      *ProtoProtocol
	reader io.Reader
	writer io.Writer
	closer io.Closer
}

type head struct {
	Length  int    `struc:"uint32,little"`
	MsgType uint16 `struc:"uint16,little"`
}

func (c *protoCodec) Receive() (interface{}, error) {
	h := &head{}
	if err := struc.Unpack(c.reader, h); err != nil {
		return nil, err
	}

	var msg interface{}
	if h.MsgType != 0 {
		if t, exists := c.p.types[h.MsgType]; exists {
			msg = reflect.New(t).Interface()
		} else {
			return nil, link.ErrUnregisteredMsg
		}
	}

	msgBytes := make([]byte, h.Length)
	if _, err := io.ReadFull(c.reader, msgBytes); err != nil {
		return nil, err
	}

	if err := proto.Unmarshal(msgBytes, msg.(proto.Message)); err != nil {
		return nil, err
	}
	return msg, nil
}

func (c *protoCodec) Send(msg interface{}) error {
	t := reflect.TypeOf(msg)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	msgType, exists := c.p.names[t]
	if !exists {
		return link.ErrUnregisteredMsg
	}

	msgBytes, err := proto.Marshal(msg.(proto.Message))
	if err != nil {
		return err
	}

	h := &head{
		MsgType: msgType,
		Length:  len(msgBytes),
	}
	if err = struc.Pack(c.writer, h); err != nil {
		return err
	}
	_, err = c.writer.Write(msgBytes)
	return err
}

func (c *protoCodec) Close() error {
	if c.closer != nil {
		return c.closer.Close()
	}
	return nil
}
