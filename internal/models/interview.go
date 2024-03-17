package models

type InterviewResults struct {
	PublicID string
	Result   []byte
}

type Interview struct {
	PublicID string `json:"public_id"`
	Result   Result `json:"result"`
}

type Question struct {
	Question       string          `json:"question"`
	Evaluation     string          `json:"evaluation"`
	Score          int             `json:"score"`
	VideoLink      string          `json:"video_link"`
	EmotionResults []EmotionResult `json:"emotion_results"`
}

type EmotionResult struct {
	Emotion   string  `json:"emotion"`
	ExactTime float64 `json:"exact_time"`
	Duration  float64 `json:"duration"`
}

type Result struct {
	Questions []Question `json:"questions"`
	Score     int        `json:"score"`
}
