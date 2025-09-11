import { Gym } from "@/types/gym";
import Header from "@/components/Header";
import { notFound } from "next/navigation";

// レビューの型定義
type Review = {
  id: number;
  userName: string;
  rating: number;
  comment: string;
  visitDate: string;
  isVerified: boolean;
};

// レスポンスの型定義
type GymDetailResponse = {
  gym: Gym;
};

type ReviewsResponse = {
  reviews: Review[];
};

// ジム詳細を取得する関数（モック）
const fetchGymDetail = async (id: string): Promise<GymDetailResponse> => {
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

// レビューを取得する関数（モック）
const fetchReviews = async (gymId: string): Promise<ReviewsResponse> => {
  const mockReviews: Review[] = [
    {
      id: 1,
      userName: "山田太郎",
      rating: 5,
      comment: "設備が非常に充実していて、スタッフの対応も素晴らしいです。24時間営業なので仕事帰りにも通いやすく、とても満足しています。フリーウェイトエリアが広く、混雑時でも快適にトレーニングできます。",
      visitDate: "2024-08-15",
      isVerified: true
    },
    {
      id: 2,
      userName: "佐藤花子",
      rating: 4,
      comment: "マシンの種類が豊富で、混雑時でも待ち時間が少ないのが良いです。更衣室も清潔で使いやすいです。女性専用エリアがあるのも安心できます。ただし、料金がもう少し安ければ嬉しいです。",
      visitDate: "2024-08-10",
      isVerified: true
    },
    {
      id: 3,
      userName: "田中健一",
      rating: 5,
      comment: "パーソナルトレーナーのアドバイスが的確で、効率よくトレーニングできています。立地も良く通いやすいです。プロテインバーのメニューも豊富で、トレーニング後の補給に助かっています。",
      visitDate: "2024-08-05",
      isVerified: false
    },
    {
      id: 4,
      userName: "鈴木美咲",
      rating: 4,
      comment: "女性専用エリアがあるので安心してトレーニングできます。プロテインバーのメニューも豊富で便利です。スタッフの方々も親切で、初心者の私にも丁寧に教えてくれました。",
      visitDate: "2024-07-28",
      isVerified: true
    },
    {
      id: 5,
      userName: "高橋誠",
      rating: 5,
      comment: "原宿駅から近くてアクセスが良いです。設備も最新で、特に有酸素マシンが充実しています。24時間営業なので自分のペースで通えるのが最高です。",
      visitDate: "2024-07-20",
      isVerified: true
    },
    {
      id: 6,
      userName: "渡辺麻衣",
      rating: 4,
      comment: "清潔感があり、設備も充実していて満足です。グループレッスンも楽しく、モチベーションが上がります。ただし、人気の時間帯は混雑するので、時間を選んで利用しています。",
      visitDate: "2024-07-15",
      isVerified: true
    }
  ];

  return { reviews: mockReviews };
};

type PageProps = {
  params: {
    id: string;
  };
};

const GymReviewsPage = async ({ params }: PageProps) => {
  const gymResponse = await fetchGymDetail(params.id);
  const reviewsResponse = await fetchReviews(params.id);

  const gym = gymResponse.gym;
  const reviews = reviewsResponse.reviews;

  if (!gym) {
    notFound();
  }

  const renderStars = (rating: number) => {
    return Array.from({ length: 5 }, (_, i) => (
      <span key={i} className={i < rating ? "text-yellow-400" : "text-gray-300"}>
        ★
      </span>
    ));
  };

  return (
    <div className="min-h-screen bg-gray-50">
      <Header />

      {/* ページヘッダー */}
      <div className="bg-white shadow">
        <div className="container mx-auto px-4 py-6">
          <nav className="flex items-center space-x-2 text-sm text-gray-500 mb-4">
            <a href="/" className="hover:text-booking-600">ホーム</a>
            <span>›</span>
            <a href={`/gym/${gym.id}`} className="hover:text-booking-600">{gym.name}</a>
            <span>›</span>
            <span className="text-gray-900">レビュー</span>
          </nav>
          
          <div className="flex items-center justify-between">
            <div>
              <h1 className="text-3xl font-bold text-gray-900 mb-2">
                {gym.name}のレビュー
              </h1>
              <div className="flex items-center space-x-4">
                <div className="flex items-center">
                  <div className="bg-booking-600 text-white px-3 py-2 rounded font-bold text-lg">
                    {gym.average_rating?.toFixed(1) || "0.0"}
                  </div>
                  <span className="ml-3 text-gray-600">
                    {gym.review_count}件のレビュー
                  </span>
                </div>
              </div>
            </div>
            
            <a
              href={`/gym/${gym.id}`}
              className="bg-booking-600 hover:bg-booking-700 text-white font-semibold py-2 px-4 rounded-lg transition-colors"
            >
              ジム詳細へ戻る
            </a>
          </div>
        </div>
      </div>

      {/* レビュー一覧 */}
      <div className="container mx-auto px-4 py-8">
        <div className="max-w-4xl mx-auto">
          <div className="bg-white rounded-lg shadow">
            <div className="p-6">
              <div className="flex items-center justify-between mb-6">
                <h2 className="text-xl font-bold text-gray-900">
                  全{reviews.length}件のレビュー
                </h2>
                <select className="border border-gray-300 rounded-lg px-3 py-2 text-sm">
                  <option>新しい順</option>
                  <option>古い順</option>
                  <option>評価が高い順</option>
                  <option>評価が低い順</option>
                </select>
              </div>

              <div className="space-y-6">
                {reviews.map((review, index) => (
                  <div key={review.id} className={`${index !== reviews.length - 1 ? 'border-b border-gray-200' : ''} pb-6`}>
                    <div className="flex items-start space-x-4">
                      {/* アバター */}
                      <div className="w-12 h-12 bg-booking-600 rounded-full flex items-center justify-center text-white font-semibold">
                        {review.userName.charAt(0)}
                      </div>
                      
                      {/* レビュー内容 */}
                      <div className="flex-1">
                        <div className="flex items-center justify-between mb-3">
                          <div>
                            <h4 className="font-semibold text-gray-900 text-lg">{review.userName}</h4>
                            {review.isVerified && (
                              <span className="text-sm text-green-600 flex items-center mt-1">
                                ✓ 認証済み
                              </span>
                            )}
                          </div>
                          <div className="text-right">
                            <div className="flex items-center text-sm mb-1">
                              {renderStars(review.rating)}
                            </div>
                            <p className="text-xs text-gray-500">
                              {new Date(review.visitDate).toLocaleDateString('ja-JP')}
                            </p>
                          </div>
                        </div>
                        
                        <p className="text-gray-700 leading-relaxed">
                          {review.comment}
                        </p>
                      </div>
                    </div>
                  </div>
                ))}
              </div>

              {/* ページネーション（今後の実装用） */}
              <div className="mt-8 pt-6 border-t border-gray-200">
                <div className="flex justify-center">
                  <p className="text-sm text-gray-600">
                    {reviews.length}件中{reviews.length}件を表示
                  </p>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default GymReviewsPage;