package yappgen

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/wesen/yapp-encl-resolver/pkg/resolver"
)

func TestHasPerSideClearance(t *testing.T) {
	tests := []struct {
		name     string
		resolved map[string]any
		want     bool
	}{
		{
			name: "per-side clearance present",
			resolved: map[string]any{
				"enclosure": map[string]any{
					"wall": map[string]any{
						"clearance": map[string]any{
							"front": 2.0,
							"back":  1.5,
							"left":  1.0,
							"right": 1.0,
						},
					},
				},
			},
			want: true,
		},
		{
			name: "uniform clearance (number)",
			resolved: map[string]any{
				"enclosure": map[string]any{
					"wall": map[string]any{
						"clearance": 1.5,
					},
				},
			},
			want: false,
		},
		{
			name:     "no clearance",
			resolved:  map[string]any{},
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := hasPerSideClearance(tt.resolved)
			if got != tt.want {
				t.Errorf("hasPerSideClearance() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHasFinalDimensions(t *testing.T) {
	tests := []struct {
		name     string
		resolved map[string]any
		want     bool
	}{
		{
			name: "both dimensions present",
			resolved: map[string]any{
				"enclosure": map[string]any{
					"dimensions": map[string]any{
						"length": 100.0,
						"width":  80.0,
					},
				},
			},
			want: true,
		},
		{
			name: "only length present",
			resolved: map[string]any{
				"enclosure": map[string]any{
					"dimensions": map[string]any{
						"length": 100.0,
					},
				},
			},
			want: true,
		},
		{
			name: "only width present",
			resolved: map[string]any{
				"enclosure": map[string]any{
					"dimensions": map[string]any{
						"width": 80.0,
					},
				},
			},
			want: true,
		},
		{
			name: "empty dimensions",
			resolved: map[string]any{
				"enclosure": map[string]any{
					"dimensions": map[string]any{},
				},
			},
			want: false,
		},
		{
			name:     "no dimensions",
			resolved:  map[string]any{},
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := hasFinalDimensions(tt.resolved)
			if got != tt.want {
				t.Errorf("hasFinalDimensions() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHasUniformClearance(t *testing.T) {
	tests := []struct {
		name     string
		resolved map[string]any
		want     bool
	}{
		{
			name: "uniform clearance (float)",
			resolved: map[string]any{
				"enclosure": map[string]any{
					"wall": map[string]any{
						"clearance": 1.5,
					},
				},
			},
			want: true,
		},
		{
			name: "uniform clearance (int)",
			resolved: map[string]any{
				"enclosure": map[string]any{
					"wall": map[string]any{
						"clearance": 2,
					},
				},
			},
			want: true,
		},
		{
			name: "per-side clearance",
			resolved: map[string]any{
				"enclosure": map[string]any{
					"wall": map[string]any{
						"clearance": map[string]any{
							"front": 2.0,
						},
					},
				},
			},
			want: false,
		},
		{
			name:     "no clearance",
			resolved:  map[string]any{},
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := hasUniformClearance(tt.resolved)
			if got != tt.want {
				t.Errorf("hasUniformClearance() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidateEnclosureDimensions(t *testing.T) {
	tests := []struct {
		name     string
		resolved map[string]any
		wantErr  bool
		errMsg   string
	}{
		{
			name: "valid per-side clearance",
			resolved: map[string]any{
				"enclosure": map[string]any{
					"wall": map[string]any{
						"clearance": map[string]any{
							"front": 2.0,
							"back":  1.5,
							"left":  1.0,
							"right": 1.0,
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "valid uniform clearance",
			resolved: map[string]any{
				"enclosure": map[string]any{
					"wall": map[string]any{
						"clearance": 1.5,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "valid final dimensions",
			resolved: map[string]any{
				"enclosure": map[string]any{
					"dimensions": map[string]any{
						"length": 100.0,
						"width":  80.0,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "error: per-side + final dimensions",
			resolved: map[string]any{
				"enclosure": map[string]any{
					"wall": map[string]any{
						"clearance": map[string]any{
							"front": 2.0,
							"back":  1.5,
							"left":  1.0,
							"right": 1.0,
						},
					},
					"dimensions": map[string]any{
						"length": 100.0,
						"width":  80.0,
					},
				},
			},
			wantErr: true,
			errMsg:  "cannot specify both clearance",
		},
		{
			name: "error: uniform + final dimensions",
			resolved: map[string]any{
				"enclosure": map[string]any{
					"wall": map[string]any{
						"clearance": 1.5,
					},
					"dimensions": map[string]any{
						"length": 100.0,
						"width":  80.0,
					},
				},
			},
			wantErr: true,
			errMsg:  "cannot specify both clearance",
		},
		{
			name: "error: incomplete per-side clearance",
			resolved: map[string]any{
				"enclosure": map[string]any{
					"wall": map[string]any{
						"clearance": map[string]any{
							"front": 2.0,
							"back":  1.5,
							// missing left and right
						},
					},
				},
			},
			wantErr: true,
			errMsg:  "missing required keys",
		},
		{
			name: "error: empty dimensions",
			resolved: map[string]any{
				"enclosure": map[string]any{
					"dimensions": map[string]any{},
				},
			},
			wantErr: true,
			errMsg:  "neither length nor width provided",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateEnclosureDimensions(tt.resolved)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateEnclosureDimensions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err != nil {
				if tt.errMsg != "" && !contains(err.Error(), tt.errMsg) {
					t.Errorf("validateEnclosureDimensions() error = %v, want error containing %q", err, tt.errMsg)
				}
			}
		})
	}
}

func TestBuildModel_PerSideClearance(t *testing.T) {
	resolved := map[string]any{
		"pcb": map[string]any{
			"length":    90.0,
			"width":     70.0,
			"thickness": 1.6,
		},
		"enclosure": map[string]any{
			"wall": map[string]any{
				"thickness": 2.4,
				"clearance": map[string]any{
					"front": 2.0,
					"back":  1.5,
					"left":  1.0,
					"right": 1.0,
				},
			},
		},
	}

	ctx := context.Background()
	trace := resolver.Trace{}
	comments := map[string][]string{}
	raw := map[string]any{}

	model, err := BuildModel(ctx, resolved, trace, comments, raw)
	if err != nil {
		t.Fatalf("BuildModel() error = %v", err)
	}

	if model.PaddingFront != 2.0 {
		t.Errorf("PaddingFront = %v, want 2.0", model.PaddingFront)
	}
	if model.PaddingBack != 1.5 {
		t.Errorf("PaddingBack = %v, want 1.5", model.PaddingBack)
	}
	if model.PaddingLeft != 1.0 {
		t.Errorf("PaddingLeft = %v, want 1.0", model.PaddingLeft)
	}
	if model.PaddingRight != 1.0 {
		t.Errorf("PaddingRight = %v, want 1.0", model.PaddingRight)
	}
}

func TestBuildModel_FinalDimensions(t *testing.T) {
	resolved := map[string]any{
		"pcb": map[string]any{
			"length":    90.0,
			"width":     70.0,
			"thickness": 1.6,
		},
		"enclosure": map[string]any{
			"wall": map[string]any{
				"thickness": 2.4,
			},
			"dimensions": map[string]any{
				"length": 100.0,
				"width":  80.0,
			},
		},
	}

	ctx := context.Background()
	trace := resolver.Trace{}
	comments := map[string][]string{}
	raw := map[string]any{}

	model, err := BuildModel(ctx, resolved, trace, comments, raw)
	if err != nil {
		t.Fatalf("BuildModel() error = %v", err)
	}

	if model.FinalLength != 100.0 {
		t.Errorf("FinalLength = %v, want 100.0", model.FinalLength)
	}
	if model.FinalWidth != 80.0 {
		t.Errorf("FinalWidth = %v, want 80.0", model.FinalWidth)
	}

	// Verify padding was computed correctly
	// Length: 100.0 - (90.0 + 2.4*2) = 100.0 - 94.8 = 5.2
	// Padding per side: 5.2 / 2 = 2.6
	expectedPaddingLength := 2.6
	if !approxEqual(model.PaddingFront, expectedPaddingLength) || !approxEqual(model.PaddingBack, expectedPaddingLength) {
		t.Errorf("PaddingFront/Back = %v/%v, want ~%v for both", model.PaddingFront, model.PaddingBack, expectedPaddingLength)
	}

	// Width: 80.0 - (70.0 + 2.4*2) = 80.0 - 74.8 = 5.2
	// Padding per side: 5.2 / 2 = 2.6
	expectedPaddingWidth := 2.6
	if !approxEqual(model.PaddingLeft, expectedPaddingWidth) || !approxEqual(model.PaddingRight, expectedPaddingWidth) {
		t.Errorf("PaddingLeft/Right = %v/%v, want ~%v for both", model.PaddingLeft, model.PaddingRight, expectedPaddingWidth)
	}
}

func TestBuildModel_UniformClearance_BackwardCompatible(t *testing.T) {
	resolved := map[string]any{
		"pcb": map[string]any{
			"length":    90.0,
			"width":     70.0,
			"thickness": 1.6,
		},
		"enclosure": map[string]any{
			"wall": map[string]any{
				"thickness": 2.4,
				"clearance": 1.5,
			},
		},
	}

	ctx := context.Background()
	trace := resolver.Trace{}
	comments := map[string][]string{}
	raw := map[string]any{}

	model, err := BuildModel(ctx, resolved, trace, comments, raw)
	if err != nil {
		t.Fatalf("BuildModel() error = %v", err)
	}

	// All padding should be equal to uniform clearance
	if model.PaddingFront != 1.5 {
		t.Errorf("PaddingFront = %v, want 1.5", model.PaddingFront)
	}
	if model.PaddingBack != 1.5 {
		t.Errorf("PaddingBack = %v, want 1.5", model.PaddingBack)
	}
	if model.PaddingLeft != 1.5 {
		t.Errorf("PaddingLeft = %v, want 1.5", model.PaddingLeft)
	}
	if model.PaddingRight != 1.5 {
		t.Errorf("PaddingRight = %v, want 1.5", model.PaddingRight)
	}
}

func TestBuildModel_MutualExclusivityError(t *testing.T) {
	resolved := map[string]any{
		"pcb": map[string]any{
			"length":    90.0,
			"width":     70.0,
			"thickness": 1.6,
		},
		"enclosure": map[string]any{
			"wall": map[string]any{
				"thickness": 2.4,
				"clearance": map[string]any{
					"front": 2.0,
					"back":  1.5,
					"left":  1.0,
					"right": 1.0,
				},
			},
			"dimensions": map[string]any{
				"length": 100.0,
				"width":  80.0,
			},
		},
	}

	ctx := context.Background()
	trace := resolver.Trace{}
	comments := map[string][]string{}
	raw := map[string]any{}

	_, err := BuildModel(ctx, resolved, trace, comments, raw)
	if err == nil {
		t.Fatal("BuildModel() expected error for conflicting configuration, got nil")
	}
	if !contains(err.Error(), "cannot specify both") {
		t.Errorf("BuildModel() error = %v, want error containing 'cannot specify both'", err)
	}
}

func TestComputePaddingFromDimensions(t *testing.T) {
	tests := []struct {
		name          string
		model         *Model
		wantErr       bool
		errMsg        string
		checkPadding  func(*Model) error
	}{
		{
			name: "both dimensions specified",
			model: &Model{
				PcbLength:    90.0,
				PcbWidth:     70.0,
				WallThickness: 2.4,
				FinalLength:  100.0,
				FinalWidth:   80.0,
			},
			wantErr: false,
			checkPadding: func(m *Model) error {
				// Length: 100.0 - (90.0 + 2.4*2) = 5.2, per side = 2.6
				if !approxEqual(m.PaddingFront, 2.6) || !approxEqual(m.PaddingBack, 2.6) {
					return fmt.Errorf("PaddingFront/Back = %v/%v, want ~2.6", m.PaddingFront, m.PaddingBack)
				}
				// Width: 80.0 - (70.0 + 2.4*2) = 5.2, per side = 2.6
				if !approxEqual(m.PaddingLeft, 2.6) || !approxEqual(m.PaddingRight, 2.6) {
					return fmt.Errorf("PaddingLeft/Right = %v/%v, want ~2.6", m.PaddingLeft, m.PaddingRight)
				}
				return nil
			},
		},
		{
			name: "only length specified",
			model: &Model{
				PcbLength:    90.0,
				PcbWidth:     70.0,
				WallThickness: 2.4,
				FinalLength:  100.0,
				FinalWidth:   0,
			},
			wantErr: false,
			checkPadding: func(m *Model) error {
				// Length padding computed
				if !approxEqual(m.PaddingFront, 2.6) || !approxEqual(m.PaddingBack, 2.6) {
					return fmt.Errorf("PaddingFront/Back = %v/%v, want ~2.6", m.PaddingFront, m.PaddingBack)
				}
				// Width padding uses default (1.0)
				if !approxEqual(m.PaddingLeft, 1.0) || !approxEqual(m.PaddingRight, 1.0) {
					return fmt.Errorf("PaddingLeft/Right = %v/%v, want ~1.0 (default)", m.PaddingLeft, m.PaddingRight)
				}
				return nil
			},
		},
		{
			name: "negative padding error",
			model: &Model{
				PcbLength:    90.0,
				PcbWidth:     70.0,
				WallThickness: 2.4,
				FinalLength:  90.0, // Too small
				FinalWidth:   80.0,
			},
			wantErr: true,
			errMsg:  "too small",
		},
		{
			name: "no dimensions specified",
			model: &Model{
				PcbLength:    90.0,
				PcbWidth:     70.0,
				WallThickness: 2.4,
				FinalLength:  0,
				FinalWidth:   0,
			},
			wantErr: true,
			errMsg:  "at least one final dimension must be specified",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.model.computePaddingFromDimensions()
			if (err != nil) != tt.wantErr {
				t.Errorf("computePaddingFromDimensions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				if tt.errMsg != "" && !contains(err.Error(), tt.errMsg) {
					t.Errorf("computePaddingFromDimensions() error = %v, want error containing %q", err, tt.errMsg)
				}
			} else {
				if tt.checkPadding != nil {
					if err := tt.checkPadding(tt.model); err != nil {
						t.Error(err)
					}
				}
			}
		})
	}
}

// Helper function
func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}

// approxEqual checks if two float64 values are approximately equal (within 0.0001)
func approxEqual(a, b float64) bool {
	const epsilon = 0.0001
	diff := a - b
	if diff < 0 {
		diff = -diff
	}
	return diff < epsilon
}

