"use client";
import { useState } from "react";
import { motion, AnimatePresence } from "framer-motion";

const categories = [
  {
    id: "languages", icon: "🔤", label: "Languages",
    packages: [
      { id: "python",  icon: "🐍", name: "Python",    desc: "3.12 LTS" },
      { id: "nodejs",  icon: "🟢", name: "Node.js",   desc: "LTS" },
      { id: "go",      icon: "🔵", name: "Go",        desc: "1.22" },
      { id: "rust",    icon: "🦀", name: "Rust",      desc: "stable" },
      { id: "java",    icon: "☕", name: "Java",      desc: "OpenJDK 21" },
      { id: "kotlin",  icon: "🎯", name: "Kotlin",    desc: "2.0" },
      { id: "php",     icon: "🐘", name: "PHP",       desc: "8.3" },
      { id: "ruby",    icon: "💎", name: "Ruby",      desc: "3.3" },
      { id: "dart",    icon: "🎯", name: "Dart",      desc: "3.x" },
      { id: "swift",   icon: "🦉", name: "Swift",     desc: "5.10" },
    ],
  },
  {
    id: "frameworks", icon: "🧩", label: "Frameworks",
    packages: [
      { id: "react",      icon: "⚛️",  name: "React",       desc: "18+" },
      { id: "nextjs",     icon: "▲",   name: "Next.js",     desc: "15" },
      { id: "vue",        icon: "💚",  name: "Vue",         desc: "3.x" },
      { id: "angular",    icon: "🔴",  name: "Angular",     desc: "18" },
      { id: "svelte",     icon: "🔥",  name: "Svelte",      desc: "5.x" },
      { id: "django",     icon: "🎸",  name: "Django",      desc: "5.x" },
      { id: "fastapi",    icon: "⚡",  name: "FastAPI",     desc: "latest" },
      { id: "laravel",    icon: "🎻",  name: "Laravel",     desc: "11" },
      { id: "springboot", icon: "🌱",  name: "Spring Boot", desc: "3.x" },
      { id: "flutter",    icon: "💙",  name: "Flutter",     desc: "3.x" },
    ],
  },
  {
    id: "databases", icon: "🗄️", label: "Databases",
    packages: [
      { id: "postgresql", icon: "🐘", name: "PostgreSQL", desc: "16" },
      { id: "mysql",      icon: "🐬", name: "MySQL",      desc: "8.x" },
      { id: "mongodb",    icon: "🍃", name: "MongoDB",    desc: "7.x" },
      { id: "redis",      icon: "🔴", name: "Redis",      desc: "7.x" },
      { id: "sqlite",     icon: "📦", name: "SQLite",     desc: "3.x" },
    ],
  },
  {
    id: "devops", icon: "⚙️", label: "DevOps",
    packages: [
      { id: "docker",     icon: "🐳",  name: "Docker",     desc: "latest" },
      { id: "k8s",        icon: "☸️",  name: "kubectl",    desc: "stable" },
      { id: "terraform",  icon: "🏗️",  name: "Terraform",  desc: "1.x" },
      { id: "nginx",      icon: "🌐",  name: "Nginx",      desc: "latest" },
      { id: "git",        icon: "📝",  name: "Git",        desc: "latest" },
    ],
  },
  {
    id: "mobile", icon: "📱", label: "Mobile",
    packages: [
      { id: "flutter",       icon: "💙",  name: "Flutter",        desc: "+ deps" },
      { id: "android-sdk",   icon: "🤖",  name: "Android SDK",    desc: "latest" },
      { id: "android-studio",icon: "🖥️",  name: "Android Studio", desc: "latest" },
      { id: "xcode",         icon: "🍎",  name: "Xcode CLI",      desc: "macOS" },
    ],
  },
  {
    id: "aiml", icon: "🤖", label: "AI / ML",
    packages: [
      { id: "pytorch",      icon: "🔥",  name: "PyTorch",      desc: "2.x" },
      { id: "tensorflow",   icon: "🧠",  name: "TensorFlow",   desc: "2.x" },
      { id: "jupyter",      icon: "📓",  name: "Jupyter",      desc: "latest" },
      { id: "ollama",       icon: "🦙",  name: "Ollama",       desc: "latest" },
      { id: "langchain",    icon: "🔗",  name: "LangChain",    desc: "latest" },
      { id: "huggingface",  icon: "🤗",  name: "HuggingFace",  desc: "latest" },
      { id: "opencv",       icon: "👁️",  name: "OpenCV",       desc: "4.x" },
    ],
  },
];

