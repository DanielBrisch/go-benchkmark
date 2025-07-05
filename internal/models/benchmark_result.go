package models

type BenchmarkResult struct {
	Language     string  `json:"language"`
	LoadTimeSec  float64 `json:"load_time_sec"`
	QueryTimeSec float64 `json:"query_time_sec"`
	PeakMemoryMB uint64  `json:"peak_memory_mb"`
	DatasetSize  string  `json:"dataset_size"`
}
