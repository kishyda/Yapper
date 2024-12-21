import { NextRequest, NextResponse } from "next/server";

export async function POST(req: NextRequest) {
    const body = await req.text();
    const json = await JSON.stringify(body);
    console.log(json);
    return new NextResponse();
}