import React, { useState, useEffect } from 'react';
import {
    View,
    Text,
    StyleSheet,
    ScrollView,
    TouchableOpacity,
    FlatList,
    RefreshControl,
    Alert,
} from 'react-native';
import { useRouter } from 'expo-router';
import { useAuthStore } from '../../../stores/authStore';
import api from '../../../services/api';
import { Product, User as UserType, Order } from '../../../types';

type TabKey = 'dashboard' | 'products' | 'staff' | 'orders';

export default function AdminScreen() {
    const router = useRouter();
    const { user, logout } = useAuthStore();
    const [activeTab, setActiveTab] = useState<TabKey>('dashboard');
    const [refreshing, setRefreshing] = useState(false);

    // Data
    const [products, setProducts] = useState<Product[]>([]);
    const [users, setUsers] = useState<UserType[]>([]);
    const [orders, setOrders] = useState<Order[]>([]);
    const [stats, setStats] = useState({
        totalOrders: 0,
        totalRevenue: 0,
        activeOrders: 0,
    });

    useEffect(() => {
        fetchData();
    }, []);

    const fetchData = async () => {
        try {
            const [prodRes, userRes, orderRes, statsRes] = await Promise.all([
                api.get('/products'),
                api.get('/users'),
                api.get('/orders'),
                api.get('/stores/store-1/stats'),
            ]);
            setProducts(prodRes.data.data || []);
            setUsers(userRes.data.data || []);
            setOrders(orderRes.data.data || []);
            setStats(statsRes.data.data || {});
        } catch (error) {
            console.error('Failed to fetch data:', error);
        }
    };

    const handleRefresh = async () => {
        setRefreshing(true);
        await fetchData();
        setRefreshing(false);
    };

    const handleLogout = () => {
        Alert.alert('Logout', 'Are you sure you want to logout?', [
            { text: 'Cancel', style: 'cancel' },
            { text: 'Logout', style: 'destructive', onPress: logout },
        ]);
    };

    const formatPrice = (price: number) => `Rp ${price.toLocaleString('id-ID')}`;

    const renderDashboard = () => (
        <ScrollView style={styles.tabContent}>
            {/* Welcome */}
            <View style={styles.welcomeCard}>
                <Text style={styles.welcomeText}>Welcome, {user?.name || 'Admin'}</Text>
                <Text style={styles.roleText}>{user?.role?.replace('_', ' ').toUpperCase()}</Text>
            </View>

            {/* Stats */}
            <Text style={styles.sectionTitle}>Today's Summary</Text>
            <View style={styles.statsGrid}>
                <View style={styles.statCard}>
                    <Text style={styles.statValue}>{stats.totalOrders}</Text>
                    <Text style={styles.statLabel}>Total Orders</Text>
                </View>
                <View style={styles.statCard}>
                    <Text style={[styles.statValue, { color: '#10B981' }]}>
                        {formatPrice(stats.totalRevenue)}
                    </Text>
                    <Text style={styles.statLabel}>Revenue</Text>
                </View>
                <View style={styles.statCard}>
                    <Text style={[styles.statValue, { color: '#F59E0B' }]}>{stats.activeOrders}</Text>
                    <Text style={styles.statLabel}>Active Orders</Text>
                </View>
            </View>

            {/* Quick Actions */}
            <Text style={styles.sectionTitle}>Quick Actions</Text>
            <View style={styles.actionsGrid}>
                <TouchableOpacity style={styles.actionCard} onPress={() => setActiveTab('products')}>
                    <Text style={styles.actionEmoji}>📦</Text>
                    <Text style={styles.actionText}>Manage Products</Text>
                </TouchableOpacity>
                <TouchableOpacity style={styles.actionCard} onPress={() => setActiveTab('staff')}>
                    <Text style={styles.actionEmoji}>👥</Text>
                    <Text style={styles.actionText}>Manage Staff</Text>
                </TouchableOpacity>
                <TouchableOpacity style={styles.actionCard} onPress={() => setActiveTab('orders')}>
                    <Text style={styles.actionEmoji}>📋</Text>
                    <Text style={styles.actionText}>Order History</Text>
                </TouchableOpacity>
                <TouchableOpacity style={[styles.actionCard, { backgroundColor: '#EF4444' }]} onPress={handleLogout}>
                    <Text style={styles.actionEmoji}>🚪</Text>
                    <Text style={styles.actionText}>Logout</Text>
                </TouchableOpacity>
            </View>
        </ScrollView>
    );

    const renderProducts = () => (
        <FlatList
            data={products}
            keyExtractor={(item) => item.id}
            contentContainerStyle={styles.listContent}
            renderItem={({ item }) => (
                <View style={styles.listItem}>
                    <View style={styles.listItemIcon}>
                        <Text style={styles.listItemEmoji}>☕</Text>
                    </View>
                    <View style={styles.listItemInfo}>
                        <Text style={styles.listItemName}>{item.name}</Text>
                        <Text style={styles.listItemMeta}>
                            {item.variants.length} variants • {formatPrice(item.base_price)}
                        </Text>
                    </View>
                    <View style={[
                        styles.availabilityBadge,
                        { backgroundColor: item.is_available ? '#10B981' : '#EF4444' }
                    ]}>
                        <Text style={styles.availabilityText}>
                            {item.is_available ? 'Available' : 'Sold Out'}
                        </Text>
                    </View>
                </View>
            )}
            refreshControl={
                <RefreshControl refreshing={refreshing} onRefresh={handleRefresh} tintColor="#8B5CF6" />
            }
            ListEmptyComponent={
                <Text style={styles.emptyText}>No products found</Text>
            }
        />
    );

    const renderStaff = () => (
        <FlatList
            data={users}
            keyExtractor={(item) => item.id}
            contentContainerStyle={styles.listContent}
            renderItem={({ item }) => (
                <View style={styles.listItem}>
                    <View style={styles.avatarCircle}>
                        <Text style={styles.avatarText}>{item.name[0]}</Text>
                    </View>
                    <View style={styles.listItemInfo}>
                        <Text style={styles.listItemName}>{item.name}</Text>
                        <Text style={styles.listItemMeta}>{item.email}</Text>
                    </View>
                    <View style={styles.roleBadge}>
                        <Text style={styles.roleText2}>{item.role.replace('_', ' ')}</Text>
                    </View>
                </View>
            )}
            refreshControl={
                <RefreshControl refreshing={refreshing} onRefresh={handleRefresh} tintColor="#8B5CF6" />
            }
            ListEmptyComponent={
                <Text style={styles.emptyText}>No staff found</Text>
            }
        />
    );

    const renderOrders = () => (
        <FlatList
            data={orders}
            keyExtractor={(item) => item.id}
            contentContainerStyle={styles.listContent}
            renderItem={({ item }) => (
                <View style={styles.listItem}>
                    <View style={styles.orderInfo}>
                        <Text style={styles.listItemName}>{item.order_number}</Text>
                        <Text style={styles.listItemMeta}>
                            {item.items.length} items • {formatPrice(item.total)}
                        </Text>
                    </View>
                    <View>
                        <View style={[styles.statusPill, { backgroundColor: getStatusColor(item.status) }]}>
                            <Text style={styles.statusPillText}>{item.status}</Text>
                        </View>
                        <Text style={styles.orderTime}>
                            {new Date(item.created_at).toLocaleTimeString('id-ID')}
                        </Text>
                    </View>
                </View>
            )}
            refreshControl={
                <RefreshControl refreshing={refreshing} onRefresh={handleRefresh} tintColor="#8B5CF6" />
            }
            ListEmptyComponent={
                <Text style={styles.emptyText}>No orders yet</Text>
            }
        />
    );

    const getStatusColor = (status: string) => {
        const colors: Record<string, string> = {
            confirmed: '#3B82F6',
            cooking: '#F59E0B',
            ready: '#10B981',
            completed: '#6B7280',
            cancelled: '#EF4444',
        };
        return colors[status] || '#666';
    };

    const tabs: { key: TabKey; label: string; icon: string }[] = [
        { key: 'dashboard', label: 'Dashboard', icon: '📊' },
        { key: 'products', label: 'Products', icon: '📦' },
        { key: 'staff', label: 'Staff', icon: '👥' },
        { key: 'orders', label: 'Orders', icon: '📋' },
    ];

    return (
        <View style={styles.container}>
            {/* Tab Bar */}
            <View style={styles.tabBar}>
                {tabs.map((tab) => (
                    <TouchableOpacity
                        key={tab.key}
                        style={[styles.tab, activeTab === tab.key && styles.tabActive]}
                        onPress={() => setActiveTab(tab.key)}
                    >
                        <Text style={styles.tabIcon}>{tab.icon}</Text>
                        <Text style={[styles.tabLabel, activeTab === tab.key && styles.tabLabelActive]}>
                            {tab.label}
                        </Text>
                    </TouchableOpacity>
                ))}
            </View>

            {/* Content */}
            {activeTab === 'dashboard' && renderDashboard()}
            {activeTab === 'products' && renderProducts()}
            {activeTab === 'staff' && renderStaff()}
            {activeTab === 'orders' && renderOrders()}
        </View>
    );
}

