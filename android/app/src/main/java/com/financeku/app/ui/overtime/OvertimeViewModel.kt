package com.financeku.app.ui.overtime

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import com.financeku.app.data.repository.Resource
import com.financeku.app.domain.model.Overtime
import com.financeku.app.domain.model.OvertimeCalculation
import com.financeku.app.domain.usecase.*
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.launch
import javax.inject.Inject

data class OvertimeListUiState(
    val isLoading: Boolean = false,
    val overtimes: List<Overtime> = emptyList(),
    val calculation: OvertimeCalculation? = null,
    val error: String? = null,
    val actionSuccess: String? = null
)

data class OvertimeFormUiState(
    val isLoading: Boolean = false,
    val success: Boolean = false,
    val error: String? = null
)

@HiltViewModel
class OvertimeViewModel @Inject constructor(
    private val getOvertimesUseCase: GetOvertimesUseCase,
    private val createOvertimeUseCase: CreateOvertimeUseCase,
    private val calculateOvertimeUseCase: CalculateOvertimeUseCase,
    private val disburseOvertimeUseCase: DisburseOvertimeUseCase,
    private val overtimeRepository: com.financeku.app.data.repository.OvertimeRepository
) : ViewModel() {

    private val _listState = MutableStateFlow(OvertimeListUiState())
    val listState: StateFlow<OvertimeListUiState> = _listState

    private val _formState = MutableStateFlow(OvertimeFormUiState())
    val formState: StateFlow<OvertimeFormUiState> = _formState

    init {
        loadOvertimes()
    }

    fun loadOvertimes() {
        viewModelScope.launch {
            _listState.value = _listState.value.copy(isLoading = true, error = null)
            when (val result = getOvertimesUseCase.fetch()) {
                is Resource.Success -> {
                    _listState.value = _listState.value.copy(
                        isLoading = false,
                        overtimes = result.data
                    )
                }
                is Resource.Error -> {
                    _listState.value = _listState.value.copy(
                        isLoading = false,
                        error = result.message
                    )
                }
                is Resource.Loading -> {}
            }
        }
    }

    fun createOvertime(date: String, hours: Double, rate: Double, description: String?) {
        viewModelScope.launch {
            _formState.value = OvertimeFormUiState(isLoading = true)
            when (val result = createOvertimeUseCase(date, hours, rate, description)) {
                is Resource.Success -> {
                    _formState.value = OvertimeFormUiState(success = true)
                    loadOvertimes()
                }
                is Resource.Error -> {
                    _formState.value = OvertimeFormUiState(error = result.message)
                }
                is Resource.Loading -> {}
            }
        }
    }

    fun updateOvertime(id: String, date: String, hours: Double, rate: Double, description: String?) {
        viewModelScope.launch {
            _formState.value = OvertimeFormUiState(isLoading = true)
            when (val result = overtimeRepository.updateOvertime(id, date, hours, rate, description)) {
                is Resource.Success -> {
                    _formState.value = OvertimeFormUiState(success = true)
                    loadOvertimes()
                }
                is Resource.Error -> {
                    _formState.value = OvertimeFormUiState(error = result.message)
                }
                is Resource.Loading -> {}
            }
        }
    }

    fun deleteOvertime(id: String) {
        viewModelScope.launch {
            when (val result = overtimeRepository.deleteOvertime(id)) {
                is Resource.Success -> {
                    _listState.value = _listState.value.copy(
                        actionSuccess = "Overtime deleted"
                    )
                    loadOvertimes()
                }
                is Resource.Error -> {
                    _listState.value = _listState.value.copy(error = result.message)
                }
                is Resource.Loading -> {}
            }
        }
    }

    fun calculateOvertime(periodStart: String, periodEnd: String) {
        viewModelScope.launch {
            val params = mapOf("period_start" to periodStart, "period_end" to periodEnd)
            when (val result = calculateOvertimeUseCase(params)) {
                is Resource.Success -> {
                    _listState.value = _listState.value.copy(calculation = result.data)
                }
                is Resource.Error -> {
                    _listState.value = _listState.value.copy(error = result.message)
                }
                is Resource.Loading -> {}
            }
        }
    }

    fun disburseOvertime(periodStart: String, periodEnd: String, walletId: String?) {
        viewModelScope.launch {
            _listState.value = _listState.value.copy(isLoading = true)
            when (val result = disburseOvertimeUseCase(periodStart, periodEnd, walletId)) {
                is Resource.Success -> {
                    _listState.value = _listState.value.copy(
                        isLoading = false,
                        actionSuccess = "Overtime disbursed successfully"
                    )
                    loadOvertimes()
                }
                is Resource.Error -> {
                    _listState.value = _listState.value.copy(
                        isLoading = false,
                        error = result.message
                    )
                }
                is Resource.Loading -> {}
            }
        }
    }

    fun clearActionSuccess() {
        _listState.value = _listState.value.copy(actionSuccess = null)
    }

    fun resetFormState() {
        _formState.value = OvertimeFormUiState()
    }
}
