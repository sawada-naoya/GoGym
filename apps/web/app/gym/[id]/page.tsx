import { Gym } from "@/types/gym";
import { GET } from "@/lib/api";
import Header from "@/components/Header";
import { notFound } from "next/navigation";
import GymPhotoGallery from "../components/GymPhotoGallery";
import GymBasicInfo from "../components/GymBasicInfo";
import GymAmenities from "../components/GymAmenities";
import GymAccessInfo from "../components/GymAccessInfo";
import GymReview from "../components/GymReview";
import GymContactSidebar from "../components/GymContactSidebar";

type PageProps = {
  params: {
    id: string;
  };
};

const fetchGym = async (id: string): Promise<Gym | null> => {
  try {
    const res = await GET(`/api/v1/gym/${id}`, {
      cache: "no-store",
    });
    if (!res.ok) {
      throw new Error(`HTTP ${res.status}`);
    }
    const data = await res.json();
    return data || null;
  } catch (error) {
    console.error("Failed to fetch gym:", error);
    return null;
  }
};

const GymPage = async ({ params }: PageProps) => {
  const gym = await fetchGym(params.id);
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

export default GymPage;
