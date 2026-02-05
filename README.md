# Kaori POS System

A Point of Sale system for cafes and restaurants with React Native mobile app and Golang backend.

---

## 🚀 QUICK START (5 Minutes)

### Step 1: Clone & Enter Directory
```bash
git clone https://github.com/Ertreiter/Kaori_Cashier.git
cd Kaori_Cashier
```

### Step 2: Start Backend
```bash
cd backend
go run cmd/server/main.go
```
**Expected output:**
```
Loaded .env.development
🚀 Kaori POS API starting on port 8080 (dummy data mode)
```

### Step 3: Start Mobile App
In a **new terminal**:
```bash
cd Kaori_Cashier/mobile
npm install
npx expo start --lan
```
Scan QR with Expo Go app on your phone.

### Step 4: Configure API URL (IMPORTANT!)
Edit `mobile/constants/config.ts` line 8:
```typescript
// Replace YOUR_IP with your computer's IP (e.g. 192.168.1.100)
return `http://YOUR_IP:8080/api`;
```

Find your IP with: `ip addr` (Linux) or `ipconfig` (Windows)

---

## 🔐 Login Credentials

| Email | Password | Role |
|-------|----------|------|
| admin@kaori.pos | admin123 | Super Admin |
| cashier@kaori.pos | cashier123 | Cashier |
| kitchen@kaori.pos | kitchen123 | Kitchen |

---

## 📋 Features

- **Cashier**: Product grid, cart, table selection, cash payment
- **Kitchen**: Real-time order queue, status progression
- **Admin**: Dashboard, products, staff, order history
- **Delivery**: GrabFood, GoFood, Shopee Food integration (mock)

---

## 🚚 Test Delivery Orders

```bash
curl -X POST http://localhost:8080/api/simulate/order \
  -H "Content-Type: application/json" \
  -d '{"source":"grabfood","customer_name":"John","items":[{"name":"Latte","quantity":2,"price":28000}]}'
```

---

## ⚠️ Troubleshooting

### "Network Error" on Login
1. Check backend is running on port 8080
2. Update `mobile/constants/config.ts` with correct IP
3. Phone and computer must be on same WiFi

### Backend Won't Start
```bash
cd backend
go mod tidy
go run cmd/server/main.go
```

### White Screen on App
```bash
cd mobile
rm -rf node_modules .expo
npm install
npx expo start --clear
```
