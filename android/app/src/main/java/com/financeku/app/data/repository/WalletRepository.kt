package com.financeku.app.data.repository

import com.financeku.app.data.api.ApiService
import com.financeku.app.data.api.model.*
import com.financeku.app.data.local.dao.WalletDao
import com.financeku.app.data.local.entity.WalletEntity
import com.financeku.app.domain.model.Wallet
import kotlinx.coroutines.flow.Flow
import kotlinx.coroutines.flow.map
import javax.inject.Inject
import javax.inject.Singleton

@Singleton
class WalletRepository @Inject constructor(
    private val apiService: ApiService,
    private val walletDao: WalletDao
) {
    fun getLocalWallets(): Flow<List<Wallet>> {
        return walletDao.getAll().map { entities ->
            entities.map { it.toDomain() }
        }
    }

    suspend fun fetchWallets(): Resource<List<Wallet>> {
        return try {
            val response = apiService.getWallets()
            if (response.isSuccessful && response.body()?.success == true) {
                val data = response.body()!!.data ?: emptyList()
                walletDao.deleteAll()
                walletDao.insertAll(data.map { it.toEntity() })
                Resource.Success(data.map { it.toDomain() })
            } else {
                Resource.Error(response.body()?.message ?: "Failed to fetch wallets")
            }
        } catch (e: Exception) {
            Resource.Error(e.message ?: "Network error")
        }
    }

    suspend fun createWallet(name: String, balance: Double?, icon: String?, color: String?, isDefault: Boolean?): Resource<Wallet> {
        return try {
            val response = apiService.createWallet(WalletRequest(name, balance, icon, color, isDefault))
            if (response.isSuccessful && response.body()?.success == true) {
                val data = response.body()!!.data!!
                walletDao.insert(data.toEntity())
                Resource.Success(data.toDomain())
            } else {
                Resource.Error(response.body()?.message ?: "Failed to create wallet")
            }
        } catch (e: Exception) {
            Resource.Error(e.message ?: "Network error")
        }
    }

    suspend fun updateWallet(id: String, name: String, balance: Double?, icon: String?, color: String?, isDefault: Boolean?): Resource<Wallet> {
        return try {
            val response = apiService.updateWallet(id, WalletRequest(name, balance, icon, color, isDefault))
            if (response.isSuccessful && response.body()?.success == true) {
                val data = response.body()!!.data!!
                walletDao.insert(data.toEntity())
                Resource.Success(data.toDomain())
            } else {
                Resource.Error(response.body()?.message ?: "Failed to update wallet")
            }
        } catch (e: Exception) {
            Resource.Error(e.message ?: "Network error")
        }
    }

    suspend fun deleteWallet(id: String): Resource<Unit> {
        return try {
            val response = apiService.deleteWallet(id)
            if (response.isSuccessful && response.body()?.success == true) {
                walletDao.deleteById(id)
                Resource.Success(Unit)
            } else {
                Resource.Error(response.body()?.message ?: "Failed to delete wallet")
            }
        } catch (e: Exception) {
            Resource.Error(e.message ?: "Network error")
        }
    }

    suspend fun transfer(fromWalletId: String, toWalletId: String, amount: Double, description: String?): Resource<Unit> {
        return try {
            val response = apiService.transferBetweenWallets(TransferRequest(fromWalletId, toWalletId, amount, description))
            if (response.isSuccessful && response.body()?.success == true) {
                fetchWallets()
                Resource.Success(Unit)
            } else {
                Resource.Error(response.body()?.message ?: "Failed to transfer")
            }
        } catch (e: Exception) {
            Resource.Error(e.message ?: "Network error")
        }
    }

    private fun WalletDto.toDomain() = Wallet(
        id = id, name = name, balance = balance,
        icon = icon, color = color, isDefault = isDefault
    )

    private fun WalletDto.toEntity() = WalletEntity(
        id = id, name = name, balance = balance,
        icon = icon, color = color, isDefault = isDefault
    )

    private fun WalletEntity.toDomain() = Wallet(
        id = id, name = name, balance = balance,
        icon = icon, color = color, isDefault = isDefault
    )
}
