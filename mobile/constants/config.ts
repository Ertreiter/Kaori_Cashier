import Constants from 'expo-constants';

const getBackendUrl = () => {
    // For Android Emulator
    // return 'http://10.0.2.2:8080/api';
    
    // For physical device, replace with your machine's IP
    // return 'http://192.168.1.X:8080/api';

    // Development fallback
    if (__DEV__) {
        // Try to infer from debugger host (only works in some Expo contexts)
        const debuggerHost = Constants.expoConfig?.hostUri;
        if (debuggerHost) {
            return `http://${debuggerHost.split(':')[0]}:8080/api`;
        }
    }
    
    // Default fallback (Android Emulator)
    return 'http://10.0.2.2:8080/api';
};

export const API_URL = getBackendUrl();
export const WS_URL = API_URL.replace('http', 'ws').replace('/api', '/ws');
