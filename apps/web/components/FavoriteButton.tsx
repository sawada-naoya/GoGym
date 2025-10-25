"use client";

type FavoriteButtonProps = {
  gymId: number;
};

const FavoriteButton = ({ gymId }: FavoriteButtonProps) => {
  const handleToggleFavorite = (e: React.MouseEvent) => {
    e.preventDefault();
    console.log("Toggle favorite for gym:", gymId);
    // TODO: お気に入り機能の実装
  };

  return (
    <button className="absolute top-3 right-3 px-3 py-1 bg-white/80 hover:bg-white rounded text-sm font-medium text-gray-800 shadow-sm transition-colors" onClick={handleToggleFavorite}>
      ♡
    </button>
  );
};

export default FavoriteButton;
