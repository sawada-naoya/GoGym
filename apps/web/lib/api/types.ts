/**
 * Server Actionの標準的な戻り値型
 */
export type ActionResult<T = void> = { success: true; data?: T } | { success: false; error: string };
