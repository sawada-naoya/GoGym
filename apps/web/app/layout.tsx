import type { Metadata } from "next";
import "./globals.css";
import React from "react";
import ClientProviders from "./ClientProviders";
import Header from "../components/Header";
import { BannerHost } from "@/components/Banner";

export const metadata: Metadata = {
  metadataBase: new URL("https://www.gogym.fitness"),
  title: {
    default: "GoGym｜筋トレ・ワークアウト記録（トレーニングノート）",
    template: "%s | GoGym",
  },
  description: "筋トレ・ワークアウトを簡単に記録できるトレーニングノート。部位・種目・重量・レップ・メモを残して成長を可視化。",
  applicationName: "GoGym",
  alternates: { canonical: "/" },
  robots: { index: true, follow: true },
  openGraph: {
    type: "website",
    url: "https://www.gogym.fitness",
    siteName: "GoGym",
    title: "GoGym｜筋トレ・ワークアウト記録（トレーニングノート）",
    description: "筋トレ・ワークアウトを簡単に記録。部位・種目・重量・レップ・メモで成長を可視化。",
    locale: "ja_JP",
  },
  twitter: {
    card: "summary",
    title: "GoGym｜筋トレ・ワークアウト記録",
    description: "トレーニングノートで記録を習慣化。",
  },
  icons: {
    icon: [{ url: "/images/favicon.ico" }, { url: "/images/favicon-16x16.png", sizes: "16x16", type: "image/png" }, { url: "/images/favicon-32x32.png", sizes: "32x32", type: "image/png" }],
    apple: "/images/apple-touch-icon.png",
  },
  manifest: "/images/site.webmanifest",
};

const RootLayout = ({ children }: { children: React.ReactNode }) => (
  <html lang="ja" className="h-full">
    <body className="h-full bg-gray-50 font-sans">
      <ClientProviders>
        <Header />
        <BannerHost />
        <main>{children}</main>
        <footer className="mt-16 border-t bg-white">
          <div className="container mx-auto px-4 py-8 text-sm text-gray-500">© {new Date().getFullYear()} GoGym</div>
        </footer>
      </ClientProviders>
    </body>
  </html>
);

export default RootLayout;
