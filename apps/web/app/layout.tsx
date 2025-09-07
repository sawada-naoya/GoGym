`use client`;

import { Metadata } from "next";
import "./globals.css";

export const metadata: Metadata = {
  title: "GoGym - ジム検索アプリ",
  description: "東京のジムを簡単に検索・比較できるアプリです",
};

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="ja" className="h-full">
      <body className="h-full bg-gray-50 font-sans">{children}</body>
    </html>
  );
}
