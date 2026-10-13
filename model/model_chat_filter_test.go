package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEndpointsDeclareChatCapable(t *testing.T) {
	assert.True(t, EndpointsDeclareChatCapable(""))
	assert.True(t, EndpointsDeclareChatCapable("   "))
	assert.True(t, EndpointsDeclareChatCapable(`["embeddings"]`), "array form is not a declaration map")
	assert.True(t, EndpointsDeclareChatCapable(`{"openai": "/v1/chat/completions"}`))
	assert.True(t, EndpointsDeclareChatCapable(`{"embeddings": "/v1/embeddings", "anthropic": {"path": "/v1/messages"}}`))
	assert.False(t, EndpointsDeclareChatCapable(`{"embeddings": "/v1/embeddings"}`))
	assert.False(t, EndpointsDeclareChatCapable(`{"jina-rerank": "/v1/rerank", "embeddings": "/v1/embeddings"}`))
}
