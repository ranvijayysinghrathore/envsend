"use client";

import { motion, useMotionTemplate, useMotionValue } from "framer-motion";
import { Shield, Lock, Share2, Timer, Key } from "lucide-react";
import { MouseEvent } from "react";

const features = [
  {
    title: "Zero Knowledge Architecture",
    description: "Your secrets are encrypted on your device before they ever touch our servers. We never see the plaintext, and we never have the keys.",
    icon: <Shield className="w-6 h-6 text-orange-400" />,
    colSpan: "md:col-span-2",
    bgPattern: "bg-grid-white/[0.02]",
    stats: "AES-256-GCM"
  },
  {
    title: "Nuclear Mode",
    description: "Shamir's Secret Sharing splits the key into multiple parts. Requires M-of-N people to unlock.",
    icon: <Key className="w-6 h-6 text-red-400" />,
    colSpan: "md:col-span-1",
    bgPattern: "bg-dot-white/[0.05]",
    stats: "M-of-N Security"
  },
  {
    title: "Ephemeral by Design",
    description: "Set strict limits on views and time. Secrets self-destruct automatically, leaving absolutely no trace left behind in our database.",
    icon: <Timer className="w-6 h-6 text-blue-400" />,
    colSpan: "md:col-span-3",
    bgPattern: "bg-grid-orange/[0.02]",
    stats: "Auto-Destruction"
  },
];

export function FeaturesBento() {
  return (
    <div className="grid grid-cols-1 md:grid-cols-3 gap-4 max-w-6xl mx-auto p-4">
      {features.map((feature, i) => (
        <SpotlightCard key={i} className={feature.colSpan}>
          <div className={`absolute inset-0 ${feature.bgPattern} [mask-image:linear-gradient(to_bottom,white,transparent)]`} />
          
          <div className="relative z-10 h-full flex flex-col justify-between p-6 md:p-8">
            <div>
              {/* <div className="w-12 h-12 rounded-xl bg-white/5 border border-white/10 flex items-center justify-center mb-4 backdrop-blur-sm">
                {feature.icon}
              </div> */}
              <h3 className="text-2xl font-bold text-white mb-2 tracking-tight">{feature.title}</h3>
              <p className="text-white/50 leading-relaxed max-w-sm">
                {feature.description}
              </p>
            </div>

            <div className="mt-8 pt-6 border-t border-white/5 flex items-center justify-between">
              <span className="text-xs font-mono text-white/30 uppercase tracking-widest">
                {feature.stats}
              </span>
              <div className="h-1.5 w-12 bg-white/10 rounded-full overflow-hidden">
                <motion.div 
                  initial={{ width: 0 }}
                  whileInView={{ width: "100%" }}
                  transition={{ duration: 1, delay: 0.5 + (i * 0.2) }}
                  className="h-full bg-white/20"
                />
              </div>
            </div>
          </div>
        </SpotlightCard>
      ))}
    </div>
  );
}

function SpotlightCard({ children, className = "" }: { children: React.ReactNode; className?: string }) {
  const mouseX = useMotionValue(0);
  const mouseY = useMotionValue(0);

  function handleMouseMove({ currentTarget, clientX, clientY }: MouseEvent) {
    const { left, top } = currentTarget.getBoundingClientRect();
    mouseX.set(clientX - left);
    mouseY.set(clientY - top);
  }

  return (
    <div
      className={`group relative border border-white/10 bg-black overflow-hidden rounded-3xl ${className}`}
      onMouseMove={handleMouseMove}
    >
      <motion.div
        className="pointer-events-none absolute -inset-px rounded-3xl opacity-0 transition duration-300 group-hover:opacity-100"
        style={{
          background: useMotionTemplate`
            radial-gradient(
              650px circle at ${mouseX}px ${mouseY}px,
              rgba(255, 85, 0, 0.15),
              transparent 80%
            )
          `,
        }}
      />
      <div className="relative h-full bg-[#050505]/90 backdrop-blur-xl">{children}</div>
    </div>
  );
}
