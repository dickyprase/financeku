package overtime_calc

import "math"

// OvertimeResult holds the calculation result
type OvertimeResult struct {
	BaseAmount float64 `json:"base_amount"`
	MealAmount float64 `json:"meal_amount"`
	Total      float64 `json:"total_amount"`
}

// Calculate calculates overtime amount based on hours, salary, meal allowance, and holiday status
// Weekday formula: multiplier based on hours bracket × (salary/173)
// Holiday formula: salary/173 × 2 × hours
func Calculate(hours float64, salary float64, mealAllowance float64, isHoliday bool) OvertimeResult {
	hourlyRate := salary / 173.0
	var baseAmount float64

	if isHoliday {
		baseAmount = hourlyRate * 2 * hours
	} else {
		baseAmount = calculateWeekday(hours, hourlyRate)
	}

	// Meal allowance is given per overtime session
	mealAmount := mealAllowance

	return OvertimeResult{
		BaseAmount: math.Round(baseAmount*100) / 100,
		MealAmount: math.Round(mealAmount*100) / 100,
		Total:      math.Round((baseAmount+mealAmount)*100) / 100,
	}
}

// calculateWeekday uses progressive multiplier system
// 1 jam → 1.5x, 1.5 jam → 2.5x, 2 jam → 3.5x, 2.5 jam → 4.5x, 3 jam → 5.5x
func calculateWeekday(hours float64, hourlyRate float64) float64 {
	multiplier := getWeekdayMultiplier(hours)
	return hourlyRate * multiplier
}

// getWeekdayMultiplier returns the total multiplier for weekday overtime
func getWeekdayMultiplier(hours float64) float64 {
	// Progressive multiplier table
	// First hour: 1.5x
	// Each subsequent 0.5 hour: adds 1.0x to total multiplier
	brackets := []struct {
		maxHours   float64
		multiplier float64
	}{
		{1.0, 1.5},
		{1.5, 2.5},
		{2.0, 3.5},
		{2.5, 4.5},
		{3.0, 5.5},
		{3.5, 6.5},
		{4.0, 7.5},
	}

	for _, b := range brackets {
		if hours <= b.maxHours {
			return b.multiplier
		}
	}

	// For hours beyond 4, extrapolate: 7.5 + (hours - 4) * 2
	return 7.5 + (hours-4.0)*2.0
}
