// app/layout.tsx
import type { Metadata } from "next";
import "./globals.css";
import React from "react";
import Providers from "./provider";
import Header from "../components/Header";

import { BannerProvider, BannerHost } from "@/components/Banner";
import GlobalBanners from "./GlobalBanners";

export const metadata: Metadata = {
  title: "GoGym - ジム検索アプリ",
  description: "東京のジムを簡単に検索・比較できるアプリです",
};

const RootLayout = ({ children }: { children: React.ReactNode }) => (
  <html lang="ja" className="h-full">
    <body className="h-full bg-gray-50 font-sans">
      <Providers>
        <BannerProvider>
          <Header />
          {/* 全ページ共通のバナー表示領域（画面上部に固定） */}
          <BannerHost />
          {/* （任意）/signup→/login用のフラッシュ読取り */}
          <GlobalBanners />
          <main>{children}</main>
          <footer className="mt-16 border-t bg-white">
            <div className="container mx-auto px-4 py-8 text-sm text-gray-500">© {new Date().getFullYear()} GoGym</div>
          </footer>
        </BannerProvider>
      </Providers>
    </body>
  </html>
);

export default RootLayout;
