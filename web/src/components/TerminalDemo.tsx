"use client";

import { useEffect, useState } from "react";

import { Terminal } from "lucide-react";
import clsx from "clsx";

const typingSequence = [
  { text: "envsend .env", delay: 1000 },
  { output: "🔒 Encrypting .env file...", delay: 500 },
  { output: "✅ Secret uploaded successfully!", delay: 800 },
  { output: "🔗 Share this link: https://envsend.io/s/abc-123#key", highlight: true, delay: 0 },
];

export function TerminalDemo() {
  const [lines, setLines] = useState<Array<{ text?: string; output?: string; highlight?: boolean }>>([]);
  const [currentTyping, setCurrentTyping] = useState("");
  const [stepIndex, setStepIndex] = useState(0);
  const [showCursor, setShowCursor] = useState(true);

  useEffect(() => {
    const interval = setInterval(() => setShowCursor((v) => !v), 500);
    return () => clearInterval(interval);
  }, []);

  useEffect(() => {
    if (stepIndex >= typingSequence.length) return;

    const step = typingSequence[stepIndex];
    let timeout: NodeJS.Timeout;

    if (step.text) {
      // Typing effect for commands
      let charIndex = 0;
      const typeChar = () => {
        if (charIndex < step.text!.length) {
          setCurrentTyping(step.text!.slice(0, charIndex + 1));
          charIndex++;
          timeout = setTimeout(typeChar, 50); // Typing speed
        } else {
          timeout = setTimeout(() => {
            setLines((prev) => [...prev, { text: step.text }]);
            setCurrentTyping("");
            setStepIndex((i) => i + 1);
          }, 300);
        }
      };
      timeout = setTimeout(typeChar, step.delay);
    } else {
      // Instant output
      timeout = setTimeout(() => {
        setLines((prev) => [...prev, step]);
        setStepIndex((i) => i + 1);
      }, step.delay);
    }

    return () => clearTimeout(timeout);
  }, [stepIndex]);

  return (
    <div className="relative rounded-xl overflow-hidden bg-[#0D0D0D] border border-white/10 shadow-2xl font-mono text-sm sm:text-base group">
      {/* Window Controls */}
      <div className="flex items-center gap-2 px-4 py-3 bg-white/5 border-b border-white/5">
        <div className="flex gap-1.5">
          <div className="w-3 h-3 rounded-full bg-red-500/20 border border-red-500/50" />
          <div className="w-3 h-3 rounded-full bg-yellow-500/20 border border-yellow-500/50" />
          <div className="w-3 h-3 rounded-full bg-green-500/20 border border-green-500/50" />
        </div>
        <div className="ml-4 text-xs text-gray-500 flex items-center gap-1.5">
          <Terminal className="w-3 h-3" />
          <span>bash — 80x24</span>
        </div>
      </div>

      {/* Terminal Content */}
      <div className="p-6 min-h-[300px] text-gray-300 relative">
        {lines.map((line, i) => (
          <div key={i} className="mb-2">
            {line.text && (
              <div className="flex gap-2 text-white">
                <span className="text-green-400">➜</span>
                <span className="text-cyan-400">~</span>
                <span>{line.text}</span>
              </div>
            )}
            {line.output && (
              <div className={clsx(
                "ml-6",
                line.highlight ? "text-cyan-300 font-semibold" : "text-gray-400"
              )}>
                {line.output}
              </div>
            )}
          </div>
        ))}

        {/* Active Line */}
        {stepIndex < typingSequence.length && typingSequence[stepIndex].text && (
          <div className="flex gap-2 text-white">
            <span className="text-green-400">➜</span>
            <span className="text-cyan-400">~</span>
            <span>
              {currentTyping}
              <span className={clsx("inline-block w-2.5 h-5 bg-gray-500 align-middle ml-1", showCursor ? "opacity-100" : "opacity-0")} />
            </span>
          </div>
        )}
      </div>
      
      {/* Glossy Overlay */}
      <div className="absolute inset-0 bg-gradient-to-br from-white/5 to-transparent pointer-events-none" />
    </div>
  );
}
