package allof

import (
	"crypto/sha256"
	"errors"
	"github.com/itsabgr/go-handy"
	"golang.org/x/crypto/nacl/secretbox"
)

func Decode(packet, topic []byte) (ttl byte, data []byte, err error) {
	if len(packet) <= 24 {
		return 0, nil, errors.New("allof: packet: invalid codec")
	}
	key := sha256.Sum256(topic)
	var nonce [24]byte
	copy(nonce[:], packet[:24])
	box := packet[24:]
	msg, ok := secretbox.Open(nil, box, &nonce, &key)
	if !ok {
		return 0, nil, errors.New("allof: packet: invalid topic")
	}

	return msg[0], msg[1:], nil
}

func Encode(ttl byte, msg, topic []byte) []byte {
	msg = append([]byte{ttl}, msg...)
	key := sha256.Sum256(topic)
	var nonce [24]byte
	copy(nonce[:], handy.Rand(len(nonce)))
	return secretbox.Seal(nonce[:], msg, &nonce, &key)
}
