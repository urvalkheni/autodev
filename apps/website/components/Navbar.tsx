"use client";
import Link from "next/link";
import { motion, AnimatePresence } from "framer-motion";
import { useEffect, useState } from "react";
import { Eye, Menu, X, Star } from "lucide-react";

export default function Navbar() {
  const [stars, setStars] = useState<number | null>(null);
  const [views, setViews] = useState<number | null>(null);
  const [isMenuOpen, setIsMenuOpen] = useState(false);

  useEffect(() => {
    // Fetch GitHub stars
    fetch("https://api.github.com/repos/HEETMEHTA18/autodev")
      .then((res) => res.json())
      .then((data) => {
        if (data && typeof data.stargazers_count === "number") {
          setStars(data.stargazers_count);
        }
      })
      .catch((err) => console.error("Failed to fetch stars:", err));

    // Fetch unique page views
    const hasVisited = localStorage.getItem("autodev_visited");
    const endpoint = hasVisited 
      ? "https://api.counterapi.dev/v1/heetmehta18-autodev/views/"
      : "https://api.counterapi.dev/v1/heetmehta18-autodev/views/up";

    fetch(endpoint)
      .then((res) => res.json())
      .then((data) => {
        if (data && typeof data.count === "number") {
          setViews(data.count);
          if (!hasVisited) {
            localStorage.setItem("autodev_visited", "true");
          }
        }
      })
      .catch((err) => console.error("Failed to fetch views:", err));
  }, []);

  const navLinks = [
    { name: "Features", href: "/#features" },
    { name: "Profiles", href: "/#profiles" },
    { name: "GitHub Scanner", href: "/#github-scanner" },
    { name: "Docs", href: "/docs" }
  ];

  return (
    <motion.nav
      initial={{ y: -60, opacity: 0 }}
      animate={{ y: 0, opacity: 1 }}
      transition={{ duration: 0.4 }}
      className="fixed top-0 left-0 right-0 z-50 border-b-2 border-[#2A2A2A] bg-black/95 backdrop-blur-sm"
    >
      <div className="max-w-7xl mx-auto px-6 h-16 flex items-center justify-between gap-4">
        {/* Logo */}
        <Link href="/" className="flex items-center gap-2 shrink-0">
          <span className="text-xl md:text-2xl font-black text-[#FFD700] tracking-tighter">
            ⚡ AUTODEV
          </span>
        </Link>

        {/* Desktop Links */}
        <div className="hidden lg:flex items-center gap-8 text-sm font-semibold">
          {navLinks.map((link) => (
            <Link
              key={link.name}
              href={link.href}
              className="text-[#888] hover:text-[#FFD700] transition-colors"
            >
              {link.name}
            </Link>
          ))}
        </div>

        {/* Desktop CTA / Stats */}
        <div className="hidden md:flex items-center gap-4">
          {/* Views Button */}
          {views !== null && (
            <div className="flex items-center gap-1.5 px-3 py-1.5 border border-[#2A2A2A] bg-[#111] text-xs font-mono font-bold text-[#888]">
              <Eye className="w-3.5 h-3.5 text-[#FFD700]" />
              <span>{views.toLocaleString()} VIEWS</span>
            </div>
          )}

          {/* GitHub Stars Button */}
          <a
            href="https://github.com/HEETMEHTA18/autodev"
            target="_blank"
            rel="noreferrer"
            className="nb-btn px-4 py-2 text-sm flex items-center gap-1.5 shrink-0"
          >
            <Star className="w-4 h-4 fill-current" />
            <span>Star</span>
            {stars !== null && (
              <span className="ml-1 bg-black/20 px-1.5 py-0.5 text-xs font-mono border-l border-black/30">
                {stars}
              </span>
            )}
          </a>
        </div>

        {/* Mobile Stats & Hamburger */}
        <div className="flex lg:hidden items-center gap-2.5">
          {/* Compact views display for mobile */}
          {views !== null && (
            <div className="flex items-center gap-1 px-2 py-1 border border-[#2A2A2A] bg-[#111] text-[10px] font-mono text-[#888]">
              <Eye className="w-3.5 h-3.5 text-[#FFD700]" />
              <span>{views}</span>
            </div>
          )}

          {/* Compact stars display for mobile */}
          <a
            href="https://github.com/HEETMEHTA18/autodev"
            target="_blank"
            rel="noreferrer"
            className="nb-btn px-2.5 py-1 text-xs flex items-center gap-1 shrink-0"
          >
            <Star className="w-3 h-3 fill-current" />
            {stars !== null && <span>{stars}</span>}
          </a>

          {/* Hamburger Menu Icon */}
          <button
            onClick={() => setIsMenuOpen(!isMenuOpen)}
            className="p-1.5 border-2 border-[#2A2A2A] bg-[#111] hover:border-[#FFD700] text-white transition-colors cursor-pointer"
            aria-label="Toggle menu"
          >
            {isMenuOpen ? <X className="w-5 h-5" /> : <Menu className="w-5 h-5" />}
          </button>
        </div>
      </div>

      {/* Mobile Drawer Menu */}
      <AnimatePresence>
        {isMenuOpen && (
          <motion.div
            initial={{ opacity: 0, height: 0 }}
            animate={{ opacity: 1, height: "auto" }}
            exit={{ opacity: 0, height: 0 }}
            transition={{ duration: 0.2 }}
            className="lg:hidden border-t-2 border-[#2A2A2A] bg-[#0A0A0A] overflow-hidden"
          >
            <div className="flex flex-col p-6 gap-4 font-mono font-bold text-sm">
              {navLinks.map((link) => (
                <Link
                  key={link.name}
                  href={link.href}
                  onClick={() => setIsMenuOpen(false)}
                  className="text-neutral-400 hover:text-[#FFD700] py-2 border-b border-[#1A1A1A] transition-colors"
                >
                  {link.name}
                </Link>
              ))}
            </div>
          </motion.div>
        )}
      </AnimatePresence>
    </motion.nav>
  );
}
