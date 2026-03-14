'use client'

import { Plus, Check } from 'lucide-react'
import { Drink, CartItem } from '@/lib/types'

interface MenuProps {
  drinks: Drink[]
  cartItems: CartItem[]
  onAddToCart: (drink: Drink) => void
}

export function Menu({ drinks, cartItems, onAddToCart }: MenuProps) {
  const isInCart = (drinkId: string) => {
    return cartItems.some(item => item.drink.id === drinkId)
  }

  const getQuantity = (drinkId: string) => {
    const item = cartItems.find(item => item.drink.id === drinkId)
    return item?.quantity || 0
  }

  return (
    <section className="py-12">
      <div className="mx-auto max-w-5xl px-6">
        <h2 className="text-lg font-medium text-foreground mb-2">Menu</h2>
        <p className="text-sm text-muted-foreground mb-8">
          Select drinks to add to your order
        </p>
        
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          {drinks.map((drink) => {
            const inCart = isInCart(drink.id)
            const quantity = getQuantity(drink.id)
            
            return (
              <button
                key={drink.id}
                onClick={() => onAddToCart(drink)}
                className={`group text-left p-6 rounded border transition-colors ${
                  inCart
                    ? 'bg-card border-foreground/30'
                    : 'bg-card/50 border-border hover:border-foreground/20'
                }`}
              >
                <div className="flex items-start justify-between gap-4">
                  <div className="flex-1 min-w-0">
                    <div className="flex items-center gap-2">
                      <h3 className="text-sm font-medium text-foreground">
                        {drink.name}
                      </h3>
                      {inCart && (
                        <span className="flex items-center gap-1 text-xs text-muted-foreground">
                          <Check className="h-3 w-3" />
                          {quantity}
                        </span>
                      )}
                    </div>
                    <p className="mt-1 text-sm text-muted-foreground leading-relaxed">
                      {drink.description}
                    </p>
                  </div>
                  <div className="flex items-center gap-3">
                    <span className="text-sm text-foreground font-medium">
                      ${drink.price.toFixed(2)}
                    </span>
                    <div className={`p-1.5 rounded border transition-colors ${
                      inCart
                        ? 'border-foreground/30 text-foreground'
                        : 'border-border text-muted-foreground group-hover:border-foreground/20 group-hover:text-foreground'
                    }`}>
                      <Plus className="h-3.5 w-3.5" />
                    </div>
                  </div>
                </div>
              </button>
            )
          })}
        </div>
      </div>
    </section>
  )
}
