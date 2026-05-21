# FinanceKu Backend

Unified Personal Finance & Overtime Tracker REST API built with Node.js, Express, and TypeScript.

Menggabungkan overtime tracker + personal finance menjadi 1 aplikasi:
- Overtime yang sudah cair → otomatis masuk sebagai income di cashflow
- Wallet balance → terupdate otomatis dari overtime disbursement
- Dashboard unified → overview lembur + keuangan dalam 1 tempat

## Tech Stack

- **Runtime**: Node.js 20+
- **Language**: TypeScript
- **Framework**: Express.js
- **Database**: PostgreSQL 16 (via node-postgres)
- **Auth**: JWT (jsonwebtoken)
- **Validation**: Zod
- **Password**: bcrypt
- **Testing**: Jest

## Fitur

- Auth (Register, Login, Refresh Token, Logout)
- Overtime Management (CRUD, auto-calculate, period Thu-Wed, disburse)
- Wallet Management (CRUD, transfer antar wallet)
- Transaction Management (income/expense, auto-update balance)
- Category Management (income/expense, budget limit)
- Goal Tracking (manual/single/multiple/all wallet tracking)
- Income Management (standalone + overtime disbursement)
- Daily Budget (manual/formula mode)
- Dashboard & Reports (summary, monthly cashflow)
- Profile & Admin Panel (user management, site settings, activity logs)
- Middleware (CORS, rate limiting, request logger)

## Struktur Project

```
backend/
├── src/
│   ├── config/            # env loader, database pool
│   ├── controllers/       # request handlers
│   ├── middleware/        # auth JWT
│   ├── models/            # TypeScript interfaces
│   ├── repositories/      # SQL queries via pg
│   ├── services/          # business logic
│   ├── routes/            # Express router
│   ├── utils/             # response, overtime calc, hash, pagination, validators
│   ├── app.ts             # Express app setup
│   ├── server.ts          # Entry point
│   └── migrate.ts         # Migration runner
├── migrations/            # SQL migration files
├── tests/                 # Unit tests
├── package.json
├── tsconfig.json
├── Dockerfile
├── docker-compose.yml
├── Makefile
└── .env.example
```

## Cara Menjalankan

### Prasyarat

- Node.js 20+
- PostgreSQL 16+
- Docker & Docker Compose (opsional)

### Opsi 1: Docker Compose (Recommended)

```bash
# Clone repository
git clone https://github.com/dickyprase/financeku.git
cd financeku/backend

# Jalankan PostgreSQL + API
docker-compose up -d

# API akan berjalan di http://localhost:8080
```

### Opsi 2: Manual

```bash
# Clone repository
git clone https://github.com/dickyprase/financeku.git
cd financeku/backend

# Install dependencies
npm install

# Copy dan edit environment file
cp .env.example .env
# Edit .env sesuai konfigurasi PostgreSQL lokal

# Jalankan migration + seed
npm run seed

# Jalankan server (development)
npm run dev
```

### Opsi 3: Build Production

```bash
npm run build
npm start
```

## Environment Variables

| Variable | Default | Keterangan |
|----------|---------|------------|
| SERVER_PORT | 8080 | Port server |
| DB_HOST | localhost | PostgreSQL host |
| DB_PORT | 5432 | PostgreSQL port |
| DB_USER | postgres | PostgreSQL user |
| DB_PASSWORD | postgres | PostgreSQL password |
| DB_NAME | financeku | Database name |
| DB_SSLMODE | disable | SSL mode |
| JWT_SECRET | (required) | Secret key untuk JWT |
| JWT_ACCESS_EXP_MINUTES | 15 | Access token expiry |
| JWT_REFRESH_EXP_DAYS | 7 | Refresh token expiry |
| ALLOWED_ORIGINS | * | CORS allowed origins |

## Default Credentials

Setelah seed:
- **Email**: admin@financeku.com
- **Password**: admin123
- **Role**: admin

## API Endpoints

### Auth
| Method | Endpoint | Keterangan |
|--------|----------|------------|
| POST | /api/v1/auth/register | Register user baru |
| POST | /api/v1/auth/login | Login, return JWT tokens |
| POST | /api/v1/auth/refresh | Refresh access token |
| POST | /api/v1/auth/logout | Logout |
| GET | /api/v1/auth/me | Get current user |

