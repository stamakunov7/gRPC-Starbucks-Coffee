'use client'

import { Minus, Plus, X } from 'lucide-react'
import { CartItem } from '@/lib/types'

interface CartProps {
  items: CartItem[]
  onUpdateQuantity: (drinkId: string, quantity: number) => void
  onRemoveItem: (drinkId: string) => void
  onPlaceOrder: () => void
  onClose: () => void
}

export function Cart({
  items,
  onUpdateQuantity,
  onRemoveItem,
  onPlaceOrder,
  onClose,
}: CartProps) {
  const total = items.reduce(
    (sum, item) => sum + item.drink.price * item.quantity,
    0
  )

  if (items.length === 0) {
    return (
      <div className="fixed inset-0 bg-background/80 backdrop-blur-sm z-50">
        <div className="fixed right-0 top-0 h-full w-full max-w-md bg-card border-l border-border">
          <div className="flex flex-col h-full">
            <div className="flex items-center justify-between p-6 border-b border-border">
              <h2 className="text-lg font-medium text-foreground">Your Order</h2>
              <button
                onClick={onClose}
                className="p-2 text-muted-foreground hover:text-foreground transition-colors"
              >
                <X className="h-4 w-4" />
              </button>
            </div>
            <div className="flex-1 flex items-center justify-center p-6">
              <p className="text-sm text-muted-foreground">Your cart is empty</p>
            </div>
          </div>
        </div>
      </div>
    )
  }

  return (
    <div className="fixed inset-0 bg-background/80 backdrop-blur-sm z-50">
      <div className="fixed right-0 top-0 h-full w-full max-w-md bg-card border-l border-border">
        <div className="flex flex-col h-full">
          <div className="flex items-center justify-between p-6 border-b border-border">
            <h2 className="text-lg font-medium text-foreground">Your Order</h2>
            <button
              onClick={onClose}
              className="p-2 text-muted-foreground hover:text-foreground transition-colors"
            >
              <X className="h-4 w-4" />
            </button>
          </div>

          <div className="flex-1 overflow-auto p-6">
            <div className="space-y-4">
              {items.map((item) => (
                <div
                  key={item.drink.id}
                  className="flex items-center justify-between py-4 border-b border-border last:border-0"
                >
                  <div className="flex-1 min-w-0">
                    <h3 className="text-sm font-medium text-foreground">
                      {item.drink.name}
                    </h3>
                    <p className="text-sm text-muted-foreground">
                      ${item.drink.price.toFixed(2)} each
                    </p>
                  </div>

                  <div className="flex items-center gap-3">
                    <div className="flex items-center gap-2">
                      <button
                        onClick={() =>
                          onUpdateQuantity(item.drink.id, item.quantity - 1)
                        }
                        className="p-1.5 border border-border rounded text-muted-foreground hover:text-foreground hover:border-foreground/20 transition-colors"
                      >
                        <Minus className="h-3 w-3" />
                      </button>
                      <span className="w-8 text-center text-sm text-foreground">
                        {item.quantity}
                      </span>
                      <button
                        onClick={() =>
                          onUpdateQuantity(item.drink.id, item.quantity + 1)
                        }
                        className="p-1.5 border border-border rounded text-muted-foreground hover:text-foreground hover:border-foreground/20 transition-colors"
                      >
                        <Plus className="h-3 w-3" />
                      </button>
                    </div>
                    <button
                      onClick={() => onRemoveItem(item.drink.id)}
                      className="p-1.5 text-muted-foreground hover:text-foreground transition-colors"
                    >
                      <X className="h-3.5 w-3.5" />
                    </button>
                  </div>
                </div>
              ))}
            </div>
          </div>

          <div className="p-6 border-t border-border">
            <div className="flex items-center justify-between mb-6">
              <span className="text-sm text-muted-foreground">Total</span>
              <span className="text-lg font-medium text-foreground">
                ${total.toFixed(2)}
              </span>
            </div>
            <button
              onClick={onPlaceOrder}
              className="w-full py-3 px-4 bg-foreground text-background text-sm font-medium rounded hover:bg-foreground/90 transition-colors"
            >
              Place Order
            </button>
          </div>
        </div>
      </div>
    </div>
  )
}
