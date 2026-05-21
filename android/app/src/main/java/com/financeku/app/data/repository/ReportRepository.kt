package com.financeku.app.data.repository

import com.financeku.app.data.api.ApiService
import com.financeku.app.data.api.model.*
import com.financeku.app.data.local.datastore.TokenDataStore
import com.financeku.app.domain.model.*
import javax.inject.Inject
import javax.inject.Singleton

@Singleton
class ReportRepository @Inject constructor(
    private val apiService: ApiService
) {
    suspend fun getDashboardReport(): Resource<DashboardReport> {
        return try {
            val response = apiService.getDashboardReport()
            if (response.isSuccessful && response.body()?.success == true) {
                val data = response.body()!!.data!!
                Resource.Success(data.toDomain())
            } else {
                Resource.Error(response.body()?.message ?: "Failed to fetch dashboard")
            }
        } catch (e: Exception) {
            Resource.Error(e.message ?: "Network error")
        }
    }

    suspend fun getCashflowReport(params: Map<String, String> = emptyMap()): Resource<CashflowReport> {
        return try {
            val response = apiService.getCashflowReport(params)
            if (response.isSuccessful && response.body()?.success == true) {
                val data = response.body()!!.data!!
                Resource.Success(data.toDomain())
            } else {
                Resource.Error(response.body()?.message ?: "Failed to fetch cashflow report")
            }
        } catch (e: Exception) {
            Resource.Error(e.message ?: "Network error")
        }
    }

    private fun DashboardReportDto.toDomain() = DashboardReport(
        totalBalance = totalBalance,
        monthlyIncome = monthlyIncome,
        monthlyExpense = monthlyExpense,
        overtimePending = overtimePending,
        todayBudget = todayBudget?.let {
            TodayBudget(it.budget, it.spent, it.remaining, it.percentage)
        },
        recentTransactions = recentTransactions?.map { it.toDomain() } ?: emptyList()
    )

    private fun CashflowReportDto.toDomain() = CashflowReport(
        totalIncome = totalIncome,
        totalExpense = totalExpense,
        netCashflow = netCashflow,
        byCategory = byCategory?.map {
            CategorySummary(it.categoryId, it.categoryName, it.amount, it.percentage)
        } ?: emptyList()
    )

    private fun TransactionDto.toDomain() = Transaction(
        id = id, amount = amount, type = type, description = description,
        date = date, walletId = walletId, categoryId = categoryId,
        walletName = walletName, categoryName = categoryName
    )
}
