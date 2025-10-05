import { Suspense } from "react";
import LoginClient from "./_components/LoginClient";

export const dynamic = "force-dynamic";

const LoginPage = () => {
  return (
    <Suspense fallback={null}>
      <LoginClient />
    </Suspense>
  );
};

export default LoginPage;
