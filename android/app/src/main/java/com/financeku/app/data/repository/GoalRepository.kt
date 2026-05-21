package com.financeku.app.data.repository

import com.financeku.app.data.api.ApiService
import com.financeku.app.data.api.model.*
import com.financeku.app.data.local.dao.GoalDao
import com.financeku.app.data.local.entity.GoalEntity
import com.financeku.app.domain.model.Goal
import com.financeku.app.domain.model.GoalProgress
import kotlinx.coroutines.flow.Flow
import kotlinx.coroutines.flow.map
import javax.inject.Inject
import javax.inject.Singleton

@Singleton
class GoalRepository @Inject constructor(
    private val apiService: ApiService,
    private val goalDao: GoalDao
) {
    fun getLocalGoals(): Flow<List<Goal>> {
        return goalDao.getAll().map { entities ->
            entities.map { it.toDomain() }
        }
    }

    suspend fun fetchGoals(): Resource<List<Goal>> {
        return try {
            val response = apiService.getGoals()
            if (response.isSuccessful && response.body()?.success == true) {
                val data = response.body()!!.data ?: emptyList()
                goalDao.deleteAll()
                goalDao.insertAll(data.map { it.toEntity() })
                Resource.Success(data.map { it.toDomain() })
            } else {
                Resource.Error(response.body()?.message ?: "Failed to fetch goals")
            }
        } catch (e: Exception) {
            Resource.Error(e.message ?: "Network error")
        }
    }

    suspend fun createGoal(name: String, targetAmount: Double, currentAmount: Double?, deadline: String?, icon: String?, color: String?): Resource<Goal> {
        return try {
            val response = apiService.createGoal(GoalRequest(name, targetAmount, currentAmount, deadline, icon, color))
            if (response.isSuccessful && response.body()?.success == true) {
                val data = response.body()!!.data!!
                goalDao.insert(data.toEntity())
                Resource.Success(data.toDomain())
            } else {
                Resource.Error(response.body()?.message ?: "Failed to create goal")
            }
        } catch (e: Exception) {
            Resource.Error(e.message ?: "Network error")
        }
    }

    suspend fun updateGoal(id: String, name: String, targetAmount: Double, currentAmount: Double?, deadline: String?, icon: String?, color: String?): Resource<Goal> {
        return try {
            val response = apiService.updateGoal(id, GoalRequest(name, targetAmount, currentAmount, deadline, icon, color))
            if (response.isSuccessful && response.body()?.success == true) {
                val data = response.body()!!.data!!
                goalDao.insert(data.toEntity())
                Resource.Success(data.toDomain())
            } else {
                Resource.Error(response.body()?.message ?: "Failed to update goal")
            }
        } catch (e: Exception) {
            Resource.Error(e.message ?: "Network error")
        }
    }

    suspend fun deleteGoal(id: String): Resource<Unit> {
        return try {
            val response = apiService.deleteGoal(id)
            if (response.isSuccessful && response.body()?.success == true) {
                goalDao.deleteById(id)
                Resource.Success(Unit)
            } else {
                Resource.Error(response.body()?.message ?: "Failed to delete goal")
            }
        } catch (e: Exception) {
            Resource.Error(e.message ?: "Network error")
        }
    }

    suspend fun getGoalProgress(id: String): Resource<GoalProgress> {
        return try {
            val response = apiService.getGoalProgress(id)
            if (response.isSuccessful && response.body()?.success == true) {
                val data = response.body()!!.data!!
                Resource.Success(data.toDomain())
            } else {
                Resource.Error(response.body()?.message ?: "Failed to get goal progress")
            }
        } catch (e: Exception) {
            Resource.Error(e.message ?: "Network error")
        }
    }

    private fun GoalDto.toDomain() = Goal(
        id = id, name = name, targetAmount = targetAmount,
        currentAmount = currentAmount, deadline = deadline,
        icon = icon, color = color, status = status
    )

    private fun GoalDto.toEntity() = GoalEntity(
        id = id, name = name, targetAmount = targetAmount,
        currentAmount = currentAmount, deadline = deadline,
        icon = icon, color = color, status = status
    )

    private fun GoalEntity.toDomain() = Goal(
        id = id, name = name, targetAmount = targetAmount,
        currentAmount = currentAmount, deadline = deadline,
        icon = icon, color = color, status = status
    )

    private fun GoalProgressDto.toDomain() = GoalProgress(
        id = id, name = name, targetAmount = targetAmount,
        currentAmount = currentAmount, percentage = percentage,
        remainingAmount = remainingAmount, daysRemaining = daysRemaining
    )
}
