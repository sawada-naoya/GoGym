import { useEffect, useRef } from "react";
import type { ExerciseRow } from "../_lib/utils";

/**
 * モバイル時に新規データのセット数を1に調整するフック
 * 既存レコード（実データあり）の場合は調整しない
 *
 * @param isMobile モバイルフラグ
 * @param exercises 種目データ配列
 * @param onChangeExercises データ変更コールバック
 * @param dataKey データを識別するキー（日付など、データが変わる度に変わる値）
 */
export const useMobileExerciseAdjustment = (
  isMobile: boolean,
  exercises: ExerciseRow[] | undefined,
  onChangeExercises: (exercises: ExerciseRow[]) => void,
  dataKey?: string,
) => {
  const initializedForRef = useRef<string>("");

  useEffect(() => {
    if (!isMobile || !exercises) return;

    // dataKeyが必須（提供されていない場合は動作しない）
    if (!dataKey) return;

    // 既にこのデータセットで初期化済みなら何もしない
    if (initializedForRef.current === dataKey) return;

    // 空データの場合はキーだけ記録して終了
    if (exercises.length === 0) {
      initializedForRef.current = dataKey;
      return;
    }

    // 既存レコード判定：セットに実データ（weight_kg or reps）が入っているか
    const hasExistingData = exercises.some((ex) =>
      ex.sets.some(
        (set) =>
          (set.weight_kg !== null &&
            set.weight_kg !== "" &&
            set.weight_kg !== 0) ||
          (set.reps !== null && set.reps !== "" && set.reps !== 0),
      ),
    );

    // 既存レコードがある場合は調整せずに全セット表示
    if (hasExistingData) {
      initializedForRef.current = dataKey;
      return;
    }

    // 新規データで複数セットある場合のみ1セットに調整
    const needsAdjustment = exercises.some((ex) => ex.sets.length > 1);
    if (needsAdjustment) {
      const adjusted = exercises.map((ex) => ({
        ...ex,
        sets: ex.sets.slice(0, 1),
      }));
      onChangeExercises(adjusted);
    }

    // 調整実行後にキーを記録
    initializedForRef.current = dataKey;
  }, [isMobile, exercises, onChangeExercises, dataKey]);
};
