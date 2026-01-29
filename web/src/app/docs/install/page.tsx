"use client";

import { Download, Terminal, Command, Monitor, Apple, Shield } from "lucide-react";
import Link from "next/link";
import { useState } from "react";
import { motion } from "framer-motion";
import clsx from "clsx";

export default function InstallPage() {
  const [activeTab, setActiveTab] = useState<"windows" | "mac" | "linux">("windows");

  return (
    <div className="min-h-screen bg-black text-white selection:bg-orange-500/30">
      
      {/* Background Ambience */}
      <div className="fixed inset-0 pointer-events-none">
        <div className="absolute top-[-10%] right-[-10%] w-[50%] h-[50%] bg-orange-500/10 blur-[120px] rounded-full mix-blend-screen opacity-50" />
        <div className="absolute bottom-[-10%] left-[-10%] w-[50%] h-[50%] bg-blue-500/10 blur-[120px] rounded-full mix-blend-screen opacity-30" />
      </div>

      <main className="relative z-10 max-w-7xl mx-auto px-6 pt-32 pb-20">
        
        {/* Hero Section */}
        <div className="text-center max-w-3xl mx-auto mb-20">
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.5 }}
          >
            <h1 className="text-3xl sm:text-5xl font-bold tracking-tight mb-6">
              <span className="text-white">Get </span>
              <span className="text-transparent bg-clip-text bg-gradient-to-r from-orange-500 to-orange-400">EnvSend</span>
            </h1>
            <p className="text-xl text-white/50 leading-relaxed">
              Production-grade secret management for your terminal.
              <br className="hidden sm:block" />
              Secure, fast, and completely zero-knowledge.
            </p>
          </motion.div>
        </div>

        {/* Installation Grid */}
        <div className="grid grid-cols-1 lg:grid-cols-12 gap-8">
          
          {/* Left Column: Platform Selection (Desktop) */}
          <div className="lg:col-span-3 space-y-2">
            <p className="text-xs font-bold text-white/30 uppercase tracking-widest mb-4 px-2">Select Platform</p>
            {[
              { id: "windows", label: "Windows", icon: <Monitor className="w-4 h-4" /> },
              { id: "mac", label: "macOS", icon: <Apple className="w-4 h-4" /> },
              { id: "linux", label: "Linux", icon: <Terminal className="w-4 h-4" /> },
            ].map((platform) => (
              <button
                key={platform.id}
                onClick={() => setActiveTab(platform.id as any)}
                className={clsx(
                  "w-full flex items-center gap-3 px-4 py-3 rounded-lg text-sm font-medium transition-all text-left group border",
                  activeTab === platform.id
                    ? "bg-white/10 border-white/10 text-white shadow-lg"
                    : "bg-transparent border-transparent text-white/40 hover:text-white hover:bg-white/5"
                )}
              >
                <div className={clsx(
                  "w-8 h-8 rounded-md flex items-center justify-center transition-colors",
                  activeTab === platform.id ? "bg-white/10 text-orange-500" : "bg-white/5 text-white/40"
                )}>
                  {platform.icon}
                </div>
                <span>{platform.label}</span>
                {activeTab === platform.id && (
                  <motion.div layoutId="active-indicator" className="ml-auto w-1.5 h-1.5 rounded-full bg-orange-500 shadow-[0_0_10px_rgba(249,115,22,0.5)]" />
                )}
              </button>
            ))}
          </div>

          {/* Center Column: Main Installation Card */}
          <div className="lg:col-span-9">
            <motion.div
              key={activeTab}
              initial={{ opacity: 0, x: 20 }}
              animate={{ opacity: 1, x: 0 }}
              exit={{ opacity: 0, x: -20 }}
              transition={{ duration: 0.3 }}
              className="relative overflow-hidden rounded-3xl border border-white/10 bg-[#0A0A0A]/50 backdrop-blur-xl p-8 sm:p-12 min-h-[500px] flex flex-col justify-center"
            >
              {activeTab === "windows" && (
                <div className="grid md:grid-cols-2 gap-12 items-center">
                  <div className="space-y-8">
                    <div>
                      <div className="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-blue-500/10 border border-blue-500/20 text-blue-400 text-xs font-medium mb-4">
                        <Shield className="w-3 h-3" />
                        <span>Recommended for Windows 10/11</span>
                      </div>
                      <h2 className="text-3xl font-bold mb-4">EnvSend Installer</h2>
                      <p className="text-white/50 leading-relaxed">
                        The easiest way to get started. Includes the CLI binary and automatically configures your system PATH for global access.
                      </p>
                    </div>

                    <div className="space-y-4">
                      <a
                        href="https://github.com/ranvijayysinghrathore/envsend/releases/latest/download/EnvSendSetup.exe"
                        className="flex items-center justify-center gap-3 w-full bg-white text-black h-14 rounded-xl font-bold text-lg hover:bg-gray-200 transition-all active:scale-95 shadow-xl shadow-white/5"
                      >
                        <Download className="w-5 h-5" />
                        Download Installer
                        <span className="text-black/30 text-sm font-normal ml-1">(.exe)</span>
                      </a>
                      <p className="text-center text-xs text-white/30">
                        Version 1.0.0 • 64-bit • 7.2 MB
                      </p>
                    </div>
                  </div>

                  {/* Visual Guide */}
                  <div className="relative">
                    <div className="absolute inset-0 bg-gradient-to-br from-orange-500/20 to-transparent blur-3xl opacity-30" />
                    <div className="relative rounded-xl border border-white/10 bg-black/80 overflow-hidden shadow-2xl">
                      <div className="flex items-center gap-1.5 px-4 py-3 border-b border-white/10 bg-white/5">
                        <div className="w-3 h-3 rounded-full bg-white/20" />
                        <div className="w-3 h-3 rounded-full bg-white/20" />
                      </div>
                      <div className="p-6 font-mono text-xs sm:text-sm space-y-4 text-white/70">
                        <div className="flex gap-3">
                          <span className="text-green-500">✔</span>
                          <span>Unpacking files...</span>
                        </div>
                        <div className="flex gap-3">
                          <span className="text-green-500">✔</span>
                          <span>Adding to <span className="text-orange-400">$PATH</span>...</span>
                        </div>
                        <div className="flex gap-3 pt-2 border-t border-white/10">
                          <span className="text-blue-400">➜</span>
                          <span className="text-white">envsend .env</span>
                        </div>
                        <div className="pl-6 text-green-400">
                          ✅ Secret uploaded successfully!
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              )}

              {activeTab === "mac" && (
                <div className="max-w-2xl">
                  <h2 className="text-3xl font-bold mb-6">Install on macOS</h2>
                  
                  <div className="space-y-8">
                    <div className="p-6 rounded-xl border border-white/10 bg-black/40">
                      <h3 className="text-sm font-medium text-white/60 mb-4 flex items-center gap-2">
                        <Terminal className="w-4 h-4" />
                        Homebrew (Coming Soon)
                      </h3>
                      <div className="font-mono text-sm text-white/30 select-all bg-white/5 p-4 rounded-lg">
                        brew install envsend
                      </div>
                    </div>

                    <div className="p-6 rounded-xl border border-white/10 bg-black/40">
                      <h3 className="text-sm font-medium text-white/60 mb-4">Manual Installation</h3>
                      <div className="font-mono text-sm text-white/80 space-y-2 bg-white/5 p-4 rounded-lg">
                        <p className="text-white/40"># Build from source (requires Go 1.21+)</p>
                        <p>git clone https://github.com/ranvijayysinghrathore/envsend</p>
                        <p>cd envsend/cli</p>
                        <p>sudo go build -o /usr/local/bin/envsend .</p>
                      </div>
                    </div>
                  </div>
                </div>
              )}

              {activeTab === "linux" && (
                <div className="max-w-2xl">
                  <h2 className="text-3xl font-bold mb-6">Install on Linux</h2>
                  
                  <div className="space-y-8">
                    <div className="p-6 rounded-xl border border-white/10 bg-black/40">
                      <h3 className="text-sm font-medium text-white/60 mb-4 flex items-center gap-2">
                        <Terminal className="w-4 h-4" />
                        Curl Script
                      </h3>
                      <div className="flex items-center justify-between font-mono text-sm text-white/80 bg-white/5 p-4 rounded-lg group hover:bg-white/10 transition-colors cursor-pointer"
                           onClick={() => navigator.clipboard.writeText("curl -fsSL https://envsend.io/install.sh | bash")}
                      >
                        <span>curl -fsSL https://envsend.io/install.sh | bash</span>
                        <Command className="w-4 h-4 text-white/20 group-hover:text-white transition-colors" />
                      </div>
                    </div>

                    <div className="p-6 rounded-xl border border-white/10 bg-black/40">
                      <h3 className="text-sm font-medium text-white/60 mb-4">Build from Source</h3>
                      <div className="font-mono text-sm text-white/80 space-y-2 bg-white/5 p-4 rounded-lg">
                        <p className="text-white/40"># Requires Go 1.21+</p>
                        <p>go install github.com/ranvijayysinghrathore/envsend/cli@latest</p>
                      </div>
                    </div>
                  </div>
                </div>
              )}
            </motion.div>
          </div>
        </div>

        {/* Footer Build info */}
        <div className="mt-20 pt-10 border-t border-white/5 text-center">
          <p className="text-white/30 text-sm">
            Prefer to build it yourself? 
            <a href="https://github.com/ranvijayysinghrathore/envsend" className="text-white hover:underline ml-2 transition-colors">
              View Source on GitHub
            </a>
          </p>
        </div>

      </main>
    </div>
  );
}
