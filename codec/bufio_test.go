package codec_test

import (
	"encoding/binary"
	"github.com/icodeface/link/codec"
	"testing"
)

func Test_Bufio(t *testing.T) {
	JsonTest(t, codec.Bufio(JsonTestProtocol(), 1024, 1024))
	JsonTest(t, codec.Bufio(codec.FixLen(JsonTestProtocol(), 2, binary.LittleEndian, 64*1024, 64*1024), 1024, 1024))
	ProtoTest(t, codec.Bufio(ProtoTestProtocol(), 1024, 1024))
	ProtoTest(t, codec.Bufio(codec.FixLen(ProtoTestProtocol(), 2, binary.LittleEndian, 64*1024, 64*1024), 1024, 1024))
}
