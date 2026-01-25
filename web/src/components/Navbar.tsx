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

        <div className="flex items-center gap-4">
          <div className="relative hidden sm:block group/star">
            <a
              href="https://github.com/ranvijayysinghrathore/envsend"
              target="_blank"
              rel="noopener noreferrer"
              className="flex h-9 items-center justify-center gap-2 rounded-full border border-white/10 bg-white/5 hover:bg-white/10 px-4 text-xs font-medium text-white transition-all hover:border-white/20"
            >
              <svg viewBox="0 0 24 24" className="w-4 h-4 fill-white" aria-hidden="true"><path d="M12 2C6.477 2 2 6.484 2 12.017c0 4.425 2.865 8.18 6.839 9.504.5.092.682-.217.682-.483 0-.237-.008-.868-.013-1.703-2.782.605-3.369-1.343-3.369-1.343-.454-1.158-1.11-1.466-1.11-1.466-.908-.62.069-.608.069-.608 1.003.07 1.531 1.032 1.531 1.032.892 1.53 2.341 1.088 2.91.832.092-.647.35-1.088.636-1.338-2.22-.253-4.555-1.113-4.555-4.951 0-1.093.39-1.988 1.029-2.688-.103-.253-.446-1.272.098-2.65 0 0 .84-.27 2.75 1.026A9.564 9.564 0 0112 6.844c.85.004 1.705.115 2.504.337 1.909-1.296 2.747-1.027 2.747-1.027.546 1.379.202 2.398.1 2.651.64.7 1.028 1.595 1.028 2.688 0 3.848-2.339 4.695-4.566 4.943.359.309.678.92.678 1.855 0 1.338-.012 2.419-.012 2.747 0 .268.18.58.688.482A10.019 10.019 0 0022 12.017C22 6.484 17.522 2 12 2z"></path></svg>
              <span>GitHub</span>
            </a>
            
            {/* Hanging Tag */}
            <div className="absolute top-full right-3 flex flex-col items-center pointer-events-none perspective-[1000px]">
              {/* Curved String */}
              <svg width="20" height="24" viewBox="0 0 20 24" fill="none" xmlns="http://www.w3.org/2000/svg" className="text-white/30 -mt-1 -mr-2">
                <path d="M2 0 Q 3 12 18 20" stroke="currentColor" strokeWidth="1" />
              </svg>
              
              {/* Tag Body (3D Flip) */}
              <div className="relative -mt-2 -mr-6 group-hover/star:rotate-y-180 transition-transform duration-500 preserve-3d">
                {/* Front Side */}
                <div className="relative bg-orange-500 text-white text-[10px] font-bold px-2 py-0.5 rounded-sm shadow-lg rotate-[-6deg] animate-now-pulse border border-orange-400/50 backface-hidden">
                  Star on GitHub
                </div>
                
                {/* Back Side */}
                <div className="absolute inset-0 bg-blue-600 text-white text-[10px] font-bold px-2 py-0.5 rounded-sm shadow-lg rotate-[-6deg] border border-blue-400/50 backface-hidden rotate-y-180 flex items-center justify-center">
                  Contribute
                </div>
              </div>
            </div>
          </div>

          {/* Ghost CTA */}
          <Link 
            href="/docs/install" 
            className="hidden sm:flex h-9 items-center justify-center rounded-full bg-white/10 hover:bg-white/20 border border-white/5 hover:border-white/10 px-5 text-sm font-medium text-white transition-all backdrop-blur-sm active:scale-95"
          >
            Get Started
          </Link>
        </div>
      </div>
    </motion.nav>
  );
}
