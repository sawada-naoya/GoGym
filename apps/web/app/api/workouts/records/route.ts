import { NextResponse } from "next/server";
import { getServerAccessToken } from "@/lib/auth-helpers";

const API_BASE =
  process.env.API_INTERNAL_URL ?? process.env.NEXT_PUBLIC_API_URL;

export const GetWorkoutRecords = async (params?: {
  date?: string;
  partId?: number;
}) => {
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

  const searchParams = new URLSearchParams();
  if (params?.date) searchParams.set("date", params.date);
  if (params?.partId) searchParams.set("part_id", String(params.partId));

  const queryString = searchParams.toString()
    ? `?${searchParams.toString()}`
    : "";
  const res = await fetch(`${API_BASE}/api/v1/workouts/records${queryString}`, {
    headers: { Authorization: `Bearer ${token}` },
    cache: "no-store",
  });

  const data = await res.json().catch(() => null);
  return NextResponse.json(data, { status: res.status });
};

export const CreateWorkoutRecord = async (body: any) => {
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

  const res = await fetch(`${API_BASE}/api/v1/workouts/records`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      Authorization: `Bearer ${token}`,
    },
    body: JSON.stringify(body),
    cache: "no-store",
  });

  const data = await res.json().catch(() => null);
  return NextResponse.json(data, { status: res.status });
};

// Route Handlers（Client Component用）
export const GET = async (req: Request) => {
  const url = new URL(req.url);
  const date = url.searchParams.get("date") || undefined;
  const partId = url.searchParams.get("part_id")
    ? Number(url.searchParams.get("part_id"))
    : undefined;

  return GetWorkoutRecords({ date, partId });
};

export const POST = async (req: Request) => {
  const body = await req.json().catch(() => null);
  return CreateWorkoutRecord(body);
};
