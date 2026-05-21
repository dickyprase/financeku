export interface OvertimeResult {
  baseAmount: number;
  mealAmount: number;
  totalAmount: number;
}

/**
 * Calculate overtime amount based on hours, salary, meal allowance, and holiday status
 * Weekday: progressive multiplier × (salary/173)
 * Holiday: salary/173 × 2 × hours
 */
export function calculateOvertime(
  hours: number,
  salary: number,
  mealAllowance: number,
  isHoliday: boolean
): OvertimeResult {
  const hourlyRate = salary / 173;
  let baseAmount: number;

  if (isHoliday) {
    baseAmount = hourlyRate * 2 * hours;
  } else {
    baseAmount = calculateWeekday(hours, hourlyRate);
  }

  const mealAmount = mealAllowance;

  return {
    baseAmount: Math.round(baseAmount * 100) / 100,
    mealAmount: Math.round(mealAmount * 100) / 100,
    totalAmount: Math.round((baseAmount + mealAmount) * 100) / 100,
  };
}

/**
 * Weekday progressive multiplier system
 * 1 jam → 1.5x, 1.5 jam → 2.5x, 2 jam → 3.5x, 2.5 jam → 4.5x, 3 jam → 5.5x
 */
function calculateWeekday(hours: number, hourlyRate: number): number {
  const multiplier = getWeekdayMultiplier(hours);
  return hourlyRate * multiplier;
}

function getWeekdayMultiplier(hours: number): number {
  const brackets = [
    { maxHours: 1.0, multiplier: 1.5 },
    { maxHours: 1.5, multiplier: 2.5 },
    { maxHours: 2.0, multiplier: 3.5 },
    { maxHours: 2.5, multiplier: 4.5 },
    { maxHours: 3.0, multiplier: 5.5 },
    { maxHours: 3.5, multiplier: 6.5 },
    { maxHours: 4.0, multiplier: 7.5 },
  ];

  for (const b of brackets) {
    if (hours <= b.maxHours) {
      return b.multiplier;
    }
  }

  // Beyond 4 hours: extrapolate
  return 7.5 + (hours - 4.0) * 2.0;
}

/**
 * Calculate overtime period (Thursday to Wednesday) for a given date
 */
export function calculatePeriod(dateStr: string): { periodStart: string; periodEnd: string } {
  const date = new Date(dateStr);
  const weekday = date.getDay(); // 0=Sun, 1=Mon, ..., 4=Thu, 6=Sat

  // Days since last Thursday
  let daysToThursday: number;
  switch (weekday) {
    case 4: daysToThursday = 0; break; // Thursday
    case 5: daysToThursday = 1; break; // Friday
    case 6: daysToThursday = 2; break; // Saturday
    case 0: daysToThursday = 3; break; // Sunday
    case 1: daysToThursday = 4; break; // Monday
    case 2: daysToThursday = 5; break; // Tuesday
    case 3: daysToThursday = 6; break; // Wednesday
    default: daysToThursday = 0;
  }

  const periodStart = new Date(date);
  periodStart.setDate(date.getDate() - daysToThursday);

  const periodEnd = new Date(periodStart);
  periodEnd.setDate(periodStart.getDate() + 6);

  return {
    periodStart: periodStart.toISOString().split('T')[0],
    periodEnd: periodEnd.toISOString().split('T')[0],
  };
}
