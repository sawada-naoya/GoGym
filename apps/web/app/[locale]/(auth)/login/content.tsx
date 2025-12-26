"use client";

import { useState } from "react";
import Link from "next/link";
import { useRouter, useSearchParams, useParams } from "next/navigation";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";

import { signIn } from "next-auth/react";
import { useBanner } from "@/components/Banner";
import { loginSchema, type LoginForm } from "@/features/auth/schemas/login";

type LoginFormContentProps = {
  showHeader?: boolean;
  showSignupLink?: boolean;
};

const normalizeCallbackUrl = (candidate: string, locale: string) => {
  // open redirect対策：外部URLは禁止。相対パスのみ許可
  if (!candidate.startsWith("/")) return `/${locale}/workout`;

  // locale未付与なら付ける（/workout -> /ja/workout）
  const hasLocale = /^\/[a-z]{2}(\/|$)/.test(candidate);
  if (hasLocale) return candidate;

  return `/${locale}${candidate}`;
};

const LoginFormContent = ({ showHeader = true, showSignupLink = true }: LoginFormContentProps = {}) => {
  const router = useRouter();
  const searchParams = useSearchParams();
  const params = useParams<{ locale: string }>();
  const locale = params?.locale ?? "ja";

  const [loading, setLoading] = useState(false);
  const { error } = useBanner();

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<LoginForm>({
    resolver: zodResolver(loginSchema),
    mode: "onSubmit",
    reValidateMode: "onBlur",
  });

  const onSubmit = async (data: LoginForm) => {
    setLoading(true);
    try {
      const rawCallbackUrl = searchParams.get("callbackUrl") ?? "/workout";
      const callbackUrl = normalizeCallbackUrl(rawCallbackUrl, locale);

      const result = await signIn("credentials", {
        email: data.email,
        password: data.password,
        redirect: false,
        callbackUrl,
      });

      if (result?.error) {
        if (result.error === "CredentialsSignin") {
          error("メールアドレスまたはパスワードが正しくありません");
        } else {
          error("ログインに失敗しました");
        }
        return;
      }

      sessionStorage.setItem("flash", JSON.stringify({ variant: "success", message: "ログインに成功しました" }));

      router.replace(result?.url ?? callbackUrl);
    } finally {
      setLoading(false);
    }
  };

  const formElement = (
    <form className={showHeader ? "mt-8 space-y-6" : "space-y-4"} onSubmit={handleSubmit(onSubmit)} noValidate>
      <div className="space-y-4">
        <div>
          <label htmlFor="email" className="form-label">
            メールアドレス
          </label>
          <input {...register("email")} id="email" type="email" autoComplete="email" className={`form-input ${errors.email ? "border-red-500 focus:border-red-500 focus:ring-red-500" : ""}`} placeholder="example@example.com" />
          {errors.email && <p className="mt-1 text-sm text-red-600">{errors.email.message}</p>}
        </div>

        <div>
          <label htmlFor="password" className="form-label">
            パスワード
          </label>
          <input {...register("password")} id="password" type="password" autoComplete="current-password" className={`form-input ${errors.password ? "border-red-500 focus:border-red-500 focus:ring-red-500" : ""}`} placeholder="パスワードを入力" />
          {errors.password && <p className="mt-1 text-sm text-red-600">{errors.password.message}</p>}
        </div>
      </div>

      <div>
        <button
          type="submit"
          disabled={loading}
          className="group relative w-full flex justify-center py-3 px-4 border border-transparent text-sm font-medium rounded-md text-white bg-booking-600 hover:bg-booking-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-booking-500 disabled:opacity-50 disabled:cursor-not-allowed transition-colors duration-200"
        >
          {loading ? "ログイン中..." : "ログイン"}
        </button>
      </div>

      {/* Googleログイン（後でちゃんと実装するまで消すのが正解。今はUIノイズ） */}
      <div className="mt-6">
        <div className="relative">
          <div className="absolute inset-0 flex items-center">
            <div className="w-full border-t border-gray-300" />
          </div>
          <div className="relative flex justify-center text-sm">
            <span className="px-2 bg-gray-50 text-gray-500">または</span>
          </div>
        </div>

        <div className="mt-6">
          <button type="button" className="w-full inline-flex justify-center py-3 px-4 border border-gray-300 rounded-md shadow-sm bg-white text-sm font-medium text-gray-500 hover:bg-gray-50 transition-colors duration-200" onClick={() => alert("Googleログイン機能は開発中です")}>
            Googleでログイン
          </button>
        </div>
      </div>

      <div className="text-center text-sm">
        <Link href={`/${locale}/forgot-password`} className="font-medium text-booking-600 hover:text-booking-500">
          パスワードを忘れた方はこちら
        </Link>
      </div>

      {showSignupLink && (
        <div className="text-center text-sm">
          <span className="text-gray-600">アカウントをお持ちでない方は </span>
          <Link href={`/${locale}/signup`} className="font-medium text-booking-600 hover:text-booking-500">
            新規登録
          </Link>
        </div>
      )}
    </form>
  );

  if (!showHeader) return formElement;

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
      <div className="max-w-md w-full space-y-8">
        <div>
          <h2 className="mt-6 text-center text-3xl font-extrabold text-gray-900">アカウントにログイン</h2>
          <p className="mt-2 text-center text-sm text-gray-600">トレーニング記録を開始しましょう</p>
        </div>
        {formElement}
      </div>
    </div>
  );
};

export default LoginFormContent;
