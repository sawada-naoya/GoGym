"use client";

import Link from "next/link";
import { useState } from "react";
import { useSession, signOut } from "next-auth/react";

const Header = () => {
  const [isMenuOpen, setIsMenuOpen] = useState(false);
  const { data: session, status } = useSession();
  const userName = session?.user?.name;

  return (
    <header className="bg-booking-700 shadow-lg relative z-50">
      <div className="container mx-auto px-4">
        <div className="flex justify-between items-center py-4">
          {/* ロゴ */}
          <Link href="/" className="flex items-center">
            <span className="text-3xl font-bold text-white">GoGym</span>
          </Link>

          {/* デスクトップナビゲーション */}
          <nav className="hidden md:flex items-center space-x-8">
            <Link href="/" className="text-white hover:text-booking-200 transition-colors">
              ホーム
            </Link>
            <Link href="/gyms" className="text-white hover:text-booking-200 transition-colors">
              ジム検索
            </Link>
            <Link href="/favorites" className="text-white hover:text-booking-200 transition-colors">
              お気に入り
            </Link>
            <Link href="/about" className="text-white hover:text-booking-200 transition-colors">
              GoGymについて
            </Link>
          </nav>

          {/* ユーザーメニュー */}
          <div className="hidden md:flex items-center space-x-4">
            {status === "loading" ? (
              <div className="h-8 w-28 animate-pulse rounded bg-white/20" />
            ) : userName ? (
              <>
                <span className="text-white">{userName}</span>
                <button onClick={() => signOut({ callbackUrl: "/login" })} className="bg-white text-booking-700 hover:bg-booking-50 transition-colors px-4 py-2 rounded-md font-medium">
                  ログアウト
                </button>
              </>
            ) : (
              <>
                <Link href="/login" className="text-white hover:text-booking-200 transition-colors px-4 py-2 rounded-md">
                  ログイン
                </Link>
                <Link href="/signup" className="bg-white text-booking-700 hover:bg-booking-50 transition-colors px-4 py-2 rounded-md font-medium">
                  新規登録
                </Link>
              </>
            )}
          </div>

          {/* モバイルメニューボタン */}
          <button className="md:hidden text-white font-medium px-3 py-2 border border-white/20 rounded" onClick={() => setIsMenuOpen(!isMenuOpen)}>
            {isMenuOpen ? "閉じる" : "メニュー"}
          </button>
        </div>

        {/* モバイルメニュー */}
        {isMenuOpen && (
          <div className="md:hidden py-4 border-t border-booking-600">
            <div className="flex flex-col space-y-4">
              <Link href="/" className="text-white hover:text-booking-200 transition-colors py-2">
                ホーム
              </Link>
              <Link href="/gyms" className="text-white hover:text-booking-200 transition-colors py-2">
                ジム検索
              </Link>
              <Link href="/favorites" className="text-white hover:text-booking-200 transition-colors py-2">
                お気に入り
              </Link>
              <Link href="/about" className="text-white hover:text-booking-200 transition-colors py-2">
                GoGymについて
              </Link>
              <div className="flex flex-col space-y-2 pt-4 border-t border-booking-600">
                <Link href="/login" className="text-white hover:text-booking-200 transition-colors py-2">
                  ログイン
                </Link>
                <Link href="/register" className="bg-white text-booking-700 hover:bg-booking-50 transition-colors px-4 py-2 rounded-md font-medium text-center">
                  新規登録
                </Link>
              </div>
            </div>
          </div>
        )}
      </div>
    </header>
  );
};

export default Header;
