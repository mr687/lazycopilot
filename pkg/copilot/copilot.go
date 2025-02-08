package copilot

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/mr687/lazycopilot/pkg/utils"
)

type GithubEndpoint struct {
	API           string `json:"api"`
	OriginTracker string `json:"origin-tracker"`
	Proxy         string `json:"proxy"`
	Telemetry     string `json:"telemetry"`
}

type GithubToken struct {
	AnnotationsEnabled           bool           `json:"annotations_enabled"`
	ChatEnabled                  bool           `json:"chat_enabled"`
	ChatJetbrainsEnabled         bool           `json:"chat_jetbrains_enabled"`
	CodeQuoteEnabled             bool           `json:"code_quote_enabled"`
	CodeReviewEnabled            bool           `json:"code_review_enabled"`
	Codesearch                   bool           `json:"codesearch"`
	CopilotignoreEnabled         bool           `json:"copilotignore_enabled"`
	Endpoints                    GithubEndpoint `json:"endpoints"`
	ExpiresAt                    int            `json:"expires_at"`
	Individual                   bool           `json:"individual"`
	LimitedUserQuotas            any            `json:"limited_user_quotas"`
	LimitedUserResetDate         any            `json:"limited_user_reset_date"`
	NesEnabled                   bool           `json:"nes_enabled"`
	Prompt8K                     bool           `json:"prompt_8k"`
	PublicSuggestions            string         `json:"public_suggestions"`
	RefreshIn                    int            `json:"refresh_in"`
	Sku                          string         `json:"sku"`
	SnippyLoadTestEnabled        bool           `json:"snippy_load_test_enabled"`
	Telemetry                    string         `json:"telemetry"`
	Token                        string         `json:"token"`
	TrackingID                   string         `json:"tracking_id"`
	TriggerCompletionAfterAccept bool           `json:"trigger_completion_after_accept"`
	VscElectronFetcherV2         bool           `json:"vsc_electron_fetcher_v2"`
	Xcode                        bool           `json:"xcode"`
	XcodeChat                    bool           `json:"xcode_chat"`
}

type ModelCapabilityLimit struct {
	MaxContextWindowTokens int `json:"max_context_window_tokens"`
	MaxOutputTokens        int `json:"max_output_tokens"`
	MaxPromptTokens        int `json:"max_prompt_tokens"`
}

type ModelCapabilities struct {
	Family   string               `json:"family"`
	Limits   ModelCapabilityLimit `json:"limits"`
	Object   string               `json:"object"`
	Supports struct {
		Streaming bool `json:"streaming"`
		ToolCalls bool `json:"tool_calls"`
	} `json:"supports"`
	Tokenizer string `json:"tokenizer"`
	Type      string `json:"type"`
}

type PromptMessage struct {
	Content string `json:"content"`
	Role    string `json:"role"`
}

type Model struct {
	Capabilities       ModelCapabilities `json:"capabilities"`
	ID                 string            `json:"id"`
	ModelPickerEnabled bool              `json:"model_picker_enabled"`
	Name               string            `json:"name"`
	Object             string            `json:"object"`
	Preview            bool              `json:"preview"`
	Vendor             string            `json:"vendor"`
	Version            string            `json:"version"`
}

type Agent struct {
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Default     bool   `json:"default"`
	Description string `json:"description"`
}

type CompletionResponse struct {
	Choices             []Choice             `json:"choices"`
	Created             int64                `json:"created"`
	ID                  string               `json:"id"`
	Model               string               `json:"model"`
	PromptFilterResults []PromptFilterResult `json:"prompt_filter_results"`
	SystemFingerprint   string               `json:"system_fingerprint"`
	Usage               Usage                `json:"usage"`
}

type Choice struct {
	ContentFilterResults ContentFilterResults `json:"content_filter_results"`
	FinishReason         string               `json:"finish_reason"`
	Index                int64                `json:"index"`
	Message              Message              `json:"message"`
}

type ContentFilterResults struct {
	Hate     Hate `json:"hate"`
	SelfHarm Hate `json:"self_harm"`
	Sexual   Hate `json:"sexual"`
	Violence Hate `json:"violence"`
}

type Hate struct {
	Filtered bool   `json:"filtered"`
	Severity string `json:"severity"`
}

type Message struct {
	Content string `json:"content"`
	Role    string `json:"role"`
}

type PromptFilterResult struct {
	ContentFilterResults ContentFilterResults `json:"content_filter_results"`
	PromptIndex          int64                `json:"prompt_index"`
}

type Usage struct {
	CompletionTokens        int64                   `json:"completion_tokens"`
	CompletionTokensDetails CompletionTokensDetails `json:"completion_tokens_details"`
	PromptTokens            int64                   `json:"prompt_tokens"`
	PromptTokensDetails     PromptTokensDetails     `json:"prompt_tokens_details"`
	TotalTokens             int64                   `json:"total_tokens"`
}

type CompletionTokensDetails struct {
	ReasoningTokens int64 `json:"reasoning_tokens"`
}

