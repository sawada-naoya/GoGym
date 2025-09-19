export type FieldErrors = Record<string, string>;
export type ApiError = { kind: "validation"; message: string; fields: FieldErrors } | { kind: "conflict"; message: string; code: "EMAIL_EXISTS" | "USERNAME_EXISTS" } | { kind: "auth"; message: string } | { kind: "notfound"; message: string } | { kind: "rate"; message: string } | { kind: "server"; message: string } | { kind: "network"; message: string };

const dictionary: Record<string, string> = {
  EMAIL_EXISTS: "このメールアドレスは既に登録されています",
  USERNAME_EXISTS: "このユーザー名は既に使われています",
  UNAUTHORIZED: "ログインが必要です",
  FORBIDDEN: "権限がありません",
  NOT_FOUND: "見つかりませんでした",
  RATE_LIMITED: "混み合っています。しばらくしてからお試しください。",
  DEFAULT: "登録に失敗しました。時間を置いて再度お試しください。",
};

export const normalizeError = async (res: Response): Promise<ApiError> => {
  let body: any = null;
  try {
    body = await res.json();
  } catch {}
  const e = body?.error ?? {};
  const code = e.code as string | undefined;

  if (res.status === 422 && e.fields) return { kind: "validation", message: dictionary.DEFAULT, fields: e.fields };
  if (res.status === 409) return { kind: "conflict", message: dictionary[code ?? "DEFAULT"], code: (code as any) ?? "EMAIL_EXISTS" };
  if (res.status === 401 || res.status === 403) return { kind: "auth", message: dictionary[code ?? "UNAUTHORIZED"] };
  if (res.status === 404) return { kind: "notfound", message: dictionary.NOT_FOUND };
  if (res.status === 429) return { kind: "rate", message: dictionary.RATE_LIMITED };
  return { kind: "server", message: dictionary.DEFAULT };
};

export interface ErrorDisplayProps {
  error?: ApiError | null;
  fieldName?: string;
  className?: string;
}

export const getFieldError = (error: ApiError | null | undefined, fieldName: string): string | null => {
  if (!error) return null;
  if (error.kind === "validation" && error.fields[fieldName]) {
    return error.fields[fieldName];
  }
  return null;
};

export const getFormError = (error: ApiError | null | undefined): string | null => {
  if (!error) return null;
  if (error.kind === "validation") return null; // フィールドエラーは個別に表示
  return error.message;
};

export const getInputClassName = (error: ApiError | null | undefined, fieldName: string, baseClassName: string = "form-input"): string => {
  const hasError = getFieldError(error, fieldName) !== null;
  if (hasError) {
    return `${baseClassName} border-red-500 focus:border-red-500 focus:ring-red-500`;
  }
  return baseClassName;
};