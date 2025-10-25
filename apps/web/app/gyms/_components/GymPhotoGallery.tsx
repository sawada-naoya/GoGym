import { Gym } from "@/types/gym";

type GymPhotoGalleryProps = {
  gym: Gym;
};

const GymPhotoGallery = ({ gym }: GymPhotoGalleryProps) => {
  const placeholderPhoto = "/placeholder-gym.jpg";
  const photos = [placeholderPhoto, placeholderPhoto, placeholderPhoto, placeholderPhoto];
  
  return (
    <div className="bg-white">
      <div className="container mx-auto px-4 py-4">
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-2 h-96">
          {/* メイン写真 */}
          <div className="md:col-span-2 lg:col-span-2 md:row-span-2">
            <img
              src={placeholderPhoto}
              alt={gym.name}
              className="w-full h-full object-cover rounded-lg hover:opacity-90 transition-opacity cursor-pointer"
            />
          </div>
          
          {/* サブ写真 */}
          {photos.slice(0, 4).map((photo, index) => (
            <div key={index} className="hidden lg:block">
              <img
                src={photo}
                alt={`${gym.name} - 写真${index + 2}`}
                className="w-full h-full object-cover rounded-lg hover:opacity-90 transition-opacity cursor-pointer"
              />
            </div>
          ))}
          
          {photos.length > 4 && (
            <div className="hidden lg:block relative">
              <img
                src={photos[4]}
                alt={`${gym.name} - その他の写真`}
                className="w-full h-full object-cover rounded-lg"
              />
              <div className="absolute inset-0 bg-black bg-opacity-50 flex items-center justify-center rounded-lg cursor-pointer hover:bg-opacity-40 transition-all">
                <span className="text-white font-semibold">
                  +{photos.length - 4}枚の写真を見る
                </span>
              </div>
            </div>
          )}
        </div>
        
        {/* モバイル用：写真を全て見るボタン */}
        <div className="lg:hidden mt-4">
          <button className="w-full py-2 px-4 border border-booking-600 text-booking-600 rounded-lg hover:bg-booking-50 transition-colors">
            すべての写真を見る ({photos.length}枚)
          </button>
        </div>
      </div>
    </div>
  );
};

export default GymPhotoGallery;