"use client";
import { useState } from "react";
import { Copy, Check } from "lucide-react";

const methods = [
  { label: "NPM",      desc: "Install globally via npm",       cmd: "npm install -g @heetmehta18/autodev",                   icon: "📦" },
  { label: "NPX",      desc: "Run without local install",      cmd: "npx @heetmehta18/autodev setup",                        icon: "🚀" },
  { label: "Shell",    desc: "Linux & macOS — recommended",  cmd: "curl -fsSL https://raw.githubusercontent.com/HEETMEHTA18/autodev/main/scripts/install.sh | bash", icon: "🐚" },
  { label: "Homebrew", desc: "macOS & Linux",                  cmd: "brew install HEETMEHTA18/tap/autodev",                   icon: "🍺" },
  { label: "Scoop",    desc: "Windows — add bucket first",      cmd: "scoop bucket add autodev https://github.com/HEETMEHTA18/scoop-bucket && scoop install autodev", icon: "🪣" },
  { label: "Docker",   desc: "No local install required",      cmd: "docker run --rm -v $(pwd):/workspace ghcr.io/heetmehta18/autodev setup", icon: "🐳" },
];

export default function InstallMethods() {
  const [copied, setCopied] = useState<string | null>(null);

  const copy = (cmd: string) => {
    navigator.clipboard.writeText(cmd);
    setCopied(cmd);
    setTimeout(() => setCopied(null), 1800);
  };

  return (
    <section id="install" className="py-24 px-6 max-w-7xl mx-auto">
      <div className="mb-16 text-center">
        <span className="text-xs text-[#FFD700] font-bold uppercase tracking-widest">Install anywhere</span>
        <h2 className="text-5xl font-black text-white mt-2 mb-4">
          GET STARTED IN SECONDS
        </h2>
        <p className="text-[#888]">Pick your preferred installation method.</p>
      </div>

      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
        {methods.map((m) => (
          <div key={m.label} className="nb-card p-5 cursor-pointer group" onClick={() => copy(m.cmd)}>
            <div className="flex items-center justify-between mb-1">
              <div className="flex items-center gap-2">
                <span className="text-xl">{m.icon}</span>
                <span className="font-bold text-white text-sm">{m.label}</span>
              </div>
              <div className="flex items-center gap-1 text-xs">
                {copied === m.cmd ? (
                  <>
                    <Check className="w-3.5 h-3.5 text-[#00FF87]" />
                    <span className="text-[10px] text-[#00FF87] font-mono">Copied!</span>
                  </>
                ) : (
                  <>
                    <Copy className="w-3.5 h-3.5 text-[#555] group-hover:text-[#FFD700] transition-colors" />
                    <span className="text-[10px] text-[#555] group-hover:text-[#FFD700] font-mono transition-colors">Copy</span>
                  </>
                )}
              </div>
            </div>
            <p className="text-[10px] text-[#555] font-sans mb-3 ml-7">{m.desc}</p>
            <div
              className="font-mono text-xs text-[#00FF87] bg-[#0D0D0D] border border-[#222] px-3 py-2 truncate"
              title={m.cmd}
            >
              {m.cmd}
            </div>
          </div>
        ))}
      </div>

      {/* OS badges */}
      <div className="mt-12 flex flex-wrap gap-3 justify-center">
        {["🐧 Linux", "🍎 macOS", "🪟 Windows", "🐳 Docker", "☁️ Cloud"].map((os) => (
          <span key={os} className="border-2 border-[#2A2A2A] text-[#888] text-sm font-semibold px-4 py-2">
            {os}
          </span>
        ))}
      </div>
    </section>
  );
}
