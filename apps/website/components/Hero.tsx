"use client";
import { motion } from "framer-motion";

const container = {
  hidden: {},
  show: { transition: { staggerChildren: 0.12 } },
};
const item = {
  hidden: { opacity: 0, y: 30 },
  show: { opacity: 1, y: 0, transition: { duration: 0.5 } },
};

export default function Hero() {
  return (
    <section className="pt-36 pb-24 px-6 max-w-7xl mx-auto">
      <motion.div variants={container} initial="hidden" animate="show">
        {/* Badge */}
        <motion.div variants={item} className="mb-8">
          <span className="inline-block border-2 border-[#FFD700] text-[#FFD700] text-xs font-bold px-3 py-1 uppercase tracking-widest">
            v0.1.0 — Open Source
          </span>
        </motion.div>

        {/* Headline */}
        <motion.h1
          variants={item}
          className="text-[clamp(3.5rem,10vw,8rem)] font-black leading-[0.9] tracking-tighter text-white mb-4"
        >
          THE APP STORE
          <br />
          <span className="text-[#FFD700]">FOR DEVELOPERS.</span>
        </motion.h1>

        {/* Sub-headline */}
        <motion.p variants={item} className="text-xl text-[#888] max-w-2xl mb-4 font-medium">
          Clone. Scan. Install. Build.
        </motion.p>
        <motion.p variants={item} className="text-[#666] max-w-xl mb-12 leading-relaxed">
          Install any language, framework, database, or DevOps tool with a single command.
          Smart dependency resolution. Cross-platform. Fully open-source.
        </motion.p>

        {/* CTAs */}
        <motion.div variants={item} className="flex flex-wrap gap-4 mb-16">
          <a href="#install" className="nb-btn px-8 py-4 text-lg inline-block">
            ⚡ GET STARTED
          </a>
          <a
            href="https://github.com/HEETMEHTA18/autodev"
            target="_blank"
            rel="noreferrer"
            className="nb-btn-outline px-8 py-4 text-lg inline-block"
          >
            View on GitHub →
          </a>
        </motion.div>

        {/* Quick install */}
        <motion.div variants={item}>
          <p className="text-xs text-[#555] mb-2 uppercase tracking-widest font-semibold">Quick install</p>
          <div className="terminal inline-block rounded-none">
            <div className="terminal-bar">
              <span className="terminal-dot bg-[#FF5F56]" />
              <span className="terminal-dot bg-[#FFBD2E]" />
              <span className="terminal-dot bg-[#27C93F]" />
              <span className="text-xs text-[#666] ml-2 font-mono">bash</span>
            </div>
            <div className="px-6 py-4 font-mono text-sm text-[#00FF87]">
              <span className="text-[#555]">$ </span>
              curl -fsSL https://raw.githubusercontent.com/HEETMEHTA18/autodev/main/scripts/install.sh | bash
            </div>
          </div>
        </motion.div>

        {/* Stats row */}
        <motion.div variants={item} className="flex flex-wrap gap-8 mt-16">
          {[
            { value: "40+", label: "Packages" },
            { value: "9", label: "Dev Profiles" },
            { value: "3", label: "Platforms" },
            { value: "100%", label: "Open Source" },
          ].map(({ value, label }) => (
            <div key={label} className="nb-card px-6 py-4 min-w-[120px]">
              <div className="text-3xl font-black text-[#FFD700]">{value}</div>
              <div className="text-xs text-[#666] mt-1 font-semibold uppercase tracking-wider">{label}</div>
            </div>
          ))}
        </motion.div>
      </motion.div>
    </section>
  );
}
