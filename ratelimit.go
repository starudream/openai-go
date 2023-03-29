package openai

import (
	"net/http"
	"strconv"
	"time"
)

type Ratelimit struct {
	LimitRequests     int
	LimitTokens       int
	RemainingRequests int
	RemainingTokens   int
	ResetRequests     time.Duration
	ResetTokens       time.Duration
}

func GetRatelimit(resp *http.Response) *Ratelimit {
	if resp == nil {
		return nil
	}
	return &Ratelimit{
		LimitRequests:     convertHeaderToInt(resp.Header, "X-Ratelimit-Limit-Requests"),
		LimitTokens:       convertHeaderToInt(resp.Header, "X-Ratelimit-Limit-Tokens"),
		RemainingRequests: convertHeaderToInt(resp.Header, "X-Ratelimit-Remaining-Requests"),
		RemainingTokens:   convertHeaderToInt(resp.Header, "X-Ratelimit-Remaining-Tokens"),
		ResetRequests:     convertHeaderToDuration(resp.Header, "X-Ratelimit-Reset-Requests"),
		ResetTokens:       convertHeaderToDuration(resp.Header, "X-Ratelimit-Reset-Tokens"),
	}
}

func convertHeaderToInt(header http.Header, key string) int {
	if header == nil {
		return 0
	}
	v := header.Get(key)
	if v == "" {
		return 0
	}
	i, err := strconv.Atoi(v)
	if err != nil {
		return 0
	}
	return i
}

func convertHeaderToDuration(header http.Header, key string) time.Duration {
	if header == nil {
		return 0
	}
	v := header.Get(key)
	if v == "" {
		return 0
	}
	i, err := time.ParseDuration(v)
	if err != nil {
		return 0
	}
	return i
}
