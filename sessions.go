package main

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Preços públicos gpt-6-luna por milhão de tokens (USD).
const (
	priceInputPerM      = 3.00
	priceOutputPerM     = 15.00
	priceCachedReadPerM = 0.75
	priceCacheWritePerM = 3.75
	priceReasoningPerM  = 15.00
)

type tokenUsagePayload struct {
	Usage struct {
		InputTokens           int `json:"input_tokens"`
		CachedInputTokens     int `json:"cached_input_tokens"`
		CacheWriteInputTokens int `json:"cache_write_input_tokens"`
		OutputTokens          int `json:"output_tokens"`
		ReasoningOutputTokens int `json:"reasoning_output_tokens"`
	} `json:"usage"`
}

type sessionEvent struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type SessionStats struct {
	InputTokens  int
	OutputTokens int
	CachedRead   int
	CacheWrite   int
	Reasoning    int
	CostUSD      float64
	Turns        int
}

func calcCost(s SessionStats) float64 {
	billableInput := s.InputTokens - s.CachedRead - s.CacheWrite
	billableInput = max(billableInput, 0)
	return float64(billableInput)*priceInputPerM/1e6 +
		float64(s.CachedRead)*priceCachedReadPerM/1e6 +
		float64(s.CacheWrite)*priceCacheWritePerM/1e6 +
		float64(s.OutputTokens-s.Reasoning)*priceOutputPerM/1e6 +
		float64(s.Reasoning)*priceReasoningPerM/1e6
}

func readSessionFile(path string) SessionStats {
	f, err := os.Open(path)
	if err != nil {
		return SessionStats{}
	}
	defer f.Close()

	var s SessionStats
	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 1<<20), 1<<20)
	for scanner.Scan() {
		var ev sessionEvent
		if err := json.Unmarshal(scanner.Bytes(), &ev); err != nil {
			continue
		}
		if ev.Type != "token_usage_record" {
			continue
		}
		var p tokenUsagePayload
		if err := json.Unmarshal(ev.Payload, &p); err != nil {
			continue
		}
		u := p.Usage
		// token_usage_record acumula — o último registro tem o total da thread
		s.InputTokens = u.InputTokens
		s.OutputTokens = u.OutputTokens
		s.CachedRead = u.CachedInputTokens
		s.CacheWrite = u.CacheWriteInputTokens
		s.Reasoning = u.ReasoningOutputTokens
		s.Turns++
	}
	s.CostUSD = calcCost(s)
	return s
}

func codexSessionsDir() string {
	if h := os.Getenv("CODEX_HOME"); h != "" {
		return filepath.Join(h, "sessions")
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".codex", "sessions")
}

// ReadTodayStats soma todos os arquivos .jsonl de hoje.
func ReadTodayStats() SessionStats {
	base := codexSessionsDir()
	today := time.Now().UTC().Format("2006/01/02")
	dir := filepath.Join(base, today)

	entries, err := os.ReadDir(dir)
	if err != nil {
		return SessionStats{}
	}

	var total SessionStats
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".jsonl") {
			continue
		}
		s := readSessionFile(filepath.Join(dir, e.Name()))
		total.InputTokens += s.InputTokens
		total.OutputTokens += s.OutputTokens
		total.CachedRead += s.CachedRead
		total.CacheWrite += s.CacheWrite
		total.Reasoning += s.Reasoning
		total.Turns += s.Turns
	}
	total.CostUSD = calcCost(total)
	return total
}

// ReadSessionStats lê o arquivo .jsonl mais recente de hoje (sessão ativa).
func ReadSessionStats() SessionStats {
	base := codexSessionsDir()
	today := time.Now().UTC().Format("2006/01/02")
	dir := filepath.Join(base, today)

	entries, err := os.ReadDir(dir)
	if err != nil {
		return SessionStats{}
	}

	var latest string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".jsonl") && e.Name() > latest {
			latest = e.Name()
		}
	}
	if latest == "" {
		return SessionStats{}
	}
	return readSessionFile(filepath.Join(dir, latest))
}

