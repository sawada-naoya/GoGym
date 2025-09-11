import { Gym } from "@/types/gym";

type GymAmenitiesProps = {
  gym: Gym;
};

const GymAmenities = ({ gym }: GymAmenitiesProps) => {
  // モックの設備データ（実際のGym型にamenitiesプロパティがないため）
  const mockAmenities = [
    "フリーウェイト", "マシン", "有酸素マシン", 
    "シャワー", "更衣室", "ロッカー", 
    "プロテインバー", "パーソナルトレーニング", "グループレッスン"
  ];

  if (!mockAmenities || mockAmenities.length === 0) {
    return null;
  }

  return (
    <div className="bg-white rounded-lg shadow p-6">
      <h2 className="text-xl font-bold text-gray-900 mb-4">設備・サービス</h2>
      <div className="grid grid-cols-2 md:grid-cols-3 gap-3">
        {mockAmenities.map((amenity, index) => (
          <div key={index} className="flex items-center space-x-2">
            <div className="w-2 h-2 bg-booking-600 rounded-full"></div>
            <span className="text-gray-700">{amenity}</span>
          </div>
        ))}
      </div>
    </div>
  );
};

export default GymAmenities;