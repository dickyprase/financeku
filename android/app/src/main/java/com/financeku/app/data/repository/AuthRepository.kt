package com.financeku.app.data.repository

import com.financeku.app.data.api.ApiService
import com.financeku.app.data.api.model.*
import com.financeku.app.data.local.datastore.TokenDataStore
import com.financeku.app.domain.model.*
import kotlinx.coroutines.flow.Flow
import kotlinx.coroutines.flow.firstOrNull
import javax.inject.Inject
import javax.inject.Singleton

sealed class Resource<out T> {
    data class Success<T>(val data: T) : Resource<T>()
    data class Error(val message: String, val code: Int? = null) : Resource<Nothing>()
    data object Loading : Resource<Nothing>()
}

@Singleton
class AuthRepository @Inject constructor(
    private val apiService: ApiService,
    private val tokenDataStore: TokenDataStore
) {
    val isLoggedIn: Flow<String?> = tokenDataStore.accessToken

    suspend fun login(email: String, password: String): Resource<User> {
        return try {
            val response = apiService.login(LoginRequest(email, password))
            if (response.isSuccessful && response.body()?.success == true) {
                val data = response.body()!!.data!!
                tokenDataStore.saveTokens(data.accessToken, data.refreshToken)
                data.user?.let {
                    tokenDataStore.saveUserInfo(it.id, it.name, it.email)
                    Resource.Success(it.toDomain())
                } ?: Resource.Success(User("", "", "", null))
            } else {
                Resource.Error(response.body()?.message ?: "Login failed")
            }
        } catch (e: Exception) {
            Resource.Error(e.message ?: "Network error")
        }
    }

    suspend fun register(name: String, email: String, password: String, passwordConfirmation: String): Resource<User> {
        return try {
            val response = apiService.register(RegisterRequest(name, email, password, passwordConfirmation))
            if (response.isSuccessful && response.body()?.success == true) {
                val data = response.body()!!.data!!
                tokenDataStore.saveTokens(data.accessToken, data.refreshToken)
                data.user?.let {
                    tokenDataStore.saveUserInfo(it.id, it.name, it.email)
                    Resource.Success(it.toDomain())
                } ?: Resource.Success(User("", name, email, null))
            } else {
                Resource.Error(response.body()?.message ?: "Registration failed")
            }
        } catch (e: Exception) {
            Resource.Error(e.message ?: "Network error")
        }
    }

    suspend fun getMe(): Resource<User> {
        return try {
            val response = apiService.getMe()
            if (response.isSuccessful && response.body()?.success == true) {
                val user = response.body()!!.data!!.toDomain()
                tokenDataStore.saveUserInfo(user.id, user.name, user.email)
                Resource.Success(user)
            } else {
                Resource.Error(response.body()?.message ?: "Failed to get user")
            }
        } catch (e: Exception) {
            Resource.Error(e.message ?: "Network error")
        }
    }

    suspend fun logout() {
        tokenDataStore.clearTokens()
    }

    private fun UserDto.toDomain() = User(
        id = id,
        name = name,
        email = email,
        avatarUrl = avatarUrl
    )
}
