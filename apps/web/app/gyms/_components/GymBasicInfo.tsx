import { Gym } from "@/types/gym";
import FavoriteButton from "@/components/FavoriteButton";

type GymBasicInfoProps = {
  gym: Gym;
};

const GymBasicInfo = ({ gym }: GymBasicInfoProps) => {
  const averageRating = gym.average_rating ? gym.average_rating.toFixed(1) : "0.0";
  const reviewCount = gym.review_count;
  return (
    <div className="bg-white rounded-lg shadow p-6">
      <div className="flex items-start justify-between mb-4">
        <div>
          <h1 className="text-2xl md:text-3xl font-bold text-gray-900 mb-2">{gym.name}</h1>
          <div className="flex items-center space-x-4">
            <div className="flex items-center">
              <div className="flex items-center bg-booking-600 text-white px-2 py-1 rounded text-sm font-semibold">⭐ {averageRating}</div>
              <span className="ml-2 text-gray-600">({reviewCount}件のレビュー)</span>
            </div>
          </div>
        </div>
        <FavoriteButton gymId={gym.id} />
      </div>

      <p className="text-gray-700 leading-relaxed mb-4">{gym.description}</p>

      {/* 基本情報 */}
      <div className="flex items-center space-x-6 text-sm text-gray-600">
        {gym.prefecture && (
          <div>
            <span className="font-medium">都道府県:</span> {gym.prefecture}
          </div>
        )}
        {gym.city && (
          <div>
            <span className="font-medium">市区町村:</span> {gym.city}
          </div>
        )}
        {gym.tags && gym.tags.length > 0 && (
          <div>
            <span className="font-medium">タグ:</span> {gym.tags.map((tag) => tag.name).join(", ")}
          </div>
        )}
      </div>
    </div>
  );
};

export default GymBasicInfo;
