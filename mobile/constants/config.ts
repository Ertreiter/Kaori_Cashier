import Constants from 'expo-constants';
import { Platform } from 'react-native';

// Get backend URL - localhost for browser testing
const getBackendUrl = (): string => {
    return 'http://localhost:8080/api';
};

export const API_URL = getBackendUrl();
export const WS_URL = API_URL.replace('http', 'ws').replace('/api', '/ws');

// App configuration
export const APP_CONFIG = {
    appName: 'Kaori POS',
    version: '1.0.0',
    taxRate: 0.11, // 11% tax
    currency: 'Rp',
    locale: 'id-ID',
};

// Order sources
export const ORDER_SOURCES = {
    CASHIER: 'cashier',
    TABLE_QR: 'table_qr',
    GRABFOOD: 'grabfood',
    GOFOOD: 'gofood',
    SHOPEE_FOOD: 'shopee_food',
} as const;

export type OrderSource = typeof ORDER_SOURCES[keyof typeof ORDER_SOURCES];
