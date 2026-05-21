package com.financeku.app.ui.navigation

import androidx.compose.animation.*
import androidx.compose.foundation.layout.*
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.hilt.navigation.compose.hiltViewModel
import androidx.lifecycle.compose.collectAsStateWithLifecycle
import androidx.navigation.NavDestination.Companion.hierarchy
import androidx.navigation.NavGraph.Companion.findStartDestination
import androidx.navigation.NavHostController
import androidx.navigation.NavType
import androidx.navigation.compose.*
import androidx.navigation.navArgument
import com.financeku.app.data.local.datastore.TokenDataStore
import com.financeku.app.ui.auth.LoginScreen
import com.financeku.app.ui.auth.RegisterScreen
import com.financeku.app.ui.cashflow.CashflowScreen
import com.financeku.app.ui.cashflow.TransactionFormScreen
import com.financeku.app.ui.cashflow.TransactionListScreen
import com.financeku.app.ui.cashflow.TransferScreen
import com.financeku.app.ui.cashflow.WalletListScreen
import com.financeku.app.ui.dashboard.DashboardScreen
import com.financeku.app.ui.goals.GoalsScreen
import com.financeku.app.ui.overtime.OvertimeFormScreen
import com.financeku.app.ui.overtime.OvertimeListScreen
import com.financeku.app.ui.profile.ChangePasswordScreen
import com.financeku.app.ui.profile.ProfileScreen
import com.financeku.app.ui.reports.ReportsScreen
import com.financeku.app.ui.theme.LocalDarkMode

@Composable
fun FinanceKuNavHost() {
    val navController = rememberNavController()
    val tokenDataStore: TokenDataStore = hiltViewModel<NavViewModel>().tokenDataStore
    val token by tokenDataStore.accessToken.collectAsStateWithLifecycle(initialValue = null)
    val startDestination = if (token.isNullOrEmpty()) Screen.Login.route else Screen.Dashboard.route

    var currentStartDest by remember { mutableStateOf(startDestination) }

    LaunchedEffect(token) {
        if (token.isNullOrEmpty()) {
            navController.navigate(Screen.Login.route) {
                popUpTo(0) { inclusive = true }
            }
        }
    }

    Scaffold(
        bottomBar = {
            val navBackStackEntry by navController.currentBackStackEntryAsState()
            val currentDestination = navBackStackEntry?.destination
            val showBottomBar = Screen.bottomNavItems.any {
                currentDestination?.hierarchy?.any { dest -> dest.route == it.route } == true
            }

            if (showBottomBar) {
                GlassBottomBar(navController = navController, currentDestination = currentDestination)
            }
        }
    ) { paddingValues ->
        NavHost(
            navController = navController,
            startDestination = startDestination,
            modifier = Modifier.padding(paddingValues)
        ) {
            // Auth
            composable(Screen.Login.route) {
                LoginScreen(
                    onNavigateToRegister = {
                        navController.navigate(Screen.Register.route)
                    },
                    onLoginSuccess = {
                        navController.navigate(Screen.Dashboard.route) {
                            popUpTo(0) { inclusive = true }
                        }
                    }
                )
            }
            composable(Screen.Register.route) {
                RegisterScreen(
                    onNavigateToLogin = { navController.popBackStack() },
                    onRegisterSuccess = {
                        navController.navigate(Screen.Dashboard.route) {
                            popUpTo(0) { inclusive = true }
                        }
                    }
                )
            }

            // Dashboard
            composable(Screen.Dashboard.route) {
                DashboardScreen(
                    onNavigateToReports = { navController.navigate(Screen.Reports.route) }
                )
            }

            // Overtime
            composable(Screen.Overtime.route) {
                OvertimeListScreen(
                    onNavigateToForm = { id ->
                        navController.navigate(Screen.OvertimeForm.createRoute(id))
                    }
                )
            }
            composable(
                route = Screen.OvertimeForm.route,
                arguments = listOf(navArgument("id") { type = NavType.StringType; defaultValue = "" })
            ) { backStackEntry ->
                val id = backStackEntry.arguments?.getString("id")?.takeIf { it.isNotEmpty() }
                OvertimeFormScreen(
                    overtimeId = id,
                    onNavigateBack = { navController.popBackStack() }
                )
            }

            // Cashflow
            composable(Screen.Cashflow.route) {
                CashflowScreen(
                    onNavigateToTransactions = { navController.navigate(Screen.TransactionList.route) },
                    onNavigateToWallets = { navController.navigate(Screen.WalletList.route) },
                    onNavigateToTransfer = { navController.navigate(Screen.Transfer.route) },
                    onNavigateToTransactionForm = { id ->
                        navController.navigate(Screen.TransactionForm.createRoute(id))
                    }
                )
            }
            composable(Screen.TransactionList.route) {
                TransactionListScreen(
                    onNavigateBack = { navController.popBackStack() },
                    onNavigateToForm = { id ->
                        navController.navigate(Screen.TransactionForm.createRoute(id))
                    }
                )
            }
            composable(Screen.WalletList.route) {
                WalletListScreen(onNavigateBack = { navController.popBackStack() })
            }
            composable(Screen.Transfer.route) {
                TransferScreen(onNavigateBack = { navController.popBackStack() })
            }
            composable(
                route = Screen.TransactionForm.route,
                arguments = listOf(navArgument("id") { type = NavType.StringType; defaultValue = "" })
            ) { backStackEntry ->
                val id = backStackEntry.arguments?.getString("id")?.takeIf { it.isNotEmpty() }
                TransactionFormScreen(
                    transactionId = id,
                    onNavigateBack = { navController.popBackStack() }
                )
            }

            // Goals
            composable(Screen.Goals.route) {
                GoalsScreen()
            }

            // Reports
            composable(Screen.Reports.route) {
                ReportsScreen(onNavigateBack = { navController.popBackStack() })
            }

            // Profile
            composable(Screen.Profile.route) {
                ProfileScreen(
                    onNavigateToChangePassword = { navController.navigate(Screen.ChangePassword.route) },
                    onLogout = {
                        navController.navigate(Screen.Login.route) {
                            popUpTo(0) { inclusive = true }
                        }
                    }
                )
            }
            composable(Screen.ChangePassword.route) {
                ChangePasswordScreen(onNavigateBack = { navController.popBackStack() })
            }
        }
    }
}

