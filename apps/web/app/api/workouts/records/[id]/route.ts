import { NextResponse } from "next/server";
import { getServerAccessToken } from "@/lib/auth-helpers";

const API_BASE =
  process.env.API_INTERNAL_URL ?? process.env.NEXT_PUBLIC_API_URL;

export const UpdateWorkoutRecord = async (id: string, body: any) => {
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

  const res = await fetch(`${API_BASE}/api/v1/workouts/records/${id}`, {
    method: "PUT",
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

// Route Handler（Client Component用）
export const PUT = async (
  req: Request,
  { params }: { params: { id: string } },
) => {
  const body = await req.json().catch(() => null);
  return UpdateWorkoutRecord(params.id, body);
};
