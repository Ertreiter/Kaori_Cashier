import React, { useState, useEffect } from "react";
import {
  View,
  Text,
  StyleSheet,
  FlatList,
  TouchableOpacity,
  ScrollView,
  Modal,
  TextInput,
  Alert,
  ActivityIndicator,
  SafeAreaView,
  StatusBar,
  Platform,
} from "react-native";
import { useRouter } from "expo-router";
import { useAuthStore } from "../../../stores/authStore";
import { useCartStore, CartItem } from "../../../stores/cartStore";
import api from "../../../services/api";
import { Category, Product, Table } from "../../../types";

export default function CashierScreen() {
  const router = useRouter();
  const { user, logout } = useAuthStore();
  const cart = useCartStore();

  // Data
  const [categories, setCategories] = useState<Category[]>([]);
  const [products, setProducts] = useState<Product[]>([]);
  const [tables, setTables] = useState<Table[]>([]);
  const [loading, setLoading] = useState(true);

  // UI State
  const [selectedCategory, setSelectedCategory] = useState<string>("all");
  const [showTableModal, setShowTableModal] = useState(false);
  const [showPaymentModal, setShowPaymentModal] = useState(false);
  const [showProductModal, setShowProductModal] = useState(false);
  const [selectedProduct, setSelectedProduct] = useState<Product | null>(null);
  const [cashAmount, setCashAmount] = useState("");
  const [processing, setProcessing] = useState(false);

  // Fetch data
  useEffect(() => {
    fetchData();
  }, []);

  const fetchData = async () => {
    try {
      setLoading(true);
      const [catRes, prodRes, tableRes] = await Promise.all([
        api.get("/categories"),
        api.get("/products"),
        api.get("/tables"),
      ]);
      setCategories(catRes.data.data || []);
      setProducts(prodRes.data.data || []);
      setTables(tableRes.data.data || []);
    } catch (error) {
      console.error("Failed to fetch data:", error);
    } finally {
      setLoading(false);
    }
  };

  // Filter products
  const filteredProducts =
    selectedCategory === "all"
      ? products.filter((p) => p.is_available)
      : products.filter(
          (p) => p.category_id === selectedCategory && p.is_available,
        );

  // Add product to cart
  const handleProductPress = (product: Product) => {
    if (product.variants.length > 0 || product.modifiers.length > 0) {
      setSelectedProduct(product);
      setShowProductModal(true);
    } else {
      cart.addItem({
        productId: product.id,
        productName: product.name,
        quantity: 1,
        unitPrice: product.base_price,
        modifiers: [],
      });
    }
  };

  // Add product with variant
  const handleAddWithVariant = (product: Product, variantId?: string) => {
    const variant = product.variants.find((v) => v.id === variantId);
    const price = product.base_price + (variant?.price_adjustment || 0);

    cart.addItem({
      productId: product.id,
      productName: product.name,
      variantId,
      variantName: variant?.name,
      quantity: 1,
      unitPrice: price,
      modifiers: [],
    });
    setShowProductModal(false);
    setSelectedProduct(null);
  };

  // Submit order
  const handleSubmitOrder = async () => {
    if (cart.items.length === 0) {
      Alert.alert("Error", "Cart is empty");
      return;
    }

    if (cart.orderType === "dine_in" && !cart.tableId) {
      setShowTableModal(true);
      return;
    }

    setShowPaymentModal(true);
  };

  // Process payment
  const handlePayment = async () => {
    const total = cart.getTotal();
    const paid = parseInt(cashAmount) || 0;

    if (paid < total) {
      Alert.alert("Error", "Insufficient payment amount");
      return;
    }

    setProcessing(true);
    try {
      // Create order
      const orderRes = await api.post("/orders", {
        order_type: cart.orderType,
        table_id: cart.tableId,
        items: cart.items.map((item) => ({
          product_id: item.productId,
          variant_id: item.variantId,
          modifiers: item.modifiers,
          quantity: item.quantity,
          notes: item.notes,
        })),
        notes: cart.notes,
      });

      const order = orderRes.data.data;

      // Process payment
      await api.post("/payments/cash", {
        order_id: order.id,
        amount_paid: paid,
      });

      const change = paid - total;
      Alert.alert(
        "✅ Payment Successful",
        `Order: ${order.order_number}\nTotal: Rp ${total.toLocaleString()}\nPaid: Rp ${paid.toLocaleString()}\nChange: Rp ${change.toLocaleString()}`,
        [{ text: "OK", onPress: () => cart.clearCart() }],
      );
      setShowPaymentModal(false);
      setCashAmount("");
    } catch (error) {
      console.error("Payment error:", error);
      Alert.alert("Error", "Failed to process payment");
    } finally {
      setProcessing(false);
    }
  };

  // Format currency
  const formatPrice = (price: number | undefined) =>
    `Rp ${(price || 0).toLocaleString("id-ID")}`;

  if (loading) {
    return (
      <SafeAreaView style={styles.loadingContainer}>
        <ActivityIndicator size="large" color="#8B5CF6" />
        <Text style={styles.loadingText}>Loading...</Text>
      </SafeAreaView>
    );
  }

  return (
    <SafeAreaView style={styles.container}>
      {/* Left Side - Products */}
      <View style={styles.leftPanel}>
        {/* Categories */}
        <ScrollView
          horizontal
          showsHorizontalScrollIndicator={false}
          style={styles.categoryBar}
        >
          <TouchableOpacity
            style={[
              styles.categoryChip,
              selectedCategory === "all" && styles.categoryChipActive,
            ]}
            onPress={() => setSelectedCategory("all")}
          >
            <Text
              style={[
                styles.categoryText,
                selectedCategory === "all" && styles.categoryTextActive,
              ]}
            >
              All
            </Text>
          </TouchableOpacity>
          {categories.map((cat) => (
            <TouchableOpacity
              key={cat.id}
              style={[
                styles.categoryChip,
                selectedCategory === cat.id && styles.categoryChipActive,
              ]}
              onPress={() => setSelectedCategory(cat.id)}
            >
              <Text
                style={[
                  styles.categoryText,
                  selectedCategory === cat.id && styles.categoryTextActive,
                ]}
              >
                {cat.name}
              </Text>
            </TouchableOpacity>
          ))}
        </ScrollView>

        {/* Product Grid */}
        <FlatList
          data={filteredProducts}
          keyExtractor={(item) => item.id}
          numColumns={3}
          contentContainerStyle={styles.productGrid}
          renderItem={({ item }) => (
            <TouchableOpacity
              style={styles.productCard}
              onPress={() => handleProductPress(item)}
            >
              <View style={styles.productImagePlaceholder}>
                <Text style={styles.productEmoji}>☕</Text>
              </View>
              <Text style={styles.productName} numberOfLines={2}>
                {item.name}
              </Text>
              <Text style={styles.productPrice}>
                {formatPrice(item.base_price)}
              </Text>
            </TouchableOpacity>
          )}
        />
      </View>

      {/* Right Side - Cart */}
      <View style={styles.rightPanel}>
        {/* Order Type Toggle */}
        <View style={styles.orderTypeContainer}>
          <TouchableOpacity
            style={[
              styles.orderTypeBtn,
              cart.orderType === "dine_in" && styles.orderTypeBtnActive,
            ]}
            onPress={() => cart.setOrderType("dine_in")}
          >
            <Text
              style={[
                styles.orderTypeText,
                cart.orderType === "dine_in" && styles.orderTypeTextActive,
              ]}
            >
              Dine In
            </Text>
          </TouchableOpacity>
          <TouchableOpacity
            style={[
              styles.orderTypeBtn,
              cart.orderType === "takeaway" && styles.orderTypeBtnActive,
            ]}
            onPress={() => cart.setOrderType("takeaway")}
          >
            <Text
              style={[
                styles.orderTypeText,
                cart.orderType === "takeaway" && styles.orderTypeTextActive,
              ]}
            >
              Takeaway
            </Text>
          </TouchableOpacity>
        </View>

        {/* Table Selection (for dine-in) */}
        {cart.orderType === "dine_in" && (
          <TouchableOpacity
            style={styles.tableSelector}
            onPress={() => setShowTableModal(true)}
          >
            <Text style={styles.tableSelectorText}>
              {cart.tableId
                ? `Table ${tables.find((t) => t.id === cart.tableId)?.number || ""}`
                : "Select Table"}
            </Text>
          </TouchableOpacity>
        )}

        {/* Cart Items */}
        <ScrollView style={styles.cartItems}>
          {cart.items.length === 0 ? (
            <Text style={styles.emptyCart}>Cart is empty</Text>
          ) : (
            cart.items.map((item) => (
              <View key={item.id} style={styles.cartItem}>
                <View style={styles.cartItemInfo}>
                  <Text style={styles.cartItemName}>{item.productName}</Text>
                  {item.variantName && (
                    <Text style={styles.cartItemVariant}>
                      {item.variantName}
                    </Text>
                  )}
                  <Text style={styles.cartItemPrice}>
                    {formatPrice(item.unitPrice)}
                  </Text>
                </View>
                <View style={styles.cartItemQty}>
                  <TouchableOpacity
                    style={styles.qtyBtn}
                    onPress={() =>
                      cart.updateQuantity(item.id, item.quantity - 1)
                    }
                  >
                    <Text style={styles.qtyBtnText}>-</Text>
                  </TouchableOpacity>
                  <Text style={styles.qtyText}>{item.quantity}</Text>
                  <TouchableOpacity
                    style={styles.qtyBtn}
                    onPress={() =>
                      cart.updateQuantity(item.id, item.quantity + 1)
                    }
                  >
                    <Text style={styles.qtyBtnText}>+</Text>
                  </TouchableOpacity>
                </View>
              </View>
            ))
          )}
        </ScrollView>

        {/* Cart Summary */}
        <View style={styles.cartSummary}>
          <View style={styles.summaryRow}>
            <Text style={styles.summaryLabel}>Subtotal</Text>
            <Text style={styles.summaryValue}>
              {formatPrice(cart.getSubtotal())}
            </Text>
          </View>
          <View style={styles.summaryRow}>
            <Text style={styles.summaryLabel}>Tax (11%)</Text>
            <Text style={styles.summaryValue}>
              {formatPrice(cart.getTax())}
            </Text>
          </View>
          <View style={[styles.summaryRow, styles.totalRow]}>
            <Text style={styles.totalLabel}>Total</Text>
            <Text style={styles.totalValue}>
              {formatPrice(cart.getTotal())}
            </Text>
          </View>
        </View>

        {/* Action Buttons */}
        <View style={styles.actionButtons}>
          <TouchableOpacity
            style={styles.clearBtn}
            onPress={() => cart.clearCart()}
          >
            <Text style={styles.clearBtnText}>Clear</Text>
          </TouchableOpacity>
          <TouchableOpacity
            style={[
              styles.payBtn,
              cart.items.length === 0 && styles.payBtnDisabled,
            ]}
            onPress={handleSubmitOrder}
            disabled={cart.items.length === 0}
          >
            <Text style={styles.payBtnText}>
              Pay ({cart.getItemCount()} items)
            </Text>
          </TouchableOpacity>
        </View>
      </View>

      {/* Table Selection Modal */}
      <Modal visible={showTableModal} transparent animationType="fade">
        <View style={styles.modalOverlay}>
          <View style={styles.modalContent}>
            <Text style={styles.modalTitle}>Select Table</Text>
            <View style={styles.tableGrid}>
              {tables
                .filter((t) => t.status === "available")
                .map((table) => (
                  <TouchableOpacity
                    key={table.id}
                    style={[
                      styles.tableBtn,
                      cart.tableId === table.id && styles.tableBtnActive,
                    ]}
                    onPress={() => {
                      cart.setTableId(table.id);
                      setShowTableModal(false);
                    }}
                  >
                    <Text
                      style={[
                        styles.tableBtnText,
                        cart.tableId === table.id && styles.tableBtnTextActive,
                      ]}
                    >
                      {table.number}
                    </Text>
                    <Text style={styles.tableCapacity}>
                      ({table.capacity} pax)
                    </Text>
                  </TouchableOpacity>
                ))}
            </View>
            <TouchableOpacity
              style={styles.modalCloseBtn}
              onPress={() => setShowTableModal(false)}
            >
              <Text style={styles.modalCloseBtnText}>Close</Text>
            </TouchableOpacity>
          </View>
        </View>
      </Modal>

      {/* Payment Modal */}
      <Modal visible={showPaymentModal} transparent animationType="fade">
        <View style={styles.modalOverlay}>
          <View style={styles.modalContent}>
            <Text style={styles.modalTitle}>Cash Payment</Text>
            <Text style={styles.paymentTotal}>
              Total: {formatPrice(cart.getTotal())}
            </Text>
            <TextInput
              style={styles.cashInput}
              placeholder="Enter cash amount"
              placeholderTextColor="#666"
              keyboardType="numeric"
              value={cashAmount}
              onChangeText={setCashAmount}
            />
            {parseInt(cashAmount) > cart.getTotal() && (
              <Text style={styles.changeText}>
                Change: {formatPrice(parseInt(cashAmount) - cart.getTotal())}
              </Text>
            )}
            <View style={styles.quickAmounts}>
              {[50000, 100000, 150000, 200000].map((amount) => (
                <TouchableOpacity
                  key={amount}
                  style={styles.quickAmountBtn}
                  onPress={() => setCashAmount(amount.toString())}
                >
                  <Text style={styles.quickAmountText}>
                    {formatPrice(amount)}
                  </Text>
                </TouchableOpacity>
              ))}
            </View>
            <View style={styles.paymentActions}>
              <TouchableOpacity
                style={styles.cancelPayBtn}
                onPress={() => {
                  setShowPaymentModal(false);
                  setCashAmount("");
                }}
              >
                <Text style={styles.cancelPayBtnText}>Cancel</Text>
              </TouchableOpacity>
              <TouchableOpacity
                style={[
                  styles.confirmPayBtn,
                  processing && styles.payBtnDisabled,
                ]}
                onPress={handlePayment}
                disabled={processing}
              >
                {processing ? (
                  <ActivityIndicator color="#fff" />
                ) : (
                  <Text style={styles.confirmPayBtnText}>Confirm Payment</Text>
                )}
              </TouchableOpacity>
            </View>
          </View>
        </View>
      </Modal>

      {/* Product Variant Modal */}
      <Modal visible={showProductModal} transparent animationType="fade">
        <View style={styles.modalOverlay}>
          <View style={styles.modalContent}>
            <Text style={styles.modalTitle}>{selectedProduct?.name}</Text>
            <Text style={styles.variantHint}>Select variant:</Text>
            {selectedProduct?.variants.map((variant) => (
              <TouchableOpacity
                key={variant.id}
                style={styles.variantBtn}
                onPress={() =>
                  handleAddWithVariant(selectedProduct!, variant.id)
                }
              >
                <Text style={styles.variantName}>{variant.name}</Text>
                <Text style={styles.variantPrice}>
                  {variant.price_adjustment > 0
                    ? `+${formatPrice(variant.price_adjustment)}`
                    : formatPrice(selectedProduct.base_price)}
                </Text>
              </TouchableOpacity>
            ))}
            {selectedProduct?.variants.length === 0 && (
              <TouchableOpacity
                style={styles.variantBtn}
                onPress={() => handleAddWithVariant(selectedProduct!)}
              >
                <Text style={styles.variantName}>Regular</Text>
                <Text style={styles.variantPrice}>
                  {formatPrice(selectedProduct.base_price)}
                </Text>
              </TouchableOpacity>
            )}
            <TouchableOpacity
              style={styles.modalCloseBtn}
              onPress={() => {
                setShowProductModal(false);
                setSelectedProduct(null);
              }}
            >
              <Text style={styles.modalCloseBtnText}>Cancel</Text>
            </TouchableOpacity>
          </View>
        </View>
      </Modal>
    </SafeAreaView>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    flexDirection: "row",
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
  loadingText: {
    color: "#fff",
    marginTop: 10,
  },

  // Left Panel - Products
  leftPanel: {
    flex: 2,
    borderRightWidth: 1,
    borderRightColor: "#222",
  },
  categoryBar: {
    maxHeight: 50,
    paddingHorizontal: 10,
    paddingVertical: 8,
    backgroundColor: "#111",
  },
  categoryChip: {
    paddingHorizontal: 16,
    paddingVertical: 8,
    marginRight: 8,
    borderRadius: 20,
    backgroundColor: "#222",
  },
  categoryChipActive: {
    backgroundColor: "#8B5CF6",
  },
  categoryText: {
    color: "#888",
    fontSize: 14,
  },
  categoryTextActive: {
    color: "#fff",
    fontWeight: "bold",
  },
  productGrid: {
    padding: 10,
  },
  productCard: {
    flex: 1,
    margin: 5,
    padding: 12,
    backgroundColor: "#1a1a1a",
    borderRadius: 12,
    alignItems: "center",
    maxWidth: "31%",
  },
  productImagePlaceholder: {
    width: 60,
    height: 60,
    borderRadius: 30,
    backgroundColor: "#333",
    justifyContent: "center",
    alignItems: "center",
    marginBottom: 8,
  },
  productEmoji: {
    fontSize: 28,
  },
  productName: {
    color: "#fff",
    fontSize: 13,
    fontWeight: "600",
    textAlign: "center",
    marginBottom: 4,
  },
  productPrice: {
    color: "#8B5CF6",
    fontSize: 12,
    fontWeight: "bold",
  },

  // Right Panel - Cart
  rightPanel: {
    flex: 1,
    backgroundColor: "#111",
  },
  orderTypeContainer: {
    flexDirection: "row",
    padding: 10,
    borderBottomWidth: 1,
    borderBottomColor: "#222",
  },
  orderTypeBtn: {
    flex: 1,
    paddingVertical: 10,
    alignItems: "center",
    borderRadius: 8,
    backgroundColor: "#222",
    marginHorizontal: 4,
  },
  orderTypeBtnActive: {
    backgroundColor: "#8B5CF6",
  },
  orderTypeText: {
    color: "#888",
    fontWeight: "600",
  },
  orderTypeTextActive: {
    color: "#fff",
  },
  tableSelector: {
    padding: 12,
    backgroundColor: "#1a1a1a",
    margin: 10,
    borderRadius: 8,
    alignItems: "center",
  },
  tableSelectorText: {
    color: "#8B5CF6",
    fontWeight: "600",
  },
  cartItems: {
    flex: 1,
    padding: 10,
  },
  emptyCart: {
    color: "#666",
    textAlign: "center",
    marginTop: 40,
  },
  cartItem: {
    flexDirection: "row",
    justifyContent: "space-between",
    alignItems: "center",
    paddingVertical: 12,
    borderBottomWidth: 1,
    borderBottomColor: "#222",
  },
  cartItemInfo: {
    flex: 1,
  },
  cartItemName: {
    color: "#fff",
    fontSize: 14,
    fontWeight: "600",
  },
  cartItemVariant: {
    color: "#888",
    fontSize: 12,
  },
  cartItemPrice: {
    color: "#8B5CF6",
    fontSize: 12,
    marginTop: 2,
  },
  cartItemQty: {
    flexDirection: "row",
    alignItems: "center",
  },
  qtyBtn: {
    width: 28,
    height: 28,
    borderRadius: 14,
    backgroundColor: "#333",
    justifyContent: "center",
    alignItems: "center",
  },
  qtyBtnText: {
    color: "#fff",
    fontSize: 16,
    fontWeight: "bold",
  },
  qtyText: {
    color: "#fff",
    fontSize: 14,
    marginHorizontal: 12,
  },
  cartSummary: {
    padding: 15,
    borderTopWidth: 1,
    borderTopColor: "#222",
  },
  summaryRow: {
    flexDirection: "row",
    justifyContent: "space-between",
    marginBottom: 8,
  },
  summaryLabel: {
    color: "#888",
    fontSize: 14,
  },
  summaryValue: {
    color: "#fff",
    fontSize: 14,
  },
  totalRow: {
    marginTop: 8,
    paddingTop: 8,
    borderTopWidth: 1,
    borderTopColor: "#333",
  },
  totalLabel: {
    color: "#fff",
    fontSize: 16,
    fontWeight: "bold",
  },
  totalValue: {
    color: "#8B5CF6",
    fontSize: 18,
    fontWeight: "bold",
  },
  actionButtons: {
    flexDirection: "row",
    padding: 10,
    gap: 10,
  },
  clearBtn: {
    flex: 1,
    paddingVertical: 14,
    backgroundColor: "#333",
    borderRadius: 10,
    alignItems: "center",
  },
  clearBtnText: {
    color: "#fff",
    fontWeight: "600",
  },
  payBtn: {
    flex: 2,
    paddingVertical: 14,
    backgroundColor: "#8B5CF6",
    borderRadius: 10,
    alignItems: "center",
  },
  payBtnDisabled: {
    backgroundColor: "#444",
  },
  payBtnText: {
    color: "#fff",
    fontWeight: "bold",
    fontSize: 15,
  },

  // Modals
  modalOverlay: {
    flex: 1,
    backgroundColor: "rgba(0,0,0,0.8)",
    justifyContent: "center",
    alignItems: "center",
  },
  modalContent: {
    backgroundColor: "#1a1a1a",
    borderRadius: 16,
    padding: 24,
    width: "80%",
    maxWidth: 400,
  },
  modalTitle: {
    color: "#fff",
    fontSize: 20,
    fontWeight: "bold",
    marginBottom: 16,
    textAlign: "center",
  },
  tableGrid: {
    flexDirection: "row",
    flexWrap: "wrap",
    justifyContent: "center",
    gap: 10,
    marginBottom: 20,
  },
  tableBtn: {
    width: 70,
    height: 70,
    borderRadius: 12,
    backgroundColor: "#222",
    justifyContent: "center",
    alignItems: "center",
  },
  tableBtnActive: {
    backgroundColor: "#8B5CF6",
  },
  tableBtnText: {
    color: "#fff",
    fontSize: 24,
    fontWeight: "bold",
  },
  tableBtnTextActive: {
    color: "#fff",
  },
  tableCapacity: {
    color: "#888",
    fontSize: 11,
  },
  modalCloseBtn: {
    paddingVertical: 12,
    alignItems: "center",
  },
  modalCloseBtnText: {
    color: "#888",
    fontSize: 16,
  },

  // Payment Modal
  paymentTotal: {
    color: "#8B5CF6",
    fontSize: 24,
    fontWeight: "bold",
    textAlign: "center",
    marginBottom: 20,
  },
  cashInput: {
    backgroundColor: "#222",
    borderRadius: 10,
    padding: 14,
    color: "#fff",
    fontSize: 18,
    textAlign: "center",
    marginBottom: 10,
  },
  changeText: {
    color: "#10B981",
    fontSize: 16,
    textAlign: "center",
    marginBottom: 15,
  },
  quickAmounts: {
    flexDirection: "row",
    flexWrap: "wrap",
    justifyContent: "center",
    gap: 8,
    marginBottom: 20,
  },
  quickAmountBtn: {
    paddingHorizontal: 16,
    paddingVertical: 10,
    backgroundColor: "#333",
    borderRadius: 8,
  },
  quickAmountText: {
    color: "#fff",
    fontSize: 14,
  },
  paymentActions: {
    flexDirection: "row",
    gap: 10,
  },
  cancelPayBtn: {
    flex: 1,
    paddingVertical: 14,
    backgroundColor: "#333",
    borderRadius: 10,
    alignItems: "center",
  },
  cancelPayBtnText: {
    color: "#fff",
    fontWeight: "600",
  },
  confirmPayBtn: {
    flex: 2,
    paddingVertical: 14,
    backgroundColor: "#10B981",
    borderRadius: 10,
    alignItems: "center",
  },
  confirmPayBtnText: {
    color: "#fff",
    fontWeight: "bold",
    fontSize: 15,
  },

  // Variant Modal
  variantHint: {
    color: "#888",
    marginBottom: 12,
  },
  variantBtn: {
    flexDirection: "row",
    justifyContent: "space-between",
    padding: 14,
    backgroundColor: "#222",
    borderRadius: 10,
    marginBottom: 8,
  },
  variantName: {
    color: "#fff",
    fontSize: 15,
  },
  variantPrice: {
    color: "#8B5CF6",
    fontSize: 15,
    fontWeight: "600",
  },
});
