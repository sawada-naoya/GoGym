import LoginFormContent from "./[locale]/(auth)/login/content";
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
        <div className="grid md:grid-cols-2 gap-12 items-center">
          {/* 左側：説明 */}
          <div className="space-y-6">
            <h1 className="text-5xl font-bold text-gray-900">GoGym</h1>
            <p className="text-2xl text-gray-700">日々のトレーニングを記録できるアプリです</p>
            <div className="space-y-4">
              <div className="flex items-start gap-3">
                <div className="flex-shrink-0 w-8 h-8 bg-booking-600 rounded-full flex items-center justify-center">
                  <svg className="w-5 h-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                  </svg>
                </div>
                <div>
                  <h3 className="font-semibold text-gray-900">トレーニング記録</h3>
                  <p className="text-gray-600">セット数、重量、回数を簡単に記録</p>
                </div>
              </div>
              <div className="flex items-start gap-3">
                <div className="flex-shrink-0 w-8 h-8 bg-booking-600 rounded-full flex items-center justify-center">
                  <svg className="w-5 h-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                  </svg>
                </div>
                <div>
                  <h3 className="font-semibold text-gray-900">進捗管理</h3>
                  <p className="text-gray-600">日々の成長を可視化して確認</p>
                </div>
              </div>
              <div className="flex items-start gap-3">
                <div className="flex-shrink-0 w-8 h-8 bg-booking-600 rounded-full flex items-center justify-center">
                  <svg className="w-5 h-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                  </svg>
                </div>
                <div>
                  <h3 className="font-semibold text-gray-900">シンプルで使いやすい</h3>
                  <p className="text-gray-600">直感的な操作で素早く記録</p>
                </div>
              </div>
            </div>
          </div>

          {/* 右側：ログインフォーム */}
          <div className="bg-white rounded-2xl shadow-xl p-8">
            <h2 className="text-2xl font-bold text-gray-900 mb-6 text-center">今すぐトレーニングを開始する</h2>

            <LoginFormContent showHeader={false} showSignupLink={true} />
          </div>
        </div>
      </div>
    </div>
  );
};

export default Home;
