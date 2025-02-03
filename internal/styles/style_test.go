package styles

import (
	"testing"
)

func TestInfer(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		wantStyleName  string
		wantScore      float64
		wantIsHeader   bool
		wantIsFooter   bool
		scoreThreshold float64
	}{
		{
			name:           "Exact match with simple header",
			input:          "----------------------------------------",
			wantStyleName:  "Simple",
			wantScore:      1.0,
			wantIsHeader:   true,
			wantIsFooter:   true,
			scoreThreshold: 0.001,
		},
		{
			name:           "Match brackets style header",
			input:          "[ License Start ]----------------------",
			wantStyleName:  "Brackets",
			wantScore:      1.0,
			wantIsHeader:   true,
			wantIsFooter:   false,
			scoreThreshold: 0.001,
		},
		{
			name:           "Match brackets style footer",
			input:          "[ License End ]------------------------",
			wantStyleName:  "Brackets",
			wantScore:      1.0,
			wantIsHeader:   false,
			wantIsFooter:   true,
			scoreThreshold: 0.001,
		},
		{
			name:           "Empty input",
			input:          "",
			wantStyleName:  "",
			wantScore:      0.0,
			wantIsHeader:   false,
			wantIsFooter:   false,
			scoreThreshold: 0.001,
		},
		{
			name:           "Non-matching input",
			input:          "This is just some random text",
			wantStyleName:  "",
			wantScore:      0.0,
			wantIsHeader:   false,
			wantIsFooter:   false,
			scoreThreshold: 0.001,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Infer(tt.input)

			// Check style name
			if tt.wantStyleName != "" {
				if got.Style.Name != tt.wantStyleName {
					t.Errorf("Infer() style = %v, want %v", got.Style.Name, tt.wantStyleName)
				}
			}

			// Check score within threshold
			if diff := got.Score - tt.wantScore; diff < -tt.scoreThreshold ||
				diff > tt.scoreThreshold {
				t.Errorf(
					"Infer() score = %v, want %v (±%v)",
					got.Score,
					tt.wantScore,
					tt.scoreThreshold,
				)
			}

			// Check header/footer flags
			if got.IsHeader != tt.wantIsHeader {
				t.Errorf("Infer() isHeader = %v, want %v", got.IsHeader, tt.wantIsHeader)
			}
			if got.IsFooter != tt.wantIsFooter {
				t.Errorf("Infer() isFooter = %v, want %v", got.IsFooter, tt.wantIsFooter)
			}
		})
	}
}

func TestCalculateSimilarity(t *testing.T) {
	tests := []struct {
		name      string
		a         string
		b         string
		want      float64
		threshold float64
	}{
		{
			name:      "Exact match",
			a:         "----------------------------------------",
			b:         "----------------------------------------",
			want:      1.0,
			threshold: 0.001,
		},
		{
			name:      "Case insensitive match",
			a:         "LICENSE",
			b:         "license",
			want:      1.0,
			threshold: 0.001,
		},
		{
			name:      "No match",
			a:         "abc",
			b:         "xyz",
			want:      0.0,
			threshold: 0.001,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := calculateSimilarity(tt.a, tt.b)
			if diff := got - tt.want; diff < -tt.threshold || diff > tt.threshold {
				t.Errorf("calculateSimilarity() = %v, want %v (±%v)", got, tt.want, tt.threshold)
			}
		})
	}
}
