'use client'

import { useState, useCallback } from 'react'
import { Header } from '@/components/header'
import { Menu } from '@/components/menu'
import { Cart } from '@/components/cart'
import { OrderStatusComponent } from '@/components/order-status'
import { Footer } from '@/components/footer'
import { drinks } from '@/lib/drinks'
import { Drink, CartItem, Order, OrderStatus } from '@/lib/types'

function generateOrderId(): string {
  return `ORD-${Date.now().toString(36).toUpperCase()}-${Math.random()
    .toString(36)
    .substring(2, 6)
    .toUpperCase()}`
}

export default function CoffeeShop() {
  const [cartItems, setCartItems] = useState<CartItem[]>([])
  const [isCartOpen, setIsCartOpen] = useState(false)
  const [currentOrder, setCurrentOrder] = useState<Order | null>(null)

  const handleAddToCart = useCallback((drink: Drink) => {
    setCartItems((prev) => {
      const existing = prev.find((item) => item.drink.id === drink.id)
      if (existing) {
        return prev.map((item) =>
          item.drink.id === drink.id
            ? { ...item, quantity: item.quantity + 1 }
            : item
        )
      }
      return [...prev, { drink, quantity: 1 }]
    })
  }, [])

  const handleUpdateQuantity = useCallback((drinkId: string, quantity: number) => {
    if (quantity <= 0) {
      setCartItems((prev) => prev.filter((item) => item.drink.id !== drinkId))
    } else {
      setCartItems((prev) =>
        prev.map((item) =>
          item.drink.id === drinkId ? { ...item, quantity } : item
        )
      )
    }
  }, [])

  const handleRemoveItem = useCallback((drinkId: string) => {
    setCartItems((prev) => prev.filter((item) => item.drink.id !== drinkId))
  }, [])

  const handlePlaceOrder = useCallback(() => {
    const total = cartItems.reduce(
      (sum, item) => sum + item.drink.price * item.quantity,
      0
    )

    const order: Order = {
      id: generateOrderId(),
      items: [...cartItems],
      total,
      status: 'received',
      createdAt: new Date(),
    }

    setCurrentOrder(order)
    setCartItems([])
    setIsCartOpen(false)
  }, [cartItems])

  const handleRefreshStatus = useCallback(() => {
    if (!currentOrder) return

    const statusFlow: OrderStatus[] = ['received', 'preparing', 'ready']
    const currentIndex = statusFlow.indexOf(currentOrder.status)

    if (currentIndex < statusFlow.length - 1) {
      setCurrentOrder((prev) =>
        prev ? { ...prev, status: statusFlow[currentIndex + 1] } : null
      )
    }
  }, [currentOrder])

  const handleCloseOrderStatus = useCallback(() => {
    setCurrentOrder(null)
  }, [])

  return (
    <div className="min-h-screen bg-background">
      <Header cartItems={cartItems} onCartClick={() => setIsCartOpen(true)} />

      {currentOrder && (
        <OrderStatusComponent
          order={currentOrder}
          onRefresh={handleRefreshStatus}
          onClose={handleCloseOrderStatus}
        />
      )}

      <main>
        <Menu
          drinks={drinks}
          cartItems={cartItems}
          onAddToCart={handleAddToCart}
        />
      </main>

      {isCartOpen && (
        <Cart
          items={cartItems}
          onUpdateQuantity={handleUpdateQuantity}
          onRemoveItem={handleRemoveItem}
          onPlaceOrder={handlePlaceOrder}
          onClose={() => setIsCartOpen(false)}
        />
      )}

      <Footer />
    </div>
  )
}
