package ads

import "errors"

type Message struct {
	NodeID      string
	VersionInfo string
	Nonce       string
	TypeUrl     string
}

type MessageChan chan Message

var (
	ErrChannelClosed = errors.New("can't send message on closed channel")
)

func (m MessageChan) Push(message Message) (err error) {
	defer func() {
		if recover() != nil {
			err = ErrChannelClosed
		}
	}()

	m <- message
	return nil
}
