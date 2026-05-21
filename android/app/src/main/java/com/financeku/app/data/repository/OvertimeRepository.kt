package com.financeku.app.data.repository

import com.financeku.app.data.api.ApiService
import com.financeku.app.data.api.model.*
import com.financeku.app.data.local.dao.OvertimeDao
import com.financeku.app.data.local.entity.OvertimeEntity
import com.financeku.app.domain.model.Overtime
import com.financeku.app.domain.model.OvertimeCalculation
import kotlinx.coroutines.flow.Flow
import kotlinx.coroutines.flow.map
import javax.inject.Inject
import javax.inject.Singleton

@Singleton
class OvertimeRepository @Inject constructor(
    private val apiService: ApiService,
    private val overtimeDao: OvertimeDao
) {
    fun getLocalOvertimes(): Flow<List<Overtime>> {
        return overtimeDao.getAll().map { entities ->
            entities.map { it.toDomain() }
        }
    }

    suspend fun fetchOvertimes(params: Map<String, String> = emptyMap()): Resource<List<Overtime>> {
        return try {
            val response = apiService.getOvertimes(params)
            if (response.isSuccessful && response.body()?.success == true) {
                val data = response.body()!!.data ?: emptyList()
                val entities = data.map { it.toEntity() }
                overtimeDao.deleteAll()
                overtimeDao.insertAll(entities)
                Resource.Success(data.map { it.toDomain() })
            } else {
                Resource.Error(response.body()?.message ?: "Failed to fetch overtimes")
            }
        } catch (e: Exception) {
            Resource.Error(e.message ?: "Network error")
        }
    }

    suspend fun createOvertime(date: String, hours: Double, rate: Double, description: String?): Resource<Overtime> {
        return try {
            val response = apiService.createOvertime(OvertimeRequest(date, hours, rate, description))
            if (response.isSuccessful && response.body()?.success == true) {
                val data = response.body()!!.data!!
                overtimeDao.insert(data.toEntity())
                Resource.Success(data.toDomain())
            } else {
                Resource.Error(response.body()?.message ?: "Failed to create overtime")
            }
        } catch (e: Exception) {
            Resource.Error(e.message ?: "Network error")
        }
    }

    suspend fun updateOvertime(id: String, date: String, hours: Double, rate: Double, description: String?): Resource<Overtime> {
        return try {
            val response = apiService.updateOvertime(id, OvertimeRequest(date, hours, rate, description))
            if (response.isSuccessful && response.body()?.success == true) {
                val data = response.body()!!.data!!
                overtimeDao.insert(data.toEntity())
                Resource.Success(data.toDomain())
            } else {
                Resource.Error(response.body()?.message ?: "Failed to update overtime")
            }
        } catch (e: Exception) {
            Resource.Error(e.message ?: "Network error")
        }
    }

    suspend fun deleteOvertime(id: String): Resource<Unit> {
        return try {
            val response = apiService.deleteOvertime(id)
            if (response.isSuccessful && response.body()?.success == true) {
                overtimeDao.deleteById(id)
                Resource.Success(Unit)
            } else {
                Resource.Error(response.body()?.message ?: "Failed to delete overtime")
            }
        } catch (e: Exception) {
            Resource.Error(e.message ?: "Network error")
        }
    }

    suspend fun calculateOvertime(params: Map<String, String>): Resource<OvertimeCalculation> {
        return try {
            val response = apiService.calculateOvertime(params)
            if (response.isSuccessful && response.body()?.success == true) {
                val data = response.body()!!.data!!
                Resource.Success(data.toDomain())
            } else {
                Resource.Error(response.body()?.message ?: "Failed to calculate overtime")
            }
        } catch (e: Exception) {
            Resource.Error(e.message ?: "Network error")
        }
    }

    suspend fun disburseOvertime(periodStart: String, periodEnd: String, walletId: String?): Resource<Unit> {
        return try {
            val response = apiService.disburseOvertime(DisburseRequest(periodStart, periodEnd, walletId))
            if (response.isSuccessful && response.body()?.success == true) {
                Resource.Success(Unit)
            } else {
                Resource.Error(response.body()?.message ?: "Failed to disburse overtime")
            }
        } catch (e: Exception) {
            Resource.Error(e.message ?: "Network error")
        }
    }

    private fun OvertimeDto.toDomain() = Overtime(
        id = id, date = date, hours = hours, rate = rate,
        amount = amount, description = description, status = status, periodId = periodId
    )

    private fun OvertimeDto.toEntity() = OvertimeEntity(
        id = id, date = date, hours = hours, rate = rate,
        amount = amount, description = description, status = status, periodId = periodId
    )

    private fun OvertimeEntity.toDomain() = Overtime(
        id = id, date = date, hours = hours, rate = rate,
        amount = amount, description = description, status = status, periodId = periodId
    )

    private fun OvertimeCalculationDto.toDomain() = OvertimeCalculation(
        totalHours = totalHours, totalAmount = totalAmount,
        periodStart = periodStart, periodEnd = periodEnd
    )
}
