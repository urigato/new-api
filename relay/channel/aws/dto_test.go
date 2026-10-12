package aws

import (
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Claude Code(2.1.x)가 보내는 anthropic-beta 중 Bedrock이 400으로 거부하는
// first-party 전용 beta는 제거하고 나머지는 순서대로 보존한다.
func TestFilterBedrockBetasDropsFirstPartyOnlyBetas(t *testing.T) {
	header := "claude-code-20250219, interleaved-thinking-2025-05-14,prompt-caching-scope-2026-01-05,context-1m-2025-08-07,oauth-2025-04-20, ,effort-2025-11-24"
	got := filterBedrockBetas(header)
	assert.Equal(t, []string{"claude-code-20250219", "interleaved-thinking-2025-05-14", "context-1m-2025-08-07", "effort-2025-11-24"}, got)
	assert.Nil(t, filterBedrockBetas(""))
	assert.Nil(t, filterBedrockBetas("prompt-caching-scope-2026-01-05"))
}

func TestFormatRequestFiltersBetaHeader(t *testing.T) {
	h := http.Header{}
	h.Set("anthropic-beta", "prompt-caching-scope-2026-01-05,interleaved-thinking-2025-05-14")
	req, err := formatRequest(strings.NewReader(`{"max_tokens":8,"messages":[{"role":"user","content":"hi"}]}`), h)
	require.NoError(t, err)
	assert.JSONEq(t, `["interleaved-thinking-2025-05-14"]`, string(req.AnthropicBeta))
	assert.Equal(t, "bedrock-2023-05-31", req.AnthropicVersion)

	h.Set("anthropic-beta", "prompt-caching-scope-2026-01-05")
	req, err = formatRequest(strings.NewReader(`{"max_tokens":8,"messages":[]}`), h)
	require.NoError(t, err)
	assert.Empty(t, req.AnthropicBeta)
}
