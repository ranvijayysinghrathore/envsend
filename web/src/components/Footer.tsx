"use client";

import { useState, useEffect } from "react";
import { motion, AnimatePresence } from "framer-motion";
import { X, Shield, Lock, FileText, Globe } from "lucide-react";

const privacyContent = {
  title: "Privacy Policy",
  content: (
    <div className="space-y-4 text-white/70 text-sm leading-relaxed">
      <p>
        At <strong>EnvSend</strong>, we believe privacy is a fundamental human right. Our architecture is designed to be Zero-Knowledge, meaning we cannot see your secrets even if we wanted to.
      </p>
      <ul className="list-disc pl-4 space-y-2">
        <li><strong>No Logs:</strong> We do not store access logs or IP addresses associated with secret retrieval.</li>
        <li><strong>Client-Side Encryption:</strong> All encryption happens in your browser using AES-256-GCM. The server never receives the encryption key.</li>
        <li><strong>Ephemeral Storage:</strong> Encrypted blobs are automatically deleted after the specified expiration time or view count.</li>
      </ul>
      <p className="pt-2">
        For a more detailed technical breakdown, please review our open-source repository.
      </p>
    </div>
  )
};

const termsContent = {
  title: "Terms of Service",
  content: (
    <div className="space-y-4 text-white/70 text-sm leading-relaxed">
      <p>
        By using <strong>EnvSend</strong>, you agree to the following terms designed to ensure the safety and utility of our service.
      </p>
      <ul className="list-disc pl-4 space-y-2">
        <li><strong>Acceptable Use:</strong> You may not use EnvSend for illegal activities, distributing malware, or hosting abusive content.</li>
        <li><strong>No Warranty:</strong> This service is provided "as is" without any warranties. While we strive for perfection, multiple layers of encryption mean you are responsible for your keys.</li>
        <li><strong>Availability:</strong> We reserve the right to limit or terminate service to protect the integrity of our infrastructure.</li>
      </ul>
      <p className="pt-2">
        Usage of the API implies consent to these terms and any future updates.
      </p>
    </div>
  )
};

export function Footer() {
  const [activeModal, setActiveModal] = useState<'privacy' | 'terms' | null>(null);

  // Close modal on Escape key
  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      if (e.key === "Escape") setActiveModal(null);
    };
    window.addEventListener("keydown", handleKeyDown);
    return () => window.removeEventListener("keydown", handleKeyDown);
  }, []);

  const openData = activeModal === 'privacy' ? privacyContent : activeModal === 'terms' ? termsContent : null;

  return (
    <>
      <footer className="py-12 border-t border-white/[0.05]">
        <div className="max-w-6xl mx-auto px-6 flex flex-col sm:flex-row justify-between items-center gap-4 text-xs text-white/30">
          <p>&copy; {new Date().getFullYear()} EnvSend. Open Source MIT.</p>
          <div className="flex gap-6">
            <button 
              onClick={() => setActiveModal('privacy')} 
              className="hover:text-white transition-colors focus:outline-none"
            >
              Privacy
            </button>
            <button 
              onClick={() => setActiveModal('terms')} 
              className="hover:text-white transition-colors focus:outline-none"
            >
              Terms
            </button>
            <a href="https://x.com/Ranvijayy_" target="_blank" rel="noopener noreferrer" className="hover:text-white transition-colors">Twitter</a>
          </div>
        </div>
      </footer>

      <AnimatePresence>
        {activeModal && openData && (
          <div className="fixed inset-0 z-[100] flex items-center justify-center px-4">
            {/* Backdrop */}
            <motion.div
              initial={{ opacity: 0 }}
              animate={{ opacity: 1 }}
              exit={{ opacity: 0 }}
              onClick={() => setActiveModal(null)}
              className="absolute inset-0 bg-black/60 backdrop-blur-sm cursor-pointer"
            />
            
            {/* Modal */}
            <motion.div
              initial={{ opacity: 0, scale: 0.9, y: 10 }}
              animate={{ opacity: 1, scale: 1, y: 0 }}
              exit={{ opacity: 0, scale: 0.9, y: 10 }}
              transition={{ type: "spring", damping: 25, stiffness: 300 }}
              className="relative w-full max-w-md bg-[#0A0A0A] border border-white/10 rounded-xl shadow-2xl overflow-hidden"
            >
              {/* Header */}
              <div className="flex items-center justify-between px-6 py-4 border-b border-white/5 bg-white/[0.02]">
                <h3 className="text-lg font-medium text-white flex items-center gap-2">
                  <Shield className="w-4 h-4 text-orange-500" />
                  {openData.title}
                </h3>
                <button
                  onClick={() => setActiveModal(null)}
                  className="p-1 rounded-full hover:bg-white/10 text-white/40 hover:text-white transition-colors"
                >
                  <X className="w-4 h-4" />
                </button>
              </div>

              {/* Content */}
              <div className="p-6 max-h-[60vh] overflow-y-auto custom-scrollbar">
                {openData.content}
              </div>

              {/* Footer Decoration */}
              <div className="h-1 w-full bg-gradient-to-r from-orange-500/20 via-orange-500/50 to-orange-500/20" />
            </motion.div>
          </div>
        )}
      </AnimatePresence>
    </>
  );
}
