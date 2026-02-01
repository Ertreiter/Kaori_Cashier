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

export interface Order {
    id: string;
    order_number: string;
    table_id?: string;
    table_number?: number;
    order_type: 'dine_in' | 'takeaway';
    order_source: 'cashier' | 'table_qr' | 'client_app';
    status: 'pending' | 'confirmed' | 'cooking' | 'ready' | 'completed' | 'cancelled';
    payment_status: 'unpaid' | 'paid';
    items: OrderItem[];
    subtotal: number;
    tax: number;
    total: number;
    notes?: string;
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
