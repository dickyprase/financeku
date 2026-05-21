package com.financeku.app.domain.usecase

import com.financeku.app.data.repository.*
import com.financeku.app.domain.model.*
import kotlinx.coroutines.flow.Flow
import javax.inject.Inject

class LoginUseCase @Inject constructor(
    private val authRepository: AuthRepository
) {
    suspend operator fun invoke(email: String, password: String): Resource<User> {
        return authRepository.login(email, password)
    }
}

class RegisterUseCase @Inject constructor(
    private val authRepository: AuthRepository
) {
    suspend operator fun invoke(name: String, email: String, password: String, confirmation: String): Resource<User> {
        return authRepository.register(name, email, password, confirmation)
    }
}

class GetMeUseCase @Inject constructor(
    private val authRepository: AuthRepository
) {
    suspend operator fun invoke(): Resource<User> {
        return authRepository.getMe()
    }
}

class LogoutUseCase @Inject constructor(
    private val authRepository: AuthRepository
) {
    suspend operator fun invoke() {
        authRepository.logout()
    }
}

class GetDashboardUseCase @Inject constructor(
    private val reportRepository: ReportRepository
) {
    suspend operator fun invoke(): Resource<DashboardReport> {
        return reportRepository.getDashboardReport()
    }
}

class GetCashflowReportUseCase @Inject constructor(
    private val reportRepository: ReportRepository
) {
    suspend operator fun invoke(params: Map<String, String> = emptyMap()): Resource<CashflowReport> {
        return reportRepository.getCashflowReport(params)
    }
}

class GetWalletsUseCase @Inject constructor(
    private val walletRepository: WalletRepository
) {
    fun getLocal(): Flow<List<Wallet>> = walletRepository.getLocalWallets()
    suspend fun fetch(): Resource<List<Wallet>> = walletRepository.fetchWallets()
}

class CreateWalletUseCase @Inject constructor(
    private val walletRepository: WalletRepository
) {
    suspend operator fun invoke(name: String, balance: Double?, icon: String?, color: String?, isDefault: Boolean?): Resource<Wallet> {
        return walletRepository.createWallet(name, balance, icon, color, isDefault)
    }
}

class TransferUseCase @Inject constructor(
    private val walletRepository: WalletRepository
) {
    suspend operator fun invoke(fromId: String, toId: String, amount: Double, description: String?): Resource<Unit> {
        return walletRepository.transfer(fromId, toId, amount, description)
    }
}

class GetTransactionsUseCase @Inject constructor(
    private val transactionRepository: TransactionRepository
) {
    fun getLocal(): Flow<List<Transaction>> = transactionRepository.getLocalTransactions()
    suspend fun fetch(params: Map<String, String> = emptyMap()): Resource<List<Transaction>> =
        transactionRepository.fetchTransactions(params)
}

class CreateTransactionUseCase @Inject constructor(
    private val transactionRepository: TransactionRepository
) {
    suspend operator fun invoke(
        amount: Double, type: String, description: String?,
        date: String, walletId: String, categoryId: String?
    ): Resource<Transaction> {
        return transactionRepository.createTransaction(amount, type, description, date, walletId, categoryId)
    }
}

class GetCategoriesUseCase @Inject constructor(
    private val categoryRepository: CategoryRepository
) {
    fun getLocal(): Flow<List<Category>> = categoryRepository.getLocalCategories()
    fun getByType(type: String): Flow<List<Category>> = categoryRepository.getLocalCategoriesByType(type)
    suspend fun fetch(): Resource<List<Category>> = categoryRepository.fetchCategories()
}

class GetOvertimesUseCase @Inject constructor(
    private val overtimeRepository: OvertimeRepository
) {
    fun getLocal(): Flow<List<Overtime>> = overtimeRepository.getLocalOvertimes()
    suspend fun fetch(params: Map<String, String> = emptyMap()): Resource<List<Overtime>> =
        overtimeRepository.fetchOvertimes(params)
}

class CreateOvertimeUseCase @Inject constructor(
    private val overtimeRepository: OvertimeRepository
) {
    suspend operator fun invoke(date: String, hours: Double, rate: Double, description: String?): Resource<Overtime> {
        return overtimeRepository.createOvertime(date, hours, rate, description)
    }
}

class CalculateOvertimeUseCase @Inject constructor(
    private val overtimeRepository: OvertimeRepository
) {
    suspend operator fun invoke(params: Map<String, String>): Resource<OvertimeCalculation> {
        return overtimeRepository.calculateOvertime(params)
    }
}

class DisburseOvertimeUseCase @Inject constructor(
    private val overtimeRepository: OvertimeRepository
) {
    suspend operator fun invoke(periodStart: String, periodEnd: String, walletId: String?): Resource<Unit> {
        return overtimeRepository.disburseOvertime(periodStart, periodEnd, walletId)
    }
}

class GetGoalsUseCase @Inject constructor(
    private val goalRepository: GoalRepository
) {
    fun getLocal(): Flow<List<Goal>> = goalRepository.getLocalGoals()
    suspend fun fetch(): Resource<List<Goal>> = goalRepository.fetchGoals()
}

class CreateGoalUseCase @Inject constructor(
    private val goalRepository: GoalRepository
) {
    suspend operator fun invoke(
        name: String, targetAmount: Double, currentAmount: Double?,
        deadline: String?, icon: String?, color: String?
    ): Resource<Goal> {
        return goalRepository.createGoal(name, targetAmount, currentAmount, deadline, icon, color)
    }
}

class UpdateGoalUseCase @Inject constructor(
    private val goalRepository: GoalRepository
) {
    suspend operator fun invoke(
        id: String, name: String, targetAmount: Double, currentAmount: Double?,
        deadline: String?, icon: String?, color: String?
    ): Resource<Goal> {
        return goalRepository.updateGoal(id, name, targetAmount, currentAmount, deadline, icon, color)
    }
}

class GetGoalProgressUseCase @Inject constructor(
    private val goalRepository: GoalRepository
) {
    suspend operator fun invoke(id: String): Resource<GoalProgress> {
        return goalRepository.getGoalProgress(id)
    }
}

class UpdateProfileUseCase @Inject constructor(
    private val profileRepository: ProfileRepository
) {
    suspend operator fun invoke(name: String?, email: String?, avatarUrl: String?): Resource<User> {
        return profileRepository.updateProfile(name, email, avatarUrl)
    }
}

class ChangePasswordUseCase @Inject constructor(
    private val profileRepository: ProfileRepository
) {
    suspend operator fun invoke(currentPassword: String, newPassword: String, confirmation: String): Resource<Unit> {
        return profileRepository.changePassword(currentPassword, newPassword, confirmation)
    }
}
