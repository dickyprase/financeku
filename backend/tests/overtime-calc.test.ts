import { calculateOvertime, calculatePeriod } from '../src/utils/overtime-calc';

describe('calculateOvertime', () => {
  const salary = 5000000;
  const meal = 30000;
  const hourlyRate = salary / 173;

  describe('weekday', () => {
    it('1 hour should use 1.5x multiplier', () => {
      const result = calculateOvertime(1, salary, meal, false);
      expect(result.baseAmount).toBeCloseTo(hourlyRate * 1.5, 2);
      expect(result.mealAmount).toBe(meal);
      expect(result.totalAmount).toBeCloseTo(hourlyRate * 1.5 + meal, 2);
    });

    it('1.5 hours should use 2.5x multiplier', () => {
      const result = calculateOvertime(1.5, salary, meal, false);
      expect(result.baseAmount).toBeCloseTo(hourlyRate * 2.5, 2);
    });

    it('2 hours should use 3.5x multiplier', () => {
      const result = calculateOvertime(2, salary, meal, false);
      expect(result.baseAmount).toBeCloseTo(hourlyRate * 3.5, 2);
    });

    it('2.5 hours should use 4.5x multiplier', () => {
      const result = calculateOvertime(2.5, salary, meal, false);
      expect(result.baseAmount).toBeCloseTo(hourlyRate * 4.5, 2);
    });

    it('3 hours should use 5.5x multiplier', () => {
      const result = calculateOvertime(3, salary, meal, false);
      expect(result.baseAmount).toBeCloseTo(hourlyRate * 5.5, 2);
    });
  });

  describe('holiday', () => {
    it('1 hour holiday should use 2x hourly rate', () => {
      const result = calculateOvertime(1, salary, meal, true);
      expect(result.baseAmount).toBeCloseTo(hourlyRate * 2 * 1, 2);
    });

    it('2 hours holiday', () => {
      const result = calculateOvertime(2, salary, meal, true);
      expect(result.baseAmount).toBeCloseTo(hourlyRate * 2 * 2, 2);
    });

    it('3 hours holiday', () => {
      const result = calculateOvertime(3, salary, meal, true);
      expect(result.baseAmount).toBeCloseTo(hourlyRate * 2 * 3, 2);
    });

    it('4 hours holiday', () => {
      const result = calculateOvertime(4, salary, meal, true);
      expect(result.baseAmount).toBeCloseTo(hourlyRate * 2 * 4, 2);
    });
  });

  describe('edge cases', () => {
    it('zero salary should return 0 base amount', () => {
      const result = calculateOvertime(2, 0, meal, false);
      expect(result.baseAmount).toBe(0);
      expect(result.mealAmount).toBe(meal);
    });
  });
});

describe('calculatePeriod', () => {
  it('Thursday should be period start', () => {
    // 2024-01-04 is a Thursday
    const { periodStart, periodEnd } = calculatePeriod('2024-01-04');
    expect(periodStart).toBe('2024-01-04');
    expect(periodEnd).toBe('2024-01-10');
  });

  it('Wednesday should be period end', () => {
    // 2024-01-10 is a Wednesday
    const { periodStart, periodEnd } = calculatePeriod('2024-01-10');
    expect(periodStart).toBe('2024-01-04');
    expect(periodEnd).toBe('2024-01-10');
  });

  it('Monday should map to correct period', () => {
    // 2024-01-08 is a Monday
    const { periodStart, periodEnd } = calculatePeriod('2024-01-08');
    expect(periodStart).toBe('2024-01-04');
    expect(periodEnd).toBe('2024-01-10');
  });
});
