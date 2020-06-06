package codec_test

import (
	"bytes"
	"github.com/icodeface/link/codec"
	"github.com/icodeface/link/codec/proto_test_message"
	"testing"

	"github.com/icodeface/link"
)

const (
	MsgTypeMyMessage1 uint16 = 0x0001
	MsgTypeMyMessage2 uint16 = 0x0002
)

func ProtoTestProtocol() *codec.ProtoProtocol {
	protocol := codec.Proto()
	protocol.Register(MsgTypeMyMessage1, proto_test_message.MyMessage1{})
	protocol.Register(MsgTypeMyMessage2, &proto_test_message.MyMessage2{})
	return protocol
}

func equlal1(a *proto_test_message.MyMessage1, b *proto_test_message.MyMessage1) bool {
	return a.Field1 == b.Field1 && a.Field2 == b.Field2
}

func equlal2(a *proto_test_message.MyMessage2, b *proto_test_message.MyMessage2) bool {
	return a.Field1 == b.Field1 && a.Field2 == b.Field2
}

func ProtoTest(t *testing.T, protocol link.Protocol) {
	var stream bytes.Buffer

	codec, _ := protocol.NewCodec(&stream)

	sendMsg1 := proto_test_message.MyMessage1{
		Field1: "abc",
		Field2: 123,
	}

	err := codec.Send(&sendMsg1)
	if err != nil {
		t.Fatal(err)
	}

	recvMsg1, err := codec.Receive()
	if err != nil {
		t.Fatal(err)
	}

	if _, ok := recvMsg1.(*proto_test_message.MyMessage1); !ok {
		t.Fatalf("message type not match: %#v", recvMsg1)
	}

	if !equlal1(&sendMsg1, recvMsg1.(*proto_test_message.MyMessage1)) {
		t.Fatalf("message not match: %v, %v", sendMsg1, recvMsg1)
	}

	sendMsg2 := proto_test_message.MyMessage2{
		Field1: 123,
		Field2: "abc",
	}

	err = codec.Send(&sendMsg2)
	if err != nil {
		t.Fatal(err)
	}

	recvMsg2, err := codec.Receive()
	if err != nil {
		t.Fatal(err)
	}

	if _, ok := recvMsg2.(*proto_test_message.MyMessage2); !ok {
		t.Fatalf("message type not match: %#v", recvMsg2)
	}

	if !equlal2(&sendMsg2, recvMsg2.(*proto_test_message.MyMessage2)) {
		t.Fatalf("message not match: %v, %v", sendMsg2, recvMsg2)
	}

	// multi messages
	for i := 0; i < 3; i++ {
		err = codec.Send(&sendMsg1)
		if err != nil {
			t.Fatal(err)
		}
		err = codec.Send(&sendMsg2)
		if err != nil {
			t.Fatal(err)
		}
	}

	for i := 0; i < 3; i++ {
		recvMsg1, err := codec.Receive()
		if err != nil {
			t.Fatal(err)
		}
		if _, ok := recvMsg1.(*proto_test_message.MyMessage1); !ok {
			t.Fatalf("message type not match: %#v", recvMsg1)
		}
		if !equlal1(&sendMsg1, recvMsg1.(*proto_test_message.MyMessage1)) {
			t.Fatalf("message not match: %v, %v", sendMsg1, recvMsg1)
		}

		recvMsg2, err := codec.Receive()
		if err != nil {
			t.Fatal(err)
		}
		if _, ok := recvMsg2.(*proto_test_message.MyMessage2); !ok {
			t.Fatalf("message type not match: %#v", recvMsg2)
		}
		if !equlal2(&sendMsg2, recvMsg2.(*proto_test_message.MyMessage2)) {
			t.Fatalf("message not match: %v, %v", sendMsg2, recvMsg2)
		}
	}
}

func Test_Proto(t *testing.T) {
	protocol := ProtoTestProtocol()
	ProtoTest(t, protocol)
}
