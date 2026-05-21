package com.financeku.app.ui.cashflow

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import com.financeku.app.data.repository.Resource
import com.financeku.app.domain.model.*
import com.financeku.app.domain.usecase.*
import com.financeku.app.data.repository.WalletRepository
import com.financeku.app.data.repository.TransactionRepository
import com.financeku.app.data.repository.CategoryRepository
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.launch
import javax.inject.Inject

data class CashflowUiState(
    val isLoading: Boolean = false,
    val wallets: List<Wallet> = emptyList(),
    val transactions: List<Transaction> = emptyList(),
    val categories: List<Category> = emptyList(),
    val error: String? = null,
    val actionSuccess: String? = null
)

data class TransferUiState(
    val isLoading: Boolean = false,
    val success: Boolean = false,
    val error: String? = null
)

data class TransactionFormUiState(
    val isLoading: Boolean = false,
    val success: Boolean = false,
    val error: String? = null
)

@HiltViewModel
class CashflowViewModel @Inject constructor(
    private val getWalletsUseCase: GetWalletsUseCase,
    private val getTransactionsUseCase: GetTransactionsUseCase,
    private val getCategoriesUseCase: GetCategoriesUseCase,
    private val createTransactionUseCase: CreateTransactionUseCase,
    private val transferUseCase: TransferUseCase,
    private val walletRepository: WalletRepository,
    private val transactionRepository: TransactionRepository,
    private val categoryRepository: CategoryRepository
) : ViewModel() {

    private val _uiState = MutableStateFlow(CashflowUiState())
    val uiState: StateFlow<CashflowUiState> = _uiState

    private val _transferState = MutableStateFlow(TransferUiState())
    val transferState: StateFlow<TransferUiState> = _transferState

    private val _transactionFormState = MutableStateFlow(TransactionFormUiState())
    val transactionFormState: StateFlow<TransactionFormUiState> = _transactionFormState

    init {
        loadAll()
    }

    fun loadAll() {
        viewModelScope.launch {
            _uiState.value = _uiState.value.copy(isLoading = true, error = null)
            val walletsResult = getWalletsUseCase.fetch()
            val transactionsResult = getTransactionsUseCase.fetch()
            val categoriesResult = getCategoriesUseCase.fetch()

            _uiState.value = _uiState.value.copy(
                isLoading = false,
                wallets = (walletsResult as? Resource.Success)?.data ?: _uiState.value.wallets,
                transactions = (transactionsResult as? Resource.Success)?.data ?: _uiState.value.transactions,
                categories = (categoriesResult as? Resource.Success)?.data ?: _uiState.value.categories,
                error = when {
                    walletsResult is Resource.Error -> walletsResult.message
                    transactionsResult is Resource.Error -> transactionsResult.message
                    categoriesResult is Resource.Error -> categoriesResult.message
                    else -> null
                }
            )
        }
    }

    fun createTransaction(
        amount: Double, type: String, description: String?,
        date: String, walletId: String, categoryId: String?
    ) {
        viewModelScope.launch {
            _transactionFormState.value = TransactionFormUiState(isLoading = true)
            when (val result = createTransactionUseCase(amount, type, description, date, walletId, categoryId)) {
                is Resource.Success -> {
                    _transactionFormState.value = TransactionFormUiState(success = true)
                    loadAll()
                }
                is Resource.Error -> {
                    _transactionFormState.value = TransactionFormUiState(error = result.message)
                }
                is Resource.Loading -> {}
            }
        }
    }

    fun updateTransaction(
        id: String, amount: Double, type: String, description: String?,
        date: String, walletId: String, categoryId: String?
    ) {
        viewModelScope.launch {
            _transactionFormState.value = TransactionFormUiState(isLoading = true)
            when (val result = transactionRepository.updateTransaction(id, amount, type, description, date, walletId, categoryId)) {
                is Resource.Success -> {
                    _transactionFormState.value = TransactionFormUiState(success = true)
                    loadAll()
                }
                is Resource.Error -> {
                    _transactionFormState.value = TransactionFormUiState(error = result.message)
                }
                is Resource.Loading -> {}
            }
        }
    }

    fun deleteTransaction(id: String) {
        viewModelScope.launch {
            when (val result = transactionRepository.deleteTransaction(id)) {
                is Resource.Success -> {
                    _uiState.value = _uiState.value.copy(actionSuccess = "Transaction deleted")
                    loadAll()
                }
                is Resource.Error -> {
                    _uiState.value = _uiState.value.copy(error = result.message)
                }
                is Resource.Loading -> {}
            }
        }
    }

    fun transfer(fromId: String, toId: String, amount: Double, description: String?) {
        viewModelScope.launch {
            _transferState.value = TransferUiState(isLoading = true)
            when (val result = transferUseCase(fromId, toId, amount, description)) {
                is Resource.Success -> {
                    _transferState.value = TransferUiState(success = true)
                    loadAll()
                }
                is Resource.Error -> {
                    _transferState.value = TransferUiState(error = result.message)
                }
                is Resource.Loading -> {}
            }
        }
    }

    fun createWallet(name: String, balance: Double?, icon: String?, color: String?, isDefault: Boolean?) {
        viewModelScope.launch {
            _uiState.value = _uiState.value.copy(isLoading = true)
            when (val result = walletRepository.createWallet(name, balance, icon, color, isDefault)) {
                is Resource.Success -> {
                    _uiState.value = _uiState.value.copy(actionSuccess = "Wallet created")
                    loadAll()
                }
                is Resource.Error -> {
                    _uiState.value = _uiState.value.copy(isLoading = false, error = result.message)
                }
                is Resource.Loading -> {}
            }
        }
    }

    fun deleteWallet(id: String) {
        viewModelScope.launch {
            when (val result = walletRepository.deleteWallet(id)) {
                is Resource.Success -> {
                    _uiState.value = _uiState.value.copy(actionSuccess = "Wallet deleted")
                    loadAll()
                }
                is Resource.Error -> {
                    _uiState.value = _uiState.value.copy(error = result.message)
                }
                is Resource.Loading -> {}
            }
        }
    }

    fun resetTransferState() {
        _transferState.value = TransferUiState()
    }

    fun resetTransactionFormState() {
        _transactionFormState.value = TransactionFormUiState()
    }

    fun clearError() {
        _uiState.value = _uiState.value.copy(error = null)
    }
}
