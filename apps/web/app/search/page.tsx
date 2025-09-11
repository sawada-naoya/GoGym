import { Suspense } from "react";
import Header from "@/components/Header";
import GymCard from "@/components/GymCard";
import SearchForm from "@/components/SearchForm";
import { SearchGymResponse, SearchGymParams } from "@/types/gym";
import { GET } from "@/lib/api";

// おすすめジム取得関数（評価順でソート）
const fetchRecommendedGyms = async (searchParams?: SearchGymParams): Promise<SearchGymResponse> => {
  try {
    // クエリパラメータを構築
    const queryParams: Record<string, any> = {
      limit: searchParams?.limit || 20,
    };

    if (searchParams?.q) queryParams.q = searchParams.q;
    if (searchParams?.lat) queryParams.lat = searchParams.lat;
    if (searchParams?.lon) queryParams.lon = searchParams.lon;
    if (searchParams?.radius_m) queryParams.radius_m = searchParams.radius_m;
    if (searchParams?.cursor) queryParams.cursor = searchParams.cursor;

    // 検索条件があるかどうかで呼び出すエンドポイントを変える
    const hasSearchCondition = searchParams?.q || searchParams?.lat || searchParams?.lon;

    if (hasSearchCondition) {
      // 検索条件がある場合は通常の検索エンドポイント
      return await GET<SearchGymResponse>("/api/v1/gyms", {
        query: queryParams,
        cache: "no-store",
      });
    } else {
      // 検索条件がない場合はおすすめジムエンドポイント
      return await GET<SearchGymResponse>("/api/v1/gyms/recommended", {
        query: { limit: queryParams.limit || 20 },
        cache: "no-store",
      });
    }
  } catch (error) {
    console.error("Failed to fetch recommended gyms:", error);
    return {
      gyms: [],
      next_cursor: null,
      has_more: false,
    };
  }
};

// ローディングコンポーネント
const SearchResultsSkeleton = () => {
  return (
    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      {Array.from({ length: 6 }).map((_, i) => (
        <div key={i} className="bg-white rounded-xl shadow-sm overflow-hidden animate-pulse">
          <div className="h-64 bg-gray-200"></div>
          <div className="p-4">
            <div className="h-4 bg-gray-200 rounded mb-2"></div>
            <div className="h-6 bg-gray-200 rounded mb-2"></div>
            <div className="h-4 bg-gray-200 rounded w-3/4"></div>
          </div>
        </div>
      ))}
    </div>
  );
};

// 検索結果コンポーネント
const SearchResults = async ({ searchParams }: { searchParams: SearchGymParams }) => {
  const { gyms, has_more } = await fetchRecommendedGyms(searchParams);

  // 検索条件の有無を判定
  const hasSearchQuery = searchParams.q || searchParams.lat || searchParams.lon;

  return (
    <div className="container mx-auto px-4 py-8">
      <div className="flex items-center justify-between mb-8">
        <div>
          <h1 className="text-3xl font-bold text-gray-900 mb-2">{hasSearchQuery ? "検索結果" : "おすすめのジム"}</h1>
          <p className="text-gray-600">{hasSearchQuery ? "条件に合うジムが見つかりました" : "評価の高い順に表示しています"}</p>
        </div>
        <div className="text-sm text-gray-600">{gyms.length}件のジムが見つかりました</div>
      </div>

      {/* フィルターバー（今後実装予定） */}
      <div className="bg-white rounded-lg shadow-sm p-4 mb-8">
        <div className="flex items-center gap-4 text-sm">
          <span className="text-gray-700 font-medium">並び替え:</span>
          <button className="px-3 py-1 bg-booking-600 text-white rounded-md text-sm">評価順</button>
          <button className="px-3 py-1 text-gray-600 hover:bg-gray-100 rounded-md text-sm">距離順</button>
          <button className="px-3 py-1 text-gray-600 hover:bg-gray-100 rounded-md text-sm">料金順</button>
        </div>
      </div>

      {/* ジム一覧 */}
      {gyms.length > 0 ? (
        <>
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {gyms.map((gym) => (
              <GymCard key={gym.id} gym={gym} />
            ))}
          </div>

          {/* もっと見るボタン */}
          {has_more && (
            <div className="text-center mt-12">
              <button className="bg-booking-600 hover:bg-booking-700 text-white font-semibold py-3 px-8 rounded-lg transition-colors duration-200">もっと見る</button>
            </div>
          )}
        </>
      ) : (
        <div className="text-center py-12 bg-white rounded-lg shadow">
          <div className="text-gray-500">
            <p className="text-lg mb-2">{hasSearchQuery ? "条件に合うジムが見つかりませんでした" : "ジムデータが見つかりませんでした"}</p>
            <p className="text-sm">{hasSearchQuery ? "検索条件を変更してお試しください" : "APIサーバーが起動しているか確認してください"}</p>
          </div>
        </div>
      )}
    </div>
  );
};

interface SearchPageProps {
  searchParams: {
    q?: string;
    lat?: string;
    lon?: string;
    radius_m?: string;
    cursor?: string;
    limit?: string;
  };
}

const SearchPage = ({ searchParams }: SearchPageProps) => {
  // URLパラメータをSearchGymParamsに変換
  const searchQuery: SearchGymParams = {
    q: searchParams.q,
    lat: searchParams.lat ? parseFloat(searchParams.lat) : undefined,
    lon: searchParams.lon ? parseFloat(searchParams.lon) : undefined,
    radius_m: searchParams.radius_m ? parseInt(searchParams.radius_m) : undefined,
    cursor: searchParams.cursor,
    limit: searchParams.limit ? parseInt(searchParams.limit) : undefined,
  };

  return (
    <div className="min-h-screen bg-gray-50">
      <Header />

      {/* 検索バー */}
      <div className="bg-white border-b border-gray-200">
        <div className="container mx-auto px-4 py-6">
          <SearchForm />
        </div>
      </div>

      {/* 検索結果 */}
      <Suspense fallback={<SearchResultsSkeleton />}>
        <SearchResults searchParams={searchQuery} />
      </Suspense>
    </div>
  );
};

export default SearchPage;
