"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import { motion, AnimatePresence } from "framer-motion";
import clsx from "clsx";
import { Menu, X, Github } from "lucide-react";
import { useState, useEffect } from "react";

const links = [
  { name: "Documentation", href: "/docs/cli" },
  { name: "Install", href: "/docs/install" },
];

export function Navbar() {
  const pathname = usePathname();
  const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false);

  // Lock body scroll when menu is open
  useEffect(() => {
    if (isMobileMenuOpen) {
      document.body.style.overflow = "hidden";
    } else {
      document.body.style.overflow = "unset";
    }
    return () => {
      document.body.style.overflow = "unset";
    };
  }, [isMobileMenuOpen]);

  return (
    <>
      <motion.nav 
        initial={{ y: -20, opacity: 0 }}
        animate={{ y: 0, opacity: 1 }}
        transition={{ duration: 0.5, ease: "easeOut" }}
        className="fixed top-0 left-0 right-0 z-50 h-20 flex items-center"
      >
        <div className="absolute inset-0 bg-gradient-to-b from-black/80 to-transparent pointer-events-none" />
        
        <div className="max-w-7xl w-full mx-auto px-6 flex items-center justify-between relative z-10">
          {/* Logo */}
          <Link href="/" className="flex items-center gap-1 text-lg font-bold tracking-tight text-white/90 hover:text-white transition-colors z-50">
            <img src="/logo.png" alt="EnvSend" className="w-10 h-10 sm:w-14 sm:h-14 rounded-lg" />
            EnvSend
          </Link>

          {/* Desktop Links */}
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

          <div className="hidden md:flex items-center gap-4">
            <div className="relative hidden sm:block group/star">
              <a
                href="https://github.com/ranvijayysinghrathore/envsend"
                target="_blank"
                rel="noopener noreferrer"
                className="flex h-9 items-center justify-center gap-2 rounded-full border border-white/10 bg-white/5 hover:bg-white/10 px-4 text-xs font-medium text-white transition-all hover:border-white/20"
              >
                <Github className="w-4 h-4 fill-white" />
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

          {/* Mobile Menu Toggle */}
          <button 
            className="md:hidden relative z-50 p-2 text-white/70 hover:text-white transition-colors"
            onClick={() => setIsMobileMenuOpen(!isMobileMenuOpen)}
          >
            {isMobileMenuOpen ? <X className="w-6 h-6" /> : <Menu className="w-6 h-6" />}
          </button>
        </div>
      </motion.nav>

      {/* Mobile Menu Overlay */}
      <AnimatePresence>
        {isMobileMenuOpen && (
          <motion.div
            initial={{ opacity: 0, y: -20 }}
            animate={{ opacity: 1, y: 0 }}
            exit={{ opacity: 0, y: -20 }}
            transition={{ duration: 0.2 }}
            className="fixed inset-0 z-40 bg-black/95 backdrop-blur-3xl md:hidden flex flex-col items-center justify-center gap-8 pt-20"
          >
            {links.map((link) => (
              <Link 
                key={link.href} 
                href={link.href}
                onClick={() => setIsMobileMenuOpen(false)}
                className="text-2xl font-medium text-white/80 hover:text-white transition-colors"
              >
                {link.name}
              </Link>
            ))}
            
            <div className="w-16 h-px bg-white/10 my-4" />

            <a
              href="https://github.com/ranvijayysinghrathore/envsend"
              target="_blank"
              rel="noopener noreferrer"
              className="flex items-center gap-2 text-white/60 hover:text-white transition-colors"
            >
              <Github className="w-5 h-5" />
              <span>GitHub</span>
            </a>

            <Link 
              href="/docs/install"
              onClick={() => setIsMobileMenuOpen(false)}
              className="px-8 py-3 rounded-full bg-white text-black font-bold mt-4 active:scale-95 transition-transform"
            >
              Get Started
            </Link>
          </motion.div>
        )}
      </AnimatePresence>
    </>
  );
}
