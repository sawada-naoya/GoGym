"use client";
import React from "react";

// ---- 外から渡すマスタ型 ----
export type PartSource = "preset" | "custom" | null;

export type PartItem = {
  id: number;
  name: string;
  source: PartSource;
};

export type ExerciseCatalogItem = {
  id: number;
  name: string;
  trainingPartId: number | null;
  isDefault: 0 | 1;
};

// ---- Props ----
export type ExerciseComboboxProps = {
  all: ExerciseCatalogItem[];
  parts: PartItem[];
  selectedPartId: number | null; // 部位フィルタ
  valueId: number | null; // 選択中 id（既存）; 新規は null
  valueName: string; // 選択中 name（新規文字列も入る）
  onPick: (picked: { id: number; name: string }) => void; // 既存を選択
  onCreateNew: (name: string) => void; // 新規作成
};

const norm = (s: string) => s.normalize("NFKC").toLowerCase();

export const ExerciseCombobox: React.FC<ExerciseComboboxProps> = ({ all, parts, selectedPartId, valueId, valueName, onPick, onCreateNew }) => {
  const [open, setOpen] = React.useState(false);
  const [q, setQ] = React.useState("");

  const filtered = React.useMemo(() => {
    let base = all;
    if (selectedPartId != null) base = base.filter((e) => e.trainingPartId === selectedPartId);
    if (q) base = base.filter((e) => norm(e.name).includes(norm(q)));
    return [...base].sort((a, b) => a.name.localeCompare(b.name, "ja"));
  }, [all, selectedPartId, q]);

  const selected = valueId != null ? all.find((e) => e.id === valueId) : null;
  const selectedPartName = selected?.trainingPartId != null ? parts.find((p) => p.id === selected.trainingPartId)?.name ?? "" : "";

  return (
    <div className="relative">
      <input
        value={q}
        onChange={(e) => {
          setQ(e.target.value);
          setOpen(true);
        }}
        onFocus={() => setOpen(true)}
        placeholder={selected?.name || valueName || "種目を検索 / 追加"}
        className="w-full px-2 py-1 border border-gray-300 rounded focus:outline-none focus:ring-1 focus:ring-booking-500"
      />

      {open && (
        <div className="absolute z-20 mt-1 w-80 max-w-[90vw] bg-white border border-gray-200 rounded shadow">
          {q && filtered.length === 0 && (
            <button
              className="w-full text-left px-3 py-2 hover:bg-gray-50"
              onMouseDown={(e) => e.preventDefault()}
              onClick={() => {
                onCreateNew(q);
                setQ("");
                setOpen(false);
              }}
            >
              「{q}」を新規追加
            </button>
          )}

          {!q && <div className="px-3 pt-2 pb-1 text-[11px] text-gray-500">{selectedPartId == null ? "すべての部位を表示中" : `部位: ${parts.find((p) => p.id === selectedPartId)?.name ?? ""}`}</div>}

          <ul className="max-h-64 overflow-auto">
            {filtered.map((e) => (
              <li key={e.id}>
                <button
                  className="w-full flex items-center justify-between px-3 py-2 hover:bg-gray-50"
                  onMouseDown={(ev) => ev.preventDefault()}
                  onClick={() => {
                    onPick({ id: e.id, name: e.name });
                    setQ("");
                    setOpen(false);
                  }}
                >
                  <span>{e.name}</span>
                  <span className="text-[10px] text-gray-400">{e.trainingPartId != null ? parts.find((p) => p.id === e.trainingPartId)?.name ?? "" : ""}</span>
                </button>
              </li>
            ))}
          </ul>
        </div>
      )}

      {selected && <div className="mt-1 text-[10px] text-gray-400">部位: {selectedPartName || "—"}</div>}
    </div>
  );
};

// 再レンダー最適化（親からの同一propsで再描画しない）
export default React.memo(ExerciseCombobox);
