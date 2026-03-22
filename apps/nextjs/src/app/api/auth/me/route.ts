import { NextResponse } from "next/server";

import { authulaClientServer } from "@/lib/authula-client-server";

export async function GET() {
  try {
    const response = await authulaClientServer.getMe();
    return NextResponse.json(response);
  } catch (error: any) {
    return NextResponse.json({ message: error?.message }, { status: 500 });
  }
}
