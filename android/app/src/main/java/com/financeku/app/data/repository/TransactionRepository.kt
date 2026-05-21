package com.financeku.app.data.repository

import com.financeku.app.data.api.ApiService
import com.financeku.app.data.api.model.*
import com.financeku.app.data.local.dao.TransactionDao
import com.financeku.app.data.local.entity.TransactionEntity
import com.financeku.app.domain.model.Transaction
import kotlinx.coroutines.flow.Flow
import kotlinx.coroutines.flow.map
import javax.inject.Inject
import javax.inject.Singleton

@Singleton
class TransactionRepository @Inject constructor(
    private val apiService: ApiService,
    private val transactionDao: TransactionDao
) {
    fun getLocalTransactions(): Flow<List<Transaction>> {
        return transactionDao.getAll().map { entities ->
            entities.map { it.toDomain() }
        }
    }

    suspend fun fetchTransactions(params: Map<String, String> = emptyMap()): Resource<List<Transaction>> {
        return try {
            val response = apiService.getTransactions(params)
            if (response.isSuccessful && response.body()?.success == true) {
                val data = response.body()!!.data ?: emptyList()
                transactionDao.deleteAll()
                transactionDao.insertAll(data.map { it.toEntity() })
                Resource.Success(data.map { it.toDomain() })
            } else {
                Resource.Error(response.body()?.message ?: "Failed to fetch transactions")
            }
        } catch (e: Exception) {
            Resource.Error(e.message ?: "Network error")
        }
    }

    suspend fun createTransaction(
        amount: Double, type: String, description: String?,
        date: String, walletId: String, categoryId: String?
    ): Resource<Transaction> {
        return try {
            val response = apiService.createTransaction(
                TransactionRequest(amount, type, description, date, walletId, categoryId)
            )
            if (response.isSuccessful && response.body()?.success == true) {
                val data = response.body()!!.data!!
                transactionDao.insert(data.toEntity())
                Resource.Success(data.toDomain())
            } else {
                Resource.Error(response.body()?.message ?: "Failed to create transaction")
            }
        } catch (e: Exception) {
            Resource.Error(e.message ?: "Network error")
        }
    }

    suspend fun updateTransaction(
        id: String, amount: Double, type: String, description: String?,
        date: String, walletId: String, categoryId: String?
    ): Resource<Transaction> {
        return try {
            val response = apiService.updateTransaction(
                id, TransactionRequest(amount, type, description, date, walletId, categoryId)
            )
            if (response.isSuccessful && response.body()?.success == true) {
                val data = response.body()!!.data!!
                transactionDao.insert(data.toEntity())
                Resource.Success(data.toDomain())
            } else {
                Resource.Error(response.body()?.message ?: "Failed to update transaction")
            }
        } catch (e: Exception) {
            Resource.Error(e.message ?: "Network error")
        }
    }

    suspend fun deleteTransaction(id: String): Resource<Unit> {
        return try {
            val response = apiService.deleteTransaction(id)
            if (response.isSuccessful && response.body()?.success == true) {
                transactionDao.deleteById(id)
                Resource.Success(Unit)
            } else {
                Resource.Error(response.body()?.message ?: "Failed to delete transaction")
            }
        } catch (e: Exception) {
            Resource.Error(e.message ?: "Network error")
        }
    }

    private fun TransactionDto.toDomain() = Transaction(
        id = id, amount = amount, type = type, description = description,
        date = date, walletId = walletId, categoryId = categoryId,
        walletName = walletName, categoryName = categoryName
    )

    private fun TransactionDto.toEntity() = TransactionEntity(
        id = id, amount = amount, type = type, description = description,
        date = date, walletId = walletId, categoryId = categoryId,
        walletName = walletName, categoryName = categoryName
    )

    private fun TransactionEntity.toDomain() = Transaction(
        id = id, amount = amount, type = type, description = description,
        date = date, walletId = walletId, categoryId = categoryId,
        walletName = walletName, categoryName = categoryName
    )
}
