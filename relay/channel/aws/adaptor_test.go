package aws

import (
	"net/http/httptest"
	"testing"

	relaycommon "github.com/QuantumNous/new-api/relay/common"
	"github.com/QuantumNous/new-api/relaykit/dto"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newSamplingTestRequest(model string) *dto.ClaudeRequest {
	temp := 0.7
	topP := 0.9
	topK := 40
	return &dto.ClaudeRequest{
		Model:       model,
		Temperature: &temp,
		TopP:        &topP,
		TopK:        &topK,
	}
}

func convertClaudeForTest(t *testing.T, model string) *dto.ClaudeRequest {
	t.Helper()
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	info := &relaycommon.RelayInfo{ChannelMeta: &relaycommon.ChannelMeta{UpstreamModelName: model}}
	out, err := (&Adaptor{}).ConvertClaudeRequest(c, info, newSamplingTestRequest(model))
	require.NoError(t, err)
	converted, ok := out.(*dto.ClaudeRequest)
	require.True(t, ok)
	return converted
}

// Opus 4.7 이후 세대(Opus 4.7/4.8, Claude 5 계열)는 샘플링 파라미터가 API에서
// 제거되어 temperature/top_p/top_k 지정 시 Bedrock이 400으로 거부한다.
func TestConvertClaudeRequestStripsSamplingForSamplingRemovedModels(t *testing.T) {
	for _, model := range []string{
		"claude-opus-5",
		"global.anthropic.claude-opus-5-20260115-v1:0",
		"claude-sonnet-5",
		"global.anthropic.claude-sonnet-5-20260203-v1:0",
		"anthropic.claude-opus-4-7",
		"anthropic.claude-opus-4-8",
	} {
		converted := convertClaudeForTest(t, model)
		assert.Nil(t, converted.Temperature, "temperature should be stripped for %s", model)
		assert.Nil(t, converted.TopP, "top_p should be stripped for %s", model)
		assert.Nil(t, converted.TopK, "top_k should be stripped for %s", model)
	}
}

// 구형 모델은 temperature를 유지하되, 동시 지정 시 top_p만 제거(ko.2 동작 유지)
func TestConvertClaudeRequestKeepsTemperatureForOlderModels(t *testing.T) {
	for _, model := range []string{
		"anthropic.claude-sonnet-4-5-20250929-v1:0",
		"global.anthropic.claude-haiku-4-5-20251001-v1:0",
		"anthropic.claude-sonnet-4-6",
		"anthropic.claude-opus-4-6-v1",
	} {
		converted := convertClaudeForTest(t, model)
		require.NotNil(t, converted.Temperature, "temperature should be kept for %s", model)
		assert.Equal(t, 0.7, *converted.Temperature)
		assert.Nil(t, converted.TopP, "top_p should be dropped when both are set for %s", model)
	}
}
