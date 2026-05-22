package com.financeku.app.data.api

import com.financeku.app.data.api.model.*
import retrofit2.Response
import retrofit2.http.*

interface ApiService {

    // Auth
    @POST("auth/login")
    suspend fun login(@Body request: LoginRequest): Response<ApiResponse<AuthResponse>>

    @POST("auth/register")
    suspend fun register(@Body request: RegisterRequest): Response<ApiResponse<AuthResponse>>

    @POST("auth/refresh")
    suspend fun refreshToken(@Body request: RefreshTokenRequest): Response<ApiResponse<TokensDto>>

    @GET("auth/me")
    suspend fun getMe(): Response<ApiResponse<UserDto>>

    // Overtime
    @GET("overtime")
    suspend fun getOvertimes(@QueryMap params: Map<String, String> = emptyMap()): Response<ApiResponse<List<OvertimeDto>>>

    @POST("overtime")
    suspend fun createOvertime(@Body request: OvertimeRequest): Response<ApiResponse<OvertimeDto>>

    @PUT("overtime/{id}")
    suspend fun updateOvertime(@Path("id") id: String, @Body request: OvertimeRequest): Response<ApiResponse<OvertimeDto>>

    @DELETE("overtime/{id}")
    suspend fun deleteOvertime(@Path("id") id: String): Response<ApiResponse<Unit>>

    @GET("overtime/calculate")
    suspend fun calculateOvertime(@QueryMap params: Map<String, String>): Response<ApiResponse<OvertimeCalculationDto>>

    @PUT("overtime/periods/disburse")
    suspend fun disburseOvertime(@Body request: DisburseRequest): Response<ApiResponse<Unit>>

    // Wallets
    @GET("wallets")
    suspend fun getWallets(): Response<ApiResponse<List<WalletDto>>>

    @POST("wallets")
    suspend fun createWallet(@Body request: WalletRequest): Response<ApiResponse<WalletDto>>

    @PUT("wallets/{id}")
    suspend fun updateWallet(@Path("id") id: String, @Body request: WalletRequest): Response<ApiResponse<WalletDto>>

    @DELETE("wallets/{id}")
    suspend fun deleteWallet(@Path("id") id: String): Response<ApiResponse<Unit>>

    @POST("wallets/transfer")
    suspend fun transferBetweenWallets(@Body request: TransferRequest): Response<ApiResponse<Unit>>

    // Categories
    @GET("categories")
    suspend fun getCategories(): Response<ApiResponse<List<CategoryDto>>>

    @POST("categories")
    suspend fun createCategory(@Body request: CategoryRequest): Response<ApiResponse<CategoryDto>>

    @PUT("categories/{id}")
    suspend fun updateCategory(@Path("id") id: String, @Body request: CategoryRequest): Response<ApiResponse<CategoryDto>>

    @DELETE("categories/{id}")
    suspend fun deleteCategory(@Path("id") id: String): Response<ApiResponse<Unit>>

    // Transactions
    @GET("transactions")
    suspend fun getTransactions(@QueryMap params: Map<String, String> = emptyMap()): Response<ApiResponse<List<TransactionDto>>>

    @POST("transactions")
    suspend fun createTransaction(@Body request: TransactionRequest): Response<ApiResponse<TransactionDto>>

    @PUT("transactions/{id}")
    suspend fun updateTransaction(@Path("id") id: String, @Body request: TransactionRequest): Response<ApiResponse<TransactionDto>>

    @DELETE("transactions/{id}")
    suspend fun deleteTransaction(@Path("id") id: String): Response<ApiResponse<Unit>>

    // Goals
    @GET("goals")
    suspend fun getGoals(): Response<ApiResponse<List<GoalDto>>>

    @POST("goals")
    suspend fun createGoal(@Body request: GoalRequest): Response<ApiResponse<GoalDto>>

    @PUT("goals/{id}")
    suspend fun updateGoal(@Path("id") id: String, @Body request: GoalRequest): Response<ApiResponse<GoalDto>>

    @DELETE("goals/{id}")
    suspend fun deleteGoal(@Path("id") id: String): Response<ApiResponse<Unit>>

    @GET("goals/{id}/progress")
    suspend fun getGoalProgress(@Path("id") id: String): Response<ApiResponse<GoalProgressDto>>

    // Incomes
    @GET("incomes")
    suspend fun getIncomes(@QueryMap params: Map<String, String> = emptyMap()): Response<ApiResponse<List<IncomeDto>>>

    @POST("incomes")
    suspend fun createIncome(@Body request: IncomeRequest): Response<ApiResponse<IncomeDto>>

    @DELETE("incomes/{id}")
    suspend fun deleteIncome(@Path("id") id: String): Response<ApiResponse<Unit>>

    // Daily Budget
    @GET("daily-budget")
    suspend fun getDailyBudget(): Response<ApiResponse<DailyBudgetDto>>

    @PUT("daily-budget")
    suspend fun updateDailyBudget(@Body request: DailyBudgetRequest): Response<ApiResponse<DailyBudgetDto>>

    @GET("daily-budget/today")
    suspend fun getTodayBudget(): Response<ApiResponse<TodayBudgetDto>>

    // Reports
    @GET("reports/dashboard")
    suspend fun getDashboardReport(): Response<ApiResponse<DashboardReportDto>>

    @GET("reports/cashflow")
    suspend fun getCashflowReport(@QueryMap params: Map<String, String> = emptyMap()): Response<ApiResponse<CashflowReportDto>>

    // Profile
    @PUT("profile")
    suspend fun updateProfile(@Body request: ProfileRequest): Response<ApiResponse<UserDto>>

    @PUT("profile/password")
    suspend fun changePassword(@Body request: ChangePasswordRequest): Response<ApiResponse<Unit>>
}
