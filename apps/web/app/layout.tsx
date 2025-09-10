import type { Metadata } from "next";
import "./globals.css";
import React from "react";

export const metadata: Metadata = {
  title: "GoGym - ジム検索アプリ",
  description: "東京のジムを簡単に検索・比較できるアプリです",
};

const RootLayout = ({ children }: { children: React.ReactNode }) => {
  return (
    <html lang="ja" className="h-full">
      <body className="h-full bg-gray-50 font-sans">
        <main>{children}</main>

        <footer className="mt-16 border-t bg-white">
          <div className="container mx-auto px-4 py-8 text-sm text-gray-500">© {new Date().getFullYear()} GoGym</div>
        </footer>
      </body>
    </html>
  );
};

export default RootLayout;
