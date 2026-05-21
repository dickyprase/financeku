package com.financeku.app.ui.goals

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
fun GoalsScreen(
    viewModel: GoalsViewModel = hiltViewModel()
) {
    val uiState by viewModel.uiState.collectAsState()
    var showCreateDialog by remember { mutableStateOf(false) }

    val isDark = LocalDarkMode.current.value
    val bgGradient = if (isDark) {
        Brush.verticalGradient(listOf(BackgroundDark, GradientDarkMiddle))
    } else {
        Brush.verticalGradient(listOf(BackgroundLight, Color(0xFFE3F2FD)))
    }

    Column(
        modifier = Modifier
            .fillMaxSize()
            .background(bgGradient)
    ) {
        GlassTopBar(
            title = "Goals",
            actions = {
                IconButton(onClick = { showCreateDialog = true }) {
                    Icon(Icons.Filled.Add, contentDescription = "Add Goal")
                }
            }
        )

        when {
            uiState.isLoading -> LoadingIndicator()
            uiState.error != null -> ErrorMessage(
                message = uiState.error!!,
                onRetry = { viewModel.loadGoals() }
            )
            uiState.goals.isEmpty() -> {
                Box(
                    modifier = Modifier.fillMaxSize(),
                    contentAlignment = Alignment.Center
                ) {
                    Column(horizontalAlignment = Alignment.CenterHorizontally) {
                        Icon(
                            Icons.Filled.Flag,
                            contentDescription = null,
                            modifier = Modifier.size(64.dp),
                            tint = MaterialTheme.colorScheme.onSurface.copy(alpha = 0.3f)
                        )
                        Spacer(modifier = Modifier.height(16.dp))
                        Text(
                            "No goals yet",
                            style = MaterialTheme.typography.bodyLarge,
                            color = MaterialTheme.colorScheme.onSurface.copy(alpha = 0.5f)
                        )
                        Spacer(modifier = Modifier.height(8.dp))
                        Text(
                            "Set a financial goal to start tracking",
                            style = MaterialTheme.typography.bodyMedium,
                            color = MaterialTheme.colorScheme.onSurface.copy(alpha = 0.4f)
                        )
                    }
                }
            }
            else -> {
                LazyColumn(
                    modifier = Modifier.fillMaxSize(),
                    contentPadding = PaddingValues(16.dp),
                    verticalArrangement = Arrangement.spacedBy(12.dp)
                ) {
                    items(uiState.goals) { goal ->
                        GoalItem(
                            goal = goal,
                            onDelete = { viewModel.deleteGoal(goal.id) },
                            onUpdate = { currentAmount ->
                                viewModel.updateGoal(
                                    goal.id, goal.name, goal.targetAmount,
                                    currentAmount, goal.deadline, goal.icon, goal.color
                                )
                            }
                        )
                    }
                    item { Spacer(modifier = Modifier.height(16.dp)) }
                }
            }
        }
    }

    if (showCreateDialog) {
        CreateGoalDialog(
            onDismiss = { showCreateDialog = false },
            onCreate = { name, target, deadline ->
                viewModel.createGoal(name, target, 0.0, deadline, null, null)
                showCreateDialog = false
            }
        )
    }
}

