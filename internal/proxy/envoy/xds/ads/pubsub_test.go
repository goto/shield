package ads_test

import (
	"errors"
	"testing"

	"github.com/goto/shield/internal/proxy/envoy/xds/ads"
	"github.com/stretchr/testify/assert"
)

func TestPush(t *testing.T) {
	message := ads.Message{
		NodeID:      "node-1",
		VersionInfo: "v1",
		Nonce:       "test",
		TypeUrl:     ads.CLUSTER_TYPE_URL,
	}
	messageChan := make(ads.MessageChan, 1)

	err := messageChan.Push(message)
	recv := <-messageChan
	assert.NoError(t, err)
	assert.Equal(t, message, recv)

	close(messageChan)
	err = messageChan.Push(message)
	assert.True(t, errors.Is(err, ads.ErrChannelClosed))
}