export default function Categories() {
  const [active, setActive] = useState(categories[0].id);
  const [selected, setSelected] = useState<Set<string>>(new Set());

  const current = categories.find((c) => c.id === active)!;

  const toggle = (id: string) => {
    const next = new Set(selected);
    next.has(id) ? next.delete(id) : next.add(id);
    setSelected(next);
  };

  return (
    <section id="features" className="py-24 px-6 max-w-7xl mx-auto">
      <div className="mb-16">
        <span className="text-xs text-[#FFD700] font-bold uppercase tracking-widest">Interactive Demo</span>
        <h2 className="text-5xl font-black text-white mt-2 mb-4">
          BROWSE & INSTALL
        </h2>
        <p className="text-[#888] max-w-xl">
          Select categories and packages. AutoDev resolves dependencies automatically.
        </p>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-[280px_1fr] gap-6">
        {/* Sidebar */}
        <div className="space-y-2">
          {categories.map((cat) => (
            <button
              key={cat.id}
              onClick={() => setActive(cat.id)}
              className={`w-full text-left px-4 py-3 font-bold text-sm flex items-center gap-3 border-2 transition-all duration-100
                ${active === cat.id
                  ? "border-[#FFD700] text-[#FFD700] bg-[#FFD70010] translate-x-1"
                  : "border-[#2A2A2A] text-[#888] hover:border-[#444] hover:text-white"
                }`}
            >
              <span>{cat.icon}</span>
              {cat.label}
            </button>
          ))}
        </div>

        {/* Package Grid */}
        <div className="border-2 border-[#2A2A2A] p-6">
          <AnimatePresence mode="wait">
            <motion.div
              key={active}
              initial={{ opacity: 0, x: 10 }}
              animate={{ opacity: 1, x: 0 }}
              exit={{ opacity: 0, x: -10 }}
              transition={{ duration: 0.15 }}
            >
              <div className="flex items-center justify-between mb-6">
                <h3 className="text-xl font-black text-white">{current.label}</h3>
                <div className="flex gap-2">
                  <button
                    onClick={() => {
                      const next = new Set(selected);
                      current.packages.forEach((p) => next.add(p.id));
                      setSelected(next);
                    }}
                    className="text-xs px-3 py-1 border border-[#FFD700] text-[#FFD700] font-bold hover:bg-[#FFD700] hover:text-black transition-colors"
                  >
                    Select All
                  </button>
                  <button
                    onClick={() => {
                      const next = new Set(selected);
                      current.packages.forEach((p) => next.delete(p.id));
                      setSelected(next);
                    }}
                    className="text-xs px-3 py-1 border border-[#444] text-[#888] font-bold hover:border-white hover:text-white transition-colors"
                  >
                    Clear
                  </button>
                </div>
              </div>

              <div className="grid grid-cols-2 sm:grid-cols-3 xl:grid-cols-4 gap-3">
                {current.packages.map((pkg) => {
                  const isSelected = selected.has(pkg.id);
                  return (
                    <motion.button
                      key={pkg.id}
                      onClick={() => toggle(pkg.id)}
                      whileTap={{ scale: 0.96 }}
                      className={`p-4 text-left border-2 font-semibold transition-all duration-100
                        ${isSelected
                          ? "border-[#FFD700] bg-[#FFD70015] shadow-[4px_4px_0_#FFD700]"
                          : "border-[#2A2A2A] hover:border-[#444]"
                        }`}
                    >
                      <div className="text-2xl mb-2">{pkg.icon}</div>
                      <div className={`text-sm font-bold ${isSelected ? "text-[#FFD700]" : "text-white"}`}>
                        {pkg.name}
                      </div>
                      <div className="text-xs text-[#666] mt-0.5">{pkg.desc}</div>
                      {isSelected && (
                        <div className="mt-2 text-xs text-[#00FF87] font-bold">✓ Selected</div>
                      )}
                    </motion.button>
                  );
                })}
              </div>
            </motion.div>
          </AnimatePresence>

          {/* Install bar */}
          {selected.size > 0 && (
            <motion.div
              initial={{ opacity: 0, y: 10 }}
              animate={{ opacity: 1, y: 0 }}
              className="mt-6 border-t-2 border-[#2A2A2A] pt-5 flex items-center justify-between"
            >
              <span className="text-sm text-[#888]">
                <span className="text-[#FFD700] font-bold">{selected.size}</span> packages selected
              </span>
              <div className="terminal inline-block">
                <div className="px-4 py-2 font-mono text-xs text-[#00FF87]">
                  $ autodev install {[...selected].join(" ")}
                </div>
              </div>
            </motion.div>
          )}
        </div>
      </div>
    </section>
  );
}
