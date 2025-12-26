"use client";

import { useTranslation } from "react-i18next";

const LandingHero = ({ children }: { children: React.ReactNode }) => {
  const { t } = useTranslation("common");

  return (
    <div className="grid md:grid-cols-2 gap-12 items-center">
      {/* 左側：説明 */}
      <div className="space-y-6">
        <h1 className="text-5xl font-bold text-gray-900">
          {t("landing.hero.title")}
        </h1>
        <p className="text-2xl text-gray-700">
          {t("landing.hero.description")}
        </p>
        <div className="space-y-4">
          <div className="flex items-start gap-3">
            <div className="flex-shrink-0 w-8 h-8 bg-booking-600 rounded-full flex items-center justify-center">
              <svg
                className="w-5 h-5 text-white"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M5 13l4 4L19 7"
                />
              </svg>
            </div>
            <div>
              <h3 className="font-semibold text-gray-900">
                {t("landing.features.trainingRecord.title")}
              </h3>
              <p className="text-gray-600">
                {t("landing.features.trainingRecord.description")}
              </p>
            </div>
          </div>
          <div className="flex items-start gap-3">
            <div className="flex-shrink-0 w-8 h-8 bg-booking-600 rounded-full flex items-center justify-center">
              <svg
                className="w-5 h-5 text-white"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M5 13l4 4L19 7"
                />
              </svg>
            </div>
            <div>
              <h3 className="font-semibold text-gray-900">
                {t("landing.features.progressTracking.title")}
              </h3>
              <p className="text-gray-600">
                {t("landing.features.progressTracking.description")}
              </p>
            </div>
          </div>
          <div className="flex items-start gap-3">
            <div className="flex-shrink-0 w-8 h-8 bg-booking-600 rounded-full flex items-center justify-center">
              <svg
                className="w-5 h-5 text-white"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M5 13l4 4L19 7"
                />
              </svg>
            </div>
            <div>
              <h3 className="font-semibold text-gray-900">
                {t("landing.features.simpleUsage.title")}
              </h3>
              <p className="text-gray-600">
                {t("landing.features.simpleUsage.description")}
              </p>
            </div>
          </div>
        </div>
      </div>

      {/* 右側：子コンポーネント */}
      {children}
    </div>
  );
};

export default LandingHero;
