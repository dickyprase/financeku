package com.financeku.app.data.api.model

import com.google.gson.annotations.SerializedName

data class ApiResponse<T>(
    val success: Boolean,
    val message: String?,
    val data: T?,
    val meta: MetaResponse?
)

data class MetaResponse(
    val page: Int?,
    @SerializedName("per_page") val perPage: Int?,
    val total: Int?,
    @SerializedName("total_pages") val totalPages: Int?
)

// Auth
data class LoginRequest(
    val email: String,
    val password: String
)

data class RegisterRequest(
    val name: String,
    val email: String,
    val password: String,
    @SerializedName("password_confirmation") val passwordConfirmation: String
)

data class RefreshTokenRequest(
    @SerializedName("refresh_token") val refreshToken: String
)

data class AuthResponse(
    val tokens: TokensDto?,
    val user: UserDto?
) {
    val accessToken: String get() = tokens?.accessToken ?: ""
    val refreshToken: String get() = tokens?.refreshToken ?: ""
}

data class TokensDto(
    @SerializedName("access_token") val accessToken: String,
    @SerializedName("refresh_token") val refreshToken: String,
    @SerializedName("expires_in") val expiresIn: Int?
)

data class UserDto(
    val id: String,
    val name: String,
    val email: String,
    @SerializedName("avatar_url") val avatarUrl: String?,
    @SerializedName("created_at") val createdAt: String?
)

// Overtime
data class OvertimeDto(
    val id: String,
    val date: String,
    val hours: Double,
    val rate: Double,
    val amount: Double,
    val description: String?,
    val status: String,
    @SerializedName("period_id") val periodId: String?,
    @SerializedName("created_at") val createdAt: String?
)

data class OvertimeRequest(
    val date: String,
    val hours: Double,
    val rate: Double,
    val description: String?
)

data class OvertimeCalculationDto(
    @SerializedName("total_hours") val totalHours: Double,
    @SerializedName("total_amount") val totalAmount: Double,
    @SerializedName("period_start") val periodStart: String,
    @SerializedName("period_end") val periodEnd: String
)

data class DisburseRequest(
    @SerializedName("period_start") val periodStart: String,
    @SerializedName("period_end") val periodEnd: String,
    @SerializedName("wallet_id") val walletId: String?
)

// Wallets
data class WalletDto(
    val id: String,
    val name: String,
    val balance: Double,
    val icon: String?,
    val color: String?,
    @SerializedName("is_default") val isDefault: Boolean,
    @SerializedName("created_at") val createdAt: String?
)

data class WalletRequest(
    val name: String,
    val balance: Double?,
    val icon: String?,
    val color: String?,
    @SerializedName("is_default") val isDefault: Boolean?
)

data class TransferRequest(
    @SerializedName("from_wallet_id") val fromWalletId: String,
    @SerializedName("to_wallet_id") val toWalletId: String,
    val amount: Double,
    val description: String?
)

// Categories
data class CategoryDto(
    val id: String,
    val name: String,
    val type: String,
    val icon: String?,
    val color: String?,
    @SerializedName("created_at") val createdAt: String?
)

data class CategoryRequest(
    val name: String,
    val type: String,
    val icon: String?,
    val color: String?
)

// Transactions
data class TransactionDto(
    val id: String,
    val amount: Double,
    val type: String,
    val description: String?,
    val date: String,
    @SerializedName("wallet_id") val walletId: String,
    @SerializedName("category_id") val categoryId: String?,
    @SerializedName("wallet_name") val walletName: String?,
    @SerializedName("category_name") val categoryName: String?,
    @SerializedName("created_at") val createdAt: String?
)

data class TransactionRequest(
    val amount: Double,
    val type: String,
    val description: String?,
    val date: String,
    @SerializedName("wallet_id") val walletId: String,
    @SerializedName("category_id") val categoryId: String?
)

// Goals
data class GoalDto(
    val id: String,
    val name: String,
    @SerializedName("target_amount") val targetAmount: Double,
    @SerializedName("current_amount") val currentAmount: Double,
    val deadline: String?,
    val icon: String?,
    val color: String?,
    val status: String,
    @SerializedName("created_at") val createdAt: String?
)

data class GoalRequest(
    val name: String,
    @SerializedName("target_amount") val targetAmount: Double,
    @SerializedName("current_amount") val currentAmount: Double?,
    val deadline: String?,
    val icon: String?,
    val color: String?
)

data class GoalProgressDto(
    val id: String,
    val name: String,
    @SerializedName("target_amount") val targetAmount: Double,
    @SerializedName("current_amount") val currentAmount: Double,
    val percentage: Double,
    @SerializedName("remaining_amount") val remainingAmount: Double,
    @SerializedName("days_remaining") val daysRemaining: Int?
)

// Incomes
data class IncomeDto(
    val id: String,
    val amount: Double,
    val source: String,
    val date: String,
    val description: String?,
    @SerializedName("wallet_id") val walletId: String?,
    @SerializedName("created_at") val createdAt: String?
)

data class IncomeRequest(
    val amount: Double,
    val source: String,
    val date: String,
    val description: String?,
    @SerializedName("wallet_id") val walletId: String?
)

// Daily Budget
data class DailyBudgetDto(
    val id: String,
    val amount: Double,
    @SerializedName("start_date") val startDate: String?,
    @SerializedName("end_date") val endDate: String?
)

data class DailyBudgetRequest(
    val amount: Double
)

data class TodayBudgetDto(
    val budget: Double,
    val spent: Double,
    val remaining: Double,
    val percentage: Double
)

// Reports
data class DashboardReportDto(
    @SerializedName("total_balance") val totalBalance: Double,
    @SerializedName("monthly_income") val monthlyIncome: Double,
    @SerializedName("monthly_expense") val monthlyExpense: Double,
    @SerializedName("overtime_pending") val overtimePending: Double,
    @SerializedName("today_budget") val todayBudget: TodayBudgetDto?,
    @SerializedName("recent_transactions") val recentTransactions: List<TransactionDto>?
)

data class CashflowReportDto(
    @SerializedName("total_income") val totalIncome: Double,
    @SerializedName("total_expense") val totalExpense: Double,
    @SerializedName("net_cashflow") val netCashflow: Double,
    @SerializedName("by_category") val byCategory: List<CategorySummaryDto>?
)

data class CategorySummaryDto(
    @SerializedName("category_id") val categoryId: String,
    @SerializedName("category_name") val categoryName: String,
    val amount: Double,
    val percentage: Double
)

// Profile
data class ProfileRequest(
    val name: String?,
    val email: String?,
    @SerializedName("avatar_url") val avatarUrl: String?
)

data class ChangePasswordRequest(
    @SerializedName("current_password") val currentPassword: String,
    @SerializedName("new_password") val newPassword: String,
    @SerializedName("new_password_confirmation") val newPasswordConfirmation: String
)
