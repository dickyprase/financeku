package com.financeku.app.data.api

import com.financeku.app.data.api.model.ApiResponse
import com.financeku.app.data.api.model.AuthResponse
import com.financeku.app.data.api.model.RefreshTokenRequest
import com.financeku.app.data.local.datastore.TokenDataStore
import kotlinx.coroutines.flow.firstOrNull
import kotlinx.coroutines.runBlocking
import okhttp3.Authenticator
import okhttp3.Request
import okhttp3.Response
import okhttp3.Route
import retrofit2.Retrofit
import retrofit2.converter.gson.GsonConverterFactory
import javax.inject.Inject
import javax.inject.Singleton

@Singleton
class TokenAuthenticator @Inject constructor(
    private val tokenDataStore: TokenDataStore
) : Authenticator {

    override fun authenticate(route: Route?, response: Response): Request? {
        // If we already tried to refresh, give up
        if (response.request.header("X-Retry") != null) {
            runBlocking { tokenDataStore.clearTokens() }
            return null
        }

        val refreshToken = runBlocking {
            tokenDataStore.refreshToken.firstOrNull()
        } ?: return null

        // Try to refresh the token
        val newToken = runBlocking {
            try {
                val baseUrl = response.request.url.scheme + "://" + response.request.url.host +
                        ":" + response.request.url.port + "/api/v1/"
                val retrofit = Retrofit.Builder()
                    .baseUrl(baseUrl)
                    .addConverterFactory(GsonConverterFactory.create())
                    .build()

                val api = retrofit.create(ApiService::class.java)
                val refreshResponse = api.refreshToken(RefreshTokenRequest(refreshToken))

                if (refreshResponse.isSuccessful) {
                    val body = refreshResponse.body()
                    if (body?.success == true && body.data != null) {
                        tokenDataStore.saveTokens(
                            accessToken = body.data.accessToken,
                            refreshToken = body.data.refreshToken
                        )
                        body.data.accessToken
                    } else {
                        tokenDataStore.clearTokens()
                        null
                    }
                } else {
                    tokenDataStore.clearTokens()
                    null
                }
            } catch (e: Exception) {
                tokenDataStore.clearTokens()
                null
            }
        }

        return if (newToken != null) {
            response.request.newBuilder()
                .header("Authorization", "Bearer $newToken")
                .header("X-Retry", "true")
                .build()
        } else {
            null
        }
    }
}
