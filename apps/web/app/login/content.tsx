"use client";

import { useState } from "react";
import Link from "next/link";
import { useRouter, useSearchParams } from "next/navigation";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { useTranslation } from "react-i18next";

import { signIn } from "next-auth/react";
import { useBanner } from "@/components/Banner";
import { loginSchema, type LoginForm } from "@/features/auth/schemas";

type LoginFormContentProps = {
  showHeader?: boolean;
  showSignupLink?: boolean;
};

const normalizeCallbackUrl = (candidate: string) => {
  // open redirect対策：外部URLは禁止。相対パスのみ許可
  if (!candidate.startsWith("/")) return "/workout";
  return candidate;
};

const LoginFormContent = ({
  showHeader = true,
  showSignupLink = true,
}: LoginFormContentProps = {}) => {
  const { t } = useTranslation("common");
  const router = useRouter();
  const searchParams = useSearchParams();

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
      const callbackUrl = normalizeCallbackUrl(rawCallbackUrl);

      const result = await signIn("credentials", {
        email: data.email,
        password: data.password,
        redirect: false,
        callbackUrl,
      });

      if (result?.error) {
        if (result.error === "CredentialsSignin") {
          error(t("auth.login.errorInvalidCredentials"));
        } else {
          error(t("auth.login.errorLoginFailed"));
        }
        return;
      }

      sessionStorage.setItem(
        "flash",
        JSON.stringify({
          variant: "success",
          message: t("auth.login.successMessage"),
        }),
      );

      router.replace(result?.url ?? callbackUrl);
    } finally {
      setLoading(false);
    }
  };

  const formElement = (
    <form
      className={showHeader ? "mt-8 space-y-6" : "space-y-4"}
      onSubmit={handleSubmit(onSubmit)}
      noValidate
    >
      <div className="space-y-4">
        <div>
          <label htmlFor="email" className="form-label">
            {t("auth.login.emailLabel")}
          </label>
          <input
            {...register("email")}
            id="email"
            type="email"
            autoComplete="email"
            className={`form-input ${errors.email ? "border-red-500 focus:border-red-500 focus:ring-red-500" : ""}`}
            placeholder={t("auth.login.emailPlaceholder")}
          />
          {errors.email && (
            <p className="mt-1 text-sm text-red-600">{errors.email.message}</p>
          )}
        </div>

        <div>
          <label htmlFor="password" className="form-label">
            {t("auth.login.passwordLabel")}
          </label>
          <input
            {...register("password")}
            id="password"
            type="password"
            autoComplete="current-password"
            className={`form-input ${errors.password ? "border-red-500 focus:border-red-500 focus:ring-red-500" : ""}`}
            placeholder={t("auth.login.passwordPlaceholder")}
          />
          {errors.password && (
            <p className="mt-1 text-sm text-red-600">
              {errors.password.message}
            </p>
          )}
        </div>
      </div>

      <div>
        <button
          type="submit"
          disabled={loading}
          className="group relative w-full flex justify-center py-3 px-4 border border-transparent text-sm font-medium rounded-md text-white bg-booking-600 hover:bg-booking-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-booking-500 disabled:opacity-50 disabled:cursor-not-allowed transition-colors duration-200"
        >
          {loading ? t("auth.login.loggingIn") : t("auth.login.loginButton")}
        </button>
      </div>

      <div className="text-center text-sm">
        <Link
          href="/forgot-password"
          className="font-medium text-booking-600 hover:text-booking-500"
        >
          {t("auth.login.forgotPassword")}
        </Link>
      </div>

      {showSignupLink && (
        <div className="text-center text-sm">
          <span className="text-gray-600">{t("auth.login.noAccount")} </span>
          <Link
            href="/signup"
            className="font-medium text-booking-600 hover:text-booking-500"
          >
            {t("auth.login.signupLink")}
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
          <h2 className="mt-6 text-center text-3xl font-extrabold text-gray-900">
            {t("auth.login.title")}
          </h2>
          <p className="mt-2 text-center text-sm text-gray-600">
            {t("auth.login.subtitle")}
          </p>
        </div>
        {formElement}
      </div>
    </div>
  );
};

export default LoginFormContent;
