package com.financeku.app.data.api

import com.financeku.app.data.local.datastore.TokenDataStore
import kotlinx.coroutines.flow.firstOrNull
import kotlinx.coroutines.runBlocking
import okhttp3.Interceptor
import okhttp3.Response
import javax.inject.Inject
import javax.inject.Singleton

@Singleton
class AuthInterceptor @Inject constructor(
    private val tokenDataStore: TokenDataStore
) : Interceptor {

    override fun intercept(chain: Interceptor.Chain): Response {
        val originalRequest = chain.request()

        // Skip auth header for login/register/refresh endpoints
        val path = originalRequest.url.encodedPath
        if (path.contains("auth/login") || path.contains("auth/register") || path.contains("auth/refresh")) {
            return chain.proceed(originalRequest)
        }

        val token = runBlocking {
            tokenDataStore.accessToken.firstOrNull()
        }

        val request = if (!token.isNullOrEmpty()) {
            originalRequest.newBuilder()
                .header("Authorization", "Bearer $token")
                .build()
        } else {
            originalRequest
        }

        return chain.proceed(request)
    }
}
