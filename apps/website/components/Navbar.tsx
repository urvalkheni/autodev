"use client";
import Link from "next/link";
import { motion } from "framer-motion";

export default function Navbar() {
  return (
    <motion.nav
      initial={{ y: -60, opacity: 0 }}
      animate={{ y: 0, opacity: 1 }}
      transition={{ duration: 0.4 }}
      className="fixed top-0 left-0 right-0 z-50 border-b-2 border-[#2A2A2A] bg-black/95 backdrop-blur-sm"
    >
      <div className="max-w-7xl mx-auto px-6 h-16 flex items-center justify-between">
        {/* Logo */}
        <Link href="/" className="flex items-center gap-3">
          <span className="text-2xl font-black text-[#FFD700] tracking-tighter">
            ⚡ AUTODEV
          </span>
        </Link>

        {/* Links */}
        <div className="hidden md:flex items-center gap-8 text-sm font-semibold">
          {["Features", "Profiles", "GitHub Scanner", "Docs"].map((item) => (
            <Link
              key={item}
              href={item === "Docs" ? "/docs" : `/#${item.toLowerCase().replace(" ", "-")}`}
              className="text-[#888] hover:text-[#FFD700] transition-colors"
            >
              {item}
            </Link>
          ))}
        </div>

        {/* CTA */}
        <a
          href="https://github.com/HEETMEHTA18/autodev"
          target="_blank"
          rel="noreferrer"
          className="nb-btn px-4 py-2 text-sm"
        >
          ★ Star on GitHub
        </a>
      </div>
    </motion.nav>
  );
}
