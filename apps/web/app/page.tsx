import Header from '@/components/Header'
import GymCard from '@/components/GymCard'
import SearchForm from '@/components/SearchForm'
import { SearchGymResponse } from '@/types/gym'
import { GET } from '@/lib/api'


// おすすめのジムを取得する関数
const fetchRecommendedGyms = async (): Promise<SearchGymResponse> => {
  try {
    return await GET<SearchGymResponse>('/api/v1/gyms/recommended', {
      query: { limit: 6 }, // トップページでは6件程度に制限
      cache: 'no-store'
    })
  } catch (error) {
    console.error('Failed to fetch recommended gyms:', error)
    // エラー時は空のレスポンスを返す
    return {
      gyms: [],
      next_cursor: null,
      has_more: false
    }
  }
}

const Home = async () => {
  const { gyms } = await fetchRecommendedGyms()

  return (
    <div className="min-h-screen bg-gray-50">
      <Header />

      <div className="bg-gradient-to-br from-booking-600 to-booking-800 relative overflow-hidden">
        <div className="absolute inset-0 bg-black/10"></div>

        <div className="relative container mx-auto px-4 py-16">
          <div className="text-center mb-12">
            <h1 className="text-4xl md:text-5xl font-bold text-white mb-4">
              あなたにぴったりのジムを見つけよう
            </h1>
            <p className="text-xl text-white/90">
              評価の高いおすすめジムを表示中
            </p>
          </div>

          <div className="max-w-4xl mx-auto">
            <SearchForm />
          </div>
        </div>
      </div>

      {/* ジム一覧 */}
      <div className="container mx-auto px-4 py-12">
        <div className="flex items-center justify-between mb-8">
          <div>
            <h2 className="text-3xl font-bold text-gray-900 mb-2">
              {gyms.length > 0 ? '⭐ おすすめのジム' : '新しいジムを追加中...'}
            </h2>
            <p className="text-gray-600">
              {gyms.length > 0 ? '評価の高い順に表示しています' : 'データベースにジム情報を追加してください'}
            </p>
          </div>
          <div className="text-sm text-gray-600">
            {gyms.length > 0 && `厳選された${gyms.length}件のジム`}
          </div>
        </div>

        {/* ジムカードグリッド */}
        {gyms.length > 0 ? (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {gyms.map((gym) => (
              <GymCard key={gym.id} gym={gym} />
            ))}
          </div>
        ) : (
          <div className="text-center py-12 bg-white rounded-lg shadow">
            <div className="text-gray-500">
              <p className="text-lg mb-2">ジムデータが見つかりませんでした</p>
              <p className="text-sm">APIサーバーが起動しているか確認してください</p>
            </div>
          </div>
        )}

        {/* もっと見る・すべて見るボタン */}
        <div className="text-center mt-12">
          {gyms.length > 0 ? (
            <a
              href="/search"
              className="inline-block bg-booking-600 hover:bg-booking-700 text-white font-semibold py-3 px-8 rounded-lg transition-colors duration-200"
            >
              すべてのおすすめジムを見る
            </a>
          ) : (
            <a
              href="/search"
              className="inline-block bg-booking-600 hover:bg-booking-700 text-white font-semibold py-3 px-8 rounded-lg transition-colors duration-200"
            >
              ジムを探す
            </a>
          )}
        </div>
      </div>
    </div>
  )
}

export default Home