package com.financeku.app.ui.cashflow

import androidx.compose.foundation.background
import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.items
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.*
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Brush
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import androidx.hilt.navigation.compose.hiltViewModel
import com.financeku.app.ui.components.*
import com.financeku.app.ui.theme.*
import java.text.NumberFormat
import java.util.Locale

@Composable
fun CashflowScreen(
    onNavigateToTransactions: () -> Unit,
    onNavigateToWallets: () -> Unit,
    onNavigateToTransfer: () -> Unit,
    onNavigateToTransactionForm: (String?) -> Unit,
    viewModel: CashflowViewModel = hiltViewModel()
) {
    val uiState by viewModel.uiState.collectAsState()
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
            title = "Cashflow",
            actions = {
                IconButton(onClick = { onNavigateToTransactionForm(null) }) {
                    Icon(Icons.Filled.Add, contentDescription = "Add Transaction")
                }
            }
        )

        when {
            uiState.isLoading -> LoadingIndicator()
            uiState.error != null -> ErrorMessage(
                message = uiState.error!!,
                onRetry = { viewModel.loadAll() }
            )
            else -> {
                LazyColumn(
                    modifier = Modifier.fillMaxSize(),
                    contentPadding = PaddingValues(16.dp),
                    verticalArrangement = Arrangement.spacedBy(12.dp)
                ) {
                    // Quick Actions
                    item {
                        Row(
                            modifier = Modifier.fillMaxWidth(),
                            horizontalArrangement = Arrangement.spacedBy(12.dp)
                        ) {
                            QuickActionCard(
                                icon = Icons.Filled.AccountBalanceWallet,
                                label = "Wallets",
                                color = GoalBlue,
                                modifier = Modifier.weight(1f),
                                onClick = onNavigateToWallets
                            )
                            QuickActionCard(
                                icon = Icons.Filled.SwapHoriz,
                                label = "Transfer",
                                color = TransferPurple,
                                modifier = Modifier.weight(1f),
                                onClick = onNavigateToTransfer
                            )
                        }
                    }

                    // Wallets Summary
                    item {
                        Text(
                            text = "Wallets",
                            style = MaterialTheme.typography.titleMedium,
                            color = MaterialTheme.colorScheme.onSurface,
                            modifier = Modifier.padding(top = 8.dp)
                        )
                    }

                    if (uiState.wallets.isEmpty()) {
                        item {
                            GlassCard(modifier = Modifier.fillMaxWidth()) {
                                Text(
                                    "No wallets yet. Create one to get started.",
                                    style = MaterialTheme.typography.bodyMedium,
                                    color = MaterialTheme.colorScheme.onSurface.copy(alpha = 0.6f)
                                )
                            }
                        }
                    } else {
                        items(uiState.wallets) { wallet ->
                            GlassCard(modifier = Modifier.fillMaxWidth()) {
                                Row(
                                    modifier = Modifier.fillMaxWidth(),
                                    horizontalArrangement = Arrangement.SpaceBetween,
                                    verticalAlignment = Alignment.CenterVertically
                                ) {
                                    Row(verticalAlignment = Alignment.CenterVertically) {
                                        Icon(
                                            Icons.Filled.AccountBalanceWallet,
                                            contentDescription = null,
                                            tint = GoalBlue,
                                            modifier = Modifier.size(24.dp)
                                        )
                                        Spacer(modifier = Modifier.width(12.dp))
                                        Column {
                                            Text(
                                                text = wallet.name,
                                                style = MaterialTheme.typography.bodyMedium,
                                                fontWeight = FontWeight.Medium
                                            )
                                            if (wallet.isDefault) {
                                                Text(
                                                    "Default",
                                                    style = MaterialTheme.typography.labelSmall,
                                                    color = MaterialTheme.colorScheme.primary
                                                )
                                            }
                                        }
                                    }
                                    Text(
                                        text = formatCurrency(wallet.balance),
                                        style = MaterialTheme.typography.titleSmall,
                                        fontWeight = FontWeight.SemiBold,
                                        color = if (wallet.balance >= 0) IncomeGreen else ExpenseRed
                                    )
                                }
                            }
                        }
                    }

                    // Recent Transactions
                    item {
                        Row(
                            modifier = Modifier.fillMaxWidth(),
                            horizontalArrangement = Arrangement.SpaceBetween,
                            verticalAlignment = Alignment.CenterVertically
                        ) {
                            Text(
                                text = "Recent Transactions",
                                style = MaterialTheme.typography.titleMedium,
                                color = MaterialTheme.colorScheme.onSurface
                            )
                            TextButton(onClick = onNavigateToTransactions) {
                                Text("See All")
                            }
                        }
                    }

                    if (uiState.transactions.isEmpty()) {
                        item {
                            GlassCard(modifier = Modifier.fillMaxWidth()) {
                                Text(
                                    "No transactions yet.",
                                    style = MaterialTheme.typography.bodyMedium,
                                    color = MaterialTheme.colorScheme.onSurface.copy(alpha = 0.6f)
                                )
                            }
                        }
                    } else {
                        items(uiState.transactions.take(10)) { transaction ->
                            GlassCard(
                                modifier = Modifier
                                    .fillMaxWidth()
                                    .clickable { onNavigateToTransactionForm(transaction.id) }
                            ) {
                                Row(
                                    modifier = Modifier.fillMaxWidth(),
                                    horizontalArrangement = Arrangement.SpaceBetween,
                                    verticalAlignment = Alignment.CenterVertically
                                ) {
                                    Row(
                                        verticalAlignment = Alignment.CenterVertically,
                                        modifier = Modifier.weight(1f)
                                    ) {
                                        Icon(
                                            if (transaction.type == "income") Icons.Filled.ArrowDownward
                                            else Icons.Filled.ArrowUpward,
                                            contentDescription = null,
                                            tint = if (transaction.type == "income") IncomeGreen else ExpenseRed,
                                            modifier = Modifier.size(20.dp)
                                        )
                                        Spacer(modifier = Modifier.width(12.dp))
                                        Column {
                                            Text(
                                                text = transaction.description ?: transaction.categoryName ?: "Transaction",
                                                style = MaterialTheme.typography.bodyMedium,
                                                maxLines = 1
                                            )
                                            Text(
                                                text = "${transaction.walletName ?: ""} • ${transaction.date}",
                                                style = MaterialTheme.typography.bodySmall,
                                                color = MaterialTheme.colorScheme.onSurface.copy(alpha = 0.5f)
                                            )
                                        }
                                    }
                                    Text(
                                        text = "${if (transaction.type == "expense") "-" else "+"}${formatCurrency(transaction.amount)}",
                                        style = MaterialTheme.typography.titleSmall,
                                        color = if (transaction.type == "expense") ExpenseRed else IncomeGreen,
                                        fontWeight = FontWeight.SemiBold
                                    )
                                }
                            }
                        }
                    }

                    item { Spacer(modifier = Modifier.height(16.dp)) }
                }
            }
        }
    }
}

@Composable
private fun QuickActionCard(
    icon: androidx.compose.ui.graphics.vector.ImageVector,
    label: String,
    color: Color,
    modifier: Modifier = Modifier,
    onClick: () -> Unit
) {
    GlassCard(
        modifier = modifier.clickable { onClick() }
    ) {
        Column(
            horizontalAlignment = Alignment.CenterHorizontally,
            modifier = Modifier.fillMaxWidth()
        ) {
            Icon(icon, contentDescription = label, tint = color, modifier = Modifier.size(28.dp))
            Spacer(modifier = Modifier.height(8.dp))
            Text(label, style = MaterialTheme.typography.labelMedium)
        }
    }
}

private fun formatCurrency(amount: Double): String {
    val format = NumberFormat.getCurrencyInstance(Locale("id", "ID"))
    return format.format(amount)
}
