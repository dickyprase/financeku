package com.financeku.app.ui.cashflow

import androidx.compose.foundation.background
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.rememberScrollState
import androidx.compose.foundation.text.KeyboardOptions
import androidx.compose.foundation.verticalScroll
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.*
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Brush
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.text.input.KeyboardType
import androidx.compose.ui.unit.dp
import androidx.hilt.navigation.compose.hiltViewModel
import com.financeku.app.ui.components.*
import com.financeku.app.ui.theme.*

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun TransactionFormScreen(
    transactionId: String?,
    onNavigateBack: () -> Unit,
    viewModel: CashflowViewModel = hiltViewModel()
) {
    val formState by viewModel.transactionFormState.collectAsState()
    val uiState by viewModel.uiState.collectAsState()

    var amount by remember { mutableStateOf("") }
    var type by remember { mutableStateOf("expense") }
    var description by remember { mutableStateOf("") }
    var date by remember { mutableStateOf("") }
    var selectedWalletId by remember { mutableStateOf("") }
    var selectedCategoryId by remember { mutableStateOf("") }
    var walletExpanded by remember { mutableStateOf(false) }
    var categoryExpanded by remember { mutableStateOf(false) }

    val isEditing = transactionId != null

    LaunchedEffect(formState.success) {
        if (formState.success) {
            viewModel.resetTransactionFormState()
            onNavigateBack()
        }
    }

    val isDark = LocalDarkMode.current.value
    val bgGradient = if (isDark) {
        Brush.verticalGradient(listOf(BackgroundDark, GradientDarkMiddle))
    } else {
        Brush.verticalGradient(listOf(BackgroundLight, Color(0xFFE8F5E9)))
    }

    Column(
        modifier = Modifier
            .fillMaxSize()
            .background(bgGradient)
    ) {
        GlassTopBar(
            title = if (isEditing) "Edit Transaction" else "Add Transaction",
            navigationIcon = {
                IconButton(onClick = onNavigateBack) {
                    Icon(Icons.Filled.ArrowBack, contentDescription = "Back")
                }
            }
        )

        Column(
            modifier = Modifier
                .fillMaxSize()
                .verticalScroll(rememberScrollState())
                .padding(16.dp),
            verticalArrangement = Arrangement.spacedBy(16.dp)
        ) {
            GlassCard(modifier = Modifier.fillMaxWidth()) {
                // Type selector
                Row(
                    modifier = Modifier.fillMaxWidth(),
                    horizontalArrangement = Arrangement.spacedBy(8.dp)
                ) {
                    FilterChip(
                        selected = type == "expense",
                        onClick = { type = "expense" },
                        label = { Text("Expense") },
                        leadingIcon = if (type == "expense") {
                            { Icon(Icons.Filled.Check, contentDescription = null, modifier = Modifier.size(16.dp)) }
                        } else null,
                        colors = FilterChipDefaults.filterChipColors(
                            selectedContainerColor = ExpenseRed.copy(alpha = 0.2f)
                        ),
                        modifier = Modifier.weight(1f)
                    )
                    FilterChip(
                        selected = type == "income",
                        onClick = { type = "income" },
                        label = { Text("Income") },
                        leadingIcon = if (type == "income") {
                            { Icon(Icons.Filled.Check, contentDescription = null, modifier = Modifier.size(16.dp)) }
                        } else null,
                        colors = FilterChipDefaults.filterChipColors(
                            selectedContainerColor = IncomeGreen.copy(alpha = 0.2f)
                        ),
                        modifier = Modifier.weight(1f)
                    )
                }

                Spacer(modifier = Modifier.height(16.dp))

                GlassTextField(
                    value = amount,
                    onValueChange = { amount = it },
                    label = "Amount",
                    placeholder = "0",
                    leadingIcon = { Icon(Icons.Filled.AttachMoney, contentDescription = null) },
                    keyboardOptions = KeyboardOptions(keyboardType = KeyboardType.Decimal),
                    modifier = Modifier.fillMaxWidth()
                )

                Spacer(modifier = Modifier.height(16.dp))

                GlassTextField(
                    value = date,
                    onValueChange = { date = it },
                    label = "Date",
                    placeholder = "YYYY-MM-DD",
                    leadingIcon = { Icon(Icons.Filled.CalendarToday, contentDescription = null) },
                    modifier = Modifier.fillMaxWidth()
                )

                Spacer(modifier = Modifier.height(16.dp))

                // Wallet dropdown
                ExposedDropdownMenuBox(
                    expanded = walletExpanded,
                    onExpandedChange = { walletExpanded = it }
                ) {
                    GlassTextField(
                        value = uiState.wallets.find { it.id == selectedWalletId }?.name ?: "",
                        onValueChange = {},
                        label = "Wallet",
                        placeholder = "Select wallet",
                        leadingIcon = { Icon(Icons.Filled.AccountBalanceWallet, contentDescription = null) },
                        modifier = Modifier
                            .fillMaxWidth()
                            .menuAnchor()
                    )
                    ExposedDropdownMenu(
                        expanded = walletExpanded,
                        onDismissRequest = { walletExpanded = false }
                    ) {
                        uiState.wallets.forEach { wallet ->
                            DropdownMenuItem(
                                text = { Text(wallet.name) },
                                onClick = {
                                    selectedWalletId = wallet.id
                                    walletExpanded = false
                                }
                            )
                        }
                    }
                }

                Spacer(modifier = Modifier.height(16.dp))

                // Category dropdown
                ExposedDropdownMenuBox(
                    expanded = categoryExpanded,
                    onExpandedChange = { categoryExpanded = it }
                ) {
                    GlassTextField(
                        value = uiState.categories.find { it.id == selectedCategoryId }?.name ?: "",
                        onValueChange = {},
                        label = "Category",
                        placeholder = "Select category",
                        leadingIcon = { Icon(Icons.Filled.Category, contentDescription = null) },
                        modifier = Modifier
                            .fillMaxWidth()
                            .menuAnchor()
                    )
                    ExposedDropdownMenu(
                        expanded = categoryExpanded,
                        onDismissRequest = { categoryExpanded = false }
                    ) {
                        uiState.categories
                            .filter { it.type == type }
                            .forEach { category ->
                                DropdownMenuItem(
                                    text = { Text(category.name) },
                                    onClick = {
                                        selectedCategoryId = category.id
                                        categoryExpanded = false
                                    }
                                )
                            }
                    }
                }

                Spacer(modifier = Modifier.height(16.dp))

                GlassTextField(
                    value = description,
                    onValueChange = { description = it },
                    label = "Description (Optional)",
                    placeholder = "Add a note",
                    leadingIcon = { Icon(Icons.Filled.Description, contentDescription = null) },
                    singleLine = false,
                    maxLines = 3,
                    modifier = Modifier.fillMaxWidth()
                )

                Spacer(modifier = Modifier.height(24.dp))

                if (formState.error != null) {
                    Text(
                        text = formState.error!!,
                        style = MaterialTheme.typography.bodySmall,
                        color = MaterialTheme.colorScheme.error,
                        modifier = Modifier.padding(bottom = 8.dp)
                    )
                }

                GlassButton(
                    onClick = {
                        val amountVal = amount.toDoubleOrNull() ?: return@GlassButton
                        if (isEditing) {
                            viewModel.updateTransaction(
                                transactionId!!, amountVal, type,
                                description.ifBlank { null }, date,
                                selectedWalletId, selectedCategoryId.ifBlank { null }
                            )
                        } else {
                            viewModel.createTransaction(
                                amountVal, type,
                                description.ifBlank { null }, date,
                                selectedWalletId, selectedCategoryId.ifBlank { null }
                            )
                        }
                    },
                    enabled = !formState.isLoading && amount.isNotBlank() && date.isNotBlank() && selectedWalletId.isNotBlank(),
                    modifier = Modifier.fillMaxWidth()
                ) {
                    if (formState.isLoading) {
                        CircularProgressIndicator(modifier = Modifier.size(20.dp), color = Color.White, strokeWidth = 2.dp)
                        Spacer(modifier = Modifier.width(8.dp))
                    }
                    Text(if (isEditing) "Update" else "Save", color = Color.White)
                }
            }
        }
    }
}