type PromptTokensDetails struct {
	CachedTokens int64 `json:"cached_tokens"`
}

type Copilot interface {
	FetchModels(ctx context.Context) (map[string]*Model, error)
	FetchAgents(ctx context.Context) (map[string]*Agent, error)

	Ask(ctx context.Context, prompt string, opts any) (string, error)
}

func NewCopilot() Copilot {
	c := &copilot{
		models:    make(map[string]*Model),
		agents:    make(map[string]*Agent),
		machineId: utils.GenerateMachineId(),
	}

	configPath := utils.GetConfigPath() + "/lazycopilot"
	c.githubToken = c.getCachedToken()
	_ = utils.LoadFileJson(configPath+"/token.json", &c.token)
	_ = utils.LoadFileJson(configPath+"/models.json", &c.models)
	_ = utils.LoadFileJson(configPath+"/agents.json", &c.agents)

	return c
}

type copilot struct {
	machineId   string
	sessionId   string
	token       *GithubToken
	agents      map[string]*Agent
	models      map[string]*Model
	histories   []PromptMessage
	githubToken *string
}

func (c *copilot) generateAskRequest(histories []PromptMessage, prompt, systemPrompt, model string, temperature float64, maxOutputToken int, stream bool) interface{} {
	isO1 := strings.HasPrefix(model, "o1")
	systemRole := "system"
	if isO1 {
		systemRole = "user"
	}
	messages := make([]PromptMessage, 0)
	if systemPrompt != "" {
		messages = append(messages, PromptMessage{
			Content: systemPrompt,
			Role:    systemRole,
		})
	}

	if len(histories) > 0 {
		messages = append(messages, histories...)
	}

	if prompt != "" {
		messages = append(messages, PromptMessage{
			Content: prompt,
			Role:    "user",
		})
	}

	body := map[string]any{
		"messages": messages,
		"model":    model,
		"stream":   stream,
		"n":        1,
	}

	if maxOutputToken > 0 {
		body["max_tokens"] = maxOutputToken
	}

	if !isO1 {
		body["temperature"] = temperature
		body["top_p"] = 1
	}

	return body
}

// Ask implements Copilot.
func (c *copilot) Ask(ctx context.Context, prompt string, opts any) (string, error) {
	prompt = strings.TrimSpace(prompt)

	systemPrompt := strings.TrimSpace(COPILOT_INSTRUCTIONS)
	temperature := 0.1

	model := "gpt-4o"
	models, err := c.FetchModels(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to fetch models: %v. Please check your network connection and try again.", err)
	}

	modelConfig, ok := models[model]
	if !ok {
		return "", fmt.Errorf("model %s not found. Please ensure the model ID is correct and try again.", model)
	}

	capabilities := modelConfig.Capabilities
	maxOutputToken := capabilities.Limits.MaxOutputTokens
	stream := false

	fullResponse := ""
	parseLine := func(line string) {
		if line == "" {
			return
		}

		var res CompletionResponse
		err := json.Unmarshal([]byte(line), &res)
		if err != nil {
			return
		}

		if len(res.Choices) == 0 {
			return
		}

		choice := res.Choices[0]
		content := choice.Message.Content
		if content != "" {
			fullResponse += content
		}
	}

	body := c.generateAskRequest(
		c.histories,
		prompt,
		systemPrompt,
		model,
		temperature,
		maxOutputToken,
		stream,
	)

	headers, err := c.generateHeaders(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to generate headers: %v. Please try again.", err)
	}

	res, err := utils.HttpRequest(ctx, utils.HttpOptions{
		Method:  http.MethodPost,
		Url:     "https://api.githubcopilot.com/chat/completions",
		Headers: headers,
		Body:    body,
	})
	if err != nil {
		return "", fmt.Errorf("failed to send request: %v. Please check your network connection and try again.", err)
	}

	if res.StatusCode != 200 {
		return "", fmt.Errorf("failed to fetch completion response (%d): %s. Please try again.", res.StatusCode, res.Status)
	}

	resBody, err := res.StringDecode()
	if err != nil {
		return "", fmt.Errorf("failed to decode response: %v. Please try again.", err)
	}

	if !stream {
		parseLine(resBody)
	}

	c.histories = append(c.histories, PromptMessage{
		Content: prompt,
		Role:    "user",
	})

	c.histories = append(c.histories, PromptMessage{
		Content: fullResponse,
		Role:    "assistant",
	})

	if fullResponse == "" {
		return "", fmt.Errorf("failed to get response. Please try again.")
	}

	return fullResponse, nil
}

