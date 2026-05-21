package overtime_calc

import (
	"math"
	"testing"
)

func TestCalculateWeekday(t *testing.T) {
	salary := 5000000.0
	meal := 30000.0
	hourlyRate := salary / 173.0

	tests := []struct {
		name     string
		hours    float64
		expected float64
	}{
		{"1 hour", 1.0, hourlyRate * 1.5},
		{"1.5 hours", 1.5, hourlyRate * 2.5},
		{"2 hours", 2.0, hourlyRate * 3.5},
		{"2.5 hours", 2.5, hourlyRate * 4.5},
		{"3 hours", 3.0, hourlyRate * 5.5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Calculate(tt.hours, salary, meal, false)

			expectedBase := math.Round(tt.expected*100) / 100
			if result.BaseAmount != expectedBase {
				t.Errorf("BaseAmount = %v, want %v", result.BaseAmount, expectedBase)
			}

			if result.MealAmount != meal {
				t.Errorf("MealAmount = %v, want %v", result.MealAmount, meal)
			}

			expectedTotal := math.Round((tt.expected+meal)*100) / 100
			if result.Total != expectedTotal {
				t.Errorf("Total = %v, want %v", result.Total, expectedTotal)
			}
		})
	}
}

func TestCalculateHoliday(t *testing.T) {
	salary := 5000000.0
	meal := 30000.0
	hourlyRate := salary / 173.0

	tests := []struct {
		name     string
		hours    float64
		expected float64
	}{
		{"1 hour holiday", 1.0, hourlyRate * 2 * 1.0},
		{"2 hours holiday", 2.0, hourlyRate * 2 * 2.0},
		{"3 hours holiday", 3.0, hourlyRate * 2 * 3.0},
		{"4 hours holiday", 4.0, hourlyRate * 2 * 4.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Calculate(tt.hours, salary, meal, true)

			expectedBase := math.Round(tt.expected*100) / 100
			if result.BaseAmount != expectedBase {
				t.Errorf("BaseAmount = %v, want %v", result.BaseAmount, expectedBase)
			}

			if result.MealAmount != meal {
				t.Errorf("MealAmount = %v, want %v", result.MealAmount, meal)
			}
		})
	}
}

func TestCalculateZeroSalary(t *testing.T) {
	result := Calculate(2.0, 0, 30000, false)
	if result.BaseAmount != 0 {
		t.Errorf("BaseAmount should be 0 for zero salary, got %v", result.BaseAmount)
	}
	if result.MealAmount != 30000 {
		t.Errorf("MealAmount should still be 30000, got %v", result.MealAmount)
	}
}
