import { NextResponse } from "next/server"

import { getOrderStatus } from "@/lib/grpc"

export const runtime = "nodejs"

export async function GET(
  _req: Request,
  { params }: { params: Promise<{ id: string }> }
) {
  try {
    const { id } = await params
    if (!id) {
      return NextResponse.json({ error: "Missing order id" }, { status: 400 })
    }
    const status = await getOrderStatus(id)
    return NextResponse.json({ orderId: id, status }, { headers: { "Cache-Control": "no-store" } })
  } catch (err: any) {
    return NextResponse.json(
      { error: err?.message ?? "Failed to fetch order status" },
      { status: 500 }
    )
  }
}