@Composable
private fun GlassBottomBar(
    navController: NavHostController,
    currentDestination: androidx.navigation.NavDestination?
) {
    val isDark = LocalDarkMode.current.value
    val containerColor = if (isDark) Color(0xDD1A1929) else Color(0xDDFFFFFF)

    NavigationBar(
        containerColor = containerColor,
        tonalElevation = 0.dp
    ) {
        Screen.bottomNavItems.forEach { screen ->
            val selected = currentDestination?.hierarchy?.any { it.route == screen.route } == true
            NavigationBarItem(
                icon = {
                    screen.icon?.let {
                        Icon(imageVector = it, contentDescription = screen.title)
                    }
                },
                label = { Text(screen.title, style = MaterialTheme.typography.labelSmall) },
                selected = selected,
                onClick = {
                    navController.navigate(screen.route) {
                        popUpTo(navController.graph.findStartDestination().id) {
                            saveState = true
                        }
                        launchSingleTop = true
                        restoreState = true
                    }
                },
                colors = NavigationBarItemDefaults.colors(
                    selectedIconColor = MaterialTheme.colorScheme.primary,
                    selectedTextColor = MaterialTheme.colorScheme.primary,
                    unselectedIconColor = MaterialTheme.colorScheme.onSurface.copy(alpha = 0.5f),
                    unselectedTextColor = MaterialTheme.colorScheme.onSurface.copy(alpha = 0.5f),
                    indicatorColor = MaterialTheme.colorScheme.primary.copy(alpha = 0.1f)
                )
            )
        }
    }
}
