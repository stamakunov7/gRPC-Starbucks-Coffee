import { NextResponse } from "next/server"

import { MenuItem, placeOrder } from "@/lib/grpc"

export const runtime = "nodejs"

type PlaceOrderBody = {
  items: Array<{ item: MenuItem; quantity: number }>
}

export async function POST(req: Request) {
  try {
    const body = (await req.json()) as PlaceOrderBody
    const expanded: MenuItem[] = []

    for (const entry of body?.items ?? []) {
      const qty = Math.max(0, Number(entry.quantity ?? 0))
      if (!entry?.item?.id || qty <= 0) continue
      for (let i = 0; i < qty; i++) expanded.push(entry.item)
    }

    if (expanded.length === 0) {
      return NextResponse.json(
        { error: "No items to order" },
        { status: 400 }
      )
    }

    const receiptId = await placeOrder(expanded)
    if (!receiptId) {
      return NextResponse.json(
        { error: "Empty receipt id from server" },
        { status: 502 }
      )
    }

    return NextResponse.json({ receiptId }, { headers: { "Cache-Control": "no-store" } })
  } catch (err: any) {
    return NextResponse.json(
      { error: err?.message ?? "Failed to place order" },
      { status: 500 }
    )
  }
}

