import { create } from 'zustand';
import { createJSONStorage, persist } from 'zustand/middleware';
import AsyncStorage from '@react-native-async-storage/async-storage';
import { router } from 'expo-router';

export interface User {
  id: string;
  email: string;
  name: string;
  role: 'cashier' | 'kitchen' | 'store_admin' | 'super_admin';
  store_id?: string;
}

interface AuthState {
  user: User | null;
  token: string | null;
  refreshToken: string | null;
  isLoading: boolean;
  login: (user: User, token: string, refreshToken: string) => void;
  logout: () => void;
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set) => ({
      user: null,
      token: null,
      refreshToken: null,
      isLoading: false,
      login: (user, token, refreshToken) => {
        set({ user, token, refreshToken });
        if (user.role === 'cashier') {
          router.replace('/(main)/cashier');
        } else if (user.role === 'kitchen') {
          router.replace('/(main)/kitchen');
        } else {
          router.replace('/(main)/admin');
        }
      },
      logout: () => {
        set({ user: null, token: null, refreshToken: null });
        router.replace('/(auth)/login');
      },
    }),
    {
      name: 'auth-storage',
      storage: createJSONStorage(() => AsyncStorage),
    }
  )
);
