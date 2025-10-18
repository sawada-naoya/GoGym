import { Gym } from "@/types/gym";
import { ReviewListResponse } from "@/types/review";
import { GET } from "@/lib/api";
import Header from "@/components/Header";
import { notFound } from "next/navigation";
import GymPhotoGallery from "../components/GymPhotoGallery";
import GymBasicInfo from "../components/GymBasicInfo";
import GymAmenities from "../components/GymAmenities";
import GymAccessInfo from "../components/GymAccessInfo";
import GymReview from "../components/GymReview";
import GymContactSidebar from "../components/GymContactSidebar";

export const dynamic = "force-dynamic";

type PageProps = {
  params: {
    id: string;
  };
};

const fetchGym = async (id: string): Promise<Gym | null> => {
  const fallback: Gym = {
    id: 0,
    name: "",
    description: null,
    location: { latitude: 0, longitude: 0 },
    address: "",
    city: null,
    prefecture: null,
    postal_code: null,
    is_active: false,
    tags: [],
    average_rating: null,
    review_count: 0,
    created_at: "",
    updated_at: "",
  };
  try {
    const res = await GET<Gym>(`/api/v1/gyms/${id}`, {
      cache: "no-store",
    });
    if (!res.ok || !res.data) return fallback;
    return res.data;
  } catch (error) {
    return null;
  }
};

const fetchGymReviews = async (id: string): Promise<ReviewListResponse | null> => {
  const fallback: ReviewListResponse = {
    reviews: [],
    next_cursor: null,
  };
  try {
    const res = await GET<ReviewListResponse>(`/api/v1/gyms/${id}/reviews`, {
      cache: "no-store",
    });
    if (!res.ok || !res.data) return fallback;
    return res.data;
  } catch (error) {
    return null;
  }
};

const GymPage = async ({ params }: PageProps) => {
  const gym = await fetchGym(params.id);
  if (!gym) {
    notFound();
  }
  const reviews = await fetchGymReviews(params.id);
  if (!reviews) {
    console.error("No reviews found for this gym");
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
            <GymReview gym={gym} reviews={reviews || null} />
          </div>
          <GymContactSidebar gym={gym} />
        </div>
      </div>
    </div>
  );
};

export default GymPage;