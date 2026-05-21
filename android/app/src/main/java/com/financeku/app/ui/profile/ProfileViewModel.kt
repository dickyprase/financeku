package com.financeku.app.ui.profile

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import com.financeku.app.data.local.datastore.TokenDataStore
import com.financeku.app.data.repository.Resource
import com.financeku.app.domain.model.User
import com.financeku.app.domain.usecase.*
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.firstOrNull
import kotlinx.coroutines.launch
import javax.inject.Inject

data class ProfileUiState(
    val isLoading: Boolean = false,
    val user: User? = null,
    val isDarkMode: Boolean = false,
    val error: String? = null,
    val actionSuccess: String? = null
)

data class ChangePasswordUiState(
    val isLoading: Boolean = false,
    val success: Boolean = false,
    val error: String? = null
)

@HiltViewModel
class ProfileViewModel @Inject constructor(
    private val getMeUseCase: GetMeUseCase,
    private val updateProfileUseCase: UpdateProfileUseCase,
    private val changePasswordUseCase: ChangePasswordUseCase,
    private val logoutUseCase: LogoutUseCase,
    private val tokenDataStore: TokenDataStore
) : ViewModel() {

    private val _uiState = MutableStateFlow(ProfileUiState())
    val uiState: StateFlow<ProfileUiState> = _uiState

    private val _changePasswordState = MutableStateFlow(ChangePasswordUiState())
    val changePasswordState: StateFlow<ChangePasswordUiState> = _changePasswordState

    init {
        loadProfile()
        observeDarkMode()
    }

    private fun observeDarkMode() {
        viewModelScope.launch {
            tokenDataStore.isDarkMode.collect { isDark ->
                _uiState.value = _uiState.value.copy(isDarkMode = isDark)
            }
        }
    }

    fun loadProfile() {
        viewModelScope.launch {
            _uiState.value = _uiState.value.copy(isLoading = true, error = null)

            // Try to load from local first
            val name = tokenDataStore.userName.firstOrNull()
            val email = tokenDataStore.userEmail.firstOrNull()
            val id = tokenDataStore.userId.firstOrNull()

            if (name != null && email != null && id != null) {
                _uiState.value = _uiState.value.copy(
                    user = User(id, name, email, null),
                    isLoading = false
                )
            }

            // Then refresh from API
            when (val result = getMeUseCase()) {
                is Resource.Success -> {
                    _uiState.value = _uiState.value.copy(
                        isLoading = false,
                        user = result.data
                    )
                }
                is Resource.Error -> {
                    _uiState.value = _uiState.value.copy(
                        isLoading = false,
                        error = if (_uiState.value.user == null) result.message else null
                    )
                }
                is Resource.Loading -> {}
            }
        }
    }

    fun updateProfile(name: String?, email: String?) {
        viewModelScope.launch {
            _uiState.value = _uiState.value.copy(isLoading = true)
            when (val result = updateProfileUseCase(name, email, null)) {
                is Resource.Success -> {
                    _uiState.value = _uiState.value.copy(
                        isLoading = false,
                        user = result.data,
                        actionSuccess = "Profile updated"
                    )
                }
                is Resource.Error -> {
                    _uiState.value = _uiState.value.copy(
                        isLoading = false,
                        error = result.message
                    )
                }
                is Resource.Loading -> {}
            }
        }
    }

    fun changePassword(currentPassword: String, newPassword: String, confirmation: String) {
        viewModelScope.launch {
            _changePasswordState.value = ChangePasswordUiState(isLoading = true)
            when (val result = changePasswordUseCase(currentPassword, newPassword, confirmation)) {
                is Resource.Success -> {
                    _changePasswordState.value = ChangePasswordUiState(success = true)
                }
                is Resource.Error -> {
                    _changePasswordState.value = ChangePasswordUiState(error = result.message)
                }
                is Resource.Loading -> {}
            }
        }
    }

    fun toggleDarkMode(isDark: Boolean) {
        viewModelScope.launch {
            tokenDataStore.setDarkMode(isDark)
        }
    }

    fun logout() {
        viewModelScope.launch {
            logoutUseCase()
        }
    }

    fun resetChangePasswordState() {
        _changePasswordState.value = ChangePasswordUiState()
    }
}