const styles = StyleSheet.create({
    container: {
        flex: 1,
        backgroundColor: '#0a0a0a',
    },

    // Tab Bar
    tabBar: {
        flexDirection: 'row',
        backgroundColor: '#111',
        borderBottomWidth: 1,
        borderBottomColor: '#222',
    },
    tab: {
        flex: 1,
        paddingVertical: 12,
        alignItems: 'center',
    },
    tabActive: {
        borderBottomWidth: 2,
        borderBottomColor: '#8B5CF6',
    },
    tabIcon: {
        fontSize: 20,
        marginBottom: 4,
    },
    tabLabel: {
        color: '#666',
        fontSize: 12,
    },
    tabLabelActive: {
        color: '#8B5CF6',
        fontWeight: 'bold',
    },

    // Tab Content
    tabContent: {
        flex: 1,
        padding: 16,
    },

    // Welcome Card
    welcomeCard: {
        backgroundColor: '#8B5CF6',
        borderRadius: 12,
        padding: 20,
        marginBottom: 20,
    },
    welcomeText: {
        color: '#fff',
        fontSize: 24,
        fontWeight: 'bold',
    },
    roleText: {
        color: 'rgba(255,255,255,0.8)',
        fontSize: 14,
        marginTop: 4,
    },

    // Section Title
    sectionTitle: {
        color: '#fff',
        fontSize: 18,
        fontWeight: 'bold',
        marginBottom: 12,
        marginTop: 8,
    },

    // Stats Grid
    statsGrid: {
        flexDirection: 'row',
        gap: 12,
        marginBottom: 20,
    },
    statCard: {
        flex: 1,
        backgroundColor: '#1a1a1a',
        borderRadius: 12,
        padding: 16,
        alignItems: 'center',
    },
    statValue: {
        color: '#8B5CF6',
        fontSize: 24,
        fontWeight: 'bold',
    },
    statLabel: {
        color: '#888',
        fontSize: 12,
        marginTop: 4,
    },

    // Actions Grid
    actionsGrid: {
        flexDirection: 'row',
        flexWrap: 'wrap',
        gap: 12,
    },
    actionCard: {
        width: '48%',
        backgroundColor: '#1a1a1a',
        borderRadius: 12,
        padding: 20,
        alignItems: 'center',
    },
    actionEmoji: {
        fontSize: 32,
        marginBottom: 8,
    },
    actionText: {
        color: '#fff',
        fontSize: 14,
        fontWeight: '600',
    },

    // List
    listContent: {
        padding: 16,
    },
    listItem: {
        flexDirection: 'row',
        alignItems: 'center',
        backgroundColor: '#1a1a1a',
        borderRadius: 12,
        padding: 14,
        marginBottom: 10,
    },
    listItemIcon: {
        width: 44,
        height: 44,
        borderRadius: 22,
        backgroundColor: '#333',
        justifyContent: 'center',
        alignItems: 'center',
        marginRight: 12,
    },
    listItemEmoji: {
        fontSize: 20,
    },
    listItemInfo: {
        flex: 1,
    },
    listItemName: {
        color: '#fff',
        fontSize: 15,
        fontWeight: '600',
    },
    listItemMeta: {
        color: '#888',
        fontSize: 12,
        marginTop: 2,
    },
    availabilityBadge: {
        paddingHorizontal: 10,
        paddingVertical: 4,
        borderRadius: 12,
    },
    availabilityText: {
        color: '#fff',
        fontSize: 11,
        fontWeight: 'bold',
    },

    // Staff
    avatarCircle: {
        width: 44,
        height: 44,
        borderRadius: 22,
        backgroundColor: '#8B5CF6',
        justifyContent: 'center',
        alignItems: 'center',
        marginRight: 12,
    },
    avatarText: {
        color: '#fff',
        fontSize: 18,
        fontWeight: 'bold',
    },
    roleBadge: {
        backgroundColor: '#333',
        paddingHorizontal: 10,
        paddingVertical: 4,
        borderRadius: 12,
    },
    roleText2: {
        color: '#8B5CF6',
        fontSize: 11,
        fontWeight: 'bold',
        textTransform: 'capitalize',
    },

    // Orders
    orderInfo: {
        flex: 1,
    },
    statusPill: {
        paddingHorizontal: 10,
        paddingVertical: 4,
        borderRadius: 12,
        alignSelf: 'flex-end',
    },
    statusPillText: {
        color: '#fff',
        fontSize: 11,
        fontWeight: 'bold',
        textTransform: 'capitalize',
    },
    orderTime: {
        color: '#666',
        fontSize: 11,
        marginTop: 4,
        textAlign: 'right',
    },

    // Empty
    emptyText: {
        color: '#666',
        textAlign: 'center',
        marginTop: 40,
    },
});
