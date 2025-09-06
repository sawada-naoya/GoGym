'use client'

import { useState } from 'react'
// アイコンは削除してシンプルにする
import Header from '@/components/Header'
import GymCard from '@/components/GymCard'

// ダミーデータ
const dummyGyms = [
  {
    id: '1',
    name: 'Gold\'s Gym 表参道店',
    location: '表参道・青山エリア',
    distance: '駅から徒歩3分',
    rating: 4.2,
    reviewCount: 128,
    price: '¥8,800',
    image: '',
    features: ['24時間営業', 'パーソナルトレーニング', 'サウナ', 'プール'],
    isFavorite: false
  },
  {
    id: '2', 
    name: 'エニタイムフィットネス 渋谷店',
    location: '渋谷エリア',
    distance: '駅から徒歩1分',
    rating: 4.0,
    reviewCount: 89,
    price: '¥6,980',
    image: '',
    features: ['24時間営業', 'マシン特化', '世界中で利用可能'],
    isFavorite: true
  },
  {
    id: '3',
    name: 'RIZAP 新宿店',
    location: '新宿エリア',
    distance: '駅から徒歩5分',
    rating: 4.7,
    reviewCount: 234,
    price: '¥29,800',
    image: '',
    features: ['パーソナル専門', '完全個室', '食事指導', '返金保証'],
    isFavorite: false
  },
  {
    id: '4',
    name: 'コナミスポーツクラブ 池袋',
    location: '池袋エリア',
    distance: '駅から徒歩2分',
    rating: 3.9,
    reviewCount: 156,
    price: '¥7,590',
    image: '',
    features: ['プール', 'スタジオ', 'サウナ', 'お風呂'],
    isFavorite: false
  },
  {
    id: '5',
    name: 'ジェクサー・フィットネスクラブ大井町',
    location: '大井町エリア',
    distance: '駅から徒歩1分',
    rating: 4.1,
    reviewCount: 98,
    price: '¥9,350',
    image: '',
    features: ['大型施設', 'プール', 'テニス', 'ゴルフ'],
    isFavorite: true
  },
  {
    id: '6',
    name: 'チョコザップ 恵比寿店',
    location: '恵比寿エリア',
    distance: '駅から徒歩30秒',
    rating: 3.8,
    reviewCount: 67,
    price: '¥2,980',
    image: '',
    features: ['コンビニジム', '24時間', 'セルフエステ', 'セルフ脱毛'],
    isFavorite: false
  }
]

export default function Home() {
  const [searchQuery, setSearchQuery] = useState('')
  const [selectedArea, setSelectedArea] = useState('')
  const [loading, setLoading] = useState(false)

  const handleSearch = async (e: React.FormEvent) => {
    e.preventDefault()
    setLoading(true)
    
    console.log('Search:', { searchQuery, selectedArea })
    
    setTimeout(() => {
      setLoading(false)
      alert(`「${searchQuery}」「${selectedArea}」で検索機能は開発中です`)
    }, 1000)
  }

  return (
    <div className="min-h-screen bg-gray-50">
      <Header />
      
      {/* ヒーローセクション */}
      <div className="bg-gradient-to-br from-booking-600 to-booking-800 relative overflow-hidden">
        <div className="absolute inset-0 bg-black/10"></div>
        
        <div className="relative container mx-auto px-4 py-16">
          <div className="text-center mb-12">
            <h1 className="text-4xl md:text-5xl font-bold text-white mb-4">
              あなたにぴったりのジムを見つけよう
            </h1>
            <p className="text-xl text-white/90 max-w-2xl mx-auto">
              東京都内の厳選されたジムから、あなたの条件に合った最適な場所を発見できます
            </p>
          </div>

          {/* 検索フォーム - Booking.com風 */}
          <div className="max-w-4xl mx-auto">
            <form onSubmit={handleSearch} className="bg-white rounded-xl shadow-lg p-6">
              <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                {/* ジム名・キーワード検索 */}
                <div className="space-y-2">
                  <label className="block text-sm font-medium text-gray-700">ジム名・キーワード</label>
                  <input
                    type="text"
                    placeholder="Gold's Gym, エニタイム..."
                    className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-booking-500 focus:border-transparent"
                    value={searchQuery}
                    onChange={(e) => setSearchQuery(e.target.value)}
                  />
                </div>

                {/* エリア選択 */}
                <div className="space-y-2">
                  <label className="block text-sm font-medium text-gray-700">エリア</label>
                  <select
                    className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-booking-500 focus:border-transparent appearance-none bg-white"
                    value={selectedArea}
                    onChange={(e) => setSelectedArea(e.target.value)}
                  >
                    <option value="">エリアを選択</option>
                    <option value="shibuya">渋谷・原宿</option>
                    <option value="shinjuku">新宿・代々木</option>
                    <option value="omotesando">表参道・青山</option>
                    <option value="ikebukuro">池袋・巣鴨</option>
                    <option value="ebisu">恵比寿・中目黒</option>
                    <option value="roppongi">六本木・赤坂</option>
                  </select>
                </div>

                {/* 検索ボタン */}
                <div className="space-y-2">
                  <label className="block text-sm font-medium text-transparent">検索</label>
                  <button
                    type="submit"
                    disabled={loading}
                    className="w-full bg-booking-600 hover:bg-booking-700 disabled:opacity-50 text-white font-semibold py-3 px-6 rounded-lg transition-colors duration-200"
                  >
                    {loading ? '検索中...' : '検索する'}
                  </button>
                </div>
              </div>

              {/* フィルターボタン */}
              <div className="flex items-center justify-center mt-4 pt-4 border-t border-gray-200">
                <button type="button" className="text-booking-600 hover:text-booking-700 font-medium">
                  詳細条件で絞り込み
                </button>
              </div>
            </form>
          </div>
        </div>
      </div>

      {/* おすすめジム一覧 */}
      <div className="container mx-auto px-4 py-12">
        <div className="flex items-center justify-between mb-8">
          <div>
            <h2 className="text-3xl font-bold text-gray-900 mb-2">おすすめのジム</h2>
            <p className="text-gray-600">東京の人気ジムをチェック</p>
          </div>
          <div className="text-sm text-gray-600">
            {dummyGyms.length}件のジムが見つかりました
          </div>
        </div>

        {/* ジムカードグリッド - Airbnb風 */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {dummyGyms.map((gym) => (
            <GymCard key={gym.id} gym={gym} />
          ))}
        </div>

        {/* もっと見るボタン */}
        <div className="text-center mt-12">
          <button className="bg-booking-600 hover:bg-booking-700 text-white font-semibold py-3 px-8 rounded-lg transition-colors duration-200">
            もっと見る
          </button>
        </div>
      </div>
    </div>
  )
}