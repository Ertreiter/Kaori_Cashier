import { create } from 'zustand';

export interface CartItem {
    id: string;
    productId: string;
    productName: string;
    variantId?: string;
    variantName?: string;
    modifiers: string[];
    quantity: number;
    unitPrice: number;
    notes?: string;
}

interface CartState {
    items: CartItem[];
    orderType: 'dine_in' | 'takeaway';
    tableId?: string;
    notes: string;

    // Actions
    addItem: (item: Omit<CartItem, 'id'>) => void;
    removeItem: (id: string) => void;
    updateQuantity: (id: string, quantity: number) => void;
    clearCart: () => void;
    setOrderType: (type: 'dine_in' | 'takeaway') => void;
    setTableId: (tableId: string | undefined) => void;
    setNotes: (notes: string) => void;

    // Computed
    getSubtotal: () => number;
    getTax: () => number;
    getTotal: () => number;
    getItemCount: () => number;
}

export const useCartStore = create<CartState>((set, get) => ({
    items: [],
    orderType: 'dine_in',
    tableId: undefined,
    notes: '',

    addItem: (item) => {
        const id = `${item.productId}-${item.variantId || 'default'}-${Date.now()}`;
        set((state) => ({
            items: [...state.items, { ...item, id }],
        }));
    },

    removeItem: (id) => {
        set((state) => ({
            items: state.items.filter((item) => item.id !== id),
        }));
    },

    updateQuantity: (id, quantity) => {
        if (quantity <= 0) {
            get().removeItem(id);
            return;
        }
        set((state) => ({
            items: state.items.map((item) =>
                item.id === id ? { ...item, quantity } : item
            ),
        }));
    },

    clearCart: () => {
        set({ items: [], tableId: undefined, notes: '' });
    },

    setOrderType: (orderType) => {
        set({ orderType });
        if (orderType === 'takeaway') {
            set({ tableId: undefined });
        }
    },

    setTableId: (tableId) => {
        set({ tableId });
    },

    setNotes: (notes) => {
        set({ notes });
    },

    getSubtotal: () => {
        return get().items.reduce((sum, item) => sum + item.unitPrice * item.quantity, 0);
    },

    getTax: () => {
        return Math.round(get().getSubtotal() * 0.11); // 11% tax
    },

    getTotal: () => {
        return get().getSubtotal() + get().getTax();
    },

    getItemCount: () => {
        return get().items.reduce((sum, item) => sum + item.quantity, 0);
    },
}));
