import { Redirect } from 'expo-router';
import { useAuthStore } from '../stores/authStore';
import { View, ActivityIndicator } from 'react-native';

export default function Index() {
    const { token, user, isLoading } = useAuthStore();

    if (isLoading) {
        return (
            <View style={{ flex: 1, justifyContent: 'center', alignItems: 'center' }}>
                <ActivityIndicator size="large" color="#4a90e2" />
            </View>
        );
    }

    if (token && user) {
        if (user.role === 'cashier') {
            return <Redirect href="/(main)/cashier" />;
        } else if (user.role === 'kitchen') {
            return <Redirect href="/(main)/kitchen" />;
        }
        return <Redirect href="/(main)/admin" />;
    }

    return <Redirect href="/(auth)/login" />;
}
