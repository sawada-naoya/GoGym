import { proxyToGo } from "@/lib/server/proxy";

export const GET = (req: Request) => proxyToGo(req, "/api/v1/workouts/records");

export const POST = async (req: Request) => {
  const body = await req.json().catch(() => null);
  return proxyToGo(req, "/api/v1/workouts/records", { method: "POST", body });
};
