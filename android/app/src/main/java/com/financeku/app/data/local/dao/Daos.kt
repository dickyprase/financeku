package com.financeku.app.data.local.dao

import androidx.room.*
import com.financeku.app.data.local.entity.*
import kotlinx.coroutines.flow.Flow

@Dao
interface WalletDao {
    @Query("SELECT * FROM wallets ORDER BY name ASC")
    fun getAll(): Flow<List<WalletEntity>>

    @Insert(onConflict = OnConflictStrategy.REPLACE)
    suspend fun insertAll(wallets: List<WalletEntity>)

    @Insert(onConflict = OnConflictStrategy.REPLACE)
    suspend fun insert(wallet: WalletEntity)

    @Delete
    suspend fun delete(wallet: WalletEntity)

    @Query("DELETE FROM wallets")
    suspend fun deleteAll()

    @Query("DELETE FROM wallets WHERE id = :id")
    suspend fun deleteById(id: String)
}

@Dao
interface CategoryDao {
    @Query("SELECT * FROM categories ORDER BY name ASC")
    fun getAll(): Flow<List<CategoryEntity>>

    @Query("SELECT * FROM categories WHERE type = :type ORDER BY name ASC")
    fun getByType(type: String): Flow<List<CategoryEntity>>

    @Insert(onConflict = OnConflictStrategy.REPLACE)
    suspend fun insertAll(categories: List<CategoryEntity>)

    @Insert(onConflict = OnConflictStrategy.REPLACE)
    suspend fun insert(category: CategoryEntity)

    @Query("DELETE FROM categories")
    suspend fun deleteAll()

    @Query("DELETE FROM categories WHERE id = :id")
    suspend fun deleteById(id: String)
}

@Dao
interface TransactionDao {
    @Query("SELECT * FROM transactions ORDER BY date DESC")
    fun getAll(): Flow<List<TransactionEntity>>

    @Query("SELECT * FROM transactions WHERE walletId = :walletId ORDER BY date DESC")
    fun getByWallet(walletId: String): Flow<List<TransactionEntity>>

    @Insert(onConflict = OnConflictStrategy.REPLACE)
    suspend fun insertAll(transactions: List<TransactionEntity>)

    @Insert(onConflict = OnConflictStrategy.REPLACE)
    suspend fun insert(transaction: TransactionEntity)

    @Query("DELETE FROM transactions")
    suspend fun deleteAll()

    @Query("DELETE FROM transactions WHERE id = :id")
    suspend fun deleteById(id: String)
}

@Dao
interface GoalDao {
    @Query("SELECT * FROM goals ORDER BY name ASC")
    fun getAll(): Flow<List<GoalEntity>>

    @Insert(onConflict = OnConflictStrategy.REPLACE)
    suspend fun insertAll(goals: List<GoalEntity>)

    @Insert(onConflict = OnConflictStrategy.REPLACE)
    suspend fun insert(goal: GoalEntity)

    @Query("DELETE FROM goals")
    suspend fun deleteAll()

    @Query("DELETE FROM goals WHERE id = :id")
    suspend fun deleteById(id: String)
}

@Dao
interface OvertimeDao {
    @Query("SELECT * FROM overtimes ORDER BY date DESC")
    fun getAll(): Flow<List<OvertimeEntity>>

    @Insert(onConflict = OnConflictStrategy.REPLACE)
    suspend fun insertAll(overtimes: List<OvertimeEntity>)

    @Insert(onConflict = OnConflictStrategy.REPLACE)
    suspend fun insert(overtime: OvertimeEntity)

    @Query("DELETE FROM overtimes")
    suspend fun deleteAll()

    @Query("DELETE FROM overtimes WHERE id = :id")
    suspend fun deleteById(id: String)
}
