import { NextResponse } from "next/server";
import { getServerAccessToken } from "@/lib/auth-helpers";

const API_BASE =
  process.env.API_INTERNAL_URL ?? process.env.NEXT_PUBLIC_API_URL;

export const GetLastExerciseRecord = async (id: string) => {
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

  const res = await fetch(`${API_BASE}/api/v1/workouts/exercises/${id}/last`, {
    headers: { Authorization: `Bearer ${token}` },
    cache: "no-store",
  });

  const data = await res.json().catch(() => null);
  return NextResponse.json(data, { status: res.status });
};

export const UpsertWorkoutExercises = async (body: any) => {
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

  const res = await fetch(`${API_BASE}/api/v1/workouts/exercises/bulk`, {
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

export const DeleteWorkoutExercise = async (id: string) => {
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

  const res = await fetch(`${API_BASE}/api/v1/workouts/exercises/${id}`, {
    method: "DELETE",
    headers: { Authorization: `Bearer ${token}` },
    cache: "no-store",
  });

  const data = await res.json().catch(() => null);
  return NextResponse.json(data, { status: res.status });
};

// Route Handlers（Client Component用）
export const GET = async (req: Request) => {
  const url = new URL(req.url);
  const id = url.searchParams.get("id");
  const action = url.searchParams.get("action");

  if (!id) {
    return NextResponse.json({ message: "id is required" }, { status: 400 });
  }
  if (action !== "last") {
    return NextResponse.json({ message: "invalid action" }, { status: 400 });
  }

  return GetLastExerciseRecord(id);
};

export const POST = async (req: Request) => {
  const body = await req.json().catch(() => null);
  return UpsertWorkoutExercises(body);
};

export const DELETE = async (req: Request) => {
  const url = new URL(req.url);
  const id = url.searchParams.get("id");

  if (!id) {
    return NextResponse.json({ message: "id is required" }, { status: 400 });
  }

  return DeleteWorkoutExercise(id);
};
