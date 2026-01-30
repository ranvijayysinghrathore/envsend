import type { Metadata } from "next";
import { Inter } from "next/font/google";
import { Navbar } from "@/components/Navbar";
import { Analytics } from "@vercel/analytics/next"
import "./globals.css";

const inter = Inter({ subsets: ["latin"], variable: "--font-inter" });

export const metadata: Metadata = {
  title: "EnvSend - Zero-Knowledge Secret Transfer",
  description: "Securely transfer .env files and secrets without exposing plaintext.",
  metadataBase: new URL("https://envsend.vercel.app"),
  openGraph: {
    title: "EnvSend - Zero-Knowledge Secret Transfer",
    description: "End-to-end encrypted secret sharing. No servers, no logs, pure client-side AES-256.",
    url: "https://envsend.vercel.app",
    siteName: "EnvSend",
    images: [
      {
        url: "/og-image.png",
        width: 1200,
        height: 630,
        alt: "EnvSend Interface",
      },
    ],
    locale: "en_US",
    type: "website",
  },
  twitter: {
    card: "summary_large_image",
    title: "EnvSend - Zero-Knowledge Secret Transfer",
    description: "End-to-end encrypted secret sharing. No servers, no logs.",
    creator: "@Ranvijayy_",
    images: ["/og-image.png"],
  },
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
        <Analytics />
      </body>
    </html>
  );
}
