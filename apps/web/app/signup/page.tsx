"use client";

import { useState } from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { z } from "zod";
import { POST } from "@/lib/api";
import { ErrorBanner, SuccessBanner } from "../../components/ui/Banner";

const SignUpSchema = z
  .object({
    name: z.string().min(1, "名前は必須です"),
    email: z.email("有効なメールアドレスを入力してください").min(1, "メールアドレスは必須です"),
    password: z.string().min(6, "パスワードは6文字以上で入力してください"),
    confirmPassword: z.string().min(1, "パスワード（確認）は必須です"),
  })
  .refine((data) => data.password === data.confirmPassword, {
    message: "パスワードが一致しません",
    path: ["confirmPassword"],
  });

type FieldErrors = Record<string, string>;

const SignUpPage = () => {
  const router = useRouter();
  const [form, setForm] = useState({
    name: "",
    email: "",
    password: "",
    confirmPassword: "",
  });
  const [loading, setLoading] = useState(false);
  const [errors, setErrors] = useState<FieldErrors>({});
  const [apiError, setApiError] = useState<string | null>(null);

  const onChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setForm({ ...form, [e.target.name]: e.target.value });
    if (errors[e.target.name]) {
      setErrors({ ...errors, [e.target.name]: "" });
    }
    // 入力開始時にエラーメッセージをクリア
    if (apiError) setApiError(null);
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setErrors({});
    setApiError(null);

    const parsed = SignUpSchema.safeParse(form);
    if (!parsed.success) {
      const fe: FieldErrors = {};
      for (const issue of parsed.error.issues) {
        const k = issue.path[0] as string;
        if (!fe[k]) fe[k] = issue.message;
      }
      setErrors(fe);
      return;
    }

    setLoading(true);
    try {
      const res = await POST("/api/v1/users", {
        body: {
          name: form.name,
          email: form.email,
          password: form.password,
        },
      });

      if (!res.ok) {
        setApiError("このメールアドレスは既に使用されています");
        return;
      }

      router.push("/login?success=signup");
    } catch (error) {
      setApiError("ネットワークエラーが発生しました。時間を置いて再度お試しください。");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center py-12 px-4 sm:px-6 lg:px-8">
      <div className="max-w-md w-full space-y-8">
        <div>
          <div className="flex justify-center">
            <h1 className="text-4xl font-bold text-gray-900">GoGym</h1>
          </div>
          <h2 className="mt-6 text-center text-3xl font-bold tracking-tight text-gray-900">新規登録</h2>
          <p className="mt-2 text-center text-sm text-gray-600">
            既にアカウントをお持ちの方は{" "}
            <Link href="/login" className="font-medium text-booking-600 hover:text-booking-500">
              ログイン
            </Link>
          </p>
        </div>

        <form className="mt-8 space-y-6" onSubmit={handleSubmit} noValidate>
          {apiError && <ErrorBanner message={apiError} />}

          <div className="space-y-4">
            <div>
              <label htmlFor="name" className="form-label">
                名前
              </label>
              <input id="name" name="name" type="text" autoComplete="name" className={`form-input ${errors.name ? "border-red-500 focus:border-red-500 focus:ring-red-500" : ""}`} placeholder="山田太郎" value={form.name} onChange={onChange} />
              {errors.name && <p className="mt-1 text-sm text-red-600">{errors.name}</p>}
            </div>

            <div>
              <label htmlFor="email" className="form-label">
                メールアドレス
              </label>
              <input id="email" name="email" type="email" autoComplete="email" className={`form-input ${errors.email ? "border-red-500 focus:border-red-500 focus:ring-red-500" : ""}`} placeholder="example@example.com" value={form.email} onChange={onChange} />
              {errors.email && <p className="mt-1 text-sm text-red-600">{errors.email}</p>}
            </div>

            <div>
              <label htmlFor="password" className="form-label">
                パスワード
              </label>
              <input id="password" name="password" type="password" autoComplete="new-password" className={`form-input ${errors.password ? "border-red-500 focus:border-red-500 focus:ring-red-500" : ""}`} placeholder="パスワードを入力" value={form.password} onChange={onChange} />
              {errors.password && <p className="mt-1 text-sm text-red-600">{errors.password}</p>}
            </div>

            <div>
              <label htmlFor="confirmPassword" className="form-label">
                パスワード（確認）
              </label>
              <input id="confirmPassword" name="confirmPassword" type="password" autoComplete="new-password" className={`form-input ${errors.confirmPassword ? "border-red-500 focus:border-red-500 focus:ring-red-500" : ""}`} placeholder="パスワードを再度入力" value={form.confirmPassword} onChange={onChange} />
              {errors.confirmPassword && <p className="mt-1 text-sm text-red-600">{errors.confirmPassword}</p>}
            </div>
          </div>

          <div>
            <button
              type="submit"
              disabled={loading}
              className="group relative w-full flex justify-center py-3 px-4 border border-transparent text-sm font-medium rounded-md text-white bg-booking-600 hover:bg-booking-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-booking-500 disabled:opacity-50 disabled:cursor-not-allowed transition-colors duration-200"
            >
              {loading ? (
                <div className="flex items-center">
                  <div className="animate-spin -ml-1 mr-3 h-5 w-5 text-white">
                    <svg className="h-5 w-5" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                      <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                      <path className="opacity-75" fill="currentColor" d="m4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                    </svg>
                  </div>
                  登録中...
                </div>
              ) : (
                "新規登録"
              )}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};

export default SignUpPage;
