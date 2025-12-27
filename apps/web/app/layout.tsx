import type { Metadata } from "next";
import "./globals.css";
import React from "react";
import ClientProviders from "./ClientProviders";
import Header from "../components/Header";
import { BannerHost } from "@/components/Banner";

export const metadata: Metadata = {
  title: "GoGym - トレーニング記録アプリ",
  description: "日々のトレーニングを記録し、成長を可視化するアプリ",
};

const RootLayout = ({ children }: { children: React.ReactNode }) => (
  <html lang="ja" className="h-full">
    <body className="h-full bg-gray-50 font-sans">
      <ClientProviders>
        <Header />
        <BannerHost />
        <main>{children}</main>
        <footer className="mt-16 border-t bg-white">
          <div className="container mx-auto px-4 py-8 text-sm text-gray-500">
            © {new Date().getFullYear()} GoGym
          </div>
        </footer>
      </ClientProviders>
    </body>
  </html>
);

export default RootLayout;
