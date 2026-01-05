"use client";

import { useState } from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { useTranslation } from "react-i18next";

import { useBanner } from "@/components/Banner";
import { signup } from "@/features/auth/actions";
import { signUpSchema, type SignUpForm } from "@/features/auth/schemas";

type SignupFormContentProps = {
  onSuccessCallback?: () => void;
  showHeader?: boolean;
  showLoginLink?: boolean;
};

const SignupFormContent = ({
  onSuccessCallback,
  showHeader = true,
  showLoginLink = true,
}: SignupFormContentProps = {}) => {
  const { t } = useTranslation("common");
  const router = useRouter();

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
        error(result.error || t("auth.signup.errorEmailTaken"));
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
          message: t("auth.signup.successMessage"),
        }),
      );

      router.push("/login");
    } catch {
      error(t("auth.signup.errorNetworkError"));
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
          <label htmlFor="name" className="form-label">
            {t("auth.signup.nameLabel")}
          </label>
          <input
            {...register("name")}
            id="name"
            type="text"
            autoComplete="name"
            className={`form-input ${errors.name ? "border-red-500 focus:border-red-500 focus:ring-red-500" : ""}`}
            placeholder={t("auth.signup.namePlaceholder")}
          />
          {errors.name && (
            <p className="mt-1 text-sm text-red-600">{errors.name.message}</p>
          )}
        </div>

        <div>
          <label htmlFor="email" className="form-label">
            {t("auth.signup.emailLabel")}
          </label>
          <input
            {...register("email")}
            id="email"
            type="email"
            autoComplete="email"
            className={`form-input ${errors.email ? "border-red-500 focus:border-red-500 focus:ring-red-500" : ""}`}
            placeholder={t("auth.signup.emailPlaceholder")}
          />
          {errors.email && (
            <p className="mt-1 text-sm text-red-600">{errors.email.message}</p>
          )}
        </div>

        <div>
          <label htmlFor="password" className="form-label">
            {t("auth.signup.passwordLabel")}
          </label>
          <input
            {...register("password")}
            id="password"
            type="password"
            autoComplete="new-password"
            className={`form-input ${errors.password ? "border-red-500 focus:border-red-500 focus:ring-red-500" : ""}`}
            placeholder={t("auth.signup.passwordPlaceholder")}
          />
          {errors.password && (
            <p className="mt-1 text-sm text-red-600">
              {errors.password.message}
            </p>
          )}
        </div>

        <div>
          <label htmlFor="confirmPassword" className="form-label">
            {t("auth.signup.confirmPasswordLabel")}
          </label>
          <input
            {...register("confirmPassword")}
            id="confirmPassword"
            type="password"
            autoComplete="new-password"
            className={`form-input ${errors.confirmPassword ? "border-red-500 focus:border-red-500 focus:ring-red-500" : ""}`}
            placeholder={t("auth.signup.confirmPasswordPlaceholder")}
          />
          {errors.confirmPassword && (
            <p className="mt-1 text-sm text-red-600">
              {errors.confirmPassword.message}
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
          {loading ? t("auth.signup.signingUp") : t("auth.signup.signupButton")}
        </button>
      </div>

      {showLoginLink && (
        <div className="text-center text-sm">
          <span className="text-gray-600">{t("auth.signup.hasAccount")} </span>
          <Link
            href="/login"
            className="font-medium text-booking-600 hover:text-booking-500"
          >
            {t("auth.signup.loginLink")}
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
          <h2 className="mt-6 text-center text-3xl font-bold tracking-tight text-gray-900">
            {t("auth.signup.title")}
          </h2>
          <p className="mt-2 text-center text-sm text-gray-600">
            {t("auth.signup.hasAccount")}{" "}
            <Link
              href="/login"
              className="font-medium text-booking-600 hover:text-booking-500"
            >
              {t("auth.signup.loginLink")}
            </Link>
          </p>
        </div>
        {formElement}
      </div>
    </div>
  );
};

export default SignupFormContent;
