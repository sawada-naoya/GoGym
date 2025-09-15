'use client'

import { useState } from 'react';
import { Gym } from '@/types/gym';
import { ReviewListResponse } from '@/types/review';
import GymReviewModal from './GymReviewModal';

type GymReviewProps = {
  gym: Gym;
  reviews: ReviewListResponse | null;
};

const GymReview = ({ gym, reviews }: GymReviewProps) => {
  const [isModalOpen, setIsModalOpen] = useState(false);

  const handleOpenModal = () => {
    setIsModalOpen(true);
  };

  const handleCloseModal = () => {
    setIsModalOpen(false);
  };

  return (
    <>
      <div className="bg-white rounded-lg shadow p-6">
        <h2 className="text-xl font-bold text-gray-900 mb-4">レビュー</h2>
        <div className="flex items-center justify-between mb-4">
          <div className="flex items-center">
            <div className="bg-booking-600 text-white px-3 py-2 rounded font-bold text-lg">
              {gym.average_rating?.toFixed(1) || "0.0"}
            </div>
            <div className="ml-3">
              <p className="font-semibold text-gray-900">素晴らしい</p>
              <p className="text-sm text-gray-600">{gym.review_count}件のレビュー</p>
            </div>
          </div>
        </div>

        <button
          onClick={handleOpenModal}
          className="w-full text-booking-600 hover:text-booking-700 font-semibold py-2 border-t border-gray-200 transition-colors"
        >
          すべてのレビューを見る
        </button>
      </div>

      <GymReviewModal
        gym={gym}
        reviews={reviews}
        isOpen={isModalOpen}
        onClose={handleCloseModal}
      />
    </>
  );
};

export default GymReview;