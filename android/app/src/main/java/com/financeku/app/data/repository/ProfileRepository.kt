package com.financeku.app.data.repository

import com.financeku.app.data.api.ApiService
import com.financeku.app.data.api.model.*
import com.financeku.app.data.local.datastore.TokenDataStore
import com.financeku.app.domain.model.User
import javax.inject.Inject
import javax.inject.Singleton

@Singleton
class ProfileRepository @Inject constructor(
    private val apiService: ApiService,
    private val tokenDataStore: TokenDataStore
) {
    suspend fun updateProfile(name: String?, email: String?, avatarUrl: String?): Resource<User> {
        return try {
            val response = apiService.updateProfile(ProfileRequest(name, email, avatarUrl))
            if (response.isSuccessful && response.body()?.success == true) {
                val data = response.body()!!.data!!
                tokenDataStore.saveUserInfo(data.id, data.name, data.email)
                Resource.Success(data.toDomain())
            } else {
                Resource.Error(response.body()?.message ?: "Failed to update profile")
            }
        } catch (e: Exception) {
            Resource.Error(e.message ?: "Network error")
        }
    }

    suspend fun changePassword(currentPassword: String, newPassword: String, confirmation: String): Resource<Unit> {
        return try {
            val response = apiService.changePassword(
                ChangePasswordRequest(currentPassword, newPassword, confirmation)
            )
            if (response.isSuccessful && response.body()?.success == true) {
                Resource.Success(Unit)
            } else {
                Resource.Error(response.body()?.message ?: "Failed to change password")
            }
        } catch (e: Exception) {
            Resource.Error(e.message ?: "Network error")
        }
    }

    private fun UserDto.toDomain() = User(
        id = id, name = name, email = email, avatarUrl = avatarUrl
    )
}
