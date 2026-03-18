'use client'

import { useState, useCallback, useEffect } from 'react'
import { Header } from '@/components/header'
import { Menu } from '@/components/menu'
import { Cart } from '@/components/cart'
import { OrderStatusComponent } from '@/components/order-status'
import { Footer } from '@/components/footer'
import { Drink, CartItem, Order, OrderStatus } from '@/lib/types'

function coerceStatus(s: string): OrderStatus {
  if (s === 'received' || s === 'preparing' || s === 'ready') return s
  return 'received'
}

export default function CoffeeShop() {
  const [menu, setMenu] = useState<Drink[]>([])
  const [menuError, setMenuError] = useState<string | null>(null)
  const [cartItems, setCartItems] = useState<CartItem[]>([])
  const [isCartOpen, setIsCartOpen] = useState(false)
  const [currentOrder, setCurrentOrder] = useState<Order | null>(null)

  const refreshOrderStatus = useCallback(async (orderId: string) => {
    const res = await fetch(`/api/orders/${encodeURIComponent(orderId)}/status`, {
      cache: 'no-store',
    })
    if (!res.ok) return
    const data = (await res.json()) as { status: string }
    const next = coerceStatus(data.status)
    setCurrentOrder((prev) => (prev ? { ...prev, status: next } : null))
  }, [])

  // Poll status until it becomes ready (near real-time without streaming).
  useEffect(() => {
    if (!currentOrder) return
    if (currentOrder.status === 'ready') return

    const id = currentOrder.id
    const t = setInterval(() => {
      void refreshOrderStatus(id)
    }, 1500)

    return () => clearInterval(t)
  }, [currentOrder, refreshOrderStatus])

  useEffect(() => {
    let cancelled = false
    ;(async () => {
      try {
        const res = await fetch('/api/menu', { cache: 'no-store' })
        if (!res.ok) throw new Error(`menu: ${res.status}`)
        const data = (await res.json()) as { items: Drink[] }
        if (!cancelled) setMenu(Array.isArray(data.items) ? data.items : [])
      } catch (e: any) {
        if (!cancelled) setMenuError(e?.message ?? 'Failed to load menu')
      }
    })()
    return () => {
      cancelled = true
    }
  }, [])

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

  const handlePlaceOrder = useCallback(async () => {
    const total = cartItems.reduce(
      (sum, item) => sum + item.drink.price * item.quantity,
      0
    )

    const res = await fetch('/api/orders', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        items: cartItems.map((ci) => ({ item: ci.drink, quantity: ci.quantity })),
      }),
    })

    if (!res.ok) {
      const msg = await res.text()
      throw new Error(msg || `order: ${res.status}`)
    }

    const data = (await res.json()) as { receiptId: string }
    const receiptId = data.receiptId
    if (!receiptId) throw new Error('missing receipt id')

    const order: Order = {
      id: receiptId,
      items: [...cartItems],
      total,
      status: 'received',
      createdAt: new Date(),
    }

    setCurrentOrder(order)
    setCartItems([])
    setIsCartOpen(false)
  }, [cartItems])

  const handleRefreshStatus = useCallback(async () => {
    if (!currentOrder) return
    await refreshOrderStatus(currentOrder.id)
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
        {menuError && (
          <div className="max-w-5xl mx-auto px-6 py-6">
            <p className="text-sm text-muted-foreground">
              Failed to load menu. ({menuError})
            </p>
          </div>
        )}
        <Menu
          drinks={menu}
          cartItems={cartItems}
          onAddToCart={handleAddToCart}
        />
      </main>

      {isCartOpen && (
        <Cart
          items={cartItems}
          onUpdateQuantity={handleUpdateQuantity}
          onRemoveItem={handleRemoveItem}
          onPlaceOrder={() => void handlePlaceOrder()}
          onClose={() => setIsCartOpen(false)}
        />
      )}

      <Footer />
    </div>
  )
}
