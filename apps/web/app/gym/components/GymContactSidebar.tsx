import { Gym } from "@/types/gym";

type GymContactSidebarProps = {
  gym: Gym;
};

const GymContactSidebar = ({ gym }: GymContactSidebarProps) => {
  // モックデータ（実際のGym型にはこれらのプロパティがないため）
  const mockWebsite = "https://www.goldsgym.jp/shop/13009";
  const mockPhoneNumber = "03-3797-4848";
  const mockPriceMin = 7700;
  const mockPriceMax = 15400;

  return (
    <div className="bg-white rounded-lg shadow p-6 sticky top-6 h-fit">
      <div className="text-center mb-4">
        <div className="text-2xl font-bold text-booking-600 mb-1">
          ¥{mockPriceMin.toLocaleString()}
          {mockPriceMax && mockPriceMax !== mockPriceMin && (
            <span className="text-lg">〜</span>
          )}
        </div>
        <p className="text-sm text-gray-600">月額料金から</p>
      </div>
      
      <div className="space-y-2">
        <a
          href={mockWebsite}
          target="_blank"
          rel="noopener noreferrer"
          className="w-full block bg-booking-600 hover:bg-booking-700 text-white font-semibold py-2.5 px-4 rounded-lg transition-colors text-center text-sm"
        >
          🌐 公式サイトで詳細を見る
        </a>
        
        <a
          href={`tel:${mockPhoneNumber}`}
          className="w-full block bg-white border-2 border-booking-600 text-booking-600 hover:bg-booking-50 font-semibold py-2.5 px-4 rounded-lg transition-colors text-center text-sm"
        >
          📞 {mockPhoneNumber}
        </a>
      </div>
      
      <div className="mt-4 pt-3 border-t border-gray-200">
        <p className="text-xs text-gray-500 text-center leading-tight">
          料金やプランの詳細は直接お問い合わせください
        </p>
      </div>
    </div>
  );
};

export default GymContactSidebar;