import { proxyToGoPublic, errorResponse } from "@/lib/server/proxy";

export const POST = async (req: Request) => {
  const body = await req.json().catch(() => null);
  if (!body) return errorResponse("Invalid request body", 400);

  return proxyToGoPublic(req, "/api/v1/users", { method: "POST", body });
};
