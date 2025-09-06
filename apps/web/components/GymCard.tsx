import Link from 'next/link'

interface GymCardProps {
  gym: {
    id: string
    name: string
    location: string
    distance: string
    rating: number
    reviewCount: number
    price: string
    image: string
    features: string[]
    isFavorite?: boolean
  }
}

export default function GymCard({ gym }: GymCardProps) {
  return (
    <Link href={`/gyms/${gym.id}`}>
      <div className="bg-white rounded-xl shadow-sm hover:shadow-md transition-all duration-300 overflow-hidden group cursor-pointer">
        {/* 画像部分 */}
        <div className="relative h-64 bg-gray-200">
          {/* プレースホルダー画像 */}
          <div className="w-full h-full bg-gradient-to-br from-booking-100 to-booking-200 flex items-center justify-center">
            <div className="text-center">
              <div className="text-4xl font-bold text-booking-600 mb-2">GYM</div>
              <div className="text-sm text-booking-500">No Image</div>
            </div>
          </div>
          
          {/* お気に入りボタン */}
          <button 
            className="absolute top-3 right-3 px-3 py-1 bg-white/80 hover:bg-white rounded text-sm font-medium text-gray-800 shadow-sm transition-colors"
            onClick={(e) => {
              e.preventDefault()
              console.log('Toggle favorite')
            }}
          >
            {gym.isFavorite ? '♥' : '♡'}
          </button>

          {/* 距離バッジ */}
          <div className="absolute bottom-3 left-3 bg-white/90 px-2 py-1 rounded-md text-sm font-medium text-gray-800">
            {gym.distance}
          </div>
        </div>

        {/* 情報部分 */}
        <div className="p-4">
          {/* エリア */}
          <div className="text-sm text-gray-600 mb-1">
            <span>{gym.location}</span>
          </div>

          {/* ジム名 */}
          <h3 className="font-semibold text-gray-900 text-lg mb-2 group-hover:text-booking-700 transition-colors">
            {gym.name}
          </h3>

          {/* レーティング */}
          <div className="flex items-center mb-2">
            <div className="flex items-center">
              <span className="text-yellow-400 mr-1">★</span>
              <span className="text-sm font-medium text-gray-900">{gym.rating}</span>
              <span className="ml-1 text-sm text-gray-600">({gym.reviewCount}件のレビュー)</span>
            </div>
          </div>

          {/* 特徴・設備 */}
          <div className="mb-3">
            <div className="flex flex-wrap gap-1">
              {gym.features.slice(0, 3).map((feature, index) => (
                <span 
                  key={index}
                  className="px-2 py-1 bg-booking-50 text-booking-700 text-xs rounded-md"
                >
                  {feature}
                </span>
              ))}
              {gym.features.length > 3 && (
                <span className="px-2 py-1 bg-gray-100 text-gray-600 text-xs rounded-md">
                  +{gym.features.length - 3}
                </span>
              )}
            </div>
          </div>

          {/* 価格 */}
          <div className="flex justify-between items-end">
            <div>
              <span className="text-lg font-bold text-gray-900">{gym.price}</span>
              <span className="text-sm text-gray-600 ml-1">/月</span>
            </div>
            <button className="text-booking-700 hover:text-booking-800 text-sm font-medium transition-colors">
              詳細を見る →
            </button>
          </div>
        </div>
      </div>
    </Link>
  )
}