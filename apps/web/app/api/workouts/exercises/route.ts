import { proxyToGo, errorResponse } from "@/lib/server/proxy";

export const GET = async (req: Request) => {
  const url = new URL(req.url);
  const id = url.searchParams.get("id");
  const action = url.searchParams.get("action");

  if (!id) return errorResponse("id is required", 400);
  if (action !== "last") return errorResponse("invalid action", 400);

  return proxyToGo(req, `/api/v1/workouts/exercises/${id}/last`);
};

export const POST = async (req: Request) => {
  const body = await req.json().catch(() => null);
  return proxyToGo(req, "/api/v1/workouts/exercises/bulk", {
    method: "POST",
    body,
  });
};

export const DELETE = async (req: Request) => {
  const url = new URL(req.url);
  const id = url.searchParams.get("id");

  if (!id) return errorResponse("id is required", 400);

  return proxyToGo(req, `/api/v1/workouts/exercises/${id}`, {
    method: "DELETE",
  });
};
