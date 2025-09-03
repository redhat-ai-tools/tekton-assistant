package types

type KnowledgeEntry struct {
	ID           string                 `json:"id"`
	Error        string                 `json:"error"`
	Context      string                 `json:"context"`
	Solution     string                 `json:"solution"`
	Reference    string                 `json:"reference"`
	Source       string                 `json:"source"`
	Metadata     map[string]any         `json:"metadata"`
	CombinedText string                 `json:"combined_text"`
	ErrorType    string                 `json:"error_type"`
	Severity     string                 `json:"severity"`
}

type IngestionMetrics struct {
	TotalEntries          int                `json:"total_entries"`
	SuccessfulIngestions  int                `json:"successful_ingestions"`
	FailedIngestions      int                `json:"failed_ingestions"`
	StartTimeUnix         int64              `json:"start_time"`
	EndTimeUnix           int64              `json:"end_time"`
	ProcessingTimeSeconds float64            `json:"processing_time"`
	ErrorTypes            map[string]int     `json:"error_types"`
	SeverityCounts        map[string]int     `json:"severity_counts"`
}

func NewIngestionMetrics() IngestionMetrics {
	return IngestionMetrics{
		ErrorTypes:     map[string]int{},
		SeverityCounts: map[string]int{},
	}
}

func (m IngestionMetrics) SuccessRate() float64 {
	if m.TotalEntries == 0 {
		return 0
	}
	return float64(m.SuccessfulIngestions) / float64(m.TotalEntries) * 100.0
}

func (m IngestionMetrics) EntriesPerSecond() float64 {
	if m.ProcessingTimeSeconds == 0 {
		return 0
	}
	return float64(m.SuccessfulIngestions) / m.ProcessingTimeSeconds
}
