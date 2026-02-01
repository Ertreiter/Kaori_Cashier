import React, { useState, useEffect, useCallback } from 'react';
import {
    View,
    Text,
    StyleSheet,
    FlatList,
    TouchableOpacity,
    RefreshControl,
    ActivityIndicator,
} from 'react-native';
import api from '../../../services/api';
import { Order } from '../../../types';

const STATUS_COLORS = {
    confirmed: '#3B82F6',
    cooking: '#F59E0B',
    ready: '#10B981',
    completed: '#6B7280',
    cancelled: '#EF4444',
};

const STATUS_LABELS = {
    confirmed: 'New',
    cooking: 'Cooking',
    ready: 'Ready',
    completed: 'Done',
    cancelled: 'Cancelled',
};

export default function KitchenScreen() {
    const [orders, setOrders] = useState<Order[]>([]);
    const [loading, setLoading] = useState(true);
    const [refreshing, setRefreshing] = useState(false);

    const fetchOrders = useCallback(async () => {
        try {
            const res = await api.get('/orders/active');
            setOrders(res.data.data || []);
        } catch (error) {
            console.error('Failed to fetch orders:', error);
        } finally {
            setLoading(false);
            setRefreshing(false);
        }
    }, []);

    useEffect(() => {
        fetchOrders();
        // Poll for updates every 5 seconds
        const interval = setInterval(fetchOrders, 5000);
        return () => clearInterval(interval);
    }, [fetchOrders]);

    const handleRefresh = () => {
        setRefreshing(true);
        fetchOrders();
    };

    const updateOrderStatus = async (orderId: string, newStatus: string) => {
        try {
            await api.patch(`/orders/${orderId}/status`, { status: newStatus });
            fetchOrders();
        } catch (error) {
            console.error('Failed to update order:', error);
        }
    };

    const getNextStatus = (currentStatus: string): string | null => {
        switch (currentStatus) {
            case 'confirmed': return 'cooking';
            case 'cooking': return 'ready';
            case 'ready': return 'completed';
            default: return null;
        }
    };

    const formatTime = (dateString: string) => {
        const date = new Date(dateString);
        return date.toLocaleTimeString('id-ID', { hour: '2-digit', minute: '2-digit' });
    };

    const getTimeSince = (dateString: string) => {
        const diff = Date.now() - new Date(dateString).getTime();
        const minutes = Math.floor(diff / 60000);
        if (minutes < 1) return 'Just now';
        if (minutes < 60) return `${minutes}m ago`;
        return `${Math.floor(minutes / 60)}h ${minutes % 60}m ago`;
    };

    if (loading) {
        return (
            <View style={styles.loadingContainer}>
                <ActivityIndicator size="large" color="#F59E0B" />
                <Text style={styles.loadingText}>Loading orders...</Text>
            </View>
        );
    }

    const renderOrder = ({ item }: { item: Order }) => {
        const nextStatus = getNextStatus(item.status);
        const statusColor = STATUS_COLORS[item.status as keyof typeof STATUS_COLORS] || '#666';

        return (
            <View style={styles.orderCard}>
                {/* Header */}
                <View style={styles.orderHeader}>
                    <View>
                        <Text style={styles.orderNumber}>{item.order_number}</Text>
                        <Text style={styles.orderTime}>{formatTime(item.created_at)}</Text>
                    </View>
                    <View style={styles.orderMeta}>
                        <View style={[styles.statusBadge, { backgroundColor: statusColor }]}>
                            <Text style={styles.statusText}>
                                {STATUS_LABELS[item.status as keyof typeof STATUS_LABELS]}
                            </Text>
                        </View>
                        <Text style={styles.orderType}>
                            {item.order_type === 'dine_in' ? `Table ${item.table_number}` : '🥡 Takeaway'}
                        </Text>
                    </View>
                </View>

                {/* Items */}
                <View style={styles.orderItems}>
                    {item.items.map((orderItem, index) => (
                        <View key={orderItem.id || index} style={styles.itemRow}>
                            <Text style={styles.itemQty}>{orderItem.quantity}x</Text>
                            <View style={styles.itemDetails}>
                                <Text style={styles.itemName}>{orderItem.product_name}</Text>
                                {orderItem.variant_name && (
                                    <Text style={styles.itemVariant}>{orderItem.variant_name}</Text>
                                )}
                                {orderItem.notes && (
                                    <Text style={styles.itemNotes}>📝 {orderItem.notes}</Text>
                                )}
                            </View>
                        </View>
                    ))}
                </View>

                {/* Timer */}
                <View style={styles.orderFooter}>
                    <Text style={styles.timerText}>⏱️ {getTimeSince(item.created_at)}</Text>

                    {/* Action Button */}
                    {nextStatus && (
                        <TouchableOpacity
                            style={[styles.actionBtn, { backgroundColor: STATUS_COLORS[nextStatus as keyof typeof STATUS_COLORS] }]}
                            onPress={() => updateOrderStatus(item.id, nextStatus)}
                        >
                            <Text style={styles.actionBtnText}>
                                {nextStatus === 'cooking' && '🔥 Start Cooking'}
                                {nextStatus === 'ready' && '✅ Mark Ready'}
                                {nextStatus === 'completed' && '📦 Complete'}
                            </Text>
                        </TouchableOpacity>
                    )}
                </View>
            </View>
        );
    };

    // Group orders by status
    const confirmedOrders = orders.filter(o => o.status === 'confirmed');
    const cookingOrders = orders.filter(o => o.status === 'cooking');
    const readyOrders = orders.filter(o => o.status === 'ready');

    return (
        <View style={styles.container}>
            {/* Header Stats */}
            <View style={styles.statsBar}>
                <View style={styles.statItem}>
                    <Text style={[styles.statNumber, { color: '#3B82F6' }]}>{confirmedOrders.length}</Text>
                    <Text style={styles.statLabel}>New</Text>
                </View>
                <View style={styles.statItem}>
                    <Text style={[styles.statNumber, { color: '#F59E0B' }]}>{cookingOrders.length}</Text>
                    <Text style={styles.statLabel}>Cooking</Text>
                </View>
                <View style={styles.statItem}>
                    <Text style={[styles.statNumber, { color: '#10B981' }]}>{readyOrders.length}</Text>
                    <Text style={styles.statLabel}>Ready</Text>
                </View>
            </View>

            {/* Order List */}
            {orders.length === 0 ? (
                <View style={styles.emptyState}>
                    <Text style={styles.emptyEmoji}>👨‍🍳</Text>
                    <Text style={styles.emptyTitle}>No Active Orders</Text>
                    <Text style={styles.emptySubtitle}>New orders will appear here</Text>
                </View>
            ) : (
                <FlatList
                    data={orders}
                    keyExtractor={(item) => item.id}
                    renderItem={renderOrder}
                    contentContainerStyle={styles.orderList}
                    numColumns={2}
                    refreshControl={
                        <RefreshControl
                            refreshing={refreshing}
                            onRefresh={handleRefresh}
                            tintColor="#F59E0B"
                        />
                    }
                />
            )}
        </View>
    );
}

