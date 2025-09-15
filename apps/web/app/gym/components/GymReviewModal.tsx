"use client";

import { useState, useEffect } from "react";
import { Gym } from "@/types/gym";
import { ReviewListResponse } from "@/types/review";

type GymReviewModalProps = {
  gym: Gym;
  reviews: ReviewListResponse | null;
  isOpen: boolean;
  onClose: () => void;
};

const GymReviewModal = ({ gym, reviews, isOpen, onClose }: GymReviewModalProps) => {
  // ESCキーでモーダルを閉じる
  useEffect(() => {
    const handleEscape = (event: KeyboardEvent) => {
      if (event.key === "Escape") {
        onClose();
      }
    };

    if (isOpen) {
      document.addEventListener("keydown", handleEscape);
      document.body.style.overflow = "hidden";
    }

    return () => {
      document.removeEventListener("keydown", handleEscape);
      document.body.style.overflow = "unset";
    };
  }, [isOpen, onClose]);

  if (!isOpen) return null;

  const renderStars = (rating: number) => {
    return Array.from({ length: 5 }, (_, i) => (
      <span key={i} className={i < rating ? "text-yellow-400" : "text-gray-300"}>
        ★
      </span>
    ));
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center">
      {/* オーバーレイ */}
      <div className="absolute inset-0 bg-black bg-opacity-50" onClick={onClose} />

      {/* モーダル本体 */}
      <div className="relative bg-white rounded-lg shadow-xl w-full max-w-4xl max-h-[90vh] m-4 flex flex-col">
        {/* ヘッダー */}
        <div className="flex items-center justify-between p-6 border-b border-gray-200">
          <div>
            <h2 className="text-2xl font-bold text-gray-900">{gym.name}のレビュー</h2>
            <div className="flex items-center mt-2">
              <div className="bg-booking-600 text-white px-3 py-1 rounded font-bold">{gym.average_rating?.toFixed(1) || "0.0"}</div>
              <span className="ml-3 text-gray-600">{gym.review_count}件のレビュー</span>
            </div>
          </div>
          <button onClick={onClose} className="text-gray-400 hover:text-gray-600 text-3xl font-light">
            ×
          </button>
        </div>

        {/* レビュー一覧 */}
        <div className="flex-1 overflow-y-auto p-6">
          <div className="space-y-6">
            {reviews?.reviews.map((review) => (
              <div key={review.id} className="border-b border-gray-200 last:border-b-0 pb-6 last:pb-0">
                <div className="flex items-start space-x-4">
                  {/* アバター */}
                  <div className="w-10 h-10 bg-booking-600 rounded-full flex items-center justify-center text-white font-semibold">{review.user?.name?.charAt(0) || "?"}</div>

                  {/* レビュー内容 */}
                  <div className="flex-1">
                    <div className="flex items-center justify-between mb-2">
                      <div>
                        <h4 className="font-semibold text-gray-900">{review.user?.name || "匿名ユーザー"}</h4>
                        <span className="text-sm text-green-600 flex items-center mt-1">✓ 認証済み</span>
                      </div>
                      <div className="text-right">
                        <div className="flex items-center text-sm">{renderStars(review.rating)}</div>
                        <p className="text-xs text-gray-500 mt-1">{new Date(review.created_at).toLocaleDateString("ja-JP")}</p>
                      </div>
                    </div>

                    <p className="text-gray-700 leading-relaxed">{review.content}</p>
                  </div>
                </div>
              </div>
            ))}
          </div>
        </div>

        {/* フッター */}
        <div className="p-6 border-t border-gray-200">
          <div className="flex justify-between items-center">
            <p className="text-sm text-gray-600">
              {reviews?.reviews.length}件中{reviews?.reviews.length}件を表示
            </p>
            <a href={`/gym/${gym.id}/reviews`} className="bg-booking-600 hover:bg-booking-700 text-white font-semibold py-2 px-4 rounded-lg transition-colors" onClick={onClose}>
              すべてのレビューページへ
            </a>
          </div>
        </div>
      </div>
    </div>
  );
};

export default GymReviewModal;
