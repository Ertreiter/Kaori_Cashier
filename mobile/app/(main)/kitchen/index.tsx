import React, { useState, useEffect, useCallback } from "react";
import {
  View,
  Text,
  StyleSheet,
  FlatList,
  TouchableOpacity,
  RefreshControl,
  ActivityIndicator,
  SafeAreaView,
  StatusBar,
  Platform,
} from "react-native";
import api from "../../../services/api";
import {
  Order,
  ORDER_SOURCE_CONFIG,
  isDeliveryOrder,
  OrderSource,
} from "../../../types";

const STATUS_COLORS: Record<string, string> = {
  pending: "#8B5CF6",
  confirmed: "#3B82F6",
  cooking: "#F59E0B",
  ready: "#10B981",
  completed: "#6B7280",
  cancelled: "#EF4444",
};

const STATUS_LABELS: Record<string, string> = {
  pending: "Pending",
  confirmed: "New",
  cooking: "Cooking",
  ready: "Ready",
  completed: "Done",
  cancelled: "Cancelled",
};

export default function KitchenScreen() {
  const [orders, setOrders] = useState<Order[]>([]);
  const [loading, setLoading] = useState(true);
  const [refreshing, setRefreshing] = useState(false);
  const [filterSource, setFilterSource] = useState<string | null>(null);

  const fetchOrders = useCallback(async () => {
    try {
      const res = await api.get("/orders/active");
      setOrders(res.data.data || []);
    } catch (error) {
      console.error("Failed to fetch orders:", error);
    } finally {
      setLoading(false);
      setRefreshing(false);
    }
  }, []);

  useEffect(() => {
    fetchOrders();
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
      console.error("Failed to update order:", error);
    }
  };

  const getNextStatus = (currentStatus: string): string | null => {
    switch (currentStatus) {
      case "pending":
        return "confirmed";
      case "confirmed":
        return "cooking";
      case "cooking":
        return "ready";
      case "ready":
        return "completed";
      default:
        return null;
    }
  };

  const formatTime = (dateString: string) => {
    const date = new Date(dateString);
    return date.toLocaleTimeString("id-ID", {
      hour: "2-digit",
      minute: "2-digit",
    });
  };

  const getTimeSince = (dateString: string) => {
    const diff = Date.now() - new Date(dateString).getTime();
    const minutes = Math.floor(diff / 60000);
    if (minutes < 1) return "Just now";
    if (minutes < 60) return `${minutes}m ago`;
    return `${Math.floor(minutes / 60)}h ${minutes % 60}m ago`;
  };

  const filteredOrders = filterSource
    ? orders.filter((o) => o.order_source === filterSource)
    : orders;

  const pendingOrders = filteredOrders.filter((o) => o.status === "pending");
  const confirmedOrders = filteredOrders.filter(
    (o) => o.status === "confirmed",
  );
  const cookingOrders = filteredOrders.filter((o) => o.status === "cooking");
  const readyOrders = filteredOrders.filter((o) => o.status === "ready");

  if (loading) {
    return (
      <SafeAreaView style={styles.loadingContainer}>
        <ActivityIndicator size="large" color="#F59E0B" />
        <Text style={styles.loadingText}>Loading orders...</Text>
      </SafeAreaView>
    );
  }

  const renderOrder = ({ item }: { item: Order }) => {
    const nextStatus = getNextStatus(item.status);
    const statusColor = STATUS_COLORS[item.status] || "#666";
    const sourceConfig =
      ORDER_SOURCE_CONFIG[item.order_source as OrderSource] ||
      ORDER_SOURCE_CONFIG.cashier;
    const isDelivery = isDeliveryOrder(item.order_source as OrderSource);

    return (
      <View style={styles.orderCard}>
        {/* Source Badge - Top Banner */}
        {isDelivery && (
          <View
            style={[
              styles.sourceBanner,
              { backgroundColor: sourceConfig.color },
            ]}
          >
            <Text style={styles.sourceBannerText}>
              {sourceConfig.emoji} {sourceConfig.label}
            </Text>
          </View>
        )}

        {/* Header */}
        <View style={styles.orderHeader}>
          <View>
            <Text style={styles.orderNumber}>{item.order_number}</Text>
            <Text style={styles.orderTime}>{formatTime(item.created_at)}</Text>
          </View>
          <View style={styles.orderMeta}>
            <View
              style={[styles.statusBadge, { backgroundColor: statusColor }]}
            >
              <Text style={styles.statusText}>
                {STATUS_LABELS[item.status]}
              </Text>
            </View>
            <Text style={styles.orderType}>
              {item.order_type === "dine_in" && `🍽️ Table ${item.table_number}`}
              {item.order_type === "takeaway" && "🥡 Takeaway"}
              {item.order_type === "delivery" && "🚚 Delivery"}
            </Text>
          </View>
        </View>

        {/* Customer Info for Delivery */}
        {isDelivery && item.customer_name && (
          <View style={styles.customerInfo}>
            <Text style={styles.customerName}>👤 {item.customer_name}</Text>
            {item.driver_name && (
              <Text style={styles.driverName}>
                🏍️ Driver: {item.driver_name}
              </Text>
            )}
          </View>
        )}

        {/* Items */}
        <View style={styles.orderItems}>
          {item.items.map((orderItem, index) => (
            <View key={orderItem.id || index} style={styles.itemRow}>
              <Text style={styles.itemQty}>{orderItem.quantity}x</Text>
              <View style={styles.itemDetails}>
                <Text style={styles.itemName}>{orderItem.product_name}</Text>
                {orderItem.variant_name && (
                  <Text style={styles.itemVariant}>
                    {orderItem.variant_name}
                  </Text>
                )}
                {orderItem.notes && (
                  <Text style={styles.itemNotes}>📝 {orderItem.notes}</Text>
                )}
              </View>
            </View>
          ))}
        </View>

        {/* Footer */}
        <View style={styles.orderFooter}>
          <Text style={styles.timerText}>
            ⏱️ {getTimeSince(item.created_at)}
          </Text>
          {nextStatus && (
            <TouchableOpacity
              style={[
                styles.actionBtn,
                { backgroundColor: STATUS_COLORS[nextStatus] },
              ]}
              onPress={() => updateOrderStatus(item.id, nextStatus)}
            >
              <Text style={styles.actionBtnText}>
                {nextStatus === "confirmed" && "✓ Accept"}
                {nextStatus === "cooking" && "🔥 Start Cooking"}
                {nextStatus === "ready" && "✅ Mark Ready"}
                {nextStatus === "completed" && "📦 Complete"}
              </Text>
            </TouchableOpacity>
          )}
        </View>
      </View>
    );
  };

  return (
    <SafeAreaView style={styles.container}>
      {/* Filter Bar */}
      <View style={styles.filterBar}>
        <TouchableOpacity
          style={[styles.filterBtn, !filterSource && styles.filterBtnActive]}
          onPress={() => setFilterSource(null)}
        >
          <Text
            style={[
              styles.filterText,
              !filterSource && styles.filterTextActive,
            ]}
          >
            All
          </Text>
        </TouchableOpacity>
        <TouchableOpacity
          style={[
            styles.filterBtn,
            filterSource === "cashier" && styles.filterBtnActive,
          ]}
          onPress={() =>
            setFilterSource(filterSource === "cashier" ? null : "cashier")
          }
        >
          <Text
            style={[
              styles.filterText,
              filterSource === "cashier" && styles.filterTextActive,
            ]}
          >
            💳 POS
          </Text>
        </TouchableOpacity>
        <TouchableOpacity
          style={[
            styles.filterBtn,
            filterSource === "grabfood" && styles.filterBtnActive,
            { borderColor: "#00B14F" },
          ]}
          onPress={() =>
            setFilterSource(filterSource === "grabfood" ? null : "grabfood")
          }
        >
          <Text
            style={[
              styles.filterText,
              filterSource === "grabfood" && styles.filterTextActive,
            ]}
          >
            🟢 Grab
          </Text>
        </TouchableOpacity>
        <TouchableOpacity
          style={[
            styles.filterBtn,
            filterSource === "gofood" && styles.filterBtnActive,
            { borderColor: "#D71920" },
          ]}
          onPress={() =>
            setFilterSource(filterSource === "gofood" ? null : "gofood")
          }
        >
          <Text
            style={[
              styles.filterText,
              filterSource === "gofood" && styles.filterTextActive,
            ]}
          >
            🔴 GoFood
          </Text>
        </TouchableOpacity>
        <TouchableOpacity
          style={[
            styles.filterBtn,
            filterSource === "shopee_food" && styles.filterBtnActive,
            { borderColor: "#EE4D2D" },
          ]}
          onPress={() =>
            setFilterSource(
              filterSource === "shopee_food" ? null : "shopee_food",
            )
          }
        >
          <Text
            style={[
              styles.filterText,
              filterSource === "shopee_food" && styles.filterTextActive,
            ]}
          >
            🧡 Shopee
          </Text>
        </TouchableOpacity>
      </View>

      {/* Stats Bar */}
      <View style={styles.statsBar}>
        <View style={styles.statItem}>
          <Text style={[styles.statNumber, { color: "#8B5CF6" }]}>
            {pendingOrders.length}
          </Text>
          <Text style={styles.statLabel}>Pending</Text>
        </View>
        <View style={styles.statItem}>
          <Text style={[styles.statNumber, { color: "#3B82F6" }]}>
            {confirmedOrders.length}
          </Text>
          <Text style={styles.statLabel}>New</Text>
        </View>
        <View style={styles.statItem}>
          <Text style={[styles.statNumber, { color: "#F59E0B" }]}>
            {cookingOrders.length}
          </Text>
          <Text style={styles.statLabel}>Cooking</Text>
        </View>
        <View style={styles.statItem}>
          <Text style={[styles.statNumber, { color: "#10B981" }]}>
            {readyOrders.length}
          </Text>
          <Text style={styles.statLabel}>Ready</Text>
        </View>
      </View>

      {/* Order List */}
      {filteredOrders.length === 0 ? (
        <View style={styles.emptyState}>
          <Text style={styles.emptyEmoji}>👨‍🍳</Text>
          <Text style={styles.emptyTitle}>No Active Orders</Text>
          <Text style={styles.emptySubtitle}>
            {filterSource
              ? `No ${ORDER_SOURCE_CONFIG[filterSource as OrderSource]?.label || filterSource} orders`
              : "New orders will appear here"}
          </Text>
        </View>
      ) : (
        <FlatList
          data={filteredOrders}
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
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: "#0a0a0a",
    paddingTop: Platform.OS === "android" ? StatusBar.currentHeight : 0,
  },
  loadingContainer: {
    flex: 1,
    justifyContent: "center",
    alignItems: "center",
    backgroundColor: "#0a0a0a",
    paddingTop: Platform.OS === "android" ? StatusBar.currentHeight : 0,
  },
  loadingText: { color: "#fff", marginTop: 10 },

  // Filter Bar
  filterBar: {
    flexDirection: "row",
    paddingHorizontal: 8,
    paddingVertical: 10,
    backgroundColor: "#111",
    gap: 6,
  },
  filterBtn: {
    paddingHorizontal: 12,
    paddingVertical: 6,
    borderRadius: 16,
    borderWidth: 1,
    borderColor: "#333",
    backgroundColor: "#1a1a1a",
  },
  filterBtnActive: {
    backgroundColor: "#333",
    borderColor: "#fff",
  },
  filterText: { color: "#888", fontSize: 12, fontWeight: "600" },
  filterTextActive: { color: "#fff" },

  // Stats Bar
  statsBar: {
    flexDirection: "row",
    justifyContent: "space-around",
    paddingVertical: 12,
    backgroundColor: "#111",
    borderBottomWidth: 1,
    borderBottomColor: "#222",
  },
  statItem: { alignItems: "center" },
  statNumber: { fontSize: 24, fontWeight: "bold" },
  statLabel: { color: "#888", fontSize: 11, marginTop: 2 },

  // Order List
  orderList: { padding: 8 },
  orderCard: {
    flex: 1,
    margin: 6,
    backgroundColor: "#1a1a1a",
    borderRadius: 12,
    overflow: "hidden",
    maxWidth: "48%",
  },

  // Source Banner
  sourceBanner: {
    paddingVertical: 6,
    alignItems: "center",
  },
  sourceBannerText: {
    color: "#fff",
    fontSize: 12,
    fontWeight: "bold",
  },

  // Order Header
  orderHeader: {
    flexDirection: "row",
    justifyContent: "space-between",
    alignItems: "flex-start",
    padding: 10,
    backgroundColor: "#222",
  },
  orderNumber: { color: "#fff", fontSize: 16, fontWeight: "bold" },
  orderTime: { color: "#888", fontSize: 11, marginTop: 2 },
  orderMeta: { alignItems: "flex-end" },
  statusBadge: {
    paddingHorizontal: 8,
    paddingVertical: 3,
    borderRadius: 10,
    marginBottom: 4,
  },
  statusText: { color: "#fff", fontSize: 10, fontWeight: "bold" },
  orderType: { color: "#888", fontSize: 11 },

  // Customer Info
  customerInfo: {
    paddingHorizontal: 10,
    paddingVertical: 6,
    backgroundColor: "#1f1f1f",
    borderBottomWidth: 1,
    borderBottomColor: "#333",
  },
  customerName: { color: "#fff", fontSize: 12, fontWeight: "600" },
  driverName: { color: "#888", fontSize: 11, marginTop: 2 },

  // Order Items
  orderItems: { padding: 10 },
  itemRow: { flexDirection: "row", marginBottom: 6 },
  itemQty: { color: "#F59E0B", fontSize: 14, fontWeight: "bold", width: 28 },
  itemDetails: { flex: 1 },
  itemName: { color: "#fff", fontSize: 13, fontWeight: "600" },
  itemVariant: { color: "#888", fontSize: 11 },
  itemNotes: {
    color: "#F59E0B",
    fontSize: 10,
    fontStyle: "italic",
    marginTop: 2,
  },

  // Footer
  orderFooter: {
    padding: 10,
    borderTopWidth: 1,
    borderTopColor: "#222",
  },
  timerText: { color: "#888", fontSize: 11, marginBottom: 6 },
  actionBtn: { paddingVertical: 10, borderRadius: 8, alignItems: "center" },
  actionBtnText: { color: "#fff", fontWeight: "bold", fontSize: 12 },

  // Empty State
  emptyState: { flex: 1, justifyContent: "center", alignItems: "center" },
  emptyEmoji: { fontSize: 60, marginBottom: 16 },
  emptyTitle: {
    color: "#fff",
    fontSize: 22,
    fontWeight: "bold",
    marginBottom: 8,
  },
  emptySubtitle: { color: "#888", fontSize: 14 },
});
