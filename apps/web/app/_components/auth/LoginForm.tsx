"use client";

import { useState, Suspense } from "react";
import Link from "next/link";
import { useRouter, useSearchParams } from "next/navigation";
import { signIn } from "next-auth/react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { useBanner } from "@/components/Banner";

const LoginSchema = z.object({
  email: z.string().min(1, "メールアドレスは必須です").email("有効なメールアドレスを入力してください"),
  password: z.string().min(1, "パスワードは必須です"),
});

type LoginForm = z.infer<typeof LoginSchema>;

type LoginFormContentProps = {
  showHeader?: boolean;
  showSignupLink?: boolean;
};

const LoginFormContent = ({ showHeader = true, showSignupLink = true }: LoginFormContentProps = {}) => {
  const router = useRouter();
  const searchParams = useSearchParams();
  const [loading, setLoading] = useState(false);
  const { error } = useBanner();

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<LoginForm>({
    resolver: zodResolver(LoginSchema),
    mode: "onSubmit",
    reValidateMode: "onBlur",
  });

  const onSubmit = async (data: LoginForm) => {
    setLoading(true);

    const callbackUrl = searchParams.get("callbackUrl") ?? "/workout";
    const res = await signIn("credentials", {
      email: data.email,
      password: data.password,
      redirect: false,
      callbackUrl,
    });

    if (res?.error) {
      setLoading(false);
      if (res.error === "CredentialsSignin") {
        error("メールアドレスまたはパスワードが正しくありません");
      } else {
        error("ログインに失敗しました");
      }
      return;
    }

    sessionStorage.setItem("flash", JSON.stringify({ variant: "success", message: "ログインに成功しました" }));

    router.replace(res?.url ?? callbackUrl);
  };

  const formElement = (
    <form className={showHeader ? "mt-8 space-y-6" : "space-y-4"} onSubmit={handleSubmit(onSubmit)} noValidate>
      <div className="space-y-4">
        {/* メールアドレス */}
        <div>
          <label htmlFor="email" className="form-label">
            メールアドレス
          </label>
          <input {...register("email")} id="email" type="email" autoComplete="email" className={`form-input ${errors.email ? "border-red-500 focus:border-red-500 focus:ring-red-500" : ""}`} placeholder="example@example.com" />
          {errors.email && <p className="mt-1 text-sm text-red-600">{errors.email.message}</p>}
        </div>

        {/* パスワード */}
        <div>
          <label htmlFor="password" className="form-label">
            パスワード
          </label>
          <input {...register("password")} id="password" type="password" autoComplete="current-password" className={`form-input ${errors.password ? "border-red-500 focus:border-red-500 focus:ring-red-500" : ""}`} placeholder="パスワードを入力" />
          {errors.password && <p className="mt-1 text-sm text-red-600">{errors.password.message}</p>}
        </div>
      </div>

      {/* ログインボタン */}
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
              ログイン中...
            </div>
          ) : (
            "ログイン"
          )}
        </button>
      </div>

      {/* Googleログイン */}
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
            <div className="flex items-center">
              {/* Google SVG Icon */}
              <svg className="h-5 w-5 mr-3" viewBox="0 0 48 48">
                <path fill="#EA4335" d="M24 9.5c3.54 0 6.71 1.22 9.21 3.6l6.85-6.85C35.9 2.38 30.47 0 24 0 14.62 0 6.51 5.38 2.56 13.22l7.98 6.19C12.43 13.72 17.74 9.5 24 9.5z" />
                <path fill="#4285F4" d="M46.98 24.55c0-1.57-.15-3.09-.38-4.55H24v9.02h12.94c-.58 2.96-2.26 5.48-4.78 7.18l7.73 6c4.51-4.18 7.09-10.36 7.09-17.65z" />
                <path fill="#FBBC05" d="M10.53 28.59c-.48-1.45-.76-2.99-.76-4.59s.27-3.14.76-4.59l-7.98-6.19C.92 16.46 0 20.12 0 24c0 3.88.92 7.54 2.56 10.78l7.97-6.19z" />
                <path fill="#34A853" d="M24 48c6.48 0 11.93-2.13 15.89-5.81l-7.73-6c-2.15 1.45-4.92 2.3-8.16 2.3-6.26 0-11.57-4.22-13.47-9.91l-7.98 6.19C6.51 42.62 14.62 48 24 48z" />
                <path fill="none" d="M0 0h48v48H0z" />
              </svg>
              <span>Googleでログイン</span>
            </div>
          </button>
        </div>
      </div>

      {/* フッターリンク */}
      <div className="text-center text-sm">
        <Link href="/forgot-password" className="font-medium text-booking-600 hover:text-booking-500">
          パスワードを忘れた方はこちら
        </Link>
      </div>

      {/* 新規登録リンク */}
      {showSignupLink && (
        <div className="text-center text-sm">
          <span className="text-gray-600">アカウントをお持ちでない方は </span>
          <Link href="/signup" className="font-medium text-booking-600 hover:text-booking-500">
            新規登録
          </Link>
        </div>
      )}
    </form>
  );

  if (!showHeader) {
    return formElement;
  }

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

export const LoginForm = (props: LoginFormContentProps) => (
  <Suspense fallback={null}>
    <LoginFormContent {...props} />
  </Suspense>
);

export const LoginClient = () => {
  return (
    <Suspense fallback={null}>
      <LoginFormContent />
    </Suspense>
  );
};
