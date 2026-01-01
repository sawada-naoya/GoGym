"use client";

import Link from "next/link";
import { useState, useEffect } from "react";
import { useSession, signOut } from "next-auth/react";
import { useTranslation } from "react-i18next";
import ContactModal from "./ContactModal";

const Header = () => {
  const [isMenuOpen, setIsMenuOpen] = useState(false);
  const [isMounted, setIsMounted] = useState(false);
  const [isContactModalOpen, setIsContactModalOpen] = useState(false);
  const { data: session, status } = useSession();
  const userName = session?.user?.name;
  const { t, i18n } = useTranslation("common");

  const currentLocale = i18n.language || "ja";

  useEffect(() => {
    setIsMounted(true);
  }, []);

  const handleSignOut = async () => {
    await signOut({ callbackUrl: "/", redirect: true });
  };

  const handleLocaleChange = (newLocale: string) => {
    i18n.changeLanguage(newLocale);
  };

  return (
    <header className="bg-booking-700 shadow relative z-50">
      <div className="container mx-auto px-3 md:px-4">
        <div className="flex justify-between items-center py-2 md:py-4">
          {/* ロゴ */}
          <Link href="/" className="flex items-center">
            <span className="text-xl md:text-3xl font-bold text-white">
              {t("header.appName")}
            </span>
          </Link>

          {/* デスクトップナビゲーション */}
          <nav className="hidden md:flex items-center space-x-8">
            <Link
              href="/"
              className="text-white hover:text-booking-200 transition-colors"
            >
              {t("header.trainingNote")}
            </Link>
            <button
              onClick={() => setIsContactModalOpen(true)}
              className="text-white hover:text-booking-200 transition-colors"
            >
              {t("header.contact")}
            </button>
          </nav>

          {/* ユーザーメニュー（デスクトップ） */}
          <div className="hidden md:flex items-center space-x-4">
            {status === "loading" ? (
              <div className="h-8 w-28 animate-pulse rounded bg-white/20" />
            ) : userName ? (
              <>
                <span className="text-white">{userName}</span>
                <button
                  onClick={handleSignOut}
                  className="bg-white text-booking-700 hover:bg-booking-50 transition-colors px-4 py-2 rounded-md font-medium"
                >
                  {t("header.logout")}
                </button>
              </>
            ) : null}

            {/* 言語切り替え */}
            <div className="flex items-center gap-2">
              <button
                onClick={() => handleLocaleChange("ja")}
                className={`px-2 py-1 rounded transition-colors ${!isMounted || currentLocale === "ja" ? "bg-white text-booking-700 font-medium" : "text-white hover:text-booking-200"}`}
              >
                日
              </button>
              <button
                onClick={() => handleLocaleChange("en")}
                className={`px-2 py-1 rounded transition-colors ${isMounted && currentLocale === "en" ? "bg-white text-booking-700 font-medium" : "text-white hover:text-booking-200"}`}
              >
                EN
              </button>
            </div>
          </div>

          {/* ハンバーガーメニューボタン（モバイル） */}
          <button
            className="md:hidden text-white p-1"
            onClick={() => setIsMenuOpen(!isMenuOpen)}
            aria-label={t("header.menu")}
          >
            <svg
              className="w-6 h-6"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              {isMenuOpen ? (
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M6 18L18 6M6 6l12 12"
                />
              ) : (
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M4 6h16M4 12h16M4 18h16"
                />
              )}
            </svg>
          </button>
        </div>

        {/* モバイルメニュー */}
        {isMenuOpen && (
          <div className="md:hidden py-3 border-t border-booking-600">
            <div className="flex flex-col space-y-3">
              <Link
                href="/"
                className="text-white hover:text-booking-200 transition-colors py-2 text-sm"
                onClick={() => setIsMenuOpen(false)}
              >
                {t("header.trainingNote")}
              </Link>
              <button
                onClick={() => {
                  setIsContactModalOpen(true);
                  setIsMenuOpen(false);
                }}
                className="text-white hover:text-booking-200 transition-colors py-2 text-sm text-left"
              >
                {t("header.contact")}
              </button>

              {/* 言語切り替え（モバイル） */}
              <div className="flex items-center gap-2 py-2">
                <button
                  onClick={() => handleLocaleChange("ja")}
                  className={`px-3 py-1.5 rounded transition-colors text-sm ${!isMounted || currentLocale === "ja" ? "bg-white text-booking-700 font-medium" : "bg-booking-600 text-white hover:bg-booking-500"}`}
                >
                  日本語
                </button>
                <button
                  onClick={() => handleLocaleChange("en")}
                  className={`px-3 py-1.5 rounded transition-colors text-sm ${isMounted && currentLocale === "en" ? "bg-white text-booking-700 font-medium" : "bg-booking-600 text-white hover:bg-booking-500"}`}
                >
                  English
                </button>
              </div>

              {userName && (
                <div className="flex flex-col space-y-2 pt-3 border-t border-booking-600">
                  <span className="text-white py-1 text-sm">{userName}</span>
                  <button
                    onClick={handleSignOut}
                    className="bg-white text-booking-700 hover:bg-booking-50 transition-colors px-3 py-1.5 rounded-md font-medium text-center text-sm"
                  >
                    {t("header.logout")}
                  </button>
                </div>
              )}
            </div>
          </div>
        )}
      </div>

      {/* お問い合わせモーダル */}
      <ContactModal
        isOpen={isContactModalOpen}
        onClose={() => setIsContactModalOpen(false)}
      />
    </header>
  );
};

export default Header;
