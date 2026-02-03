# Kaori POS System

A Point of Sale system for cafes and restaurants with React Native mobile app and Golang backend.

## 🚀 Quick Start

### Prerequisites
- Node.js 18+
- Go 1.21+
- Expo Go app on your phone

### Setup

#### 1. Backend
```bash
cd backend
cp .env.example .env
# Edit .env with your settings (or use defaults for dummy data)
go run cmd/server/main.go
```

The backend will start on `http://localhost:8080`

#### 2. Mobile App
```bash
cd mobile
npm install
npx expo start --lan
```

Scan the QR code with Expo Go (make sure phone and computer are on same WiFi).

### 📱 Login Credentials

| Email | Password | Role |
|-------|----------|------|
| admin@kaori.pos | admin123 | Super Admin |
| store@kaori.pos | store123 | Store Admin |
| cashier@kaori.pos | cashier123 | Cashier |
| kitchen@kaori.pos | kitchen123 | Kitchen |

### ⚠️ Important: Configure API URL

Before running the mobile app, update the backend URL in `mobile/constants/config.ts`:

```typescript
export const API_URL = 'http://YOUR_COMPUTER_IP:8080/api';
```

Replace `YOUR_COMPUTER_IP` with your computer's local IP address (e.g., `192.168.1.100`).

## Features

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

## Tech Stack

- **Backend**: Go, Gin, PostgreSQL (with dummy data mode)
- **Mobile**: React Native, Expo Router, Zustand
- **Real-time**: WebSocket for kitchen updates
