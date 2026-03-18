import { NextResponse } from "next/server"

import { fetchMenu } from "@/lib/grpc"

export const runtime = "nodejs"

export async function GET() {
  try {
    const items = await fetchMenu()
    return NextResponse.json(
      { items },
      {
        headers: {
          "Cache-Control": "no-store",
        },
      }
    )
  } catch (err: any) {
    return NextResponse.json(
      { error: err?.message ?? "Failed to fetch menu" },
      { status: 500 }
    )
  }
}