@Composable
private fun GoalItem(
    goal: com.financeku.app.domain.model.Goal,
    onDelete: () -> Unit,
    onUpdate: (Double) -> Unit
) {
    var showDeleteDialog by remember { mutableStateOf(false) }
    var showUpdateDialog by remember { mutableStateOf(false) }
    val progress = if (goal.targetAmount > 0) (goal.currentAmount / goal.targetAmount).toFloat().coerceIn(0f, 1f) else 0f
    val percentage = (progress * 100).toInt()

    GlassCard(modifier = Modifier.fillMaxWidth()) {
        Row(
            modifier = Modifier.fillMaxWidth(),
            horizontalArrangement = Arrangement.SpaceBetween,
            verticalAlignment = Alignment.Top
        ) {
            Column(modifier = Modifier.weight(1f)) {
                Text(
                    text = goal.name,
                    style = MaterialTheme.typography.titleSmall,
                    fontWeight = FontWeight.SemiBold,
                    color = MaterialTheme.colorScheme.onSurface
                )
                if (!goal.deadline.isNullOrEmpty()) {
                    Text(
                        text = "Deadline: ${goal.deadline}",
                        style = MaterialTheme.typography.bodySmall,
                        color = MaterialTheme.colorScheme.onSurface.copy(alpha = 0.5f)
                    )
                }
            }
            Row {
                IconButton(onClick = { showUpdateDialog = true }, modifier = Modifier.size(28.dp)) {
                    Icon(Icons.Filled.Add, contentDescription = "Add funds", modifier = Modifier.size(16.dp), tint = IncomeGreen)
                }
                IconButton(onClick = { showDeleteDialog = true }, modifier = Modifier.size(28.dp)) {
                    Icon(Icons.Filled.Delete, contentDescription = "Delete", modifier = Modifier.size(16.dp), tint = ExpenseRed)
                }
            }
        }

        Spacer(modifier = Modifier.height(12.dp))

        LinearProgressIndicator(
            progress = { progress },
            modifier = Modifier
                .fillMaxWidth()
                .height(10.dp),
            color = when {
                percentage >= 100 -> IncomeGreen
                percentage >= 50 -> GoalBlue
                else -> OvertimeOrange
            },
            trackColor = MaterialTheme.colorScheme.onSurface.copy(alpha = 0.1f)
        )

        Spacer(modifier = Modifier.height(8.dp))

        Row(
            modifier = Modifier.fillMaxWidth(),
            horizontalArrangement = Arrangement.SpaceBetween
        ) {
            Text(
                text = formatCurrency(goal.currentAmount),
                style = MaterialTheme.typography.bodySmall,
                color = MaterialTheme.colorScheme.onSurface.copy(alpha = 0.7f)
            )
            Text(
                text = "$percentage%",
                style = MaterialTheme.typography.labelMedium,
                fontWeight = FontWeight.Bold,
                color = when {
                    percentage >= 100 -> IncomeGreen
                    percentage >= 50 -> GoalBlue
                    else -> OvertimeOrange
                }
            )
            Text(
                text = formatCurrency(goal.targetAmount),
                style = MaterialTheme.typography.bodySmall,
                color = MaterialTheme.colorScheme.onSurface.copy(alpha = 0.7f)
            )
        }

        if (goal.status == "completed") {
            Spacer(modifier = Modifier.height(8.dp))
            AssistChip(
                onClick = {},
                label = { Text("Completed", style = MaterialTheme.typography.labelSmall) },
                leadingIcon = { Icon(Icons.Filled.CheckCircle, contentDescription = null, modifier = Modifier.size(14.dp)) },
                colors = AssistChipDefaults.assistChipColors(containerColor = IncomeGreen.copy(alpha = 0.1f))
            )
        }
    }

    if (showDeleteDialog) {
        AlertDialog(
            onDismissRequest = { showDeleteDialog = false },
            title = { Text("Delete Goal") },
            text = { Text("Are you sure you want to delete '${goal.name}'?") },
            confirmButton = {
                TextButton(onClick = { onDelete(); showDeleteDialog = false }) {
                    Text("Delete", color = ExpenseRed)
                }
            },
            dismissButton = {
                TextButton(onClick = { showDeleteDialog = false }) { Text("Cancel") }
            }
        )
    }

    if (showUpdateDialog) {
        var addAmount by remember { mutableStateOf("") }
        AlertDialog(
            onDismissRequest = { showUpdateDialog = false },
            title = { Text("Add to Goal") },
            text = {
                OutlinedTextField(
                    value = addAmount,
                    onValueChange = { addAmount = it },
                    label = { Text("Amount to add") },
                    singleLine = true,
                    modifier = Modifier.fillMaxWidth()
                )
            },
            confirmButton = {
                TextButton(onClick = {
                    val add = addAmount.toDoubleOrNull() ?: 0.0
                    onUpdate(goal.currentAmount + add)
                    showUpdateDialog = false
                }) {
                    Text("Add")
                }
            },
            dismissButton = {
                TextButton(onClick = { showUpdateDialog = false }) { Text("Cancel") }
            }
        )
    }
}

@Composable
private fun CreateGoalDialog(
    onDismiss: () -> Unit,
    onCreate: (String, Double, String?) -> Unit
) {
    var name by remember { mutableStateOf("") }
    var target by remember { mutableStateOf("") }
    var deadline by remember { mutableStateOf("") }

    AlertDialog(
        onDismissRequest = onDismiss,
        title = { Text("Create Goal") },
        text = {
            Column(verticalArrangement = Arrangement.spacedBy(12.dp)) {
                OutlinedTextField(
                    value = name,
                    onValueChange = { name = it },
                    label = { Text("Goal Name") },
                    singleLine = true,
                    modifier = Modifier.fillMaxWidth()
                )
                OutlinedTextField(
                    value = target,
                    onValueChange = { target = it },
                    label = { Text("Target Amount") },
                    singleLine = true,
                    modifier = Modifier.fillMaxWidth()
                )
                OutlinedTextField(
                    value = deadline,
                    onValueChange = { deadline = it },
                    label = { Text("Deadline (YYYY-MM-DD, optional)") },
                    singleLine = true,
                    modifier = Modifier.fillMaxWidth()
                )
            }
        },
        confirmButton = {
            TextButton(
                onClick = {
                    val targetVal = target.toDoubleOrNull() ?: return@TextButton
                    onCreate(name, targetVal, deadline.ifBlank { null })
                },
                enabled = name.isNotBlank() && target.isNotBlank()
            ) {
                Text("Create")
            }
        },
        dismissButton = {
            TextButton(onClick = onDismiss) { Text("Cancel") }
        }
    )
}

private fun formatCurrency(amount: Double): String {
    val format = NumberFormat.getCurrencyInstance(Locale("id", "ID"))
    return format.format(amount)
}