// FetchAgents implements Copilot.
func (c *copilot) FetchAgents(ctx context.Context) (map[string]*Agent, error) {
	if len(c.agents) > 0 {
		return c.agents, nil
	}

	headers, err := c.generateHeaders(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to generate headers: %v. Please try again.", err)
	}

	res, err := utils.HttpRequest(ctx, utils.HttpOptions{
		Method:  http.MethodGet,
		Url:     "https://api.githubcopilot.com/agents",
		Headers: headers,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v. Please check your network connection and try again.", err)
	}

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("failed to fetch agents (%d): %s. Please try again.", res.StatusCode, res.Status)
	}

	var restResponse struct {
		Agents []*Agent `json:"agents"`
	}

	err = res.JsonDecode(&restResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response: %v. Please try again.", err)
	}

	for _, a := range restResponse.Agents {
		c.agents[a.Slug] = a
	}

	c.agents["copilot"] = &Agent{
		Name:        "copilot",
		Slug:        "copilot",
		Default:     true,
		Description: "Default noop agent",
	}

	err = utils.SaveFile(utils.GetConfigPath()+"/lazycopilot/agents.json", c.agents)
	if err != nil {
		return nil, fmt.Errorf("failed to save agents to file: %v. Please check your file permissions and try again.", err)
	}

	return c.agents, nil
}

// FetchModels implements Copilot.
func (c *copilot) FetchModels(ctx context.Context) (map[string]*Model, error) {
	if len(c.models) > 0 {
		return c.models, nil
	}

	headers, err := c.generateHeaders(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to generate headers: %v. Please try again.", err)
	}

	res, err := utils.HttpRequest(ctx, utils.HttpOptions{
		Method:  http.MethodGet,
		Url:     "https://api.githubcopilot.com/models",
		Headers: headers,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v. Please check your network connection and try again.", err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("failed to fetch models (%d): %s. Please try again.", res.StatusCode, res.Status)
	}

	var results struct {
		Data []*Model `json:"data"`
	}
	err = res.JsonDecode(&results)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response: %v. Please try again.", err)
	}

	for _, model := range results.Data {
		if model.ModelPickerEnabled {
			c.models[model.ID] = model
		}
	}

	err = utils.SaveFile(utils.GetConfigPath()+"/lazycopilot/models.json", c.models)
	if err != nil {
		return nil, fmt.Errorf("failed to save models to file: %v. Please check your file permissions and try again.", err)
	}

	return c.models, nil
}

func (c *copilot) authenticate(ctx context.Context) error {
	if c.githubToken == nil || *c.githubToken == "" {
		return fmt.Errorf("GitHub token not set. Please authenticate using 'auth login' command.")
	}

	if c.token == nil || (c.token.ExpiresAt != 0 && c.token.ExpiresAt <= int(time.Now().Unix())) {
		sessionId := uuid.New().String() + "-" + fmt.Sprint(time.Now().UnixMicro())
		res, err := utils.HttpRequest(ctx, utils.HttpOptions{
			Method: http.MethodGet,
			Url:    "https://api.github.com/copilot_internal/v2/token",
			Headers: &utils.Headers{
				"Authorization": "Bearer " + *c.githubToken,
				"Accept":        "application/json",
			},
		})
		if err != nil {
			return fmt.Errorf("failed to request token: %v. Please check your network connection and try again.", err)
		}
		if res.StatusCode != 200 {
			return fmt.Errorf("failed to authenticate (%d): %s. Please try again.", res.StatusCode, res.Status)
		}
		var token GithubToken
		err = res.JsonDecode(&token)
		if err != nil {
			return fmt.Errorf("failed to decode token response: %v. Please try again.", err)
		}
		c.sessionId = sessionId
		c.token = &token

		err = utils.SaveFile(utils.GetConfigPath()+"/lazycopilot/token.json", c.token)
		if err != nil {
			return fmt.Errorf("failed to save token to file: %v. Please check your file permissions and try again.", err)
		}
	}

	if c.sessionId == "" {
		c.sessionId = uuid.New().String() + "-" + fmt.Sprint(time.Now().UnixMicro())
	}

	return nil
}

func (c *copilot) generateHeaders(ctx context.Context) (*utils.Headers, error) {
	err := c.authenticate(ctx)
	if err != nil {
		return nil, err
	}
	headers := &utils.Headers{
		"authorization":          "Bearer " + c.token.Token,
		"x-request-id":           uuid.New().String(),
		"vscode-sessionid":       c.sessionId,
		"vscode-machineid":       c.machineId,
		"copilot-integration-id": "vscode-chat",
		"openai-organization":    "github-copilot",
		"openai-intent":          "conversation-panel",
		"content-type":           "application/json",
	}
	return headers, nil
}

func (c *copilot) getCachedToken() *string {
	configPath := utils.GetConfigPath()

	// token can be stored in apps.json or hosts.json
	filePaths := []string{
		configPath + "/github-copilot/hosts.json",
		configPath + "/github-copilot/apps.json",
	}

	for _, path := range filePaths {
		if utils.IsFileExists(path) {
			f, err := os.Open(path)
			if err != nil {
				return nil
			}
			defer f.Close()
			var token map[string]utils.CachedToken
			if err := json.NewDecoder(f).Decode(&token); err != nil {
				return nil
			}
			for key, value := range token {
				if strings.HasPrefix(key, "github.com") {
					return &value.OauthToken
				}
			}
		}
	}
	return nil
}
