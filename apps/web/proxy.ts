import { auth } from "@/features/auth/nextauth/auth";
import { NextResponse } from "next/server";

export default auth((req) => {
  // 認証が必要なページで未認証の場合はログインページへ
  const isAuthRequired = req.nextUrl.pathname.startsWith("/workout");

  if (isAuthRequired && !req.auth) {
    const url = new URL("/", req.url);
    return NextResponse.redirect(url);
  }

  return NextResponse.next();
});

export const config = {
  matcher: [
    /*
     * Match all request paths except for the ones starting with:
     * - api/auth (authentication endpoints)
     * - _next/static (static files)
     * - _next/image (image optimization files)
     * - favicon.ico (favicon file)
     * - public folder
     */
    "/((?!api/auth|_next/static|_next/image|favicon.ico|.*\\.png$).*)",
  ],
};
