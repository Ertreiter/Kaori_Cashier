import Constants from 'expo-constants';
import { Platform } from 'react-native';

// Get backend URL based on environment
const getBackendUrl = (): string => {
    // In development, try to auto-detect from Expo host
    if (__DEV__) {
        const hostUri = Constants.expoConfig?.hostUri;
        if (hostUri) {
            const host = hostUri.split(':')[0];
            return `http://${host}:8080/api`;
        }

        // Android emulator uses special IP
        if (Platform.OS === 'android') {
            return 'http://10.0.2.2:8080/api';
        }

        // iOS simulator uses localhost
        return 'http://localhost:8080/api';
    }

    // Production - change this to your server URL
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
