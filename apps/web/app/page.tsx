import Link from "next/link";
import SearchForm from "@/components/SearchForm";
import GymCard from "@/components/GymCard";
import type { Gym } from "@/types/gym";
import { GET } from "@/lib/api";

export const dynamic = "force-dynamic";

// おすすめのジムを取得する関数
const fetchRecommendedGyms = async (): Promise<Gym[]> => {
  console.log("🔍 SSR fetch start: ", process.env.NEXT_PUBLIC_API_URL);

  const res = await GET<Gym[]>("api/v1/gyms/recommended", {
    query: { limit: 6 },
    cache: "no-store",
  });

  console.log("🔍 SSR fetch result:", {
    ok: res.ok,
    status: res.status,
    dataLength: res.ok && res.data ? res.data.length : 0,
  });

  return res.ok && res.data ? res.data : [];
};

const Home = async () => {
  const gyms = await fetchRecommendedGyms();

  return (
    <div className="min-h-screen">
      {/* ヒーロー */}
      <div className="bg-gradient-to-br from-booking-600 to-booking-800 relative overflow-hidden -mt-4">
        <div className="absolute inset-0 bg-black/10" />
        <div className="relative container mx-auto px-4 py-16 pt-20">
          <div className="text-center mb-12">
            <h1 className="text-4xl md:text-5xl font-bold text-white mb-4">あなたにぴったりのジムを見つけよう</h1>
            <p className="text-xl text-white/90">評価の高いおすすめジムを表示中</p>
          </div>
          <div className="max-w-4xl mx-auto">
            <SearchForm />
          </div>
        </div>
      </div>

      {/* ジム一覧 */}
      <div className="bg-gray-50 container mx-auto px-4 py-12">
        <div className="flex items-center justify-between mb-8">
          <div>
            <h2 className="text-3xl font-bold text-gray-900 mb-2">{gyms.length > 0 ? "⭐ おすすめのジム" : "新しいジムを追加中..."}</h2>
            <p className="text-gray-600">{gyms.length > 0 ? "評価の高い順に表示しています" : "データベースにジム情報を追加してください"}</p>
          </div>
          {gyms.length > 0 && <div className="text-sm text-gray-600">厳選された{gyms.length}件のジム</div>}
        </div>

        {gyms.length > 0 ? (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {gyms.map((gym) => (
              <GymCard key={gym.id} gym={gym} />
            ))}
          </div>
        ) : (
          <div className="text-center py-12 bg-white rounded-lg shadow">
            <p className="text-lg mb-2 text-gray-700">ジムデータが見つかりませんでした</p>
            <p className="text-sm text-gray-500">APIサーバーが起動しているか確認してください</p>
          </div>
        )}

        <div className="text-center mt-12">
          <Link href="/search" className="inline-block bg-booking-600 hover:bg-booking-700 text-white font-semibold py-3 px-8 rounded-lg transition-colors duration-200">
            {gyms.length > 0 ? "すべてのおすすめジムを見る" : "ジムを探す"}
          </Link>
        </div>
      </div>
    </div>
  );
};

export default Home;
