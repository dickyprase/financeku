package com.financeku.app.ui.theme

import android.app.Activity
import androidx.compose.foundation.isSystemInDarkTheme
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.graphics.toArgb
import androidx.compose.ui.platform.LocalView
import androidx.core.view.WindowCompat
import androidx.lifecycle.compose.collectAsStateWithLifecycle
import com.financeku.app.data.local.datastore.TokenDataStore
import androidx.hilt.navigation.compose.hiltViewModel
import androidx.lifecycle.ViewModel
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.flow.Flow
import javax.inject.Inject

private val LightColorScheme = lightColorScheme(
    primary = PrimaryLight,
    onPrimary = OnPrimaryLight,
    secondary = SecondaryLight,
    onSecondary = OnSecondaryLight,
    background = BackgroundLight,
    onBackground = OnBackgroundLight,
    surface = SurfaceLight,
    onSurface = OnSurfaceLight,
    error = ErrorLight,
    surfaceVariant = Color(0xFFF0EEFF),
    onSurfaceVariant = Color(0xFF49454F)
)

private val DarkColorScheme = darkColorScheme(
    primary = PrimaryDark,
    onPrimary = OnPrimaryDark,
    secondary = SecondaryDark,
    onSecondary = OnSecondaryDark,
    background = BackgroundDark,
    onBackground = OnBackgroundDark,
    surface = SurfaceDark,
    onSurface = OnSurfaceDark,
    error = ErrorDark,
    surfaceVariant = Color(0xFF252438),
    onSurfaceVariant = Color(0xFFCAC4D0)
)

val LocalDarkMode = compositionLocalOf { mutableStateOf(false) }

@HiltViewModel
class ThemeViewModel @Inject constructor(
    private val tokenDataStore: TokenDataStore
) : ViewModel() {
    val isDarkMode: Flow<Boolean> = tokenDataStore.isDarkMode

    suspend fun toggleDarkMode(isDark: Boolean) {
        tokenDataStore.setDarkMode(isDark)
    }
}

@Composable
fun FinanceKuTheme(
    content: @Composable () -> Unit
) {
    val themeViewModel: ThemeViewModel = hiltViewModel()
    val isDarkMode by themeViewModel.isDarkMode.collectAsStateWithLifecycle(initialValue = isSystemInDarkTheme())
    val darkModeState = remember { mutableStateOf(isDarkMode) }

    LaunchedEffect(isDarkMode) {
        darkModeState.value = isDarkMode
    }

    val colorScheme = if (darkModeState.value) DarkColorScheme else LightColorScheme

    val view = LocalView.current
    if (!view.isInEditMode) {
        SideEffect {
            val window = (view.context as Activity).window
            window.statusBarColor = colorScheme.background.toArgb()
            WindowCompat.getInsetsController(window, view).isAppearanceLightStatusBars = !darkModeState.value
        }
    }

    CompositionLocalProvider(LocalDarkMode provides darkModeState) {
        MaterialTheme(
            colorScheme = colorScheme,
            typography = FinanceKuTypography,
            content = content
        )
    }
}
