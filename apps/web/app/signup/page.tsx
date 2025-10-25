import { Suspense } from "react";
import SignupClient from "../_components/auth/SignupForm";

export const dynamic = "force-dynamic";

const SignUpPage = () => {
  return (
    <Suspense fallback={null}>
      <SignupClient />
    </Suspense>
  );
};

export default SignUpPage;