const styles = StyleSheet.create({
    container: {
        flex: 1,
        backgroundColor: '#0a0a0a',
    },
    loadingContainer: {
        flex: 1,
        justifyContent: 'center',
        alignItems: 'center',
        backgroundColor: '#0a0a0a',
    },
    loadingText: {
        color: '#fff',
        marginTop: 10,
    },

    // Stats Bar
    statsBar: {
        flexDirection: 'row',
        justifyContent: 'space-around',
        padding: 16,
        backgroundColor: '#111',
        borderBottomWidth: 1,
        borderBottomColor: '#222',
    },
    statItem: {
        alignItems: 'center',
    },
    statNumber: {
        fontSize: 28,
        fontWeight: 'bold',
    },
    statLabel: {
        color: '#888',
        fontSize: 12,
        marginTop: 2,
    },

    // Order List
    orderList: {
        padding: 8,
    },
    orderCard: {
        flex: 1,
        margin: 8,
        backgroundColor: '#1a1a1a',
        borderRadius: 12,
        overflow: 'hidden',
        maxWidth: '48%',
    },
    orderHeader: {
        flexDirection: 'row',
        justifyContent: 'space-between',
        alignItems: 'flex-start',
        padding: 12,
        backgroundColor: '#222',
    },
    orderNumber: {
        color: '#fff',
        fontSize: 18,
        fontWeight: 'bold',
    },
    orderTime: {
        color: '#888',
        fontSize: 12,
        marginTop: 2,
    },
    orderMeta: {
        alignItems: 'flex-end',
    },
    statusBadge: {
        paddingHorizontal: 10,
        paddingVertical: 4,
        borderRadius: 12,
        marginBottom: 4,
    },
    statusText: {
        color: '#fff',
        fontSize: 11,
        fontWeight: 'bold',
    },
    orderType: {
        color: '#888',
        fontSize: 12,
    },

    // Order Items
    orderItems: {
        padding: 12,
    },
    itemRow: {
        flexDirection: 'row',
        marginBottom: 8,
    },
    itemQty: {
        color: '#F59E0B',
        fontSize: 16,
        fontWeight: 'bold',
        width: 30,
    },
    itemDetails: {
        flex: 1,
    },
    itemName: {
        color: '#fff',
        fontSize: 14,
        fontWeight: '600',
    },
    itemVariant: {
        color: '#888',
        fontSize: 12,
    },
    itemNotes: {
        color: '#F59E0B',
        fontSize: 11,
        fontStyle: 'italic',
        marginTop: 2,
    },

    // Footer
    orderFooter: {
        padding: 12,
        borderTopWidth: 1,
        borderTopColor: '#222',
    },
    timerText: {
        color: '#888',
        fontSize: 12,
        marginBottom: 8,
    },
    actionBtn: {
        paddingVertical: 12,
        borderRadius: 8,
        alignItems: 'center',
    },
    actionBtnText: {
        color: '#fff',
        fontWeight: 'bold',
        fontSize: 14,
    },

    // Empty State
    emptyState: {
        flex: 1,
        justifyContent: 'center',
        alignItems: 'center',
    },
    emptyEmoji: {
        fontSize: 60,
        marginBottom: 16,
    },
    emptyTitle: {
        color: '#fff',
        fontSize: 22,
        fontWeight: 'bold',
        marginBottom: 8,
    },
    emptySubtitle: {
        color: '#888',
        fontSize: 14,
    },
});
