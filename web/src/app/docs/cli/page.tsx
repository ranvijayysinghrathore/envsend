"use client";

import { Terminal, Shield, Lock, Clock, Share2, Key, Command, Copy, Check, ArrowRight, Eye, FileLock, Download } from "lucide-react";
import Link from "next/link";
import { useState } from "react";
import { motion } from "framer-motion";
import clsx from "clsx";

export default function CliDocsPage() {
  return (
    <div className="min-h-screen bg-black text-white selection:bg-cyan-500/30">
      
      {/* Ambient Background */}
      <div className="fixed inset-0 pointer-events-none">
        <div className="absolute top-[10%] left-[-10%] w-[40%] h-[40%] bg-cyan-500/10 blur-[120px] rounded-full mix-blend-screen opacity-30" />
        <div className="absolute bottom-[10%] right-[-10%] w-[40%] h-[40%] bg-blue-500/10 blur-[120px] rounded-full mix-blend-screen opacity-30" />
      </div>

      <main className="relative z-10 max-w-7xl mx-auto px-6 pt-32 pb-20">
        
        {/* Header */}
        <div className="max-w-3xl mb-20">
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.5 }}
          >
            <div className="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-white/5 border border-white/10 text-xs font-mono text-cyan-400 mb-6">
              <Terminal className="w-3 h-3" />
              <span>CLI Reference v1.0</span>
            </div>
            <h1 className="text-5xl sm:text-6xl font-bold tracking-tight mb-6">
              Master the <span className="text-transparent bg-clip-text bg-gradient-to-r from-cyan-400 to-blue-500">Command Line</span>
            </h1>
            <p className="text-xl text-white/50 leading-relaxed max-w-2xl">
              Simple commands for complex security. Send, receive, and manage secrets directly from your terminal with zero friction.
            </p>
          </motion.div>
        </div>

        {/* Core Workflows Grid */}
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6 mb-24">
          
          {/* Send Workflow */}
          <motion.div 
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: 0.1 }}
            className="group relative overflow-hidden rounded-3xl border border-white/10 bg-[#0A0A0A]/50 backdrop-blur-xl p-8 hover:border-white/20 transition-all"
          >
            <div className="absolute top-0 right-0 p-8 opacity-50 group-hover:opacity-100 transition-opacity">
              <div className="w-12 h-12 rounded-xl bg-orange-500/10 flex items-center justify-center text-orange-500">
                <Share2 className="w-6 h-6" />
              </div>
            </div>
            <h2 className="text-2xl font-bold mb-2">Send Secrets</h2>
            <p className="text-white/50 mb-8 max-w-sm">
              Encrypt and upload a file instantly. Returns a secure, one-time link.
            </p>
            
            <CommandBlock 
              command="envsend .env" 
              label="Basic Usage" 
              description="Sends your .env file with default settings (10m expiry, 1 view)."
            />
          </motion.div>

          {/* Receive Workflow */}
          <motion.div 
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ delay: 0.2 }}
            className="group relative overflow-hidden rounded-3xl border border-white/10 bg-[#0A0A0A]/50 backdrop-blur-xl p-8 hover:border-white/20 transition-all"
          >
            <div className="absolute top-0 right-0 p-8 opacity-50 group-hover:opacity-100 transition-opacity">
              <div className="w-12 h-12 rounded-xl bg-cyan-500/10 flex items-center justify-center text-cyan-500">
                <Download className="w-6 h-6" />
              </div>
            </div>
            <h2 className="text-2xl font-bold mb-2">Receive Secrets</h2>
            <p className="text-white/50 mb-8 max-w-sm">
              Securely download and decrypt secrets directly to your machine.
            </p>
            
            <div className="space-y-4">
              <CommandBlock 
                command='envsend "/s/19201f19-5d9b..."' 
                label="Smart Receive" 
                description="Paste the full link (in quotes) to auto-download."
              />
              <CommandBlock 
                command='envsend receive "https://envsend.io/s/..."' 
                label="Manual Receive"
                description="For scripts or when auto-detection fails. Requires the full link including the key after #." 
              />
            </div>
          </motion.div>
        </div>

        {/* Security Section (Full Width Hero) */}
        <motion.div 
          initial={{ opacity: 0, scale: 0.98 }}
          whileInView={{ opacity: 1, scale: 1 }}
          viewport={{ once: true }}
          className="relative rounded-3xl border border-green-500/20 bg-green-500/5 backdrop-blur-xl p-8 md:p-12 mb-24 overflow-hidden"
        >
          <div className="absolute -right-20 -top-20 w-96 h-96 bg-green-500/10 blur-[100px] rounded-full pointer-events-none" />
          
          <div className="grid md:grid-cols-2 gap-12 items-center relative z-10">
            <div>
              <div className="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-green-500/10 border border-green-500/20 text-green-400 text-xs font-bold uppercase tracking-wider mb-6">
                <Shield className="w-3 h-3" />
                Recommended Security
              </div>
              <h2 className="text-3xl md:text-4xl font-bold mb-6">Impossible to Breach</h2>
              <p className="text-white/60 text-lg leading-relaxed mb-8">
                For sensitive production secrets, encryption alone isn't enough. 
                Use <code className="text-green-400">--require-passphrase</code> to add a second factor. 
                Even if the link is intercepted, the data remains mathematically inaccessible without your password.
              </p>
              
              <ul className="space-y-3 mb-8">
                {[
                  "Argon2id Key Derivation",
                  "AES-256-GCM Encryption",
                  "Zero-Knowledge Architecture"
                ].map((item, i) => (
                  <li key={i} className="flex items-center gap-3 text-white/70">
                    <div className="w-5 h-5 rounded-full bg-green-500/20 flex items-center justify-center text-green-400 shrink-0">
                      <Check className="w-3 h-3" />
                    </div>
                    {item}
                  </li>
                ))}
              </ul>
            </div>

            <div className="bg-black/40 rounded-2xl border border-white/10 p-6 md:p-8 backdrop-blur-sm">
              <CommandBlock 
                command="envsend .env --require-passphrase" 
                label="Maximum Security Mode"
                variant="success"
              />
              <div className="mt-6 font-mono text-sm space-y-2 text-white/40">
                <div className="flex gap-2">
                  <span className="text-green-500">➜</span>
                  <span>Enter passphrase: ****</span>
                </div>
                <div className="flex gap-2">
                  <span className="text-green-500">➜</span>
                  <span>Confirm passphrase: ****</span>
                </div>
                <div className="flex gap-2 pt-2 text-white/60">
                  <span>🔒 Secret secured with high-entropy key</span>
                </div>
              </div>
            </div>
          </div>
        </motion.div>

        {/* Advanced Features Grid */}

        {/* Advanced Features - Expanded */}
        <div className="space-y-24 mb-20">
          
          {/* Nuclear Mode (Shamir) */}
          <section id="shamir">
            <div className="flex items-center gap-3 mb-6">
              <div className="p-2 rounded-lg bg-red-500/10 text-red-400">
                <FileLock className="w-6 h-6" />
              </div>
              <div>
                <h2 className="text-3xl font-bold">Nuclear Mode</h2>
                <p className="text-sm text-red-400 font-mono mt-1">SHAMIR'S SECRET SHARING</p>
              </div>
            </div>
            <p className="text-white/50 mb-8 max-w-2xl leading-relaxed">
              Split a secret into multiple parts (shares). To decrypt, a minimum number of shares (threshold) must be combined. Perfect for "break glass" scenarios or team consensus.
            </p>

            <div className="grid md:grid-cols-2 gap-8">
              <div className="bg-[#0A0A0A]/50 border border-white/10 rounded-2xl p-6">
                <h3 className="text-lg font-semibold text-white mb-4">How to Send</h3>
                <CommandBlock 
                  command="envsend .env --shamir-shares 5 --shamir-threshold 3" 
                  label="Create 5 Shares, Require 3 to Unlock"
                />
                <ul className="mt-4 space-y-2 text-sm text-white/40 list-disc list-inside">
                  <li>Generates 5 unique share links</li>
                  <li>Any 3 links can combine to decrypt</li>
                  <li>Ideal for distributed team access</li>
                </ul>
              </div>

              <div className="bg-[#0A0A0A]/50 border border-white/10 rounded-2xl p-6">
                <h3 className="text-lg font-semibold text-white mb-4">How to Receive</h3>
                <p className="text-sm text-white/40 mb-4">
                  Run the receive command with your share. The CLI will interactively ask for the remaining shares.
                </p>
                <div className="font-mono text-sm space-y-3 bg-black/40 p-4 rounded-lg border border-white/5 text-white/60">
                  <div className="flex gap-2">
                    <span className="text-green-500">$</span>
                    <span className="text-white">envsend "s1-..."</span>
                  </div>
                  <div className="text-yellow-400">☢️ SHAMIR SECRET SHARING DETECTED</div>
                  <div className="flex gap-2">
                    <span>Need more shares. Enter next share:</span>
                    <span className="text-white">s3-...</span>
                  </div>
                  <div className="text-green-400">✅ Decrypted successfully!</div>
                </div>
              </div>
            </div>
          </section>

          {/* SSH Encryption */}
          <section id="ssh">
            <div className="flex items-center gap-3 mb-6">
              <div className="p-2 rounded-lg bg-purple-500/10 text-purple-400">
                <Key className="w-6 h-6" />
              </div>
              <h2 className="text-3xl font-bold">Encrypted for Recipient</h2>
            </div>
            <p className="text-white/50 mb-8 max-w-2xl leading-relaxed">
              Encrypt a secret specifically for a GitHub or GitLab user using their public SSH keys. Only they can decrypt it.
            </p>

            <div className="bg-[#0A0A0A]/50 border border-white/10 rounded-2xl p-8">
              <CommandBlock 
                command="envsend .env --ssh github:username" 
                label="Target Specific User"
              />
              <div className="grid md:grid-cols-3 gap-6 mt-8">
                <div className="bg-white/5 p-4 rounded-xl">
                  <div className="text-white font-medium mb-1">1. Fetch Keys</div>
                  <div className="text-xs text-white/40">CLI automatically fetches public keys from GitHub/GitLab API.</div>
                </div>
                <div className="bg-white/5 p-4 rounded-xl">
                  <div className="text-white font-medium mb-1">2. Local Encrypt</div>
                  <div className="text-xs text-white/40">Uses X25519 (Curve25519) to encrypt data specifically for those keys.</div>
                </div>
                <div className="bg-white/5 p-4 rounded-xl">
                  <div className="text-white font-medium mb-1">3. Only They Open</div>
                  <div className="text-xs text-white/40">Recipient authenticates with their private SSH key to decrypt.</div>
                </div>
              </div>
            </div>
          </section>

          {/* Piping & Automation */}
          <section id="piping">
            <div className="flex items-center gap-3 mb-6">
              <div className="p-2 rounded-lg bg-blue-500/10 text-blue-400">
                <Terminal className="w-6 h-6" />
              </div>
              <h2 className="text-3xl font-bold">Piping & Automation</h2>
            </div>
            <p className="text-white/50 mb-8 max-w-2xl leading-relaxed">
              EnvSend reads from <code>stdin</code> and writes to `stdout`, making it a powerful citizen in your shell pipelines.
            </p>

            <div className="space-y-6">
              <CommandBlock 
                command="cat .env | grep -v SECRET_KEY | envsend" 
                label="Filter & Send" 
                description="Sanitize your env file before sending."
              />
              <CommandBlock 
                command='echo "MY_API_KEY=123" | envsend --expires 1h' 
                label="Send String" 
                description="Quickly share a single value."
              />
              <CommandBlock 
                command='envsend "https://..." > .env.production' 
                label="Save to File" 
                description="Decrypt and immediately write to a specific file."
              />
            </div>
          </section>

        </div>

      </main>
    </div>
  );
}

