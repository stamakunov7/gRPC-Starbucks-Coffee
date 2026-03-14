'use client'

import { useState, useEffect } from 'react'
import { RefreshCw, X, Check } from 'lucide-react'
import { Order, OrderStatus as OrderStatusType } from '@/lib/types'

interface OrderStatusProps {
  order: Order
  onRefresh: () => void
  onClose: () => void
}

const statusLabels: Record<OrderStatusType, string> = {
  received: 'Order Received',
  preparing: 'Preparing',
  ready: 'Ready for Pickup',
}

const statusDescriptions: Record<OrderStatusType, string> = {
  received: 'Your order has been received and is in queue',
  preparing: 'Your drinks are being prepared',
  ready: 'Your order is ready. Please pick it up at the counter.',
}

export function OrderStatusComponent({ order, onRefresh, onClose }: OrderStatusProps) {
  const statusIndex = ['received', 'preparing', 'ready'].indexOf(order.status)
  const [isPickedUp, setIsPickedUp] = useState(false)
  const [isClosing, setIsClosing] = useState(false)

  const handlePickUp = () => {
    setIsPickedUp(true)
    setTimeout(() => {
      setIsClosing(true)
      setTimeout(() => {
        onClose()
      }, 300)
    }, 1500)
  }

  useEffect(() => {
    if (isPickedUp) {
      setIsClosing(false)
    }
  }, [isPickedUp])

  return (
    <div
      className={`w-full bg-card border-b border-border transition-all duration-300 ease-out overflow-hidden ${
        isClosing ? 'max-h-0 opacity-0 border-b-0' : 'max-h-48 opacity-100'
      }`}
    >
      <div className="max-w-5xl mx-auto px-6 py-6">
        {isPickedUp ? (
          <div className="flex items-center justify-center py-2">
            <div className="flex items-center gap-3">
              <Check className="h-5 w-5 text-foreground" />
              <p className="text-foreground font-medium">
                Thank you! See you again!
              </p>
            </div>
          </div>
        ) : (
          <div className="flex flex-col gap-6 md:flex-row md:items-start md:justify-between">
            {/* Left: Order Info */}
            <div className="flex-1">
              <div className="flex items-center gap-3 mb-3">
                <p className="text-sm font-medium text-foreground">
                  {statusLabels[order.status]}
                </p>
                {order.status === 'ready' && (
                  <Check className="h-4 w-4 text-foreground" />
                )}
              </div>
              <p className="text-sm text-muted-foreground mb-3">
                {statusDescriptions[order.status]}
              </p>
              <div className="flex items-center gap-4 text-xs text-muted-foreground">
                <span className="font-mono">{order.id}</span>
                <span className="text-border">|</span>
                <span>${order.total.toFixed(2)}</span>
              </div>
            </div>

            {/* Center: Progress Bar */}
            <div className="flex-1 max-w-xs">
              <div className="flex items-center gap-2 mb-2">
                {(['received', 'preparing', 'ready'] as OrderStatusType[]).map(
                  (status, index) => (
                    <div key={status} className="flex-1">
                      <div
                        className={`h-1 rounded-full transition-colors ${
                          index <= statusIndex
                            ? 'bg-foreground'
                            : 'bg-border'
                        }`}
                      />
                    </div>
                  )
                )}
              </div>
              <div className="flex items-center justify-between text-xs text-muted-foreground">
                <span>Received</span>
                <span>Preparing</span>
                <span>Ready</span>
              </div>
            </div>

            {/* Right: Actions */}
            <div className="flex items-center gap-2">
              {order.status === 'ready' ? (
                <button
                  onClick={handlePickUp}
                  className="flex items-center gap-2 py-2 px-4 bg-foreground text-background text-xs font-medium rounded hover:bg-foreground/90 transition-colors"
                >
                  <Check className="h-3.5 w-3.5" />
                  Pick Up
                </button>
              ) : (
                <button
                  onClick={onRefresh}
                  className="flex items-center gap-2 py-2 px-3 border border-border text-xs text-foreground rounded hover:border-foreground/30 transition-colors"
                >
                  <RefreshCw className="h-3.5 w-3.5" />
                  Refresh
                </button>
              )}
              <button
                onClick={onClose}
                className="p-2 text-muted-foreground hover:text-foreground transition-colors"
                aria-label="Dismiss order status"
              >
                <X className="h-4 w-4" />
              </button>
            </div>
          </div>
        )}
      </div>
    </div>
  )
}
