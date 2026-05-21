package com.financeku.app.ui.reports

import androidx.compose.foundation.background
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
fun ReportsScreen(
    onNavigateBack: () -> Unit,
    viewModel: ReportsViewModel = hiltViewModel()
) {
    val uiState by viewModel.uiState.collectAsState()
    val isDark = LocalDarkMode.current.value
    val bgGradient = if (isDark) {
        Brush.verticalGradient(listOf(BackgroundDark, GradientDarkMiddle))
    } else {
        Brush.verticalGradient(listOf(BackgroundLight, Color(0xFFE8EAF6)))
    }

    Column(
        modifier = Modifier
            .fillMaxSize()
            .background(bgGradient)
    ) {
        GlassTopBar(
            title = "Reports",
            navigationIcon = {
                IconButton(onClick = onNavigateBack) {
                    Icon(Icons.Filled.ArrowBack, contentDescription = "Back")
                }
            }
        )

        when {
            uiState.isLoading -> LoadingIndicator()
            uiState.error != null -> ErrorMessage(
                message = uiState.error!!,
                onRetry = { viewModel.loadReport() }
            )
            else -> {
                val report = uiState.report

                LazyColumn(
                    modifier = Modifier.fillMaxSize(),
                    contentPadding = PaddingValues(16.dp),
                    verticalArrangement = Arrangement.spacedBy(12.dp)
                ) {
                    // Summary Cards
                    item {
                        Text(
                            text = "Monthly Summary",
                            style = MaterialTheme.typography.titleMedium,
                            color = MaterialTheme.colorScheme.onSurface
                        )
                    }

                    item {
                        Row(
                            modifier = Modifier.fillMaxWidth(),
                            horizontalArrangement = Arrangement.spacedBy(12.dp)
                        ) {
                            GlassCard(modifier = Modifier.weight(1f)) {
                                Icon(
                                    Icons.Filled.TrendingUp,
                                    contentDescription = null,
                                    tint = IncomeGreen,
                                    modifier = Modifier.size(24.dp)
                                )
                                Spacer(modifier = Modifier.height(8.dp))
                                Text(
                                    "Income",
                                    style = MaterialTheme.typography.bodySmall,
                                    color = MaterialTheme.colorScheme.onSurface.copy(alpha = 0.6f)
                                )
                                Text(
                                    text = formatCurrency(report?.totalIncome ?: 0.0),
                                    style = MaterialTheme.typography.titleSmall,
                                    color = IncomeGreen,
                                    fontWeight = FontWeight.Bold
                                )
                            }
                            GlassCard(modifier = Modifier.weight(1f)) {
                                Icon(
                                    Icons.Filled.TrendingDown,
                                    contentDescription = null,
                                    tint = ExpenseRed,
                                    modifier = Modifier.size(24.dp)
                                )
                                Spacer(modifier = Modifier.height(8.dp))
                                Text(
                                    "Expense",
                                    style = MaterialTheme.typography.bodySmall,
                                    color = MaterialTheme.colorScheme.onSurface.copy(alpha = 0.6f)
                                )
                                Text(
                                    text = formatCurrency(report?.totalExpense ?: 0.0),
                                    style = MaterialTheme.typography.titleSmall,
                                    color = ExpenseRed,
                                    fontWeight = FontWeight.Bold
                                )
                            }
                        }
                    }

                    // Net Cashflow
                    item {
                        GlassCard(modifier = Modifier.fillMaxWidth()) {
                            Row(
                                modifier = Modifier.fillMaxWidth(),
                                horizontalArrangement = Arrangement.SpaceBetween,
                                verticalAlignment = Alignment.CenterVertically
                            ) {
                                Column {
                                    Text(
                                        "Net Cashflow",
                                        style = MaterialTheme.typography.bodyMedium,
                                        color = MaterialTheme.colorScheme.onSurface.copy(alpha = 0.6f)
                                    )
                                    Spacer(modifier = Modifier.height(4.dp))
                                    val netCashflow = report?.netCashflow ?: 0.0
                                    Text(
                                        text = "${if (netCashflow >= 0) "+" else ""}${formatCurrency(netCashflow)}",
                                        style = MaterialTheme.typography.headlineSmall,
                                        color = if (netCashflow >= 0) IncomeGreen else ExpenseRed,
                                        fontWeight = FontWeight.Bold
                                    )
                                }
                                Icon(
                                    if ((report?.netCashflow ?: 0.0) >= 0) Icons.Filled.TrendingUp else Icons.Filled.TrendingDown,
                                    contentDescription = null,
                                    tint = if ((report?.netCashflow ?: 0.0) >= 0) IncomeGreen else ExpenseRed,
                                    modifier = Modifier.size(40.dp)
                                )
                            }
                        }
                    }

                    // By Category
                    if (!report?.byCategory.isNullOrEmpty()) {
                        item {
                            Text(
                                text = "By Category",
                                style = MaterialTheme.typography.titleMedium,
                                color = MaterialTheme.colorScheme.onSurface,
                                modifier = Modifier.padding(top = 8.dp)
                            )
                        }

                        items(report!!.byCategory) { category ->
                            GlassCard(modifier = Modifier.fillMaxWidth()) {
                                Row(
                                    modifier = Modifier.fillMaxWidth(),
                                    horizontalArrangement = Arrangement.SpaceBetween,
                                    verticalAlignment = Alignment.CenterVertically
                                ) {
                                    Column(modifier = Modifier.weight(1f)) {
                                        Text(
                                            text = category.categoryName,
                                            style = MaterialTheme.typography.bodyMedium,
                                            fontWeight = FontWeight.Medium
                                        )
                                        Spacer(modifier = Modifier.height(4.dp))
                                        LinearProgressIndicator(
                                            progress = { (category.percentage / 100.0).toFloat().coerceIn(0f, 1f) },
                                            modifier = Modifier
                                                .fillMaxWidth(0.7f)
                                                .height(6.dp),
                                            color = MaterialTheme.colorScheme.primary,
                                            trackColor = MaterialTheme.colorScheme.onSurface.copy(alpha = 0.1f)
                                        )
                                    }
                                    Column(horizontalAlignment = Alignment.End) {
                                        Text(
                                            text = formatCurrency(category.amount),
                                            style = MaterialTheme.typography.titleSmall,
                                            fontWeight = FontWeight.SemiBold
                                        )
                                        Text(
                                            text = "${String.format("%.1f", category.percentage)}%",
                                            style = MaterialTheme.typography.bodySmall,
                                            color = MaterialTheme.colorScheme.onSurface.copy(alpha = 0.5f)
                                        )
                                    }
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

private fun formatCurrency(amount: Double): String {
    val format = NumberFormat.getCurrencyInstance(Locale("id", "ID"))
    return format.format(amount)
}
