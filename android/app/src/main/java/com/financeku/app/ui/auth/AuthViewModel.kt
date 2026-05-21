package com.financeku.app.ui.auth

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import com.financeku.app.data.repository.Resource
import com.financeku.app.domain.model.User
import com.financeku.app.domain.usecase.LoginUseCase
import com.financeku.app.domain.usecase.RegisterUseCase
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.launch
import javax.inject.Inject

sealed class AuthUiState {
    data object Idle : AuthUiState()
    data object Loading : AuthUiState()
    data class Success(val user: User) : AuthUiState()
    data class Error(val message: String) : AuthUiState()
}

@HiltViewModel
class AuthViewModel @Inject constructor(
    private val loginUseCase: LoginUseCase,
    private val registerUseCase: RegisterUseCase
) : ViewModel() {

    private val _loginState = MutableStateFlow<AuthUiState>(AuthUiState.Idle)
    val loginState: StateFlow<AuthUiState> = _loginState

    private val _registerState = MutableStateFlow<AuthUiState>(AuthUiState.Idle)
    val registerState: StateFlow<AuthUiState> = _registerState

    fun login(email: String, password: String) {
        viewModelScope.launch {
            _loginState.value = AuthUiState.Loading
            when (val result = loginUseCase(email, password)) {
                is Resource.Success -> _loginState.value = AuthUiState.Success(result.data)
                is Resource.Error -> _loginState.value = AuthUiState.Error(result.message)
                is Resource.Loading -> {}
            }
        }
    }

    fun register(name: String, email: String, password: String, confirmation: String) {
        viewModelScope.launch {
            _registerState.value = AuthUiState.Loading
            when (val result = registerUseCase(name, email, password, confirmation)) {
                is Resource.Success -> _registerState.value = AuthUiState.Success(result.data)
                is Resource.Error -> _registerState.value = AuthUiState.Error(result.message)
                is Resource.Loading -> {}
            }
        }
    }
}
