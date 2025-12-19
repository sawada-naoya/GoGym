import { proxyToGo } from "@/lib/server/proxy";

export const PUT = async (
  req: Request,
  { params }: { params: { id: string } },
) => {
  const body = await req.json().catch(() => null);
  return proxyToGo(req, `/api/v1/workouts/records/${params.id}`, {
    method: "PUT",
    body,
  });
};
