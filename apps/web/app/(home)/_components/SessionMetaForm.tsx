"use client";
import React from "react";
import type { SessionMeta } from "./types";

const SessionMetaForm = ({ value, onChange }: { value: SessionMeta; onChange: (v: SessionMeta) => void }) => {
  return (
    <div className="bg-white rounded-lg shadow mb-6 p-6">
      <div className="flex flex-wrap items-start justify-start gap-3 text-left">
        {/* 時間 */}
        <div className="flex items-center gap-1">
          <label className="mr-2 text-sm font-medium text-gray-700">時間</label>
          <input type="time" className="w-24 px-2 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-booking-500" value={value.startedAt?.slice(11, 16) || ""} onChange={(e) => onChange({ ...value, startedAt: e.target.value ? e.target.value : null })} />
          <span className="text-gray-500">〜</span>
          <input type="time" className="w-24 px-2 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-booking-500" value={value.endedAt?.slice(11, 16) || ""} onChange={(e) => onChange({ ...value, endedAt: e.target.value ? e.target.value : null })} />
        </div>

        {/* 場所 */}
        <div className="flex items-center gap-2">
          <label className="text-sm font-medium text-gray-700">場所</label>
          <input type="text" value={value.place ?? ""} onChange={(e) => onChange({ ...value, place: e.target.value })} placeholder="自宅 / ジム名など" className="w-48 px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-booking-500" />
        </div>

        {/* 部位 */}
        <div className="flex items-center gap-2">
          <label className="text-sm font-medium text-gray-700">部位</label>
          <select value={value.muscle ?? ""} onChange={(e) => onChange({ ...value, muscle: e.target.value as any })} className="w-36 px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-booking-500">
            <option value="">選択</option>
            <option value="腕">腕</option>
            <option value="胸">胸</option>
            <option value="脚">脚</option>
            <option value="肩">肩</option>
            <option value="背中">背中</option>
            <option value="体幹">体幹</option>
          </select>
        </div>

        {/* コンディション */}
        <div className="flex items-center gap-2">
          <label className="text-sm font-medium text-gray-700">コンディション</label>
          <div className="flex gap-1">
            {[1, 2, 3, 4, 5].map((n) => (
              <button key={n} type="button" onClick={() => onChange({ ...value, condition: n as 1 | 2 | 3 | 4 | 5 })} className={`w-8 h-8 rounded border text-sm font-medium transition-colors ${value.condition === n ? "bg-booking-600 text-white border-booking-600" : "border-gray-300 text-gray-700 hover:bg-gray-100"}`}>
                {n}
              </button>
            ))}
          </div>
        </div>
      </div>
    </div>
  );
};

export default SessionMetaForm;
