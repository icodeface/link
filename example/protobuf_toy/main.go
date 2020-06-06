package main

import (
	"github.com/icodeface/link"
	"github.com/icodeface/link/codec"
	"github.com/icodeface/link/codec/proto_test_message"
	"log"
)

const (
	MsgTypeAddReq uint16 = 0x0001
	MsgTypeAddRsq uint16 = 0x0002
)

func main() {
	proto := codec.Proto()
	proto.Register(MsgTypeAddReq, proto_test_message.AddReq{})
	proto.Register(MsgTypeAddRsq, proto_test_message.AddRsp{})

	server, err := link.Listen("tcp", "0.0.0.0:0", codec.Bufio(proto, 1024, 1024), 0 /* sync send */, link.HandlerFunc(serverSessionLoop))
	checkErr(err)
	addr := server.Listener().Addr().String()
	go server.Serve()

	client, err := link.Dial("tcp", addr, proto, 0)
	checkErr(err)
	clientSessionLoop(client)
}

func serverSessionLoop(session *link.Session) {
	log.Printf("recv req from %s", session.RemoteAddr())
	for {
		req, err := session.Receive()
		checkErr(err)

		err = session.Send(&proto_test_message.AddRsp{
			C: req.(*proto_test_message.AddReq).A + req.(*proto_test_message.AddReq).B,
		})
		checkErr(err)
	}
}

func clientSessionLoop(session *link.Session) {
	for i := 0; i < 10; i++ {
		err := session.Send(&proto_test_message.AddReq{
			A: int32(i),
			B: int32(i),
		})
		checkErr(err)
		log.Printf("Send: %d + %d", i, i)

		rsp, err := session.Receive()
		checkErr(err)
		log.Printf("Receive: %d", rsp.(*proto_test_message.AddRsp).C)
	}
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
