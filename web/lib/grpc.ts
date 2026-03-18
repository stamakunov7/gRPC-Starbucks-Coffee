import path from "path"

import * as grpc from "@grpc/grpc-js"
import * as protoLoader from "@grpc/proto-loader"

export type MenuItem = {
  id: string
  name: string
  description: string
  price: number
}

type Receipt = { id: string }
type OrderStatus = { orderId?: string; order_id?: string; status: string }

type CoffeeShopClient = grpc.Client & {
  GetMenu: (
    req: Record<string, never>,
    metadata?: grpc.Metadata
  ) => grpc.ClientReadableStream<{ items: MenuItem[] }>
  PlaceOrder: (
    req: { items: MenuItem[] },
    cb: (err: grpc.ServiceError | null, resp?: Receipt) => void
  ) => void
  GetOrderStatus: (
    req: Receipt,
    cb: (err: grpc.ServiceError | null, resp?: OrderStatus) => void
  ) => void
}

let cachedClient: CoffeeShopClient | null = null

function getGrpcAddr() {
  return process.env.GRPC_ADDR || "localhost:9001"
}

function getProtoPath() {
  // Next runs from /web, so go up to repo root → proto/coffeeshop.proto
  return path.join(process.cwd(), "..", "proto", "coffeeshop.proto")
}

export function getCoffeeShopClient(): CoffeeShopClient {
  if (cachedClient) return cachedClient

  const packageDef = protoLoader.loadSync(getProtoPath(), {
    keepCase: false,
    longs: String,
    enums: String,
    defaults: true,
    oneofs: true,
  })

  const loaded = grpc.loadPackageDefinition(packageDef) as any
  const svc = loaded.coffeeshop?.CoffeeShop
  if (!svc) {
    throw new Error("Failed to load coffeeshop.CoffeeShop service from proto")
  }

  const client = new svc(
    getGrpcAddr(),
    grpc.credentials.createInsecure()
  ) as CoffeeShopClient

  cachedClient = client
  return client
}

export async function fetchMenu(): Promise<MenuItem[]> {
  const client = getCoffeeShopClient()
  const stream = client.GetMenu({})

  return await new Promise<MenuItem[]>((resolve, reject) => {
    const byId = new Map<string, MenuItem>()

    stream.on("data", (msg: any) => {
      const items: any[] = msg?.items || []
      for (const it of items) {
        if (!it?.id) continue
        byId.set(it.id, {
          id: String(it.id),
          name: String(it.name ?? ""),
          description: String(it.description ?? ""),
          price: Number(it.price ?? 0),
        })
      }
    })

    stream.on("error", (err: any) => reject(err))
    stream.on("end", () => resolve([...byId.values()]))
  })
}

export async function placeOrder(items: MenuItem[]): Promise<string> {
  const client = getCoffeeShopClient()
  return await new Promise<string>((resolve, reject) => {
    client.PlaceOrder({ items }, (err, resp) => {
      if (err) return reject(err)
      resolve(resp?.id || "")
    })
  })
}

export async function getOrderStatus(receiptId: string): Promise<string> {
  const client = getCoffeeShopClient()
  return await new Promise<string>((resolve, reject) => {
    client.GetOrderStatus({ id: receiptId }, (err, resp) => {
      if (err) return reject(err)
      resolve(resp?.status || "unknown")
    })
  })
}

