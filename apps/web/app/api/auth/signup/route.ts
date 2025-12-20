import { NextResponse } from "next/server";

const API_BASE =
  process.env.API_INTERNAL_URL ?? process.env.NEXT_PUBLIC_API_URL;

export const SignupUser = async (body: any) => {
  if (!API_BASE) {
    return NextResponse.json(
      { message: "API base url is not configured" },
      { status: 500 },
    );
  }

  if (!body) {
    return NextResponse.json(
      { message: "Invalid request body" },
      { status: 400 },
    );
  }

  const res = await fetch(`${API_BASE}/api/v1/users`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(body),
    cache: "no-store",
  });

  const data = await res.json().catch(() => null);
  return NextResponse.json(data, { status: res.status });
};

// Route Handler（Client Component用）
export const POST = async (req: Request) => {
  const body = await req.json().catch(() => null);
  return SignupUser(body);
};
