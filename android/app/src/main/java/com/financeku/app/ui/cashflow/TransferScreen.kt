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
fun TransferScreen(
    onNavigateBack: () -> Unit,
    viewModel: CashflowViewModel = hiltViewModel()
) {
    val uiState by viewModel.uiState.collectAsState()
    val transferState by viewModel.transferState.collectAsState()

    var fromWalletId by remember { mutableStateOf("") }
    var toWalletId by remember { mutableStateOf("") }
    var amount by remember { mutableStateOf("") }
    var description by remember { mutableStateOf("") }
    var fromExpanded by remember { mutableStateOf(false) }
    var toExpanded by remember { mutableStateOf(false) }

    LaunchedEffect(transferState.success) {
        if (transferState.success) {
            viewModel.resetTransferState()
            onNavigateBack()
        }
    }

    val isDark = LocalDarkMode.current.value
    val bgGradient = if (isDark) {
        Brush.verticalGradient(listOf(BackgroundDark, GradientDarkMiddle))
    } else {
        Brush.verticalGradient(listOf(BackgroundLight, Color(0xFFF3E5F5)))
    }

    Column(
        modifier = Modifier
            .fillMaxSize()
            .background(bgGradient)
    ) {
        GlassTopBar(
            title = "Transfer",
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
                Text(
                    text = "Transfer Between Wallets",
                    style = MaterialTheme.typography.titleMedium,
                    color = MaterialTheme.colorScheme.onSurface
                )

                Spacer(modifier = Modifier.height(20.dp))

                // From Wallet
                ExposedDropdownMenuBox(
                    expanded = fromExpanded,
                    onExpandedChange = { fromExpanded = it }
                ) {
                    GlassTextField(
                        value = uiState.wallets.find { it.id == fromWalletId }?.name ?: "",
                        onValueChange = {},
                        label = "From Wallet",
                        placeholder = "Select source wallet",
                        leadingIcon = { Icon(Icons.Filled.ArrowUpward, contentDescription = null, tint = ExpenseRed) },
                        modifier = Modifier
                            .fillMaxWidth()
                            .menuAnchor()
                    )
                    ExposedDropdownMenu(
                        expanded = fromExpanded,
                        onDismissRequest = { fromExpanded = false }
                    ) {
                        uiState.wallets.forEach { wallet ->
                            DropdownMenuItem(
                                text = { Text("${wallet.name} (${formatCurrency(wallet.balance)})") },
                                onClick = {
                                    fromWalletId = wallet.id
                                    fromExpanded = false
                                }
                            )
                        }
                    }
                }

                Spacer(modifier = Modifier.height(12.dp))

                // Swap icon
                Row(
                    modifier = Modifier.fillMaxWidth(),
                    horizontalArrangement = Arrangement.Center
                ) {
                    Icon(
                        Icons.Filled.SwapVert,
                        contentDescription = "Swap",
                        tint = TransferPurple,
                        modifier = Modifier.size(32.dp)
                    )
                }

                Spacer(modifier = Modifier.height(12.dp))

                // To Wallet
                ExposedDropdownMenuBox(
                    expanded = toExpanded,
                    onExpandedChange = { toExpanded = it }
                ) {
                    GlassTextField(
                        value = uiState.wallets.find { it.id == toWalletId }?.name ?: "",
                        onValueChange = {},
                        label = "To Wallet",
                        placeholder = "Select destination wallet",
                        leadingIcon = { Icon(Icons.Filled.ArrowDownward, contentDescription = null, tint = IncomeGreen) },
                        modifier = Modifier
                            .fillMaxWidth()
                            .menuAnchor()
                    )
                    ExposedDropdownMenu(
                        expanded = toExpanded,
                        onDismissRequest = { toExpanded = false }
                    ) {
                        uiState.wallets
                            .filter { it.id != fromWalletId }
                            .forEach { wallet ->
                                DropdownMenuItem(
                                    text = { Text("${wallet.name} (${formatCurrency(wallet.balance)})") },
                                    onClick = {
                                        toWalletId = wallet.id
                                        toExpanded = false
                                    }
                                )
                            }
                    }
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
                    value = description,
                    onValueChange = { description = it },
                    label = "Description (Optional)",
                    placeholder = "Transfer note",
                    leadingIcon = { Icon(Icons.Filled.Description, contentDescription = null) },
                    modifier = Modifier.fillMaxWidth()
                )

                Spacer(modifier = Modifier.height(24.dp))

                if (transferState.error != null) {
                    Text(
                        text = transferState.error!!,
                        style = MaterialTheme.typography.bodySmall,
                        color = MaterialTheme.colorScheme.error,
                        modifier = Modifier.padding(bottom = 8.dp)
                    )
                }

                GlassButton(
                    onClick = {
                        val amountVal = amount.toDoubleOrNull() ?: return@GlassButton
                        viewModel.transfer(fromWalletId, toWalletId, amountVal, description.ifBlank { null })
                    },
                    enabled = !transferState.isLoading &&
                            fromWalletId.isNotBlank() && toWalletId.isNotBlank() &&
                            amount.isNotBlank() && fromWalletId != toWalletId,
                    modifier = Modifier.fillMaxWidth()
                ) {
                    if (transferState.isLoading) {
                        CircularProgressIndicator(modifier = Modifier.size(20.dp), color = Color.White, strokeWidth = 2.dp)
                        Spacer(modifier = Modifier.width(8.dp))
                    }
                    Icon(Icons.Filled.SwapHoriz, contentDescription = null, tint = Color.White, modifier = Modifier.size(18.dp))
                    Spacer(modifier = Modifier.width(8.dp))
                    Text("Transfer", color = Color.White)
                }
            }
        }
    }
}

private fun formatCurrency(amount: Double): String {
    val format = java.text.NumberFormat.getCurrencyInstance(java.util.Locale("id", "ID"))
    return format.format(amount)
}
