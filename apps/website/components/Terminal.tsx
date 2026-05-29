"use client";
import { useState, useEffect } from "react";
import { motion } from "framer-motion";

const lines = [
  { delay: 0,    text: "$ autodev",                          color: "#00FF87" },
  { delay: 600,  text: "",                                   color: "#888" },
  { delay: 800,  text: "  ⚡ AUTODEV · Clone. Scan. Install. Build.",  color: "#FFD700" },
  { delay: 1000, text: "",                                   color: "#888" },
  { delay: 1100, text: "  What do you want to install?",    color: "#F0F0F0" },
  { delay: 1300, text: "",                                   color: "#888" },
  { delay: 1400, text: "  1. Languages",                    color: "#4A90E2" },
  { delay: 1500, text: "  2. Frameworks",                   color: "#4A90E2" },
  { delay: 1600, text: "  3. Databases",                    color: "#4A90E2" },
  { delay: 1700, text: "  4. DevOps",                       color: "#4A90E2" },
  { delay: 1800, text: "  5. Mobile Development",           color: "#4A90E2" },
  { delay: 1900, text: "  6. AI / ML",                      color: "#4A90E2" },
  { delay: 2000, text: "  7. Install by Profile",           color: "#FFD700" },
  { delay: 2100, text: "",                                   color: "#888" },
  { delay: 2300, text: "> [User selects: Languages → Python, Go, Node.js]", color: "#888" },
  { delay: 2800, text: "",                                   color: "#888" },
  { delay: 2900, text: "  Installing 🐍 Python...",         color: "#FFD700" },
  { delay: 3400, text: "  ✓ Python 3.12 installed",         color: "#00FF87" },
  { delay: 3600, text: "  Installing 🔵 Go...",             color: "#FFD700" },
  { delay: 4100, text: "  ✓ Go 1.22 installed",             color: "#00FF87" },
  { delay: 4300, text: "  Installing 🟢 Node.js...",        color: "#FFD700" },
  { delay: 4800, text: "  ✓ Node.js LTS installed",         color: "#00FF87" },
  { delay: 5000, text: "",                                   color: "#888" },
  { delay: 5100, text: "  ✓ Setup complete! Run 'autodev doctor' to verify.", color: "#00FF87" },
];

export default function Terminal() {
  const [visible, setVisible] = useState(0);

  useEffect(() => {
    const timers = lines.map((l, i) =>
      setTimeout(() => setVisible(i + 1), l.delay)
    );
    return () => timers.forEach(clearTimeout);
  }, []);

  return (
    <section className="py-20 px-6 bg-[#0D0D0D] border-y-2 border-[#2A2A2A]">
      <div className="max-w-5xl mx-auto">
        <div className="text-center mb-12">
          <h2 className="text-4xl font-black text-white mb-3">
            SEE IT IN ACTION
          </h2>
          <p className="text-[#666]">Run <code className="text-[#FFD700] font-mono">autodev</code> and the interactive installer opens.</p>
        </div>

        <div className="terminal">
          {/* Title bar */}
          <div className="terminal-bar">
            <span className="terminal-dot bg-[#FF5F56]" />
            <span className="terminal-dot bg-[#FFBD2E]" />
            <span className="terminal-dot bg-[#27C93F]" />
            <span className="text-xs text-[#555] ml-3 font-mono">autodev — bash</span>
          </div>

          {/* Output */}
          <div className="px-6 py-5 space-y-1 min-h-[400px]">
            {lines.slice(0, visible).map((line, i) => (
              <motion.div
                key={i}
                initial={{ opacity: 0 }}
                animate={{ opacity: 1 }}
                transition={{ duration: 0.15 }}
                className="font-mono text-sm leading-6"
                style={{ color: line.color }}
              >
                {line.text || "\u00A0"}
              </motion.div>
            ))}
            {/* Blinking cursor */}
            {visible < lines.length && (
              <span className="inline-block w-2 h-4 bg-[#FFD700] animate-pulse" />
            )}
          </div>
        </div>
      </div>
    </section>
  );
}
