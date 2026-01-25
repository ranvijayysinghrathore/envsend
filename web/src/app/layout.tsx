import type { Metadata } from "next";
import { Inter } from "next/font/google";
import { Navbar } from "@/components/Navbar";
import "./globals.css";

const inter = Inter({ subsets: ["latin"], variable: "--font-inter" });

export const metadata: Metadata = {
  title: "EnvSend - Zero-Knowledge Secret Transfer",
  description: "Securely transfer .env files and secrets without exposing plaintext.",
  icons: {
    icon: "/icon.png",
    apple: "/icon.png",
  },
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body
        suppressHydrationWarning
        className={`${inter.variable} font-sans antialiased bg-[#050505] text-white selection:bg-orange-500/30`}
      >
        <div className="bg-noise" />
        <Navbar />
        {children}
      </body>
    </html>
  );
}
