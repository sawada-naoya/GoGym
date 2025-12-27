/**
 * 日付を YYYY-MM-DD 形式の文字列にフォーマット
 * タイムゾーンの影響を受けない安全な実装
 */
export function formatDate(year: number, month: number, day: number): string {
  return `${year}-${String(month).padStart(2, "0")}-${String(day).padStart(2, "0")}`;
}

/**
 * Date オブジェクトを YYYY-MM-DD 形式の文字列にフォーマット
 */
export function formatDateFromDateObject(date: Date): string {
  const year = date.getFullYear();
  const month = date.getMonth() + 1;
  const day = date.getDate();
  return formatDate(year, month, day);
}

/**
 * 指定された年月の日数を取得
 */
export function getDaysInMonth(year: number, month: number): number {
  return new Date(year, month, 0).getDate();
}

/**
 * 日付が有効な範囲内かチェックし、必要に応じて調整
 * 例: 2025年2月31日 → 2025年2月28日
 */
export function validateDay(year: number, month: number, day: number): number {
  const daysInMonth = getDaysInMonth(year, month);
  return day > daysInMonth ? daysInMonth : day;
}

/**
 * YYYY-MM-DD 形式の文字列を { year, month, day } に分解
 */
export function parseDateString(dateString: string): {
  year: number;
  month: number;
  day: number;
} {
  const [year, month, day] = dateString.split("-").map(Number);
  return { year, month, day };
}

/**
 * HH:mm 形式の時刻文字列を作成
 */
export function formatTime(hours: number, minutes: number): string {
  return `${String(hours).padStart(2, "0")}:${String(minutes).padStart(2, "0")}`;
}

/**
 * Date オブジェクトから HH:mm 形式の時刻文字列を作成
 */
export function formatTimeFromDate(date: Date): string {
  return formatTime(date.getHours(), date.getMinutes());
}
