import Link from 'next/link'
import { Gym } from '@/types/gym'

interface GymCardProps {
  gym: Gym
}

const GymCard = ({ gym }: GymCardProps) => {
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
            ♡
          </button>

          {/* エリアバッジ */}
          <div className="absolute bottom-3 left-3 bg-white/90 px-2 py-1 rounded-md text-sm font-medium text-gray-800">
            {gym.city || gym.prefecture || '東京'}
          </div>
        </div>

        {/* 情報部分 */}
        <div className="p-4">
          {/* アドレス */}
          <div className="text-sm text-gray-600 mb-1">
            <span>{gym.address}</span>
          </div>

          {/* ジム名 */}
          <h3 className="font-semibold text-gray-900 text-lg mb-2 group-hover:text-booking-700 transition-colors">
            {gym.name}
          </h3>

          {/* レーティング */}
          <div className="flex items-center mb-2">
            <div className="flex items-center">
              <span className="text-yellow-400 mr-1">★</span>
              <span className="text-sm font-medium text-gray-900">
                {gym.average_rating?.toFixed(1) || 'N/A'}
              </span>
              <span className="ml-1 text-sm text-gray-600">({gym.review_count}件のレビュー)</span>
            </div>
          </div>

          {/* 説明文 */}
          {gym.description && (
            <div className="mb-3 text-sm text-gray-600 line-clamp-2">
              {gym.description}
            </div>
          )}

          {/* タグ */}
          <div className="mb-3">
            <div className="flex flex-wrap gap-1">
              {gym.tags.slice(0, 3).map((tag) => (
                <span 
                  key={tag.id}
                  className="px-2 py-1 bg-booking-50 text-booking-700 text-xs rounded-md"
                >
                  {tag.name}
                </span>
              ))}
              {gym.tags.length > 3 && (
                <span className="px-2 py-1 bg-gray-100 text-gray-600 text-xs rounded-md">
                  +{gym.tags.length - 3}
                </span>
              )}
            </div>
          </div>

          {/* アクションボタン */}
          <div className="flex justify-end">
            <button className="text-booking-700 hover:text-booking-800 text-sm font-medium transition-colors">
              詳細を見る →
            </button>
          </div>
        </div>
      </div>
    </Link>
  )
}

export default GymCard