# FinanceKu Android App

Aplikasi Android untuk FinanceKu — Unified Personal Finance & Overtime Tracker, dibangun dengan Kotlin dan Jetpack Compose.

## Tech Stack

| Komponen | Teknologi |
|----------|-----------|
| Language | Kotlin |
| UI | Jetpack Compose + Material 3 |
| Architecture | MVVM + Clean Architecture |
| DI | Hilt |
| Network | Retrofit + OkHttp |
| Local DB | Room |
| Auth Storage | DataStore |
| Navigation | Compose Navigation |

## Fitur

- Login & Register dengan JWT auto-refresh
- Dashboard (total balance, income/expense bulanan, overtime pending)
- Overtime Management (input, list, kalkulasi preview, disburse period)
- Cashflow (wallets, transactions, transfer antar wallet, daily budget)
- Goals (tracking progress, link ke wallet)
- Reports (monthly summary)
- Profile (edit data, change password, theme toggle)
- Glass Morphism UI dengan dark/light mode
- Offline cache via Room database

## Struktur Project

```
app/src/main/java/com/financeku/app/
├── data/
│   ├── api/                    # Retrofit ApiService, AuthInterceptor, TokenAuthenticator
│   │   └── model/             # Request/Response DTOs
│   ├── local/
│   │   ├── dao/               # Room DAOs
│   │   ├── datastore/         # DataStore (token, preferences)
│   │   └── entity/            # Room entities
│   └── repository/            # Repository implementations
├── di/                         # Hilt modules (Network, Database, Repository)
├── domain/
│   ├── model/                 # Domain models
│   └── usecase/               # Use cases
├── ui/
│   ├── auth/                  # Login, Register screens + ViewModel
│   ├── cashflow/              # Transaction, Wallet, Transfer screens + ViewModel
│   ├── components/            # Glass morphism reusable components
│   ├── dashboard/             # Dashboard screen + ViewModel
│   ├── goals/                 # Goals screen + ViewModel
│   ├── navigation/            # NavGraph, Screen sealed class
│   ├── overtime/              # Overtime list, form screens + ViewModel
│   ├── profile/               # Profile, Change Password screens + ViewModel
│   ├── reports/               # Reports screen + ViewModel
│   └── theme/                 # Color, Typography, Theme (dark/light)
├── FinanceKuApp.kt            # @HiltAndroidApp Application class
└── MainActivity.kt            # @AndroidEntryPoint entry point
```

## Prasyarat

- Android Studio Hedgehog (2023.1.1) atau lebih baru
- JDK 17
- Android SDK 34 (minimum SDK 26 / Android 8.0)
- Kotlin 1.9+

## Cara Build

### 1. Clone Repository

```bash
git clone https://github.com/dickyprase/financeku.git
cd financeku/android
```

### 2. Buka di Android Studio

1. Buka Android Studio
2. Pilih **File → Open**
3. Arahkan ke folder `financeku/android`
4. Tunggu Gradle sync selesai

### 3. Konfigurasi API Base URL

Edit file `app/src/main/java/com/financeku/app/di/NetworkModule.kt`:

```kotlin
// Ganti dengan URL backend kamu
private const val BASE_URL = "http://10.0.2.2:8080/" // untuk emulator (localhost)
// atau
private const val BASE_URL = "http://192.168.x.x:8080/" // untuk device fisik
```

> **Note**: `10.0.2.2` adalah alias untuk `localhost` host machine dari Android Emulator.

### 4. Build APK (Debug)

**Via Android Studio:**
- Klik **Build → Build Bundle(s) / APK(s) → Build APK(s)**
- APK akan tersedia di `app/build/outputs/apk/debug/app-debug.apk`

**Via Command Line:**

```bash
# Linux/macOS
./gradlew assembleDebug

# Windows
gradlew.bat assembleDebug
```

### 5. Build APK (Release)

```bash
# Generate signed APK
./gradlew assembleRelease
```

Untuk release build, kamu perlu setup signing config di `app/build.gradle.kts`:

```kotlin
android {
    signingConfigs {
        create("release") {
            storeFile = file("path/to/keystore.jks")
            storePassword = "your-store-password"
            keyAlias = "your-key-alias"
            keyPassword = "your-key-password"
        }
    }
    buildTypes {
        release {
            signingConfig = signingConfigs.getByName("release")
        }
    }
}
```

### 6. Run di Emulator/Device

**Via Android Studio:**
- Pilih device/emulator di toolbar
- Klik **Run** (▶️) atau tekan `Shift+F10`

**Via Command Line:**

```bash
# Install dan run di connected device
./gradlew installDebug
```

## Build Commands

| Command | Keterangan |
|---------|------------|
| `./gradlew assembleDebug` | Build debug APK |
| `./gradlew assembleRelease` | Build release APK |
| `./gradlew installDebug` | Install debug ke device |
| `./gradlew test` | Run unit tests |
| `./gradlew lint` | Run lint checks |
| `./gradlew clean` | Clean build artifacts |
| `./gradlew dependencies` | List all dependencies |

## Konfigurasi Backend

Pastikan backend sudah berjalan sebelum menjalankan app:

```bash
cd ../backend
npm install
npm run seed
npm run dev
# Server berjalan di http://localhost:8080
```

## Design System

### Glass Morphism

Komponen UI menggunakan efek glass morphism:
- Semi-transparent background dengan alpha
- Gradient overlay untuk depth effect
- Rounded corners
- Subtle border

### Dark/Light Mode

- Toggle via Profile screen
- Preference disimpan di DataStore
- Otomatis apply ke seluruh app

### Komponen Reusable

| Komponen | Kegunaan |
|----------|----------|
| `GlassCard` | Card container dengan glass effect |
| `GlassButton` | Button dengan glass background |
| `GlassTextField` | Input field dengan glass styling |
| `GlassTopBar` | Top app bar dengan glass effect |
| `LoadingIndicator` | Loading state indicator |
| `ErrorMessage` | Error state display |

## Minimum Requirements

- Android 8.0 (API 26) ke atas
- RAM minimal 2GB
- Koneksi internet untuk sync data

## License

MIT
