// Types for API responses

export interface Category {
    id: string;
    name: string;
    description: string;
    sort_order: number;
}

export interface Variant {
    id: string;
    name: string;
    price_adjustment: number;
}

export interface Modifier {
    id: string;
    name: string;
    price: number;
    max_qty: number;
}

export interface Product {
    id: string;
    category_id: string;
    name: string;
    description: string;
    base_price: number;
    image_url: string;
    is_available: boolean;
    variants: Variant[];
    modifiers: Modifier[];
}

export interface Table {
    id: string;
    number: number;
    capacity: number;
    status: 'available' | 'occupied' | 'reserved';
    qr_code: string;
}

export interface OrderItem {
    id: string;
    product_id: string;
    product_name: string;
    variant_id?: string;
    variant_name?: string;
    modifiers?: string[];
    quantity: number;
    unit_price: number;
    subtotal: number;
    notes?: string;
}

// Order sources including delivery platforms
export type OrderSource = 'cashier' | 'table_qr' | 'grabfood' | 'gofood' | 'shopee_food';
export type OrderType = 'dine_in' | 'takeaway' | 'delivery';
export type OrderStatus = 'pending' | 'confirmed' | 'cooking' | 'ready' | 'completed' | 'cancelled';

export interface Order {
    id: string;
    order_number: string;
    external_order_id?: string; // GrabFood/GoFood/Shopee order ID
    table_id?: string;
    table_number?: number;
    order_type: OrderType;
    order_source: OrderSource;
    status: OrderStatus;
    payment_status: 'unpaid' | 'paid';
    items: OrderItem[];
    subtotal: number;
    tax: number;
    total: number;
    notes?: string;
    // Delivery fields
    customer_name?: string;
    customer_phone?: string;
    delivery_address?: string;
    driver_name?: string;
    // Timestamps
    created_at: string;
    updated_at: string;
    cashier_id?: string;
    cashier_name?: string;
}

export interface User {
    id: string;
    email: string;
    name: string;
    role: 'super_admin' | 'store_admin' | 'cashier' | 'kitchen';
}

// Helper for order source display
export const ORDER_SOURCE_CONFIG: Record<OrderSource, { label: string; color: string; emoji: string }> = {
    cashier: { label: 'Cashier', color: '#6B7280', emoji: '💳' },
    table_qr: { label: 'Table QR', color: '#8B5CF6', emoji: '📱' },
    grabfood: { label: 'GrabFood', color: '#00B14F', emoji: '🟢' },
    gofood: { label: 'GoFood', color: '#D71920', emoji: '🔴' },
    shopee_food: { label: 'Shopee', color: '#EE4D2D', emoji: '🧡' },
};

export const isDeliveryOrder = (source: OrderSource): boolean => {
    return ['grabfood', 'gofood', 'shopee_food'].includes(source);
};
