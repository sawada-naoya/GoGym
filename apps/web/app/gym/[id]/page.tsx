import { Gym } from "@/types/gym";
import Header from "@/components/Header";
import { notFound } from "next/navigation";
import GymPhotoGallery from "../components/GymPhotoGallery";
import GymBasicInfo from "../components/GymBasicInfo";
import GymAmenities from "../components/GymAmenities";
import GymAccessInfo from "../components/GymAccessInfo";
import GymReview from "../components/GymReview";
import GymContactSidebar from "../components/GymContactSidebar";

// レスポンスの型定義
type GymDetailResponse = {
  gym: Gym;
};

// ジム詳細を取得する関数（現在はモックデータ）
const fetchGymDetail = async (id: string): Promise<GymDetailResponse> => {
  // TODO: 後でAPIを実装
  // 現在はモックデータを返す
  const mockGym: Gym = {
    id: parseInt(id),
    name: "ゴールドジム 原宿東京",
    description: "本格的なウェイトトレーニングの聖地。豊富なフリーウェイトとマシンで、初心者から上級者まで満足できる設備。24時間営業で忙しい方にも最適。",
    address: "東京都渋谷区神宮前6-31-17 ベロックスビル B1F",
    city: "渋谷区",
    prefecture: "東京都",
    postal_code: "150-0001",
    location: {
      latitude: 35.670500,
      longitude: 139.702600
    },
    average_rating: 4.5,
    review_count: 128,
    tags: [
      { id: 1, name: "フリーウェイト" },
      { id: 2, name: "マシン" },
      { id: 3, name: "24時間営業" }
    ],
    created_at: "2024-01-01T00:00:00Z",
    updated_at: "2024-01-01T00:00:00Z"
  };

  return { gym: mockGym };
};

type PageProps = {
  params: {
    id: string;
  };
};

const GymDetailPage = async ({ params }: PageProps) => {
  const response = await fetchGymDetail(params.id);
  const gym = response.gym;

  if (!gym) {
    notFound();
  }

  return (
    <div className="min-h-screen bg-gray-50">
      <Header />
      <GymPhotoGallery gym={gym} />
      <div className="container mx-auto px-4 py-6">
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
          <div className="lg:col-span-2 space-y-6">
            <GymBasicInfo gym={gym} />
            <GymAmenities gym={gym} />
            <GymAccessInfo gym={gym} />
            <GymReview gym={gym} />
          </div>
          <GymContactSidebar gym={gym} />
        </div>
      </div>
    </div>
  );
};

export default GymDetailPage;