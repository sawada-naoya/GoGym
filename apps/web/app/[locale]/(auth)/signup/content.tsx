"use client";

import { useState } from "react";
import Link from "next/link";
import { useRouter, useParams } from "next/navigation";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";

import { useBanner } from "@/components/Banner";
import { signup } from "@/features/auth/actions";
import { signUpSchema, type SignUpForm } from "@/features/auth/schemas/signup";

type SignupFormContentProps = {
  onSuccessCallback?: () => void;
  showHeader?: boolean;
  showLoginLink?: boolean;
};

const SignupFormContent = ({ onSuccessCallback, showHeader = true, showLoginLink = true }: SignupFormContentProps = {}) => {
  const router = useRouter();
  const params = useParams<{ locale: string }>();
  const locale = params?.locale ?? "ja";

  const [loading, setLoading] = useState(false);
  const { error } = useBanner();

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<SignUpForm>({
    resolver: zodResolver(signUpSchema),
    mode: "onBlur",
  });

  const onSubmit = async (data: SignUpForm) => {
    setLoading(true);
    try {
      const result = await signup({
        name: data.name,
        email: data.email,
        password: data.password,
      });

      if (!result.success) {
        error(result.error || "このメールアドレスは既に使用されています");
        return;
      }

      if (onSuccessCallback) {
        onSuccessCallback();
        return;
      }

      sessionStorage.setItem(
        "flash",
        JSON.stringify({
          variant: "success",
          message: "アカウントの作成に成功しました。ログインしてください。",
        })
      );

      // i18n前提: "/" へ戻すのはやめろ
      router.push(`/${locale}/login`);
    } catch {
      error("ネットワークエラーが発生しました。時間を置いて再度お試しください。");
    } finally {
      setLoading(false);
    }
  };

  const formElement = (
    <form className={showHeader ? "mt-8 space-y-6" : "space-y-4"} onSubmit={handleSubmit(onSubmit)} noValidate>
      <div className="space-y-4">
        <div>
          <label htmlFor="name" className="form-label">
            名前
          </label>
          <input {...register("name")} id="name" type="text" autoComplete="name" className={`form-input ${errors.name ? "border-red-500 focus:border-red-500 focus:ring-red-500" : ""}`} placeholder="山田太郎" />
          {errors.name && <p className="mt-1 text-sm text-red-600">{errors.name.message}</p>}
        </div>

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
          <input {...register("password")} id="password" type="password" autoComplete="new-password" className={`form-input ${errors.password ? "border-red-500 focus:border-red-500 focus:ring-red-500" : ""}`} placeholder="パスワードを入力" />
          {errors.password && <p className="mt-1 text-sm text-red-600">{errors.password.message}</p>}
        </div>

        <div>
          <label htmlFor="confirmPassword" className="form-label">
            パスワード（確認）
          </label>
          <input {...register("confirmPassword")} id="confirmPassword" type="password" autoComplete="new-password" className={`form-input ${errors.confirmPassword ? "border-red-500 focus:border-red-500 focus:ring-red-500" : ""}`} placeholder="パスワードを再度入力" />
          {errors.confirmPassword && <p className="mt-1 text-sm text-red-600">{errors.confirmPassword.message}</p>}
        </div>
      </div>

      <div>
        <button
          type="submit"
          disabled={loading}
          className="group relative w-full flex justify-center py-3 px-4 border border-transparent text-sm font-medium rounded-md text-white bg-booking-600 hover:bg-booking-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-booking-500 disabled:opacity-50 disabled:cursor-not-allowed transition-colors duration-200"
        >
          {loading ? "登録中..." : "新規登録"}
        </button>
      </div>

      {showLoginLink && (
        <div className="text-center text-sm">
          <span className="text-gray-600">既にアカウントをお持ちの方は </span>
          <Link href={`/${locale}/login`} className="font-medium text-booking-600 hover:text-booking-500">
            ログイン
          </Link>
        </div>
      )}
    </form>
  );

  if (!showHeader) return formElement;

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
            <Link href={`/${locale}/login`} className="font-medium text-booking-600 hover:text-booking-500">
              ログイン
            </Link>
          </p>
        </div>
        {formElement}
      </div>
    </div>
  );
};

export default SignupFormContent;
