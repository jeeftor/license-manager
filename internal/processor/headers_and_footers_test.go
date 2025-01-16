package processor

// import (
// 	"testing"
// )

// func TestInferHeaderAndFooterStyle(t *testing.T) {
// 	tests := []struct {
// 		name           string
// 		input          string
// 		wantStyleName  string
// 		wantScore      float64
// 		wantIsHeader   bool
// 		wantIsFooter   bool
// 		scoreThreshold float64 // How close the score needs to be
// 	}{
// 		{
// 			name:           "Exact match with simple header",
// 			input:          "----------------------------------------",
// 			wantStyleName:  "simple",
// 			wantScore:      1.0,
// 			wantIsHeader:   true,
// 			wantIsFooter:   true, // Both true because header and footer are identical
// 			scoreThreshold: 0.001,
// 		},
// 		{
// 			name:           "Match with extra spaces",
// 			input:          "-------- -------- -------- --------",
// 			wantStyleName:  "simple",
// 			wantScore:      0.9,
// 			wantIsHeader:   true,
// 			wantIsFooter:   true,
// 			scoreThreshold: 0.001,
// 		},
// 		{
// 			name:           "Match hash style",
// 			input:          "######################################",
// 			wantStyleName:  "hash",
// 			wantScore:      1.0,
// 			wantIsHeader:   true,
// 			wantIsFooter:   true,
// 			scoreThreshold: 0.001,
// 		},
// 		{
// 			name:           "Match box style header",
// 			input:          "+------------------------------------+",
// 			wantStyleName:  "box",
// 			wantScore:      1.0,
// 			wantIsHeader:   true,
// 			wantIsFooter:   true,
// 			scoreThreshold: 0.001,
// 		},
// 		{
// 			name:           "Match brackets style header",
// 			input:          "[ License Start ]----------------------",
// 			wantStyleName:  "brackets",
// 			wantScore:      1.0,
// 			wantIsHeader:   true,
// 			wantIsFooter:   false,
// 			scoreThreshold: 0.001,
// 		},
// 		{
// 			name:           "Match brackets style footer",
// 			input:          "[ License End ]------------------------",
// 			wantStyleName:  "brackets",
// 			wantScore:      1.0,
// 			wantIsHeader:   false,
// 			wantIsFooter:   true,
// 			scoreThreshold: 0.001,
// 		},
// 		{
// 			name:           "Similar pattern but different characters",
// 			input:          "========================================",
// 			wantStyleName:  "equals",
// 			wantScore:      1.0,
// 			wantIsHeader:   true,
// 			wantIsFooter:   true,
// 			scoreThreshold: 0.001,
// 		},
// 		{
// 			name:           "Partial match with pattern",
// 			input:          "-------- LICENSE --------",
// 			wantStyleName:  "simple",
// 			wantScore:      0.8,
// 			wantIsHeader:   true,
// 			wantIsFooter:   true,
// 			scoreThreshold: 0.1, // More lenient due to pattern matching
// 		},
// 		{
// 			name:           "Empty input",
// 			input:          "",
// 			wantStyleName:  "",
// 			wantScore:      0.0,
// 			wantIsHeader:   false,
// 			wantIsFooter:   false,
// 			scoreThreshold: 0.001,
// 		},
// 		{
// 			name:           "Non-matching input",
// 			input:          "This is just some random text",
// 			wantStyleName:  "",
// 			wantScore:      0.0,
// 			wantIsHeader:   false,
// 			wantIsFooter:   false,
// 			scoreThreshold: 0.3, // More lenient as it might find some pattern matches
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			got := InferHeaderAndFooterStyle(tt.input)

// 			// Check style name
// 			if tt.wantStyleName != "" {
// 				if got.Style.Name != tt.wantStyleName {
// 					t.Errorf("InferHeaderAndFooterStyle() style = %v, want %v", got.Style.Name, tt.wantStyleName)
// 				}
// 			}

// 			// Check score within threshold
// 			if diff := got.Score - tt.wantScore; diff < -tt.scoreThreshold || diff > tt.scoreThreshold {
// 				t.Errorf("InferHeaderAndFooterStyle() score = %v, want %v (±%v)", got.Score, tt.wantScore, tt.scoreThreshold)
// 			}

// 			// Check header/footer flags
// 			if got.IsHeader != tt.wantIsHeader {
// 				t.Errorf("InferHeaderAndFooterStyle() isHeader = %v, want %v", got.IsHeader, tt.wantIsHeader)
// 			}
// 			if got.IsFooter != tt.wantIsFooter {
// 				t.Errorf("InferHeaderAndFooterStyle() isFooter = %v, want %v", got.IsFooter, tt.wantIsFooter)
// 			}
// 		})
// 	}
// }

// func TestCalculateSimilarity(t *testing.T) {
// 	tests := []struct {
// 		name      string
// 		a         string
// 		b         string
// 		want      float64
// 		threshold float64
// 	}{
// 		{
// 			name:      "Exact match",
// 			a:         "----------------------------------------",
// 			b:         "----------------------------------------",
// 			want:      1.0,
// 			threshold: 0.001,
// 		},
// 		{
// 			name:      "Case insensitive match",
// 			a:         "LICENSE",
// 			b:         "license",
// 			want:      0.9,
// 			threshold: 0.001,
// 		},
// 		{
// 			name:      "Space difference",
// 			a:         "- - - -",
// 			b:         "----",
// 			want:      0.9,
// 			threshold: 0.001,
// 		},
// 		{
// 			name:      "Pattern match",
// 			a:         "====",
// 			b:         "----",
// 			want:      0.8,
// 			threshold: 0.1,
// 		},
// 		{
// 			name:      "Partial pattern match",
// 			a:         "===LICENSE===",
// 			b:         "---LICENSE---",
// 			want:      0.7,
// 			threshold: 0.1,
// 		},
// 		{
// 			name:      "No match",
// 			a:         "abc",
// 			b:         "xyz",
// 			want:      0.0,
// 			threshold: 0.1,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			got := calculateSimilarity(tt.a, tt.b)
// 			if diff := got - tt.want; diff < -tt.threshold || diff > tt.threshold {
// 				t.Errorf("calculateSimilarity() = %v, want %v (±%v)", got, tt.want, tt.threshold)
// 			}
// 		})
// 	}
// }
