"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import { motion } from "framer-motion";
import clsx from "clsx";

const links = [
  { name: "Documentation", href: "/docs/cli" },
  { name: "Install", href: "/docs/install" },
];

export function Navbar() {
  const pathname = usePathname();

  return (
    <motion.nav 
      initial={{ y: -20, opacity: 0 }}
      animate={{ y: 0, opacity: 1 }}
      transition={{ duration: 0.5, ease: "easeOut" }}
      className="fixed top-0 left-0 right-0 z-50 h-20 flex items-center"
    >
      <div className="absolute inset-0 bg-gradient-to-b from-black/80 to-transparent pointer-events-none" />
      
      <div className="max-w-7xl w-full mx-auto px-6 flex items-center justify-between relative z-10">
        {/* Logo */}
        <Link href="/" className="flex items-center gap-1 text-lg font-bold tracking-tight text-white/90 hover:text-white transition-colors">
          <img src="/logo.png" alt="EnvSend" className="w-14 h-14 rounded-lg" />
          EnvSend
        </Link>

        {/* Center Links */}
        <div className="hidden md:flex items-center gap-8">
          {links.map((link) => (
            <Link 
              key={link.href} 
              href={link.href}
              className={clsx(
                "text-sm font-medium transition-colors hover:text-white",
                pathname === link.href ? "text-white" : "text-white/60"
              )}
            >
              {link.name}
            </Link>
          ))}
        </div>

        {/* Ghost CTA */}
        <Link 
          href="/docs/install" 
          className="hidden sm:flex h-9 items-center justify-center rounded-full bg-white/10 hover:bg-white/20 border border-white/5 hover:border-white/10 px-5 text-sm font-medium text-white transition-all backdrop-blur-sm active:scale-95"
        >
          Get Started
        </Link>
      </div>
    </motion.nav>
  );
}
