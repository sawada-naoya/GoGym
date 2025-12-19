import { proxyToGo } from "@/lib/server/proxy";

export const GET = (req: Request, { params }: { params: { id: string } }) =>
  proxyToGo(req, `/api/v1/gyms/${params.id}`);
