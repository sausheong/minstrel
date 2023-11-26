package main

import "time"

// Ollama completion request
type CompletionRequest struct {
	Model   string  `json:"model"`
	Prompt  string  `json:"prompt"`
	Format  string  `json:"format,omitempty"`
	Options Options `json:"options,omitempty"`
	System  string  `json:"system,omitempty"`
	Context []byte  `json:"context,omitempty"`
	Stream  bool    `json:"stream"`
}

// Ollama completion response
type CompletionResponse struct {
	Model              string        `json:"model"`
	CreatedAt          time.Time     `json:"created_at"`
	Response           string        `json:"response"`
	Done               bool          `json:"done"`
	Context            []int         `json:"context,omitempty"`
	TotalDuration      time.Duration `json:"total_duration,omitempty"`
	LoadDuration       time.Duration `json:"load_duration,omitempty"`
	PromptEvalCount    int           `json:"prompt_eval_count,omitempty"`
	PromptEvalDuration time.Duration `json:"prompt_eval_duration,omitempty"`
	EvalCount          int           `json:"eval_count,omitempty"`
	EvalDuration       time.Duration `json:"eval_duration,omitempty"`
}

// Ollama completion request options
type Options struct {
	NumKeep            int      `json:"num_keep,omitempty"`
	Seed               int      `json:"seed,omitempty"`
	NumPredict         int      `json:"num_predict,omitempty"`
	TopK               int      `json:"top_k,omitempty"`
	TopP               float64  `json:"top_p,omitempty"`
	TfsZ               float64  `json:"tfs_z,omitempty"`
	TypicalP           float64  `json:"typical_p,omitempty"`
	RepeatLastN        int      `json:"repeat_last_n,omitempty"`
	Temperature        float64  `json:"temperature,omitempty"`
	RepeatPenalty      float64  `json:"repeat_penalty,omitempty"`
	PresencePenalty    float64  `json:"presence_penalty,omitempty"`
	FrequencyPenalty   float64  `json:"frequency_penalty,omitempty"`
	Mirostat           int      `json:"mirostat,omitempty"`
	MirostatTau        float64  `json:"mirostat_tau,omitempty"`
	MirostatEta        float64  `json:"mirostat_eta,omitempty"`
	PenalizeNewline    bool     `json:"penalize_newline,omitempty"`
	Stop               []string `json:"stop,omitempty"`
	Numa               bool     `json:"numa,omitempty"`
	NumCtx             int      `json:"num_ctx,omitempty"`
	NumBatch           int      `json:"num_batch,omitempty"`
	NumGqa             int      `json:"num_gqa,omitempty"`
	NumGpu             int      `json:"num_gpu,omitempty"`
	MainGpu            int      `json:"main_gpu,omitempty"`
	LowVram            bool     `json:"low_vram,omitempty"`
	F16Kv              bool     `json:"f16_kv,omitempty"`
	LogitsAll          bool     `json:"logits_all,omitempty"`
	VocabOnly          bool     `json:"vocab_only,omitempty"`
	UseMmap            bool     `json:"use_mmap,omitempty"`
	UseMlock           bool     `json:"use_mlock,omitempty"`
	EmbeddingOnly      bool     `json:"embedding_only,omitempty"`
	RopeFrequencyBase  float64  `json:"rope_frequency_base,omitempty"`
	RopeFrequencyScale float64  `json:"rope_frequency_scale,omitempty"`
	NumThread          int      `json:"num_thread,omitempty"`
}

type ReplicateSDXLRequest struct {
	Version string `json:"version,omitempty"`
	Input   Input  `json:"input,omitempty"`
}

type ReplicateSDXLResponse struct {
	CompletedAt time.Time `json:"completed_at,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	Error       any       `json:"error,omitempty"`
	ID          string    `json:"id,omitempty"`
	Input       Input     `json:"input,omitempty"`
	Logs        string    `json:"logs,omitempty"`
	Metrics     Metrics   `json:"metrics,omitempty"`
	Output      []string  `json:"output,omitempty"`
	StartedAt   time.Time `json:"started_at,omitempty"`
	Status      string    `json:"status,omitempty"`
	Urls        Urls      `json:"urls,omitempty"`
	Version     string    `json:"version,omitempty"`
}
type Input struct {
	Width             int     `json:"width,omitempty"`
	Height            int     `json:"height,omitempty"`
	Prompt            string  `json:"prompt,omitempty"`
	Refine            string  `json:"refine,omitempty"`
	Scheduler         string  `json:"scheduler,omitempty"`
	LoraScale         float64 `json:"lora_scale,omitempty"`
	NumOutputs        int     `json:"num_outputs,omitempty"`
	GuidanceScale     float64 `json:"guidance_scale,omitempty"`
	ApplyWatermark    bool    `json:"apply_watermark,omitempty"`
	HighNoiseFrac     float64 `json:"high_noise_frac,omitempty"`
	NegativePrompt    string  `json:"negative_prompt,omitempty"`
	PromptStrength    float64 `json:"prompt_strength,omitempty"`
	NumInferenceSteps int     `json:"num_inference_steps,omitempty"`
}
type Metrics struct {
	PredictTime float64 `json:"predict_time,omitempty"`
}
type Urls struct {
	Get    string `json:"get,omitempty"`
	Cancel string `json:"cancel,omitempty"`
}

type Story struct {
	Title       string `json:"title"`
	Plot        string `json:"plot"`
	AuthorStyle string `json:"author_style"`
	Genre       string `json:"genre"`
}

type OpenAIResponse struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int      `json:"created"`
	Model   string   `json:"model"`
	Usage   Usage    `json:"usage"`
	Choices []Choice `json:"choices"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type FinishDetails struct {
	Type string `json:"type"`
	Stop string `json:"stop"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Choice struct {
	Message       Message       `json:"message"`
	FinishDetails FinishDetails `json:"finish_details"`
	Index         int           `json:"index"`
}
