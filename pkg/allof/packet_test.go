package allof

import (
	"bytes"
	"github.com/itsabgr/go-handy"
	"testing"
)

func TestPacket(t *testing.T) {
	for range handy.N(10) {
		topic := handy.Rand(10)
		msg := handy.Rand(300)
		msg2, err := Decode(Encode(msg, topic), topic)
		if err != nil {
			t.Fatal(err)
		}

		if !bytes.Equal(msg, msg2) {
			t.Fatal()
		}
		_, err = Decode(Encode(msg, handy.Rand(10)), topic)
		if err == nil {
			t.Fatal()
		}
		msg = Encode(msg, topic)
		msg2 = Encode(msg, topic)
		if bytes.Equal(msg, msg2) {
			t.Fatal()
		}
		if len(msg) > 400 {
			t.Fatal()
		}
		if len(msg2) > 400 {
			t.Fatal()
		}
	}
}
