"use client";
import { useState, useEffect } from "react";
import { useTranslation } from "react-i18next";
import MonthlyStrip from "./MonthlyStrip";
import { WorkoutFormDTO, GymDTO } from "@/types/workout";
import { useFormContext } from "react-hook-form";
import { useWorkoutDate } from "@/features/workout/hooks/useWorkoutDate";

const WorkoutMetadataEditor = () => {
  const { t } = useTranslation("common");
  const form = useFormContext<WorkoutFormDTO>();
  const { year, month, day, setYear, setMonth, setDay } = useWorkoutDate();

  const [isMetadataOpen, setIsMetadataOpen] = useState(false);
  const [gyms, setGyms] = useState<GymDTO[]>([]);
  const [gymInputValue, setGymInputValue] = useState("");
  const [filteredGyms, setFilteredGyms] = useState<GymDTO[]>([]);
  const [showSuggestions, setShowSuggestions] = useState(false);

  // ジムリストを取得
  useEffect(() => {
    const fetchGyms = async () => {
      // TODO: 実際のAPIエンドポイントから取得
      // const res = await fetch("/api/gyms/my");
      // if (res.ok) {
      //   const data = await res.json();
      //   setGyms(data);
      // }
      setGyms([]); // 一旦空配列
    };
    fetchGyms();
  }, []);

  // gym_name から直接ジム名を取得
  useEffect(() => {
    const gymName = form.watch("gym_name");
    if (gymName) {
      setGymInputValue(gymName);
    } else {
      setGymInputValue("");
    }
  }, [form.watch("gym_name")]);

  // 入力値に応じてフィルタリング
  useEffect(() => {
    if (gymInputValue) {
      const filtered = gyms.filter((gym) =>
        gym.name.toLowerCase().includes(gymInputValue.toLowerCase()),
      );
      setFilteredGyms(filtered);
    } else {
      setFilteredGyms([]);
    }
  }, [gymInputValue, gyms]);

  const handleDayChange = (day: number) => {
    setDay(day);
  };

  return (
    <div className="bg-white rounded-lg shadow mb-4 md:mb-6 p-3 md:p-6">
      {/* 年月選択 */}
      <div className="flex items-center justify-center gap-2 md:gap-4 mb-3 md:mb-6">
        <select
          value={year}
          onChange={(e) => {
            setYear(Number(e.target.value));
          }}
          className="w-20 md:w-28 text-center border border-gray-300 rounded-md px-2 md:px-3 py-1.5 md:py-2 text-sm md:text-base focus:outline-none focus:ring-2 focus:ring-booking-500 bg-white"
        >
          {Array.from({ length: 10 }, (_, i) => {
            const year = new Date().getFullYear() - 5 + i;
            return (
              <option key={year} value={year}>
                {year}
              </option>
            );
          })}
        </select>
        <span className="text-sm md:text-lg font-medium">
          {t("workout.metadata.year")}
        </span>

        <select
          value={month}
          onChange={(e) => {
            setMonth(Number(e.target.value));
          }}
          className="w-16 md:w-20 text-center border border-gray-300 rounded-md px-2 md:px-3 py-1.5 md:py-2 text-sm md:text-base focus:outline-none focus:ring-2 focus:ring-booking-500 bg-white"
        >
          {Array.from({ length: 12 }, (_, i) => {
            const month = i + 1;
            return (
              <option key={month} value={month}>
                {month}
              </option>
            );
          })}
        </select>
        <span className="text-sm md:text-lg font-medium">
          {t("workout.metadata.month")}
        </span>
      </div>

      {/* 日付選択 */}
      <MonthlyStrip
        year={year}
        month={month}
        selectedDay={day}
        onSelectDay={handleDayChange}
      />

      {/* モバイル: 折りたたみトグル */}
      <div className="mt-2 md:hidden">
        <button
          type="button"
          onClick={() => setIsMetadataOpen(!isMetadataOpen)}
          className="w-full flex items-center justify-between px-3 py-1.5 text-xs font-medium text-gray-700 bg-gray-50 rounded-md hover:bg-gray-100 transition-colors"
        >
          <span>{t("workout.metadata.detailsButton")}</span>
          <svg
            className={`w-4 h-4 transition-transform ${isMetadataOpen ? "rotate-180" : ""}`}
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d="M19 9l-7 7-7-7"
            />
          </svg>
        </button>
      </div>

      {/* 区切り線（デスクトップのみ） */}
      <div className="hidden md:block my-6 border-t border-gray-200"></div>

      {/* トレーニング詳細（時刻・場所・コンディション） */}
      <div
        className={`flex flex-col md:flex-row md:flex-wrap items-start justify-start gap-2 md:gap-3 text-left ${isMetadataOpen ? "mt-2" : "hidden"} md:flex`}
      >
        {/* 時間 */}
        <div className="flex items-center gap-1 w-full md:w-auto">
          <label className="text-xs md:text-sm font-medium text-gray-700 min-w-[44px] md:min-w-[60px]">
            {t("workout.metadata.timeLabel")}
          </label>
          <input
            type="time"
            className="w-24 md:w-28 px-1.5 md:px-2 py-1.5 md:py-2 text-sm border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-booking-500"
            value={form.watch("started_at") ?? ""}
            onChange={(e) =>
              form.setValue("started_at", e.target.value || null, {
                shouldDirty: true,
              })
            }
          />
          <span className="text-xs text-gray-500">〜</span>
          <input
            type="time"
            className="w-24 md:w-28 px-1.5 md:px-2 py-1.5 md:py-2 text-sm border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-booking-500"
            value={form.watch("ended_at") ?? ""}
            onChange={(e) =>
              form.setValue("ended_at", e.target.value || null, {
                shouldDirty: true,
              })
            }
          />
        </div>

        {/* 場所（ジム予測入力） */}
        <div className="flex items-center gap-2 w-full md:w-auto relative">
          <label className="text-xs md:text-sm font-medium text-gray-700 min-w-[44px] md:min-w-[60px]">
            {t("workout.metadata.locationLabel")}
          </label>
          <div className="flex-1 md:w-80 relative">
            <input
              type="text"
              value={gymInputValue}
              onChange={(e) => {
                const value = e.target.value;
                setGymInputValue(value);
                setShowSuggestions(true);
                // gym_name を form に反映（バックエンドで使用）
                form.setValue("gym_name", value.trim() || null, {
                  shouldDirty: true,
                });
                if (!value) {
                  form.setValue("gym_id", null, { shouldDirty: true });
                }
              }}
              onFocus={() => setShowSuggestions(true)}
              onBlur={() => {
                // 少し遅延させてクリックイベントを処理できるようにする
                setTimeout(() => {
                  setShowSuggestions(false);
                  // blur時にも最終的な値を反映
                  form.setValue("gym_name", gymInputValue.trim() || null, {
                    shouldDirty: true,
                  });
                }, 200);
              }}
              placeholder={t("workout.metadata.locationPlaceholder")}
              className="w-full px-2 md:px-3 py-1.5 md:py-2 text-xs md:text-sm border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-booking-500"
            />
            {showSuggestions && filteredGyms.length > 0 && (
              <ul className="absolute z-10 w-full mt-1 bg-white border border-gray-300 rounded-md shadow-lg max-h-60 overflow-auto">
                {filteredGyms.map((gym) => (
                  <li
                    key={gym.id}
                    onClick={() => {
                      setGymInputValue(gym.name);
                      form.setValue("gym_id", gym.id, { shouldDirty: true });
                      form.setValue("gym_name", gym.name.trim(), {
                        shouldDirty: true,
                      });
                      setShowSuggestions(false);
                    }}
                    className="px-3 py-2 text-xs md:text-sm hover:bg-gray-100 cursor-pointer"
                  >
                    {gym.name}
                  </li>
                ))}
              </ul>
            )}
          </div>
        </div>

        {/* コンディション（1〜5） */}
        <div className="flex items-center gap-2 w-full md:w-auto">
          <label className="text-xs md:text-sm font-medium text-gray-700 min-w-[44px] md:min-w-[60px]">
            {t("workout.metadata.conditionLabel")}
          </label>
          <div className="flex gap-1">
            {[1, 2, 3, 4, 5].map((n) => (
              <button
                key={n}
                type="button"
                onClick={() =>
                  form.setValue("condition_level", n as 1 | 2 | 3 | 4 | 5, {
                    shouldDirty: true,
                  })
                }
                className={`w-8 h-8 md:w-10 md:h-10 rounded border text-xs md:text-sm font-medium transition-colors ${form.watch("condition_level") === n ? "bg-booking-600 text-white border-booking-600" : "border-gray-300 text-gray-700 hover:bg-gray-100"}`}
              >
                {n}
              </button>
            ))}
          </div>
        </div>
      </div>
    </div>
  );
};

export default WorkoutMetadataEditor;
