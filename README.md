# FinanceKu

Unified Personal Finance & Overtime Tracker — monorepo containing backend, web app, and android app.

## Overview

Menggabungkan overtime tracker + personal finance menjadi 1 aplikasi:
- Overtime yang sudah cair → otomatis masuk sebagai income di cashflow
- Wallet balance → terupdate otomatis dari overtime disbursement
- Dashboard unified → overview lembur + keuangan dalam 1 tempat

## Project Structure

```
financeku/
├── backend/       # Node.js + Express + TypeScript REST API (Phase 1) ✅
├── web/           # Laravel + Livewire (Phase 2) 🔜
└── android/       # Kotlin + Jetpack Compose (Phase 3) 🔜
```

## Tech Stack

| Layer | Technology |
|-------|-----------|
| Backend | Node.js 20+ (Express + TypeScript) + PostgreSQL |
| Web | Laravel 11 + Livewire 3 + Tailwind CSS |
| Android | Kotlin + Jetpack Compose + Material 3 |
| Auth | JWT (access + refresh token) |

## Quick Start

### Backend

```bash
cd backend

# Docker (recommended)
docker-compose up -d

# Manual
cp .env.example .env
npm install
npm run seed
npm run dev
```

API berjalan di `http://localhost:8080`

### Default Credentials

- Email: admin@financeku.com
- Password: admin123

## Documentation

- [Backend README](backend/README.md) — API endpoints, environment variables, overtime calculation
- Web README — coming soon
- Android README — coming soon

## License

MIT
