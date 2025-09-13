import { Gym } from "@/types/gym";

type GymAccessInfoProps = {
  gym: Gym;
};

const GymAccessInfo = ({ gym }: GymAccessInfoProps) => {
  return (
    <div className="bg-white rounded-lg shadow p-6">
      <h2 className="text-xl font-bold text-gray-900 mb-4">アクセス情報</h2>
      <div className="space-y-4">
        <div>
          <h3 className="font-semibold text-gray-900 mb-2">📍 住所</h3>
          <p className="text-gray-700">
            {[gym.address, gym.city, gym.prefecture].filter(Boolean).join(" ")}
          </p>
          {gym.postal_code && (
            <p className="text-sm text-gray-600">〒{gym.postal_code}</p>
          )}
        </div>
        
        <div>
          <h3 className="font-semibold text-gray-900 mb-2">🗺️ 位置情報</h3>
          <p className="text-gray-700">
            緯度: {gym.location.Latitude}, 経度: {gym.location.Longitude}
          </p>
        </div>
      </div>
    </div>
  );
};

export default GymAccessInfo;