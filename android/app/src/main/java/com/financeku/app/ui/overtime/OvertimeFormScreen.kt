package com.financeku.app.ui.overtime

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

@Composable
fun OvertimeFormScreen(
    overtimeId: String?,
    onNavigateBack: () -> Unit,
    viewModel: OvertimeViewModel = hiltViewModel()
) {
    val formState by viewModel.formState.collectAsState()
    var date by remember { mutableStateOf("") }
    var hours by remember { mutableStateOf("") }
    var rate by remember { mutableStateOf("") }
    var description by remember { mutableStateOf("") }

    val isEditing = overtimeId != null

    LaunchedEffect(formState.success) {
        if (formState.success) {
            viewModel.resetFormState()
            onNavigateBack()
        }
    }

    val isDark = LocalDarkMode.current.value
    val bgGradient = if (isDark) {
        Brush.verticalGradient(listOf(BackgroundDark, GradientDarkMiddle))
    } else {
        Brush.verticalGradient(listOf(BackgroundLight, Color(0xFFFFF3E0)))
    }

    Column(
        modifier = Modifier
            .fillMaxSize()
            .background(bgGradient)
    ) {
        GlassTopBar(
            title = if (isEditing) "Edit Overtime" else "Add Overtime",
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
                GlassTextField(
                    value = date,
                    onValueChange = { date = it },
                    label = "Date",
                    placeholder = "YYYY-MM-DD",
                    leadingIcon = { Icon(Icons.Filled.CalendarToday, contentDescription = null) },
                    modifier = Modifier.fillMaxWidth()
                )

                Spacer(modifier = Modifier.height(16.dp))

                GlassTextField(
                    value = hours,
                    onValueChange = { hours = it },
                    label = "Hours",
                    placeholder = "e.g. 2.5",
                    leadingIcon = { Icon(Icons.Filled.AccessTime, contentDescription = null) },
                    keyboardOptions = KeyboardOptions(keyboardType = KeyboardType.Decimal),
                    modifier = Modifier.fillMaxWidth()
                )

                Spacer(modifier = Modifier.height(16.dp))

                GlassTextField(
                    value = rate,
                    onValueChange = { rate = it },
                    label = "Rate per Hour",
                    placeholder = "e.g. 50000",
                    leadingIcon = { Icon(Icons.Filled.AttachMoney, contentDescription = null) },
                    keyboardOptions = KeyboardOptions(keyboardType = KeyboardType.Number),
                    modifier = Modifier.fillMaxWidth()
                )

                Spacer(modifier = Modifier.height(16.dp))

                GlassTextField(
                    value = description,
                    onValueChange = { description = it },
                    label = "Description (Optional)",
                    placeholder = "What did you work on?",
                    leadingIcon = { Icon(Icons.Filled.Description, contentDescription = null) },
                    singleLine = false,
                    maxLines = 3,
                    modifier = Modifier.fillMaxWidth()
                )

                // Preview calculation
                val hoursVal = hours.toDoubleOrNull()
                val rateVal = rate.toDoubleOrNull()
                if (hoursVal != null && rateVal != null) {
                    Spacer(modifier = Modifier.height(16.dp))
                    GlassCard(modifier = Modifier.fillMaxWidth()) {
                        Text(
                            text = "Preview",
                            style = MaterialTheme.typography.labelMedium,
                            color = MaterialTheme.colorScheme.onSurface.copy(alpha = 0.6f)
                        )
                        Spacer(modifier = Modifier.height(4.dp))
                        Text(
                            text = "Amount: Rp ${String.format("%,.0f", hoursVal * rateVal)}",
                            style = MaterialTheme.typography.titleMedium,
                            color = OvertimeOrange
                        )
                    }
                }

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
                        val h = hours.toDoubleOrNull() ?: return@GlassButton
                        val r = rate.toDoubleOrNull() ?: return@GlassButton
                        if (isEditing) {
                            viewModel.updateOvertime(overtimeId!!, date, h, r, description.ifBlank { null })
                        } else {
                            viewModel.createOvertime(date, h, r, description.ifBlank { null })
                        }
                    },
                    enabled = !formState.isLoading && date.isNotBlank() && hours.isNotBlank() && rate.isNotBlank(),
                    modifier = Modifier.fillMaxWidth()
                ) {
                    if (formState.isLoading) {
                        CircularProgressIndicator(
                            modifier = Modifier.size(20.dp),
                            color = Color.White,
                            strokeWidth = 2.dp
                        )
                        Spacer(modifier = Modifier.width(8.dp))
                    }
                    Text(
                        if (isEditing) "Update Overtime" else "Save Overtime",
                        color = Color.White
                    )
                }
            }
        }
    }
}
