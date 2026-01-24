"use client";

import { Download, CheckCircle, Terminal, Command, Monitor } from "lucide-react";
import Link from "next/link";
import { useState } from "react";
import clsx from "clsx";

export default function InstallPage() {
  const [os, setOs] = useState<"windows" | "mac" | "linux">("windows");

  return (
    <div className="min-h-screen bg-black text-white selection:bg-cyan-500/30">
      <main className="max-w-4xl mx-auto px-6 py-20 sm:py-32">
        <div className="mb-12">
          {/* <Link href="/" className="text-gray-400 hover:text-white transition-colors mb-4 inline-block text-sm">
            ← Back to Home
          </Link> */}
          <h1 className="text-4xl sm:text-5xl font-bold tracking-tight mb-4">
            Install EnvSend
          </h1>
          <p className="text-xl text-gray-400">
            Get the production-grade CLI on your machine in seconds.
          </p>
        </div>

        {/* OS Switcher */}
        <div className="flex gap-2 p-1 bg-white/5 rounded-xl border border-white/10 w-fit mb-12">
          {(["windows", "mac", "linux"] as const).map((tab) => (
            <button
              key={tab}
              onClick={() => setOs(tab)}
              className={clsx(
                "px-6 py-2 rounded-lg text-sm font-medium transition-all capitalize",
                os === tab 
                  ? "bg-white text-black shadow-lg shadow-white/10" 
                  : "text-gray-400 hover:text-white hover:bg-white/5"
              )}
            >
              {tab}
            </button>
          ))}
        </div>

        {/* Windows Installation */}
        {os === "windows" && (
          <div className="space-y-12 animate-fade-in">
            {/* Step 1: Download */}
            <div className="relative group">
              <div className="absolute -left-4 top-0 bottom-0 w-0.5 bg-gradient-to-b from-cyan-500/50 to-transparent" />
              <div className="pl-8">
                <div className="flex items-center gap-3 mb-4">
                  <div className="flex items-center justify-center w-8 h-8 rounded-full bg-cyan-900/30 text-cyan-400 border border-cyan-500/30 font-bold shrink-0">
                    1
                  </div>
                  <h2 className="text-2xl font-semibold">Download Installer</h2>
                </div>
                <p className="text-gray-400 mb-6">
                  Get the official highly-compressed installer for Windows 10/11.
                </p>
                <a 
                  href="/download/EnvSendSetup.exe"
                  className="inline-flex items-center gap-2 px-8 py-4 rounded-xl bg-white text-black font-bold hover:bg-gray-200 transition-all active:scale-95"
                >
                  <Download className="w-5 h-5" />
                  Download EnvSendSetup.exe
                </a>
              </div>
            </div>

            {/* Step 2: Run Installer */}
            <div className="relative group">
              <div className="absolute -left-4 top-0 bottom-0 w-0.5 bg-gradient-to-b from-purple-500/50 to-transparent" />
              <div className="pl-8">
                <div className="flex items-center gap-3 mb-4">
                  <div className="flex items-center justify-center w-8 h-8 rounded-full bg-purple-900/30 text-purple-400 border border-purple-500/30 font-bold shrink-0">
                    2
                  </div>
                  <h2 className="text-2xl font-semibold">Run the Wizard</h2>
                </div>
                <p className="text-gray-400 mb-6">
                  Double-click the installer. It will guide you through the setup.
                </p>
                
                {/* Critical Checkbox Alert */}
                <div className="p-6 rounded-xl bg-yellow-500/10 border border-yellow-500/20 mb-6">
                  <div className="flex items-start gap-3">
                    <Monitor className="w-6 h-6 text-yellow-500 shrink-0 mt-1" />
                    <div>
                      <h3 className="text-yellow-400 font-semibold mb-1">Make sure to verify PATH</h3>
                      <p className="text-yellow-200/60 text-sm">
                        The installer automatically adds EnvSend to your system PATH. 
                        If prompted, ensure functionality is enabled so you can run <code className="bg-yellow-900/40 px-1 py-0.5 rounded text-yellow-200">envsend</code> from any terminal.
                      </p>
                    </div>
                  </div>
                </div>

                <div className="rounded-xl overflow-hidden border border-white/10 bg-black/40">
                    <div className="flex items-center gap-2 px-4 py-2 border-b border-white/5 bg-white/5">
                        <span className="text-xs text-gray-500">Installation Wizard</span>
                    </div>
                    <div className="p-8 flex flex-col items-center justify-center gap-4 text-center">
                         <div className="w-16 h-16 bg-gradient-to-br from-cyan-500 to-blue-600 rounded-xl mb-2 shadow-lg shadow-cyan-500/20"></div>
                         <h4 className="text-lg font-medium">Installation Complete</h4>
                         <div className="flex items-center gap-2 px-4 py-2 bg-white/5 rounded border border-white/10 w-full max-w-sm text-left">
                            <div className="w-5 h-5 rounded border border-cyan-500 bg-cyan-500/20 flex items-center justify-center text-cyan-400">
                                <CheckCircle className="w-3.5 h-3.5" />
                            </div>
                            <span className="text-sm text-gray-300">Launch EnvSend (Open Terminal)</span>
                         </div>
                    </div>
                </div>
              </div>
            </div>

            {/* Step 3: Verify */}
            <div className="relative group">
              <div className="absolute -left-4 top-0 bottom-0 w-0.5 bg-gradient-to-b from-green-500/50 to-transparent" />
              <div className="pl-8">
                <div className="flex items-center gap-3 mb-4">
                  <div className="flex items-center justify-center w-8 h-8 rounded-full bg-green-900/30 text-green-400 border border-green-500/30 font-bold shrink-0">
                    3
                  </div>
                  <h2 className="text-2xl font-semibold">Verify Installation</h2>
                </div>
                <div className="rounded-xl bg-[#0D0D0D] border border-white/10 p-4 font-mono text-sm shadow-xl">
                  <div className="flex gap-2 text-gray-400 mb-2">
                    <span>PowerShell</span>
                  </div>
                  <div className="flex gap-2">
                    <span className="text-green-400">➜</span>
                    <span className="text-white">envsend --version</span>
                  </div>
                  <div className="mt-2 text-gray-400">
                    EnvSend v1.0.0 (windows/amd64)
                  </div>
                </div>
              </div>
            </div>
          </div>
        )}

        {/* Mac/Linux Placeholder */}
        {os !== "windows" && (
          <div className="p-12 text-center rounded-2xl bg-white/5 border border-white/10 border-dashed animate-fade-in">
            <Terminal className="w-12 h-12 text-gray-600 mx-auto mb-4" />
            <h3 className="text-xl font-semibold mb-2">Coming Soon</h3>
            <p className="text-gray-400 max-w-md mx-auto">
              Native packages for macOS (Homebrew) and Linux (Keep/DEB/RPM) are being built. 
              <br />
              For now, build from source using Go 1.21+.
            </p>
          </div>
        )}
      </main>
    </div>
  );
}
