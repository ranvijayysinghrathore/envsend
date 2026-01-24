"use client";

import { Terminal, Shield, Lock, Clock, Share2, Key, Command, Copy, Check } from "lucide-react";
import Link from "next/link";
import { useState } from "react";
import clsx from "clsx";

export default function CliDocsPage() {
  return (
    <div className="min-h-screen bg-black text-white selection:bg-cyan-500/30">
      <main className="max-w-5xl mx-auto px-6 py-20 sm:py-32">
        <div className="mb-12">
          {/* <Link href="/" className="text-gray-400 hover:text-white transition-colors mb-4 inline-block text-sm">
            ← Back to Home
          </Link> */}
          <h1 className="text-4xl sm:text-5xl font-bold tracking-tight mb-4 text-transparent bg-clip-text bg-gradient-to-r from-cyan-400 to-blue-500">
            CLI Reference
          </h1>
          <p className="text-xl text-gray-400">
            Master the EnvSend command line interface.
          </p>
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-4 gap-12">
          {/* Sidebar Navigation */}
          <aside className="hidden lg:block space-y-8 sticky top-32 h-fit">
            <div>
              <h3 className="text-sm font-semibold text-gray-200 uppercase tracking-wider mb-4">Commands</h3>
              <ul className="space-y-3 border-l border-white/10 ml-1">
                {["send", "receive", "flags", "piping"].map((item) => (
                  <li key={item}>
                    <a href={`#${item}`} className="block pl-4 text-gray-400 hover:text-white hover:border-l-white border-l border-transparent transition-all capitalize">
                      {item}
                    </a>
                  </li>
                ))}
              </ul>
            </div>
          </aside>

          {/* Main Content */}
          <div className="lg:col-span-3 space-y-16">
            
            {/* Send Command */}
            <section id="send" className="scroll-mt-32">
              <div className="flex items-center gap-3 mb-6">
                <div className="p-2 rounded-lg bg-cyan-500/10 text-cyan-400">
                  {/* <Share2 className="w-6 h-6" /> */}
                </div>
                <h2 className="text-3xl font-bold">Send Secrets</h2>
              </div>
              <p className="text-gray-400 mb-6">
                Encrypt and upload a file or text string. The server never sees the plaintext.
              </p>
              
              <CommandBlock command="envsend .env" label="Basic Usage" />
              <div className="mt-4 text-sm text-gray-500">
                Default: Expires in 10 minutes, allows 1 view.
              </div>
            </section>

            {/* Flags */}
            <section id="flags" className="scroll-mt-32">
               <div className="flex items-center gap-3 mb-6">
                <div className="p-2 rounded-lg bg-purple-500/10 text-purple-400">
                  {/* <Command className="w-6 h-6" /> */}
                </div>
                <h2 className="text-3xl font-bold">Power Flags</h2>
              </div>
              
              <div className="space-y-8">
                <div className="p-6 rounded-xl bg-white/5 border border-white/5">
                  <h3 className="text-xl font-semibold mb-2 flex items-center gap-2">
                    <Clock className="w-5 h-5 text-yellow-400" />
                    Custom Expiry
                  </h3>
                  <CommandBlock command="envsend .env --expires 1h --max-views 5" />
                  <p className="mt-2 text-gray-400 text-sm">Set custom TTL and view limits. Max expiry is usually 7 days.</p>
                </div>

                <div className="p-6 rounded-xl bg-white/5 border border-white/5">
                  <h3 className="text-xl font-semibold mb-2 flex items-center gap-2">
                    <Key className="w-5 h-5 text-green-400" />
                    Passphrase Protection
                  </h3>
                  <CommandBlock command="envsend .env --require-passphrase" />
                  <p className="mt-2 text-gray-400 text-sm">Interactive prompt to set a password. Uses Argon2id for key derivation.</p>
                </div>

                <div className="p-6 rounded-xl bg-white/5 border border-white/5 border-l-4 border-l-red-500/50">
                   <h3 className="text-xl font-semibold mb-2 flex items-center gap-2 text-red-400">
                    <Shield className="w-5 h-5" />
                    Nuclear Mode (Shamir)
                  </h3>
                   <div className="mb-4 text-sm text-gray-300">
                    Split the decryption key into <strong>N</strong> pieces, requiring <strong>M</strong> people to combine them.
                   </div>
                  <CommandBlock command='echo "LAUNCH CODES" | envsend --shamir-shares 5 --shamir-threshold 3' />
                </div>
              </div>
            </section>

            {/* Receive Command */}
            <section id="receive" className="scroll-mt-32">
              <div className="flex items-center gap-3 mb-6">
                <div className="p-2 rounded-lg bg-green-500/10 text-green-400">
                  <Terminal className="w-6 h-6" />
                </div>
                <h2 className="text-3xl font-bold">Receive Secrets</h2>
              </div>
              <p className="text-gray-400 mb-6">
                Use "Smart Mode" to automatically detect download links.
              </p>
              
              <CommandBlock command='envsend "https://envsend.io/s/..." > .env' />
              
              <div className="mt-8 p-4 rounded-lg bg-blue-500/10 border border-blue-500/20 text-blue-200 text-sm">
                <strong>Tip:</strong> You can redirect the output to any file using standard shell operators (`&gt;`).
              </div>
            </section>
            
             {/* Pipe Support */}
            <section id="piping" className="scroll-mt-32">
              <div className="flex items-center gap-3 mb-6">
                <div className="p-2 rounded-lg bg-gray-500/10 text-gray-400">
                  <Terminal className="w-6 h-6" />
                </div>
                <h2 className="text-3xl font-bold">Pipe Support</h2>
              </div>
              <p className="text-gray-400 mb-6">
                EnvSend integrates seamlessly with other UNIX tools.
              </p>
              
              <CommandBlock command="cat .env.production | grep -v SECRET | envsend" />
            </section>

          </div>
        </div>
      </main>
    </div>
  );
}

function CommandBlock({ command, label }: { command: string; label?: string }) {
  const [copied, setCopied] = useState(false);

  const copy = () => {
    navigator.clipboard.writeText(command);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  return (
    <div className="relative group">
      {label && <div className="text-xs text-gray-500 mb-1.5 font-medium uppercase tracking-wide">{label}</div>}
      <div className="flex items-center justify-between p-4 rounded-lg bg-[#0D0D0D] border border-white/10 font-mono text-sm text-gray-300 shadow-lg hover:border-white/20 transition-colors">
        <div className="flex gap-2">
           <span className="text-green-500 hidden sm:inline">➜</span>
           <span>{command}</span>
        </div>
        <button 
          onClick={copy}
          className="p-2 rounded hover:bg-white/10 text-gray-500 hover:text-white transition-colors"
          title="Copy command"
        >
          {copied ? <Check className="w-4 h-4 text-green-400" /> : <Copy className="w-4 h-4" />}
        </button>
      </div>
    </div>
  );
}
