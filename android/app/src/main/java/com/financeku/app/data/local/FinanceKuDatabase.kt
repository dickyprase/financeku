package com.financeku.app.data.local

import androidx.room.Database
import androidx.room.RoomDatabase
import com.financeku.app.data.local.dao.*
import com.financeku.app.data.local.entity.*

@Database(
    entities = [
        WalletEntity::class,
        CategoryEntity::class,
        TransactionEntity::class,
        GoalEntity::class,
        OvertimeEntity::class
    ],
    version = 1,
    exportSchema = false
)
abstract class FinanceKuDatabase : RoomDatabase() {
    abstract fun walletDao(): WalletDao
    abstract fun categoryDao(): CategoryDao
    abstract fun transactionDao(): TransactionDao
    abstract fun goalDao(): GoalDao
    abstract fun overtimeDao(): OvertimeDao
}
