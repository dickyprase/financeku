import { Router } from 'express';
import { authController } from '../controllers/auth.controller';
import { overtimeController } from '../controllers/overtime.controller';
import { walletController } from '../controllers/wallet.controller';
import { transactionController } from '../controllers/transaction.controller';
import { categoryController } from '../controllers/category.controller';
import { goalController } from '../controllers/goal.controller';
import { incomeController } from '../controllers/income.controller';
import { dashboardController, dailyBudgetController } from '../controllers/dashboard.controller';
import { profileController, adminController } from '../controllers/profile-admin.controller';
import { requireAuth, requireAdmin } from '../middleware/auth';

const router = Router();

// Auth routes (public)
router.post('/auth/register', authController.register);
router.post('/auth/login', authController.login);
router.post('/auth/refresh', authController.refresh);

// Auth routes (protected)
router.post('/auth/logout', requireAuth, authController.logout);
router.get('/auth/me', requireAuth, authController.me);

// Overtime routes
router.get('/overtime', requireAuth, overtimeController.list);
router.post('/overtime', requireAuth, overtimeController.create);
router.get('/overtime/calculate', requireAuth, overtimeController.calculate);
router.get('/overtime/:id', requireAuth, overtimeController.getById);
router.put('/overtime/:id', requireAuth, overtimeController.update);
router.delete('/overtime/:id', requireAuth, overtimeController.delete);
router.put('/overtime/periods/disburse', requireAuth, overtimeController.disburse);

// Wallet routes
router.get('/wallets', requireAuth, walletController.list);
router.post('/wallets', requireAuth, walletController.create);
router.get('/wallets/:id', requireAuth, walletController.getById);
router.put('/wallets/:id', requireAuth, walletController.update);
router.delete('/wallets/:id', requireAuth, walletController.delete);
router.post('/wallets/transfer', requireAuth, walletController.transfer);

// Category routes
router.get('/categories', requireAuth, categoryController.list);
router.post('/categories', requireAuth, categoryController.create);
router.put('/categories/:id', requireAuth, categoryController.update);
router.delete('/categories/:id', requireAuth, categoryController.delete);

// Transaction routes
router.get('/transactions', requireAuth, transactionController.list);
router.post('/transactions', requireAuth, transactionController.create);
router.get('/transactions/:id', requireAuth, transactionController.getById);
router.put('/transactions/:id', requireAuth, transactionController.update);
router.delete('/transactions/:id', requireAuth, transactionController.delete);

// Goal routes
router.get('/goals', requireAuth, goalController.list);
router.post('/goals', requireAuth, goalController.create);
router.get('/goals/:id', requireAuth, goalController.getById);
router.put('/goals/:id', requireAuth, goalController.update);
router.delete('/goals/:id', requireAuth, goalController.delete);
router.get('/goals/:id/progress', requireAuth, goalController.getProgress);

// Income routes
router.get('/incomes', requireAuth, incomeController.list);
router.post('/incomes', requireAuth, incomeController.create);
router.delete('/incomes/:id', requireAuth, incomeController.delete);

// Daily budget routes
router.get('/daily-budget', requireAuth, dailyBudgetController.get);
router.put('/daily-budget', requireAuth, dailyBudgetController.update);
router.get('/daily-budget/today', requireAuth, dailyBudgetController.getToday);

// Reports/Dashboard routes
router.get('/reports/dashboard', requireAuth, dashboardController.getSummary);
router.get('/reports/cashflow', requireAuth, dashboardController.getCashflowReport);

// Profile routes
router.put('/profile', requireAuth, profileController.update);
router.put('/profile/password', requireAuth, profileController.changePassword);

// Admin routes
router.get('/admin/users', requireAuth, requireAdmin, adminController.listUsers);
router.post('/admin/users', requireAuth, requireAdmin, adminController.createUser);
router.put('/admin/users/:id', requireAuth, requireAdmin, adminController.updateUser);
router.delete('/admin/users/:id', requireAuth, requireAdmin, adminController.deleteUser);
router.post('/admin/users/:id/reset-password', requireAuth, requireAdmin, adminController.resetPassword);
router.get('/admin/settings', requireAuth, requireAdmin, adminController.getSettings);
router.put('/admin/settings', requireAuth, requireAdmin, adminController.updateSettings);
router.get('/admin/activity-logs', requireAuth, requireAdmin, adminController.getActivityLogs);

export default router;
