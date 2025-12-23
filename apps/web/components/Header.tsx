"use client";

import Link from "next/link";
import { useState } from "react";
import { useSession, signOut } from "next-auth/react";

const Header = () => {
  const [isMenuOpen, setIsMenuOpen] = useState(false);
  const { data: session, status } = useSession();
  const userName = session?.user?.name;

  const handleSignOut = async () => {
    await signOut({ callbackUrl: "/", redirect: true });
  };

  return (
    <header className="bg-booking-700 shadow relative z-50">
      <div className="container mx-auto px-3 md:px-4">
        <div className="flex justify-between items-center py-2 md:py-4">
          {/* ロゴ */}
          <Link href="/" className="flex items-center">
            <span className="text-xl md:text-3xl font-bold text-white">
              GoGym
            </span>
          </Link>

          {/* デスクトップナビゲーション */}
          <nav className="hidden md:flex items-center space-x-8">
            <Link
              href="/"
              className="text-white hover:text-booking-200 transition-colors"
            >
              トレーニングノート
            </Link>
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
                  ログアウト
                </button>
              </>
            ) : null}
          </div>

          {/* ハンバーガーメニューボタン（モバイル） */}
          <button
            className="md:hidden text-white p-1"
            onClick={() => setIsMenuOpen(!isMenuOpen)}
            aria-label="メニュー"
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
                トレーニングノート
              </Link>
              {userName && (
                <div className="flex flex-col space-y-2 pt-3 border-t border-booking-600">
                  <span className="text-white py-1 text-sm">{userName}</span>
                  <button
                    onClick={handleSignOut}
                    className="bg-white text-booking-700 hover:bg-booking-50 transition-colors px-3 py-1.5 rounded-md font-medium text-center text-sm"
                  >
                    ログアウト
                  </button>
                </div>
              )}
            </div>
          </div>
        )}
      </div>
    </header>
  );
};

export default Header;
