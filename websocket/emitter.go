package websocket

const (
	All       = ""
	Broadcast = ";to;all;except;me;"
)

type Emitter interface {
	EmitMessage([]byte) error
	Emit(string, interface{}) error
}

type emitter struct {
	conn *connection
	to   string
}

var _ Emitter = &emitter{}

func newEmitter(c *connection, to string) *emitter {
	return &emitter{conn: c, to: to}
}

func (e *emitter) EmitMessage(nativeMessage []byte) error {
	e.conn.server.emitMessage(e.conn.id, e.to, nativeMessage)
	return nil
}

func (e *emitter) Emit(event string, data interface{}) error {
	message, err := e.conn.server.messageSerializer.serialize(event, data)
	if err != nil {
		return err
	}
	e.EmitMessage(message)
	return nil
}
