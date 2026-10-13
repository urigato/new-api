package model

import (
	"strings"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/constant"
)

// chatCapableEndpointTypes are the endpoint types a chat playground can drive.
var chatCapableEndpointTypes = map[string]struct{}{
	string(constant.EndpointTypeOpenAI):                {},
	string(constant.EndpointTypeOpenAIResponse):        {},
	string(constant.EndpointTypeOpenAIResponseCompact): {},
	string(constant.EndpointTypeAnthropic):             {},
	string(constant.EndpointTypeGemini):                {},
}

// EndpointsDeclareChatCapable reports whether a model metadata endpoints value
// still allows chat. An empty or unparsable value keeps the channel-inferred
// behaviour (chat allowed). An explicit map that declares only non-chat
// endpoints (e.g. {"embeddings": "/v1/embeddings"}) marks the model as
// non-chat so UIs such as the playground can hide it.
func EndpointsDeclareChatCapable(endpoints string) bool {
	if strings.TrimSpace(endpoints) == "" {
		return true
	}
	var raw map[string]any
	if err := common.Unmarshal([]byte(endpoints), &raw); err != nil || len(raw) == 0 {
		return true
	}
	for key := range raw {
		if _, ok := chatCapableEndpointTypes[strings.ToLower(strings.TrimSpace(key))]; ok {
			return true
		}
	}
	return false
}

// FilterChatCapableModels drops models whose metadata declares only non-chat
// endpoints. Models without metadata are kept.
func FilterChatCapableModels(names []string) []string {
	if len(names) == 0 {
		return names
	}
	var records []Model
	if err := DB.Where("model_name IN ?", names).Find(&records).Error; err != nil {
		common.SysError("FilterChatCapableModels query error: " + err.Error())
		return names
	}
	resolved := resolveModelMetadata(records, names)
	kept := make([]string, 0, len(names))
	for _, name := range names {
		if meta := resolved[name]; meta != nil && !EndpointsDeclareChatCapable(meta.Endpoints) {
			continue
		}
		kept = append(kept, name)
	}
	return kept
}
