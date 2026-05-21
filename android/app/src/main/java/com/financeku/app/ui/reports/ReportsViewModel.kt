package com.financeku.app.ui.reports

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import com.financeku.app.data.repository.Resource
import com.financeku.app.domain.model.CashflowReport
import com.financeku.app.domain.usecase.GetCashflowReportUseCase
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.launch
import javax.inject.Inject

data class ReportsUiState(
    val isLoading: Boolean = false,
    val report: CashflowReport? = null,
    val error: String? = null
)

@HiltViewModel
class ReportsViewModel @Inject constructor(
    private val getCashflowReportUseCase: GetCashflowReportUseCase
) : ViewModel() {

    private val _uiState = MutableStateFlow(ReportsUiState())
    val uiState: StateFlow<ReportsUiState> = _uiState

    init {
        loadReport()
    }

    fun loadReport(params: Map<String, String> = emptyMap()) {
        viewModelScope.launch {
            _uiState.value = _uiState.value.copy(isLoading = true, error = null)
            when (val result = getCashflowReportUseCase(params)) {
                is Resource.Success -> {
                    _uiState.value = _uiState.value.copy(
                        isLoading = false,
                        report = result.data
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
}
