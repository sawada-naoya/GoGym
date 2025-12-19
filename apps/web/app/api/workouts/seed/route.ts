import { proxyToGo } from "@/lib/server/proxy";

export const POST = (req: Request) =>
  proxyToGo(req, "/api/v1/workouts/seed", { method: "POST" });
