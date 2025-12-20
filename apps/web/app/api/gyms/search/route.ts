import { NextResponse } from "next/server";
import { getServerAccessToken } from "@/lib/auth-helpers";

const API_BASE =
  process.env.API_INTERNAL_URL ?? process.env.NEXT_PUBLIC_API_URL;

export const SearchGyms = async (searchParams?: URLSearchParams) => {
  const token = await getServerAccessToken();
  if (!token) {
    return NextResponse.json({ message: "unauthorized" }, { status: 401 });
  }
  if (!API_BASE) {
    return NextResponse.json(
      { message: "API base url is not configured" },
      { status: 500 },
    );
  }

  const queryString = searchParams?.toString()
    ? `?${searchParams.toString()}`
    : "";
  const res = await fetch(`${API_BASE}/api/v1/gyms/search${queryString}`, {
    headers: { Authorization: `Bearer ${token}` },
    cache: "no-store",
  });

  const data = await res.json().catch(() => null);
  return NextResponse.json(data, { status: res.status });
};

// Route Handler（Client Component用）
export const GET = async (req: Request) => {
  const url = new URL(req.url);
  return SearchGyms(url.searchParams);
};
