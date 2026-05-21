package com.financeku.app.di

import com.financeku.app.data.api.ApiService
import com.financeku.app.data.local.dao.*
import com.financeku.app.data.local.datastore.TokenDataStore
import com.financeku.app.data.repository.*
import dagger.Module
import dagger.Provides
import dagger.hilt.InstallIn
import dagger.hilt.components.SingletonComponent
import javax.inject.Singleton

@Module
@InstallIn(SingletonComponent::class)
object RepositoryModule {

    @Provides
    @Singleton
    fun provideAuthRepository(
        apiService: ApiService,
        tokenDataStore: TokenDataStore
    ): AuthRepository {
        return AuthRepository(apiService, tokenDataStore)
    }

    @Provides
    @Singleton
    fun provideWalletRepository(
        apiService: ApiService,
        walletDao: WalletDao
    ): WalletRepository {
        return WalletRepository(apiService, walletDao)
    }

    @Provides
    @Singleton
    fun provideCategoryRepository(
        apiService: ApiService,
        categoryDao: CategoryDao
    ): CategoryRepository {
        return CategoryRepository(apiService, categoryDao)
    }

    @Provides
    @Singleton
    fun provideTransactionRepository(
        apiService: ApiService,
        transactionDao: TransactionDao
    ): TransactionRepository {
        return TransactionRepository(apiService, transactionDao)
    }

    @Provides
    @Singleton
    fun provideGoalRepository(
        apiService: ApiService,
        goalDao: GoalDao
    ): GoalRepository {
        return GoalRepository(apiService, goalDao)
    }

    @Provides
    @Singleton
    fun provideOvertimeRepository(
        apiService: ApiService,
        overtimeDao: OvertimeDao
    ): OvertimeRepository {
        return OvertimeRepository(apiService, overtimeDao)
    }

    @Provides
    @Singleton
    fun provideProfileRepository(
        apiService: ApiService,
        tokenDataStore: TokenDataStore
    ): ProfileRepository {
        return ProfileRepository(apiService, tokenDataStore)
    }

    @Provides
    @Singleton
    fun provideReportRepository(
        apiService: ApiService
    ): ReportRepository {
        return ReportRepository(apiService)
    }
}
