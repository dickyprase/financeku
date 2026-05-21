package com.financeku.app.ui.goals

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import com.financeku.app.data.repository.Resource
import com.financeku.app.domain.model.Goal
import com.financeku.app.domain.usecase.*
import com.financeku.app.data.repository.GoalRepository
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.launch
import javax.inject.Inject

data class GoalsUiState(
    val isLoading: Boolean = false,
    val goals: List<Goal> = emptyList(),
    val error: String? = null,
    val actionSuccess: String? = null
)

@HiltViewModel
class GoalsViewModel @Inject constructor(
    private val getGoalsUseCase: GetGoalsUseCase,
    private val createGoalUseCase: CreateGoalUseCase,
    private val updateGoalUseCase: UpdateGoalUseCase,
    private val goalRepository: GoalRepository
) : ViewModel() {

    private val _uiState = MutableStateFlow(GoalsUiState())
    val uiState: StateFlow<GoalsUiState> = _uiState

    init {
        loadGoals()
    }

    fun loadGoals() {
        viewModelScope.launch {
            _uiState.value = _uiState.value.copy(isLoading = true, error = null)
            when (val result = getGoalsUseCase.fetch()) {
                is Resource.Success -> {
                    _uiState.value = _uiState.value.copy(
                        isLoading = false,
                        goals = result.data
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

    fun createGoal(name: String, targetAmount: Double, currentAmount: Double?, deadline: String?, icon: String?, color: String?) {
        viewModelScope.launch {
            _uiState.value = _uiState.value.copy(isLoading = true)
            when (val result = createGoalUseCase(name, targetAmount, currentAmount, deadline, icon, color)) {
                is Resource.Success -> {
                    _uiState.value = _uiState.value.copy(actionSuccess = "Goal created")
                    loadGoals()
                }
                is Resource.Error -> {
                    _uiState.value = _uiState.value.copy(isLoading = false, error = result.message)
                }
                is Resource.Loading -> {}
            }
        }
    }

    fun updateGoal(id: String, name: String, targetAmount: Double, currentAmount: Double?, deadline: String?, icon: String?, color: String?) {
        viewModelScope.launch {
            _uiState.value = _uiState.value.copy(isLoading = true)
            when (val result = updateGoalUseCase(id, name, targetAmount, currentAmount, deadline, icon, color)) {
                is Resource.Success -> {
                    _uiState.value = _uiState.value.copy(actionSuccess = "Goal updated")
                    loadGoals()
                }
                is Resource.Error -> {
                    _uiState.value = _uiState.value.copy(isLoading = false, error = result.message)
                }
                is Resource.Loading -> {}
            }
        }
    }

    fun deleteGoal(id: String) {
        viewModelScope.launch {
            when (val result = goalRepository.deleteGoal(id)) {
                is Resource.Success -> {
                    _uiState.value = _uiState.value.copy(actionSuccess = "Goal deleted")
                    loadGoals()
                }
                is Resource.Error -> {
                    _uiState.value = _uiState.value.copy(error = result.message)
                }
                is Resource.Loading -> {}
            }
        }
    }
}
