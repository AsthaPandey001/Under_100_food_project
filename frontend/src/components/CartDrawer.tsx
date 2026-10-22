import React from "react";
import { X, Plus, Minus, Trash2, ShoppingBag, ArrowRight, ShieldCheck } from "lucide-react";
import { CartItem } from "../hooks/useCart";

interface CartDrawerProps {
  isOpen: boolean;
  onClose: () => void;
  cartItems: CartItem[];
  subtotal: number;
  totalCount: number;
  onUpdateQuantity: (id: string, delta: number) => void;
  onRemoveItem: (id: string) => void;
  onClearCart: () => void;
}

export const CartDrawer: React.FC<CartDrawerProps> = ({
  isOpen,
  onClose,
  cartItems,
  subtotal,
  totalCount,
  onUpdateQuantity,
  onRemoveItem,
  onClearCart
}) => {
  if (!isOpen) return null;

  return (
    <div style={{
      position: "fixed",
      inset: 0,
      zIndex: 200,
      display: "flex",
      justifyContent: "flex-end"
    }}>
      {/* Backdrop */}
      <div
        onClick={onClose}
        style={{
          position: "fixed",
          inset: 0,
          backgroundColor: "rgba(23, 28, 36, 0.45)",
          backdropFilter: "blur(4px)"
        }}
      />

      {/* Drawer Panel */}
      <div style={{
        position: "relative",
        width: "100%",
        maxWidth: "420px",
        height: "100%",
        backgroundColor: "#ffffff",
        boxShadow: "var(--shadow-elevated)",
        display: "flex",
        flexDirection: "column",
        zIndex: 201
      }}>
        {/* Header */}
        <div style={{
          padding: "20px 24px",
          borderBottom: "1px solid var(--color-border)",
          display: "flex",
          alignItems: "center",
          justifyContent: "space-between"
        }}>
          <div style={{ display: "flex", alignItems: "center", gap: "10px" }}>
            <div style={{
              width: "36px",
              height: "36px",
              borderRadius: "10px",
              backgroundColor: "var(--color-soft-orange-bg)",
              color: "var(--color-primary-orange)",
              display: "flex",
              alignItems: "center",
              justifyContent: "center"
            }}>
              <ShoppingBag size={20} />
            </div>
            <div>
              <h3 style={{ fontSize: "18px", fontWeight: 800, color: "var(--color-dark-text)" }}>
                Your Meal Cart
              </h3>
              <p style={{ fontSize: "12px", color: "var(--color-secondary-text)" }}>
                {totalCount} {totalCount === 1 ? "item" : "items"} selected
              </p>
            </div>
          </div>
          <button
            type="button"
            onClick={onClose}
            style={{
              padding: "6px",
              borderRadius: "50%",
              color: "var(--color-secondary-text)",
              cursor: "pointer"
            }}
            aria-label="Close cart"
          >
            <X size={20} />
          </button>
        </div>

        {/* Cart Items List */}
        <div style={{
          flex: 1,
          overflowY: "auto",
          padding: "20px 24px"
        }}>
          {cartItems.length > 0 ? (
            <div style={{ display: "flex", flexDirection: "column", gap: "16px" }}>
              {cartItems.map(({ item, quantity }) => (
                <div
                  key={item.id}
                  style={{
                    display: "flex",
                    alignItems: "center",
                    gap: "12px",
                    padding: "12px",
                    borderRadius: "var(--radius-md)",
                    backgroundColor: "var(--color-main-bg)",
                    border: "1px solid var(--color-border)"
                  }}
                >
                  <img
                    src={item.image}
                    alt={item.name}
                    style={{
                      width: "60px",
                      height: "60px",
                      borderRadius: "8px",
                      objectFit: "cover"
                    }}
                    onError={(e) => {
                      (e.target as HTMLImageElement).src = "/crispy_dosa.jpg";
                    }}
                  />
                  <div style={{ flex: 1, minWidth: 0 }}>
                    <div style={{
                      fontSize: "14px",
                      fontWeight: 700,
                      color: "var(--color-dark-text)",
                      whiteSpace: "nowrap",
                      overflow: "hidden",
                      textOverflow: "ellipsis"
                    }}>
                      {item.name}
                    </div>
                    <div style={{ fontSize: "12px", color: "var(--color-secondary-text)", marginBottom: "4px" }}>
                      {item.restaurant}
                    </div>
                    <div style={{ fontSize: "14px", fontWeight: 800, color: "var(--color-dark-text)" }}>
                      ₹{item.price * quantity} <span style={{ fontSize: "11px", fontWeight: 500, color: "var(--color-secondary-text)" }}>(₹{item.price} each)</span>
                    </div>
                  </div>

                  {/* Quantity Controls */}
                  <div style={{
                    display: "flex",
                    alignItems: "center",
                    backgroundColor: "#ffffff",
                    border: "1px solid var(--color-border)",
                    borderRadius: "var(--radius-full)"
                  }}>
                    <button
                      type="button"
                      onClick={() => onUpdateQuantity(item.id, -1)}
                      style={{ padding: "4px 8px", color: "var(--color-dark-text)" }}
                      aria-label="Decrease quantity"
                    >
                      <Minus size={13} />
                    </button>
                    <span style={{ fontSize: "12px", fontWeight: 700, minWidth: "16px", textAlign: "center" }}>
                      {quantity}
                    </span>
                    <button
                      type="button"
                      onClick={() => onUpdateQuantity(item.id, 1)}
                      style={{ padding: "4px 8px", color: "var(--color-dark-text)" }}
                      aria-label="Increase quantity"
                    >
                      <Plus size={13} />
                    </button>
                  </div>

                  <button
                    type="button"
                    onClick={() => onRemoveItem(item.id)}
                    style={{ padding: "6px", color: "#94A3B8" }}
                    aria-label="Remove item"
                  >
                    <Trash2 size={16} />
                  </button>
                </div>
              ))}

              <div style={{ display: "flex", justifyContent: "flex-end" }}>
                <button
                  type="button"
                  onClick={onClearCart}
                  style={{
                    fontSize: "12px",
                    fontWeight: 600,
                    color: "var(--color-secondary-text)",
                    textDecoration: "underline",
                    cursor: "pointer"
                  }}
                >
                  Clear all items
                </button>
              </div>
            </div>
          ) : (
            /* Empty Cart */
            <div style={{ textAlign: "center", padding: "60px 0" }}>
              <div style={{
                width: "64px",
                height: "64px",
                borderRadius: "50%",
                backgroundColor: "#F8F6F0",
                display: "flex",
                alignItems: "center",
                justifyContent: "center",
                margin: "0 auto 16px",
                color: "var(--color-secondary-text)"
              }}>
                <ShoppingBag size={28} />
              </div>
              <h4 style={{ fontSize: "17px", fontWeight: 800, color: "var(--color-dark-text)", marginBottom: "6px" }}>
                Your cart is empty
              </h4>
              <p style={{ fontSize: "14px", color: "var(--color-secondary-text)", maxWidth: "260px", margin: "0 auto 20px" }}>
                Explore Bengaluru meals under ₹100 and add your favourites.
              </p>
              <button
                type="button"
                onClick={onClose}
                style={{
                  padding: "10px 20px",
                  borderRadius: "var(--radius-full)",
                  backgroundColor: "var(--color-primary-orange)",
                  color: "#ffffff",
                  fontSize: "13px",
                  fontWeight: 700
                }}
              >
                Start Browsing
              </button>
            </div>
          )}
        </div>

        {/* Footer with Subtotal & Demo Checkout */}
        {cartItems.length > 0 && (
          <div style={{
            padding: "20px 24px",
            borderTop: "1px solid var(--color-border)",
            backgroundColor: "#F8F6F0"
          }}>
            <div style={{
              display: "flex",
              justifyContent: "space-between",
              marginBottom: "8px",
              fontSize: "14px"
            }}>
              <span style={{ color: "var(--color-secondary-text)" }}>Item Total</span>
              <span style={{ fontWeight: 700, color: "var(--color-dark-text)" }}>₹{subtotal}</span>
            </div>
            <div style={{
              display: "flex",
              justifyContent: "space-between",
              marginBottom: "14px",
              fontSize: "14px"
            }}>
              <span style={{ color: "var(--color-secondary-text)" }}>Platform Fee</span>
              <span style={{ fontWeight: 700, color: "var(--color-deep-green)" }}>FREE (Demo)</span>
            </div>

            <div style={{
              height: "1px",
              backgroundColor: "var(--color-border)",
              marginBottom: "14px"
            }} />

            <div style={{
              display: "flex",
              justifyContent: "space-between",
              alignItems: "baseline",
              marginBottom: "16px"
            }}>
              <span style={{ fontSize: "16px", fontWeight: 800, color: "var(--color-dark-text)" }}>To Pay</span>
              <span style={{ fontSize: "22px", fontWeight: 900, color: "var(--color-primary-orange)" }}>₹{subtotal}</span>
            </div>

            {/* Demo Checkout Button */}
            <button
              type="button"
              onClick={() => alert(`Frontend Demo Checkout: Order summary for ${totalCount} items (Total ₹${subtotal}) received! No real payment processed.`)}
              style={{
                width: "100%",
                padding: "14px",
                borderRadius: "var(--radius-md)",
                backgroundColor: "var(--color-deep-green)",
                color: "#ffffff",
                fontSize: "15px",
                fontWeight: 800,
                display: "flex",
                alignItems: "center",
                justifyContent: "center",
                gap: "8px",
                boxShadow: "0 4px 12px rgba(8, 116, 67, 0.2)"
              }}
            >
              <span>Demo Checkout • ₹{subtotal}</span>
              <ArrowRight size={18} />
            </button>

            <div style={{
              display: "flex",
              alignItems: "center",
              justifyContent: "center",
              gap: "6px",
              marginTop: "12px",
              fontSize: "12px",
              color: "var(--color-secondary-text)"
            }}>
              <ShieldCheck size={14} color="var(--color-deep-green)" />
              <span>Transparent item pricing • Frontend Demo</span>
            </div>
          </div>
        )}
      </div>
    </div>
  );
};
