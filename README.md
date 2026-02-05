# Kaori POS System

A Point of Sale system for cafes and restaurants with React Native mobile app and Golang backend.

## 🚀 Quick Start

### Prerequisites
- **Node.js 18+**
- **Go 1.21+**
- **Android Studio** (for building APK)

---

## 📱 Mobile App Setup

### 1. Install Dependencies
```bash
cd mobile
npm install
```

### 2. Configure Backend URL
Edit `mobile/constants/config.ts`:
```typescript
// Replace with your backend server IP
export const API_URL = 'http://YOUR_IP_ADDRESS:8080/api';
```

### 3. Build APK (Android)

#### Option A: Build with Android Studio
1. Open Android Studio
2. Open project: `mobile/android`
3. Build > Build Bundle(s) / APK(s) > Build APK(s)
4. APK location: `mobile/android/app/build/outputs/apk/release/`

#### Option B: Build via Command Line
```bash
cd mobile

# Generate native Android project
npx expo prebuild --platform android

# Build the APK
cd android
./gradlew assembleRelease

# APK will be at:
# android/app/build/outputs/apk/release/app-release.apk
```

> **Note:** Make sure `ANDROID_HOME` is set to your Android SDK path.

### 4. Install APK on Device
```bash
adb install mobile/android/app/build/outputs/apk/release/app-release.apk
```

---

## 🖥️ Backend Setup

### 1. Configure Environment
```bash
cd backend
cp .env.example .env
```

Edit `.env` with your settings (defaults work for dummy data mode).

### 2. Run Backend
```bash
cd backend
go run cmd/server/main.go
```

The server will start on `http://localhost:8080`

---

## � Login Credentials

| Email | Password | PIN | Role |
|-------|----------|-----|------|
| admin@kaori.pos | admin123 | 1234 | Super Admin |
| store@kaori.pos | store123 | 5678 | Store Admin |
| cashier@kaori.pos | cashier123 | 1111 | Cashier |
| kitchen@kaori.pos | kitchen123 | 2222 | Kitchen |

---

## 📋 Features

### Cashier Module
- Product grid with categories
- Cart management
- Table selection (dine-in/takeaway)
- Cash payment processing

### Kitchen Module
- Real-time order queue
- Status progression (New → Cooking → Ready → Done)
- Order timing display

### Admin Module
- Dashboard with stats
- Products management
- Staff management
- Order history

---

## 🛠️ Tech Stack

- **Backend**: Go, Gin, PostgreSQL (with dummy data mode)
- **Mobile**: React Native, Expo, Zustand
- **Real-time**: WebSocket for kitchen updates

---

## 🚚 Delivery Platform Integration

The system supports incoming orders from food delivery platforms:

| Platform | Webhook Endpoint | Order Prefix |
|----------|------------------|--------------|
| GrabFood | `POST /api/webhooks/grabfood` | GRAB-xxxx |
| GoFood | `POST /api/webhooks/gofood` | GOFOOD-xxxx |
| Shopee Food | `POST /api/webhooks/shopee` | SHOPEE-xxxx |

### Testing Delivery Orders

Use the simulate endpoint to test incoming delivery orders:

```bash
curl -X POST http://localhost:8080/api/simulate/order \
  -H "Content-Type: application/json" \
  -d '{
    "source": "grabfood",
    "customer_name": "John Doe",
    "items": [
      {"name": "Latte", "quantity": 2, "price": 28000},
      {"name": "Croissant", "quantity": 1, "price": 25000}
    ]
  }'
```

Sources: `cashier`, `table_qr`, `grabfood`, `gofood`, `shopee_food`

---

## ⚠️ Troubleshooting

### White Screen on App Start
```bash
cd mobile
rm -rf node_modules
npm install
npx expo start --clear
```

### Gradle Build Fails
1. Make sure Android SDK is installed
2. Set `ANDROID_HOME` environment variable:
   ```bash
   export ANDROID_HOME=$HOME/Android/Sdk
   export PATH=$PATH:$ANDROID_HOME/tools:$ANDROID_HOME/platform-tools
   ```

### Cannot Connect to Backend
- Ensure backend is running (`go run cmd/server/main.go`)
- Check the API_URL in `mobile/constants/config.ts` matches your backend IP
- Phone and computer must be on same network
