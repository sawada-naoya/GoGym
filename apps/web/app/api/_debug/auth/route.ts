import { NextResponse } from "next/server";
import { getServerAccessTokenWithDiag } from "@/features/auth/server/accessToken";

export const GET = async () => {
  if (process.env.DEBUG_AUTH !== "1") {
    return NextResponse.json({ message: "Not Found" }, { status: 404 });
  }
  const diag = await getServerAccessTokenWithDiag();

  // accessToken は返さない（metaだけ返す）
  return NextResponse.json({
    reason: diag.reason,
    meta: diag.meta,
    error: diag.error ?? null,
  });
};
