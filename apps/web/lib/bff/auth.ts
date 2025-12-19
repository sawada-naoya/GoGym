/**
 * ユーザー登録（サインアップ）
 */
export const signup = async (data: {
  name: string;
  email: string;
  password: string;
}): Promise<{ ok: boolean; error: string | null }> => {
  const res = await fetch("/api/auth/signup", {
    method: "POST",
    headers: { "content-type": "application/json" },
    body: JSON.stringify(data),
    cache: "no-store",
  });

  if (!res.ok) {
    const errorData = await res.json().catch(() => ({}));
    return {
      ok: false,
      error: errorData.message || "このメールアドレスは既に使用されています",
    };
  }

  return { ok: true, error: null };
};
