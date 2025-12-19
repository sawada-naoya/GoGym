import { proxyToGo } from "@/lib/server/proxy";

export const GET = (req: Request) => proxyToGo(req, "/api/v1/workouts/parts");
