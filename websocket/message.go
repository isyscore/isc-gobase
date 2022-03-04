package websocket

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/valyala/bytebufferpool"
)

type messageType uint8

func (m messageType) String() string {
	return strconv.Itoa(int(m))
}

func (m messageType) Name() string {
	switch m {
	case messageTypeString:
		return "string"
	case messageTypeInt:
		return "int"
	case messageTypeBool:
		return "bool"
	case messageTypeBytes:
		return "[]byte"
	case messageTypeJSON:
		return "json"
	default:
		return "Invalid(" + m.String() + ")"
	}
}

const (
	messageTypeString messageType = iota
	messageTypeInt
	messageTypeBool
	messageTypeBytes
	messageTypeJSON

	messageSeparator = ";"
)

var messageSeparatorByte = messageSeparator[0]

type messageSerializer struct {
	prefix []byte

	prefixLen       int
	separatorLen    int
	prefixAndSepIdx int
	prefixIdx       int
	separatorIdx    int

	buf *bytebufferpool.Pool
}

func newMessageSerializer(messagePrefix []byte) *messageSerializer {
	return &messageSerializer{
		prefix:          messagePrefix,
		prefixLen:       len(messagePrefix),
		separatorLen:    len(messageSeparator),
		prefixAndSepIdx: len(messagePrefix) + len(messageSeparator) - 1,
		prefixIdx:       len(messagePrefix) - 1,
		separatorIdx:    len(messageSeparator) - 1,

		buf: new(bytebufferpool.Pool),
	}
}

var (
	boolTrueB  = []byte("true")
	boolFalseB = []byte("false")
)

func (ms *messageSerializer) serialize(event string, data any) ([]byte, error) {
	b := ms.buf.Get()
	_, _ = b.Write(ms.prefix)
	_, _ = b.WriteString(event)
	_ = b.WriteByte(messageSeparatorByte)

	switch v := data.(type) {
	case string:
		_, _ = b.WriteString(messageTypeString.String())
		_ = b.WriteByte(messageSeparatorByte)
		_, _ = b.WriteString(v)
	case int:
		_, _ = b.WriteString(messageTypeInt.String())
		_ = b.WriteByte(messageSeparatorByte)
		_ = binary.Write(b, binary.LittleEndian, v)
	case bool:
		_, _ = b.WriteString(messageTypeBool.String())
		_ = b.WriteByte(messageSeparatorByte)
		if v {
			_, _ = b.Write(boolTrueB)
		} else {
			_, _ = b.Write(boolFalseB)
		}
	case []byte:
		_, _ = b.WriteString(messageTypeBytes.String())
		_ = b.WriteByte(messageSeparatorByte)
		_, _ = b.Write(v)
	default:
		//we suppose is json
		res, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}
		_, _ = b.WriteString(messageTypeJSON.String())
		_ = b.WriteByte(messageSeparatorByte)
		_, _ = b.Write(res)
	}

	message := b.Bytes()
	ms.buf.Put(b)

	return message, nil
}

var errInvalidTypeMessage = "Type %s is invalid for message: %s"

func (ms *messageSerializer) deserialize(event []byte, websocketMessage []byte) (any, error) {
	dataStartIdx := ms.prefixAndSepIdx + len(event) + 3
	if len(websocketMessage) <= dataStartIdx {
		return nil, errors.New("websocket invalid message: " + string(websocketMessage))
	}

	typ, err := strconv.Atoi(string(websocketMessage[ms.prefixAndSepIdx+len(event)+1 : ms.prefixAndSepIdx+len(event)+2])) // in order to iris-websocket-message;user;-> 4
	if err != nil {
		return nil, err
	}

	data := websocketMessage[dataStartIdx:]

	switch messageType(typ) {
	case messageTypeString:
		return string(data), nil
	case messageTypeInt:
		msg, err := strconv.Atoi(string(data))
		if err != nil {
			return nil, err
		}
		return msg, nil
	case messageTypeBool:
		if bytes.Equal(data, boolTrueB) {
			return true, nil
		}
		return false, nil
	case messageTypeBytes:
		return data, nil
	case messageTypeJSON:
		var msg any
		err := json.Unmarshal(data, &msg)
		return msg, err
	default:
		return nil, errors.New(fmt.Sprintf(errInvalidTypeMessage, messageType(typ).Name(), websocketMessage))
	}
}

func (ms *messageSerializer) getWebsocketCustomEvent(websocketMessage []byte) []byte {
	if len(websocketMessage) < ms.prefixAndSepIdx {
		return nil
	}
	s := websocketMessage[ms.prefixAndSepIdx:]
	evt := s[:bytes.IndexByte(s, messageSeparatorByte)]
	return evt
}
