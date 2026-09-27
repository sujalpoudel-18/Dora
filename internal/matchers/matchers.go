package matchers

import (
	"regexp"
	"strings"
	"time"

	"github.com/dora/dora/internal/protocols/http"
	"github.com/dora/dora/pkg/template"
)

// Evaluator checks whether a response matches template matchers
type Evaluator struct{}

func NewEvaluator() *Evaluator {
	return &Evaluator{}
}

// MatchHTTP evaluates a list of matchers against an HTTP response
func (e *Evaluator) MatchHTTP(matchers []template.Matcher, condition string, resp *http.Response) bool {
	if len(matchers) == 0 {
		return false
	}

	isAND := strings.ToLower(condition) == "and" || condition == ""
	results := make([]bool, len(matchers))

	for i, m := range matchers {
		matched := e.evalSingleMatcher(&m, resp)
		if m.Negative {
			matched = !matched
		}
		results[i] = matched

		if isAND && !matched {
			return false
		}
		if !isAND && matched {
			return true
		}
	}

	if isAND {
		for _, r := range results {
			if !r {
				return false
			}
		}
		return true
	}

	return false
}

func (e *Evaluator) evalSingleMatcher(m *template.Matcher, resp *http.Response) bool {
	targetText := resp.Body
	if strings.ToLower(m.Part) == "header" || strings.ToLower(m.Part) == "headers" {
		var sb strings.Builder
		for k, v := range resp.Headers {
			sb.WriteString(k + ": " + v + "\n")
		}
		targetText = sb.String()
	}

	switch strings.ToLower(m.Type) {
	case "status":
		for _, s := range m.Status {
			if resp.StatusCode == s {
				return true
			}
		}
		return false

	case "word", "words":
		isAND := strings.ToLower(m.Condition) == "and"
		for _, w := range m.Words {
			found := strings.Contains(targetText, w)
			if isAND && !found {
				return false
			}
			if !isAND && found {
				return true
			}
		}
		return isAND

	case "regex":
		for _, r := range m.Regex {
			re, err := regexp.Compile(r)
			if err == nil && re.MatchString(targetText) {
				return true
			}
		}
		return false

	case "size":
		for _, sz := range m.Size {
			if resp.ResponseSize == int64(sz) {
				return true
			}
		}
		return false

	case "response-time":
		// Match if response time is higher than specified threshold
		if len(m.Size) > 0 {
			if resp.ResponseTime >= time.Duration(m.Size[0])*time.Millisecond {
				return true
			}
		}
		return false

	default:
		return false
	}
}
