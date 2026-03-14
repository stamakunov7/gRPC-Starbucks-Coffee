'use client'

import { ShoppingCart } from 'lucide-react'
import { CartItem } from '@/lib/types'

interface HeaderProps {
  cartItems: CartItem[]
  onCartClick: () => void
}

export function Header({ cartItems, onCartClick }: HeaderProps) {
  const totalItems = cartItems.reduce((sum, item) => sum + item.quantity, 0)

  return (
    <header className="border-b border-border">
      <div className="mx-auto max-w-5xl px-6 py-6 flex items-center justify-between">
        <h1 className="text-xl font-medium tracking-tight text-foreground">
          Coffee Shop
        </h1>
        <button
          onClick={onCartClick}
          className="flex items-center gap-2 px-4 py-2 text-sm text-muted-foreground hover:text-foreground transition-colors border border-border rounded"
        >
          <ShoppingCart className="h-4 w-4" />
          <span>Cart</span>
          {totalItems > 0 && (
            <span className="ml-1 text-foreground">{totalItems}</span>
          )}
        </button>
      </div>
    </header>
  )
}
