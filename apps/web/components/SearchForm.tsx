"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";

const SearchForm = () => {
  const [searchQuery, setSearchQuery] = useState("");
  const [selectedArea, setSelectedArea] = useState("");
  const [loading, setLoading] = useState(false);
  const router = useRouter();

  const handleSearch = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);

    // URLパラメータを構築
    const params = new URLSearchParams();

    if (searchQuery.trim()) {
      params.append("q", searchQuery.trim());
    }

    // エリア選択に基づいて座標を設定（ダミー座標）
    if (selectedArea) {
      const areaCoordinates: Record<string, { lat: number; lon: number }> = {
        shibuya: { lat: 35.6598, lon: 139.7036 },
        shinjuku: { lat: 35.6938, lon: 139.7034 },
        omotesando: { lat: 35.6657, lon: 139.7116 },
        ikebukuro: { lat: 35.7295, lon: 139.7109 },
        ebisu: { lat: 35.6468, lon: 139.7102 },
        roppongi: { lat: 35.6627, lon: 139.7371 },
      };

      const coords = areaCoordinates[selectedArea];
      if (coords) {
        params.append("lat", coords.lat.toString());
        params.append("lon", coords.lon.toString());
        params.append("radius_m", "5000"); // 5km圏内
      }
    }

    // 検索ページへ遷移
    const searchUrl = `/search?${params.toString()}`;
    router.push(searchUrl);

    setLoading(false);
  };

  return (
    <form onSubmit={handleSearch} className="bg-white rounded-xl shadow-lg p-6">
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        {/* ジム名・キーワード検索 */}
        <div className="space-y-2">
          <label className="block text-sm font-medium text-gray-700">
            ジム名・キーワード
          </label>
          <input
            type="text"
            placeholder="エニタイム..."
            className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-booking-500 focus:border-transparent"
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
          />
        </div>

        {/* エリア選択 */}
        <div className="space-y-2">
          <label className="block text-sm font-medium text-gray-700">
            エリア
          </label>
          <select
            className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-booking-500 focus:border-transparent appearance-none bg-white"
            value={selectedArea}
            onChange={(e) => setSelectedArea(e.target.value)}
          >
            <option value="">エリアを選択</option>
            <option value="shibuya">渋谷・原宿</option>
            <option value="shinjuku">新宿・代々木</option>
            <option value="omotesando">表参道・青山</option>
            <option value="ikebukuro">池袋・巣鴨</option>
            <option value="ebisu">恵比寿・中目黒</option>
            <option value="roppongi">六本木・赤坂</option>
          </select>
        </div>

        {/* 検索ボタン */}
        <div className="space-y-2">
          <label className="block text-sm font-medium text-transparent">
            検索
          </label>
          <button
            type="submit"
            disabled={loading}
            className="w-full bg-booking-600 hover:bg-booking-700 disabled:opacity-50 text-white font-semibold py-3 px-6 rounded-lg transition-colors duration-200"
          >
            {loading ? "検索中..." : "検索する"}
          </button>
        </div>
      </div>

      {/* フィルターボタン */}
      <div className="flex items-center justify-center mt-4 pt-4 border-t border-gray-200">
        <button
          type="button"
          className="text-booking-600 hover:text-booking-700 font-medium"
        >
          詳細条件で絞り込み
        </button>
      </div>
    </form>
  );
};

export default SearchForm;
