import LoginFormContent from "./login/content";
import LandingHero from "./LandingHero";
import { redirect } from "next/navigation";
import { auth } from "@/features/auth/nextauth/auth";

const Home = async () => {
  const session = await auth();

  // ログイン済みの場合はworkoutページにリダイレクト
  if (session?.user) {
    redirect("/workout");
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-booking-50 to-blue-50">
      {/* ヒーローセクション */}
      <div className="container mx-auto px-4 py-16">
        <LandingHero>
          {/* 右側：ログインフォーム */}
          <div className="bg-white rounded-2xl shadow-xl p-8">
            <LoginFormContent showHeader={false} showSignupLink={true} />
          </div>
        </LandingHero>
      </div>
    </div>
  );
};

export default Home;