### Overtime
| Method | Endpoint | Keterangan |
|--------|----------|------------|
| GET | /api/v1/overtime | List overtime records |
| POST | /api/v1/overtime | Create overtime record |
| GET | /api/v1/overtime/calculate | Preview calculation |
| GET | /api/v1/overtime/:id | Get by ID |
| PUT | /api/v1/overtime/:id | Update |
| DELETE | /api/v1/overtime/:id | Delete |
| PUT | /api/v1/overtime/periods/disburse | Disburse period |

### Wallets
| Method | Endpoint | Keterangan |
|--------|----------|------------|
| GET | /api/v1/wallets | List wallets |
| POST | /api/v1/wallets | Create wallet |
| GET | /api/v1/wallets/:id | Get by ID |
| PUT | /api/v1/wallets/:id | Update |
| DELETE | /api/v1/wallets/:id | Delete (soft) |
| POST | /api/v1/wallets/transfer | Transfer antar wallet |

### Categories
| Method | Endpoint | Keterangan |
|--------|----------|------------|
| GET | /api/v1/categories | List categories |
| POST | /api/v1/categories | Create category |
| PUT | /api/v1/categories/:id | Update |
| DELETE | /api/v1/categories/:id | Delete |

### Transactions
| Method | Endpoint | Keterangan |
|--------|----------|------------|
| GET | /api/v1/transactions | List (filter: wallet, category, type, date) |
| POST | /api/v1/transactions | Create |
| GET | /api/v1/transactions/:id | Get by ID |
| PUT | /api/v1/transactions/:id | Update |
| DELETE | /api/v1/transactions/:id | Delete |

### Goals
| Method | Endpoint | Keterangan |
|--------|----------|------------|
| GET | /api/v1/goals | List goals |
| POST | /api/v1/goals | Create goal |
| GET | /api/v1/goals/:id | Get by ID |
| PUT | /api/v1/goals/:id | Update |
| DELETE | /api/v1/goals/:id | Delete |
| GET | /api/v1/goals/:id/progress | Get progress |

### Incomes
| Method | Endpoint | Keterangan |
|--------|----------|------------|
| GET | /api/v1/incomes | List incomes |
| POST | /api/v1/incomes | Create income |
| DELETE | /api/v1/incomes/:id | Delete |

### Daily Budget
| Method | Endpoint | Keterangan |
|--------|----------|------------|
| GET | /api/v1/daily-budget | Get settings |
| PUT | /api/v1/daily-budget | Update settings |
| GET | /api/v1/daily-budget/today | Get today's budget |

### Reports
| Method | Endpoint | Keterangan |
|--------|----------|------------|
| GET | /api/v1/reports/dashboard | Dashboard summary |
| GET | /api/v1/reports/cashflow | Monthly cashflow report |

### Profile
| Method | Endpoint | Keterangan |
|--------|----------|------------|
| PUT | /api/v1/profile | Update profile |
| PUT | /api/v1/profile/password | Change password |

### Admin (requires admin role)
| Method | Endpoint | Keterangan |
|--------|----------|------------|
| GET | /api/v1/admin/users | List users |
| POST | /api/v1/admin/users | Create user |
| PUT | /api/v1/admin/users/:id | Update user |
| DELETE | /api/v1/admin/users/:id | Delete user |
| POST | /api/v1/admin/users/:id/reset-password | Reset password |
| GET | /api/v1/admin/settings | Get site settings |
| PUT | /api/v1/admin/settings | Update settings |
| GET | /api/v1/admin/activity-logs | Activity logs |

## Overtime Calculation

### Weekday (Progressive Multiplier)
Basis: `salary / 173`

| Jam | Multiplier | Formula |
|-----|-----------|---------|
| 1.0 | 1.5x | hourly_rate × 1.5 + meal |
| 1.5 | 2.5x | hourly_rate × 2.5 + meal |
| 2.0 | 3.5x | hourly_rate × 3.5 + meal |
| 2.5 | 4.5x | hourly_rate × 4.5 + meal |
| 3.0 | 5.5x | hourly_rate × 5.5 + meal |

### Holiday
```
amount = (salary / 173) × 2 × hours + meal
```

### Period
- Period: Kamis - Rabu
- Payment: period_end + 9 hari

## Testing

```bash
# Run all tests
npm test

# Watch mode
npm run test:watch
```

## Scripts

```bash
npm run dev          # Development server (hot reload)
npm run build        # Build TypeScript
npm start            # Start production server
npm run migrate      # Run migrations
npm run seed         # Run migrations + seed
npm test             # Run tests
npm run lint         # Lint code
```

## License

MIT
