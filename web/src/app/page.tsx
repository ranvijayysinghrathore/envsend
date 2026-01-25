"use client";

import { TerminalDemo } from "@/components/TerminalDemo";
import { FeaturesBento } from "@/components/FeaturesBento";
import { Download, Shield, Lock, Share2, ArrowRight } from "lucide-react";
import Link from "next/link";
import { motion } from "framer-motion";

export default function Home() {
  return (
    <div className="relative min-h-screen flex flex-col">
      <main className="flex-grow relative z-10 pt-32 pb-20 px-6">
        
        {/* Planet Rise Glow & SVG Horizon */}
        <div className="absolute top-[4%] left-1/2 -translate-x-1/2 w-[140%] h-[800px] pointer-events-none select-none">
          {/* Main Glow */}
          <div className="absolute bottom-0 left-1/2 -translate-x-1/2 w-[80%] h-[300px] bg-orange-500/20 blur-[120px] rounded-full mix-blend-screen" />
          
          {/* Sharp SVG Horizon Line */}
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            transition={{ duration: 1.5 }}
            className="absolute bottom-[225px] left-0 right-0 flex justify-center items-end"
          >
            <svg 
              width="100%" 
              height="200" 
              viewBox="0 0 1440 200" 
              fill="none" 
              xmlns="http://www.w3.org/2000/svg"
              className="w-full text-orange-500 [mask-image:linear-gradient(to_top,black,transparent)]"
            >
              {/* Animated Arc */}
              <motion.path
                d="M-100 200 C400 0 1040 0 1540 200"
                stroke="url(#glow-gradient)"
                strokeWidth="2"
                fill="none"
                initial={{ pathLength: 0, opacity: 0 }}
                animate={{ 
                  pathLength: 1, 
                  opacity: [0.4, 0.7, 0.4],
                  y: [0, -3, 0]
                }}
                transition={{ 
                  pathLength: { duration: 1.5, ease: "easeInOut" },
                  opacity: { duration: 8, repeat: Infinity, ease: "easeInOut" },
                  y: { duration: 10, repeat: Infinity, ease: "easeInOut" }
                }}
              />
              <defs>
                <linearGradient id="glow-gradient" x1="0" y1="0" x2="1440" y2="0" gradientUnits="userSpaceOnUse">
                  <stop offset="0%" stopColor="#FF5500" stopOpacity="0" />
                  <stop offset="20%" stopColor="#FF5500" stopOpacity="0.3" />
                  <stop offset="50%" stopColor="#FF9900" stopOpacity="0.8" />
                  <stop offset="80%" stopColor="#FF5500" stopOpacity="0.3" />
                  <stop offset="100%" stopColor="#FF5500" stopOpacity="0" />
                </linearGradient>
              </defs>
            </svg>
          </motion.div>
        </div>

        {/* Hero Content */}
        <section className="relative max-w-5xl mx-auto text-center mb-32">
          
          {/* Announcement Pill */}
          <motion.div 
            initial={{ opacity: 0, y: 10 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.5 }}
            className="inline-flex items-center gap-2 px-3 py-1 mb-8 rounded-full bg-white/[0.03] border border-white/[0.08] backdrop-blur-sm text-xs font-medium text-white/60 hover:text-white/90 hover:bg-white/[0.06] transition-colors cursor-default"
          >
            <span className="w-1.5 h-1.5 rounded-full bg-orange-500 animate-pulse" />
            <span>New Version 1.0 Release</span>
            <ArrowRight className="w-3 h-3 ml-1 text-white/40" />
          </motion.div>

          {/* Headline */}
          <motion.h1 
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.6, delay: 0.1 }}
            className="text-5xl sm:text-7xl md:text-8xl font-bold tracking-tight mb-8 leading-[1.1] text-balance"
          >
            <span className="text-white">Share Secrets.</span>
            <br />
            <span className="text-white/40">Zero Knowledge.</span>
          </motion.h1>

          {/* Subtext */}
          <motion.p 
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.6, delay: 0.2 }}
            className="text-lg sm:text-xl text-white/50 max-w-2xl mx-auto mb-12 leading-relaxed font-light text-balance"
          >
            Production-grade encryption for your .env files. <br className="hidden sm:block" />
            No servers. No logs. Just pure AES-256-GCM.
          </motion.p>

          {/* CTAs */}
          <motion.div 
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.6, delay: 0.3 }}
            className="flex flex-col sm:flex-row items-center justify-center gap-4 mb-24"
          >
            <Link 
              href="/docs/install"
              className="h-12 px-8 rounded-full bg-gradient-to-b from-orange-500 to-orange-600 text-white font-medium flex items-center justify-center gap-2 hover:shadow-[0_0_20px_-5px_rgba(255,85,0,0.4)] transition-all active:scale-95"
            >
              Install for Free
            </Link>
            <button className="h-12 px-8 rounded-full text-white/70 hover:text-white font-medium flex items-center justify-center gap-2 transition-colors group">
              <span className="border-b border-white/0 group-hover:border-white/40 transition-all">Watch Demo</span>
              <div className="w-6 h-6 rounded-full bg-white/10 flex items-center justify-center group-hover:bg-white/20 transition-colors">
                <div className="w-0 h-0 border-t-[3px] border-t-transparent border-l-[6px] border-l-white border-b-[3px] border-b-transparent ml-0.5" />
              </div>
            </button>
          </motion.div>

          {/* Glass Terminal Container */}
          <motion.div 
            initial={{ opacity: 0, scale: 0.95 }}
            animate={{ opacity: 1, scale: 1 }}
            transition={{ duration: 0.8, delay: 0.4 }}
            className="relative mx-auto max-w-4xl"
          >
            <div className="absolute -inset-px bg-gradient-to-b from-white/10 to-transparent rounded-xl pointer-events-none" />
            <div className="bg-[#0A0A0A]/80 backdrop-blur-xl rounded-xl border border-white/[0.08] shadow-2xl p-2">
              <TerminalDemo />
            </div>
            
            {/* Ambient Reflection */}
            <div className="absolute top-0 left-1/2 -translate-x-1/2 w-3/4 h-px bg-gradient-to-r from-transparent via-white/20 to-transparent" />
          </motion.div>
        </section>

        {/* Trusted By (Muted) */}
        <section className="text-center mb-40">
          <p className="text-xs font-semibold text-white/20 uppercase tracking-widest mb-10">Trusted by security teams at</p>
          <div className="flex flex-wrap justify-center gap-12 sm:gap-20 opacity-30 grayscale hover:grayscale-0 hover:opacity-60 transition-all duration-700">
            {['Acme Corp', 'GlobalBank', 'Nebula', 'LightBox', 'Vercel'].map(brand => (
              <span key={brand} className="text-lg font-bold flex items-center gap-2 text-white">
                <div className="w-5 h-5 rounded-full bg-current opacity-20" />
                {brand}
              </span>
            ))}
          </div>
        </section>

        {/* Feature Grid */}
        <section id="features" className="max-w-6xl mx-auto mb-32">
          <FeaturesBento />
        </section>

      </main>
      
      <footer className="py-12 border-t border-white/[0.05]">
        <div className="max-w-6xl mx-auto px-6 flex flex-col sm:flex-row justify-between items-center gap-4 text-xs text-white/30">
          <p>&copy; 2026 EnvSend. Open Source MIT.</p>
          <div className="flex gap-6">
            <a href="#" className="hover:text-white transition-colors">Privacy</a>
            <a href="#" className="hover:text-white transition-colors">Terms</a>
            <a href="https://x.com/Ranvijayy_" className="hover:text-white transition-colors">Twitter</a>
          </div>
        </div>
      </footer>
    </div>
  );
}

const features = [
  {
    title: "Zero Knowledge",
    description: "Encryption happens on your device. The server never sees the key or the plaintext data. Period.",
    icon: <Shield className="w-5 h-5" />,
  },
  {
    title: "Nuclear Mode",
    description: "Shamir's Secret Sharing splits the key into parts. Requires M-of-N people to decrypt.",
    icon: <Lock className="w-5 h-5" />,
  },
  {
    title: "Ephemeral Links",
    description: "Secrets self-destruct after one view or a set time limit. No traces left behind.",
    icon: <Share2 className="w-5 h-5" />,
  },
];
