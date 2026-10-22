import { useState, useEffect } from "react";
import { FoodItem } from "../data/foodItems";

export interface CartItem {
  item: FoodItem;
  quantity: number;
}

export function useCart() {
  const [cartItems, setCartItems] = useState<CartItem[]>(() => {
    try {
      const saved = localStorage.getItem("under100_cart");
      return saved ? JSON.parse(saved) : [];
    } catch {
      return [];
    }
  });

  const [isDrawerOpen, setIsDrawerOpen] = useState(false);

  useEffect(() => {
    try {
      localStorage.setItem("under100_cart", JSON.stringify(cartItems));
    } catch {
      // Ignore storage error
    }
  }, [cartItems]);

  const addToCart = (item: FoodItem) => {
    setCartItems(prev => {
      const existing = prev.find(ci => ci.item.id === item.id);
      if (existing) {
        return prev.map(ci =>
          ci.item.id === item.id ? { ...ci, quantity: ci.quantity + 1 } : ci
        );
      }
      return [...prev, { item, quantity: 1 }];
    });
  };

  const updateQuantity = (itemId: string, delta: number) => {
    setCartItems(prev => {
      return prev
        .map(ci => {
          if (ci.item.id === itemId) {
            const newQty = ci.quantity + delta;
            return newQty > 0 ? { ...ci, quantity: newQty } : null;
          }
          return ci;
        })
        .filter((ci): ci is CartItem => ci !== null);
    });
  };

  const removeFromCart = (itemId: string) => {
    setCartItems(prev => prev.filter(ci => ci.item.id !== itemId));
  };

  const clearCart = () => {
    setCartItems([]);
  };

  const getItemQuantity = (itemId: string): number => {
    const found = cartItems.find(ci => ci.item.id === itemId);
    return found ? found.quantity : 0;
  };

  const totalCount = cartItems.reduce((sum, ci) => sum + ci.quantity, 0);
  const subtotal = cartItems.reduce((sum, ci) => sum + ci.item.price * ci.quantity, 0);

  return {
    cartItems,
    addToCart,
    updateQuantity,
    removeFromCart,
    clearCart,
    getItemQuantity,
    totalCount,
    subtotal,
    isDrawerOpen,
    setIsDrawerOpen
  };
}
