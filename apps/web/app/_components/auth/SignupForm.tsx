"use client";

import { useState, Suspense } from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { POST } from "@/lib/api";
import { useBanner } from "@/components/Banner";

const SignUpSchema = z
  .object({
    name: z.string().min(1, "名前は必須です"),
    email: z.string().min(1, "メールアドレスは必須です").email("有効なメールアドレスを入力してください"),
    password: z
      .string()
      .min(8, "パスワードは8文字以上で入力してください")
      .regex(/[A-Z]/, "パスワードには大文字を1文字以上含めてください")
      .regex(/[a-z]/, "パスワードには小文字を1文字以上含めてください")
      .regex(/[0-9]/, "パスワードには数字を1文字以上含めてください"),
    confirmPassword: z.string().min(1, "パスワード（確認）は必須です"),
  })
  .refine((data) => data.password === data.confirmPassword, {
    message: "パスワードが一致しません",
    path: ["confirmPassword"],
  });

type SignUpForm = z.infer<typeof SignUpSchema>;

type SignupFormContentProps = {
  onSuccessCallback?: () => void;
  showHeader?: boolean;
  showLoginLink?: boolean;
};

const SignupFormContent = ({ onSuccessCallback, showHeader = true, showLoginLink = true }: SignupFormContentProps = {}) => {
  const router = useRouter();
  const [loading, setLoading] = useState(false);
  const { error } = useBanner();

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<SignUpForm>({
    resolver: zodResolver(SignUpSchema),
    mode: "onBlur",
  });

  const clearApiError = () => {};

  const onSubmit = async (data: SignUpForm) => {
    setLoading(true);

    try {
      const res = await POST("/api/v1/users", {
        body: {
          name: data.name,
          email: data.email,
          password: data.password,
        },
      });

      if (!res.ok) {
        error("このメールアドレスは既に使用されています");
        return;
      }

      if (onSuccessCallback) {
        onSuccessCallback();
      } else {
        sessionStorage.setItem(
          "flash",
          JSON.stringify({
            variant: "success",
            message: "アカウントの作成に成功しました。ログインしてください。",
          })
        );
        router.push("/");
      }
    } catch (e) {
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
          <input {...register("name", { onChange: clearApiError })} id="name" type="text" autoComplete="name" className={`form-input ${errors.name ? "border-red-500 focus:border-red-500 focus:ring-red-500" : ""}`} placeholder="山田太郎" />
          {errors.name && <p className="mt-1 text-sm text-red-600">{errors.name.message}</p>}
        </div>

        <div>
          <label htmlFor="email" className="form-label">
            メールアドレス
          </label>
          <input {...register("email", { onChange: clearApiError })} id="email" type="email" autoComplete="email" className={`form-input ${errors.email ? "border-red-500 focus:border-red-500 focus:ring-red-500" : ""}`} placeholder="example@example.com" />
          {errors.email && <p className="mt-1 text-sm text-red-600">{errors.email.message}</p>}
        </div>

        <div>
          <label htmlFor="password" className="form-label">
            パスワード
          </label>
          <input {...register("password", { onChange: clearApiError })} id="password" type="password" autoComplete="new-password" className={`form-input ${errors.password ? "border-red-500 focus:border-red-500 focus:ring-red-500" : ""}`} placeholder="パスワードを入力" />
          {errors.password && <p className="mt-1 text-sm text-red-600">{errors.password.message}</p>}
        </div>

        <div>
          <label htmlFor="confirmPassword" className="form-label">
            パスワード（確認）
          </label>
          <input {...register("confirmPassword", { onChange: clearApiError })} id="confirmPassword" type="password" autoComplete="new-password" className={`form-input ${errors.confirmPassword ? "border-red-500 focus:border-red-500 focus:ring-red-500" : ""}`} placeholder="パスワードを再度入力" />
          {errors.confirmPassword && <p className="mt-1 text-sm text-red-600">{errors.confirmPassword.message}</p>}
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

      {/* ログインリンク */}
      {showLoginLink && (
        <div className="text-center text-sm">
          <span className="text-gray-600">既にアカウントをお持ちの方は </span>
          <Link href="/" className="font-medium text-booking-600 hover:text-booking-500">
            ログイン
          </Link>
        </div>
      )}
    </form>
  );

  if (!showHeader) {
    return formElement;
  }

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
            <Link href="/" className="font-medium text-booking-600 hover:text-booking-500">
              ログイン
            </Link>
          </p>
        </div>
        {formElement}
      </div>
    </div>
  );
};

export const SignupForm = (props: SignupFormContentProps) => {
  return (
    <Suspense fallback={null}>
      <SignupFormContent {...props} />
    </Suspense>
  );
};

const SignupClient = () => {
  return (
    <Suspense fallback={null}>
      <SignupFormContent />
    </Suspense>
  );
};

export default SignupClient;