function CommandBlock({ command, label, description, variant = "default" }: { command: string; label?: string; description?: string, variant?: "default" | "success" }) {
  const [copied, setCopied] = useState(false);

  const copy = () => {
    navigator.clipboard.writeText(command);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  return (
    <div className="w-full">
      {label && (
        <div className={clsx(
          "text-xs font-bold uppercase tracking-wider mb-2 flex items-center gap-2",
          variant === "success" ? "text-green-400" : "text-white/40"
        )}>
          {label}
        </div>
      )}
      <div className={clsx(
        "relative flex items-center justify-between p-4 rounded-xl font-mono text-sm transition-all group",
        variant === "success" 
          ? "bg-green-500/10 border border-green-500/20 text-green-100" 
          : "bg-[#1A1A1A] border border-white/10 text-gray-300 hover:border-white/20"
      )}>
        <div className="flex gap-3 overflow-x-auto scrollbar-none">
           <span className={clsx("select-none", variant === "success" ? "text-green-400" : "text-cyan-500")}>$</span>
           <span className="whitespace-nowrap">{command}</span>
        </div>
        <button 
          onClick={copy}
          className={clsx(
            "p-2 rounded-lg transition-colors ml-4 shrink-0",
            variant === "success" ? "hover:bg-green-500/20 text-green-400" : "hover:bg-white/10 text-white/40 hover:text-white"
          )}
          title="Copy command"
        >
          {copied ? <Check className="w-4 h-4" /> : <Copy className="w-4 h-4" />}
        </button>  
      </div>
      {description && <p className="mt-2 text-xs text-white/40 leading-relaxed ml-1">{description}</p>}
    </div>
  );
}

function FeatureCard({ icon, title, command, description }: any) {
  return (
    <div className="p-6 rounded-2xl bg-white/[0.03] border border-white/5 hover:bg-white/[0.05] transition-colors">
      <div className="flex items-center gap-3 mb-4">
        <div className="p-2 rounded-lg bg-white/5">{icon}</div>
        <h4 className="font-semibold">{title}</h4>
      </div>
      <code className="block w-fit px-2 py-1 rounded bg-black/50 border border-white/10 text-xs font-mono text-cyan-400 mb-3">
        {command}
      </code>
      <p className="text-sm text-white/40 leading-relaxed">
        {description}
      </p>
    </div>
  )
}
