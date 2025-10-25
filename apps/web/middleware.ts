import { withAuth } from "next-auth/middleware";

export default withAuth({
  pages: {
    signIn: "/",
  },
});

export const config = {
  matcher: ["/:user_id/workout/:path*", "/dashboard/:path*", "/settings/:path*"],
};
