package com.financeku.app.di

import android.content.Context
import androidx.room.Room
import com.financeku.app.data.local.FinanceKuDatabase
import com.financeku.app.data.local.dao.*
import dagger.Module
import dagger.Provides
import dagger.hilt.InstallIn
import dagger.hilt.android.qualifiers.ApplicationContext
import dagger.hilt.components.SingletonComponent
import javax.inject.Singleton

@Module
@InstallIn(SingletonComponent::class)
object DatabaseModule {

    @Provides
    @Singleton
    fun provideDatabase(@ApplicationContext context: Context): FinanceKuDatabase {
        return Room.databaseBuilder(
            context,
            FinanceKuDatabase::class.java,
            "financeku_db"
        ).fallbackToDestructiveMigration().build()
    }

    @Provides
    fun provideWalletDao(db: FinanceKuDatabase): WalletDao = db.walletDao()

    @Provides
    fun provideCategoryDao(db: FinanceKuDatabase): CategoryDao = db.categoryDao()

    @Provides
    fun provideTransactionDao(db: FinanceKuDatabase): TransactionDao = db.transactionDao()

    @Provides
    fun provideGoalDao(db: FinanceKuDatabase): GoalDao = db.goalDao()

    @Provides
    fun provideOvertimeDao(db: FinanceKuDatabase): OvertimeDao = db.overtimeDao()
}
