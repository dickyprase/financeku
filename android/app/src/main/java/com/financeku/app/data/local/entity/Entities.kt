package com.financeku.app.data.local.entity

import androidx.room.Entity
import androidx.room.PrimaryKey

@Entity(tableName = "wallets")
data class WalletEntity(
    @PrimaryKey val id: String,
    val name: String,
    val balance: Double,
    val icon: String?,
    val color: String?,
    val isDefault: Boolean,
    val updatedAt: Long = System.currentTimeMillis()
)

@Entity(tableName = "categories")
data class CategoryEntity(
    @PrimaryKey val id: String,
    val name: String,
    val type: String,
    val icon: String?,
    val color: String?,
    val updatedAt: Long = System.currentTimeMillis()
)

@Entity(tableName = "transactions")
data class TransactionEntity(
    @PrimaryKey val id: String,
    val amount: Double,
    val type: String,
    val description: String?,
    val date: String,
    val walletId: String,
    val categoryId: String?,
    val walletName: String?,
    val categoryName: String?,
    val updatedAt: Long = System.currentTimeMillis()
)

@Entity(tableName = "goals")
data class GoalEntity(
    @PrimaryKey val id: String,
    val name: String,
    val targetAmount: Double,
    val currentAmount: Double,
    val deadline: String?,
    val icon: String?,
    val color: String?,
    val status: String,
    val updatedAt: Long = System.currentTimeMillis()
)

@Entity(tableName = "overtimes")
data class OvertimeEntity(
    @PrimaryKey val id: String,
    val date: String,
    val hours: Double,
    val rate: Double,
    val amount: Double,
    val description: String?,
    val status: String,
    val periodId: String?,
    val updatedAt: Long = System.currentTimeMillis()
)
