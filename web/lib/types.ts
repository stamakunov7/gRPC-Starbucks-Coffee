export interface Drink {
  id: string
  name: string
  description: string
  price: number
}

export interface CartItem {
  drink: Drink
  quantity: number
}

export type OrderStatus = 'received' | 'preparing' | 'ready'

export interface Order {
  id: string
  items: CartItem[]
  total: number
  status: OrderStatus
  createdAt: Date
}
