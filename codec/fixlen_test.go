package codec_test

import (
	"encoding/binary"
	"github.com/icodeface/link/codec"
	"testing"
)

func Test_FixLen(t *testing.T) {
	base := JsonTestProtocol()
	protocol := codec.FixLen(base, 2, binary.LittleEndian, 1024, 1024)
	JsonTest(t, protocol)
}
