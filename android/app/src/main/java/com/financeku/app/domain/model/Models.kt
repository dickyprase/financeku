package com.financeku.app.domain.model

data class User(
    val id: String,
    val name: String,
    val email: String,
    val avatarUrl: String?
)

data class Wallet(
    val id: String,
    val name: String,
    val balance: Double,
    val icon: String?,
    val color: String?,
    val isDefault: Boolean
)

data class Category(
    val id: String,
    val name: String,
    val type: String,
    val icon: String?,
    val color: String?
)

data class Transaction(
    val id: String,
    val amount: Double,
    val type: String,
    val description: String?,
    val date: String,
    val walletId: String,
    val categoryId: String?,
    val walletName: String?,
    val categoryName: String?
)

data class Goal(
    val id: String,
    val name: String,
    val targetAmount: Double,
    val currentAmount: Double,
    val deadline: String?,
    val icon: String?,
    val color: String?,
    val status: String
)

data class GoalProgress(
    val id: String,
    val name: String,
    val targetAmount: Double,
    val currentAmount: Double,
    val percentage: Double,
    val remainingAmount: Double,
    val daysRemaining: Int?
)

data class Overtime(
    val id: String,
    val date: String,
    val hours: Double,
    val rate: Double,
    val amount: Double,
    val description: String?,
    val status: String,
    val periodId: String?
)

data class OvertimeCalculation(
    val totalHours: Double,
    val totalAmount: Double,
    val periodStart: String,
    val periodEnd: String
)

data class DashboardReport(
    val totalBalance: Double,
    val monthlyIncome: Double,
    val monthlyExpense: Double,
    val overtimePending: Double,
    val todayBudget: TodayBudget?,
    val recentTransactions: List<Transaction>
)

data class TodayBudget(
    val budget: Double,
    val spent: Double,
    val remaining: Double,
    val percentage: Double
)

data class CashflowReport(
    val totalIncome: Double,
    val totalExpense: Double,
    val netCashflow: Double,
    val byCategory: List<CategorySummary>
)

data class CategorySummary(
    val categoryId: String,
    val categoryName: String,
    val amount: Double,
    val percentage: Double
)

data class Income(
    val id: String,
    val amount: Double,
    val source: String,
    val date: String,
    val description: String?,
    val walletId: String?
)

data class DailyBudget(
    val id: String,
    val amount: Double,
    val startDate: String?,
    val endDate: String?
)
