package printtac

import (
	"fmt"
	"strconv"
	"strings"
)

type Rec struct {
	Title, Body string
	Tags        []string
}

const MaxTAC = 320

func Sample() Rec {
	return Rec{Title: "WO-1042-plate", Body: "C=80 M=60 Y=55 K=70", Tags: []string{"C", "M", "Y", "K"}}
}

func Seed() []Rec {
	return []Rec{
		Sample(),
		{Title: "WO-1042-spot", Body: "C=10 M=90 Y=80 K=0 Pantone185=40", Tags: []string{"C", "M", "Y", "K", "Pantone185"}},
	}
}

func AfterWrite(getMin func() (string, error), setMin func(string) error, body string) error {
	return nil
}

func Steps() []string { return []string{"tac-check", "index-plates", "export-inkkeys"} }

func Enforce(title, body string, tags []string) error {
	if strings.TrimSpace(title) == "" {
		return fmt.Errorf("work order title required")
	}
	tac, err := TotalCoverage(body)
	if err != nil {
		return err
	}
	if tac > MaxTAC {
		return fmt.Errorf("TAC %d exceeds %d", tac, MaxTAC)
	}
	if len(tags) == 0 {
		return fmt.Errorf("ink channel tags required")
	}
	return nil
}

func TotalCoverage(body string) (int, error) {
	sum := 0
	found := 0
	for _, part := range strings.Fields(body) {
		k, v, ok := strings.Cut(part, "=")
		if !ok || k == "" {
			continue
		}
		n, err := strconv.Atoi(strings.TrimSuffix(v, "%"))
		if err != nil {
			return 0, fmt.Errorf("coverage %q: %w", part, err)
		}
		if n < 0 || n > 100 {
			return 0, fmt.Errorf("coverage %s out of 0..100", k)
		}
		sum += n
		found++
	}
	if found == 0 {
		return 0, fmt.Errorf("body must list channel coverage like C=80")
	}
	return sum, nil
}
