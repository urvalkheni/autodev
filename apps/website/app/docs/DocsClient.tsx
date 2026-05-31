"use client";
import { useState } from "react";
import Navbar from "@/components/Navbar";
import Footer from "@/components/Footer";
import { Copy, Check } from "lucide-react";

interface DocSection {
  id: string;
  title: string;
  category: string;
  content: React.ReactNode;
}

const Callout = ({ type, title, children }: { type: "info" | "warning" | "tip"; title: string; children: React.ReactNode }) => {
  const styles = {
    info: { border: "border-l-4 border-l-[#4A90E2] border-[#2A2A2A]", bg: "bg-[#4A90E2]/5", icon: "ℹ️", color: "text-[#4A90E2]" },
    warning: { border: "border-l-4 border-l-[#FF4444] border-[#2A2A2A]", bg: "bg-[#FF4444]/5", icon: "⚠️", color: "text-[#FF4444]" },
    tip: { border: "border-l-4 border-l-[#00FF87] border-[#2A2A2A]", bg: "bg-[#00FF87]/5", icon: "⚡", color: "text-[#00FF87]" },
  };

  const currentStyle = styles[type];

  return (
    <div className={`p-4 my-6 rounded border ${currentStyle.border} ${currentStyle.bg}`}>
      <div className="flex items-center gap-2 mb-2 font-bold font-mono text-sm uppercase tracking-wider">
        <span className="text-base">{currentStyle.icon}</span>
        <span className={currentStyle.color}>{title}</span>
      </div>
      <div className="text-neutral-400 text-sm font-sans">{children}</div>
    </div>
  );
};

export default function DocsClient() {
  const [activeSection, setActiveSection] = useState("introduction");
  const [copiedId, setCopiedId] = useState<string | null>(null);

  const handleCopy = (text: string, id: string) => {
    navigator.clipboard.writeText(text);
    setCopiedId(id);
    setTimeout(() => setCopiedId(null), 2000);
  };

  const renderCodeBlock = (id: string, code: string) => (
    <div key={id} className="docs-code-block relative group border-2 border-[#2A2A2A] bg-[#050505] p-4 font-mono text-sm my-4 rounded-md overflow-x-auto pr-24">
      <button
        onClick={() => handleCopy(code, id)}
        className="docs-code-copy absolute top-2.5 right-2.5 text-xs bg-[#111] hover:bg-[#FFD700] hover:text-black hover:border-[#FFD700] border border-[#2A2A2A] px-2.5 py-1.5 rounded transition-all text-neutral-400 font-mono flex items-center gap-1.5 cursor-pointer"
      >
        {copiedId === id ? (
          <>
            <Check className="w-3.5 h-3.5 text-[#00FF87]" />
            <span className="text-[#00FF87]">Copied!</span>
          </>
        ) : (
          <>
            <Copy className="w-3.5 h-3.5" />
            <span>Copy</span>
          </>
        )}
      </button>
      <pre className="text-neutral-200 bg-transparent border-0 p-0 m-0">{code}</pre>
    </div>
  );

  const sections: DocSection[] = [
    {
      id: "introduction",
      category: "Getting Started",
      title: "Introduction",
      content: (
        <>
          <h1 className="text-4xl font-black text-white mb-6 uppercase tracking-tight">Introduction</h1>
          <p className="text-neutral-300 leading-relaxed mb-4">
            AutoDev is an open-source, cross-platform developer environment bootstrapper. It acts as an
            <strong> &quot;App Store for Developers,&quot;</strong> simplifying complex toolchain setups through intelligent, profile-based automation.
          </p>
          <p className="text-neutral-300 leading-relaxed mb-4">
            Setting up local environments often involves manually downloading compilers, installing dependencies, configuring paths, and matching versions.
            AutoDev handles this end-to-end: it scans your project files, detects the programming languages, frameworks, package managers, and databases used,
            and installs any missing pieces with a single command.
          </p>
          <Callout type="tip" title="Why use AutoDev?">
            Bring team onboarding down from hours to minutes. Run <code className="text-[#00FF87] font-mono font-bold bg-[#111] px-1.5 py-0.5 rounded">autodev setup</code> in any newly cloned repository to configure the exact toolchains required.
          </Callout>
        </>
      ),
    },
    {
      id: "quick-start",
      category: "Getting Started",
      title: "Quick Start",
      content: (
        <>
          <h1 className="text-4xl font-black text-white mb-6 uppercase tracking-tight">Quick Start</h1>
          <p className="text-neutral-300 leading-relaxed mb-4">
            Get up and running with AutoDev instantly. Run the bootstrap script for your platform, run it directly via Node.js (npx), or run it inside Docker without installation.
          </p>
          <h2 className="text-2xl font-bold text-[#FFD700] mt-8 mb-4">One-line Shell Install</h2>
          <p className="text-neutral-300 mb-2">For Linux and macOS, install using curl:</p>
          {renderCodeBlock("qs-curl", "curl -fsSL https://raw.githubusercontent.com/HEETMEHTA18/autodev/main/scripts/install.sh | bash")}

          <h2 className="text-2xl font-bold text-[#FFD700] mt-8 mb-4">Run via NPX</h2>
          <p className="text-neutral-300 mb-2">Run the interactive bootstrapper setup instantly using Node.js:</p>
          {renderCodeBlock("qs-npx", "npx @heetmehta18/autodev setup")}

          <h2 className="text-2xl font-bold text-[#FFD700] mt-8 mb-4">Interactive Setup</h2>
          <p className="text-neutral-300 mb-2">Once installed via shell/package manager, enter your project&apos;s directory and run:</p>
          {renderCodeBlock("qs-setup", "autodev setup")}

          <p className="text-neutral-300 leading-relaxed mt-4">
            AutoDev will analyze your codebase, display a list of recommended tools and libraries, and prompt you to install them.
          </p>
        </>
      ),
    },
    {
      id: "installation",
      category: "Getting Started",
      title: "Installation Methods",
      content: (
        <>
          <h1 className="text-4xl font-black text-white mb-6 uppercase tracking-tight">Installation Methods</h1>
          <p className="text-neutral-300 leading-relaxed mb-6">
            AutoDev is available via NPM, Shell scripts, Homebrew, Scoop, and Docker. Pick the method that fits your workflow.
          </p>

          <Callout type="warning" title="Important Npm Naming Notice">
            The unscoped <code className="font-mono text-neutral-300">autodev</code> package on npm is owned by an unrelated author. 
            Do <strong>not</strong> install <code className="font-mono text-neutral-300">autodev</code>. Instead, always use the scoped package: <strong><code className="font-mono text-neutral-200">@heetmehta18/autodev</code></strong>.
          </Callout>

          <div className="space-y-6">
            <div>
              <h3 className="text-lg font-bold text-white mb-1">NPM (Global Install)</h3>
              <p className="text-neutral-400 text-sm mb-2">Install the official CLI globally via npm:</p>
              {renderCodeBlock("inst-npm", "npm install -g @heetmehta18/autodev")}
            </div>

            <div>
              <h3 className="text-lg font-bold text-white mb-1">NPX (Run instantly)</h3>
              <p className="text-neutral-400 text-sm mb-2">Run the interactive setup on the fly without permanent installation:</p>
              {renderCodeBlock("inst-npx", "npx @heetmehta18/autodev setup")}
            </div>

            <div>
              <h3 className="text-lg font-bold text-white mb-1">Shell Script (Linux & macOS) — Recommended</h3>
              <p className="text-neutral-400 text-sm mb-2">Downloads and installs the pre-compiled binary directly from GitHub:</p>
              {renderCodeBlock("inst-curl", "curl -fsSL https://raw.githubusercontent.com/HEETMEHTA18/autodev/main/scripts/install.sh | bash")}
            </div>

            <div>
              <h3 className="text-lg font-bold text-white mb-1">Homebrew (macOS / Linux)</h3>
              <p className="text-neutral-400 text-sm mb-2">Install globally via Brew tap:</p>
              {renderCodeBlock("inst-brew", "brew install HEETMEHTA18/tap/autodev")}
            </div>

            <div>
              <h3 className="text-lg font-bold text-white mb-1">Scoop (Windows)</h3>
              <p className="text-neutral-400 text-sm mb-2">Add the AutoDev bucket, then install:</p>
              {renderCodeBlock("inst-scoop-bucket", "scoop bucket add autodev https://github.com/HEETMEHTA18/scoop-bucket")}
              {renderCodeBlock("inst-scoop", "scoop install autodev")}
            </div>

            <div>
              <h3 className="text-lg font-bold text-white mb-1">Docker Container</h3>
              <p className="text-neutral-400 text-sm mb-2">Run inside a container mapping your repository directory:</p>
              {renderCodeBlock("inst-docker", "docker run --rm -v $(pwd):/workspace ghcr.io/heetmehta18/autodev setup")}
            </div>
          </div>
        </>
      ),
    },

    {
      id: "cmd-doctor",
      category: "Commands",
      title: "autodev doctor",
      content: (
        <>
          <h1 className="text-4xl font-black text-white mb-6 uppercase tracking-tight">autodev doctor</h1>
          <p className="text-neutral-300 leading-relaxed mb-4">
            The <code className="text-[#FFD700] font-mono bg-[#111] px-1 py-0.5 rounded text-sm">doctor</code> command inspects your system&apos;s specifications, active package managers, and checks the installation status of managed runtimes.
          </p>
          {renderCodeBlock("cmd-doc-run", "autodev doctor")}
          <p className="text-neutral-300 leading-relaxed mb-4">
            It performs a diagnostic check of languages, framework execution environments, databases, container configurations, DevOps infrastructure, and mobile development SDKs.
          </p>
          <Callout type="info" title="Diagnostic Output">
            It lists:
            <ul className="list-disc pl-5 mt-2 space-y-1 text-neutral-400 font-sans">
              <li>System OS, Kernel, CPU Core Count, Memory Capacity</li>
              <li>Status of compilers & runtimes (e.g. Go, Python, Node, Java, Docker, Rust)</li>
              <li>Missing components, coupled with recommendations on how to install them</li>
            </ul>
          </Callout>
        </>
      ),
    },
    {
      id: "cmd-scan",
      category: "Commands",
      title: "autodev scan",
      content: (
        <>
          <h1 className="text-4xl font-black text-white mb-6 uppercase tracking-tight">autodev scan</h1>
          <p className="text-neutral-300 leading-relaxed mb-4">
            The <code className="text-[#FFD700] font-mono bg-[#111] px-1 py-0.5 rounded text-sm">scan</code> command analyzes your current working directory for configuration markers, lockfiles, and structures.
          </p>
          {renderCodeBlock("cmd-scan-run", "autodev scan")}
          <p className="text-neutral-300 leading-relaxed mb-4">
            Instead of simply looking at file extensions, AutoDev queries package registries (like <code className="font-mono text-neutral-400">package.json</code>, <code className="font-mono text-neutral-400">go.mod</code>, or <code className="font-mono text-neutral-400">requirements.txt</code>) to build a precise map of frameworks, databases, and secondary development tools.
          </p>
        </>
      ),
    },
    {
      id: "cmd-audit",
      category: "Commands",
      title: "autodev audit",
      content: (
        <>
          <h1 className="text-4xl font-black text-white mb-6 uppercase tracking-tight">autodev audit</h1>
          <p className="text-neutral-300 leading-relaxed mb-4">
            The <code className="text-[#FFD700] font-mono bg-[#111] px-1 py-0.5 rounded text-sm">audit</code> command scans your project&apos;s lockfiles and dependency lists and queries the OSV (Open Source Vulnerabilities) database.
          </p>
          {renderCodeBlock("cmd-audit-run", "autodev audit")}
          <p className="text-neutral-300 leading-relaxed mb-4">
            It checks Python, Node.js, and Go dependencies for known supply-chain vulnerabilities, compromised packages, and security advisories, listing the CVE and severity level for any threats.
          </p>
        </>
      ),
    },
    {
      id: "cmd-setup",
      category: "Commands",
      title: "autodev setup",
      content: (
        <>
          <h1 className="text-4xl font-black text-white mb-6 uppercase tracking-tight">autodev setup</h1>
          <p className="text-neutral-300 leading-relaxed mb-4">
            The <code className="text-[#FFD700] font-mono bg-[#111] px-1 py-0.5 rounded text-sm">setup</code> command scans the project and aligns your local development environment.
          </p>
          {renderCodeBlock("cmd-setup-run", "autodev setup")}
          <p className="text-neutral-300 leading-relaxed mb-4">
            It performs download caching, automated binary extraction, environment variable linking (PATH configuration), and triggers package-manager installs (like running <code className="font-mono text-neutral-400">npm install</code>, <code className="font-mono text-neutral-400">go mod download</code>, or setting up virtual environments).
          </p>
          <Callout type="warning" title="Sudo Permissions">
            Depending on your installation path (e.g. installing global binaries to <code className="font-mono text-neutral-300">/usr/local/bin</code>), AutoDev might request elevated execution permissions (sudo) to register links.
          </Callout>
        </>
      ),
    },
    {
      id: "cmd-profile",
      category: "Commands",
      title: "autodev profile",
      content: (
        <>
          <h1 className="text-4xl font-black text-white mb-6 uppercase tracking-tight">autodev profile</h1>
          <p className="text-neutral-300 leading-relaxed mb-4">
            Bootstrap a pre-defined set of developer tools based on your job profile or team role. This is useful for configuring new machines.
          </p>
          {renderCodeBlock("cmd-prof-web", "autodev profile web-dev")}
          <p className="text-neutral-300 mb-4">Supported profiles include:</p>
          <ul className="list-disc pl-5 mb-4 space-y-2 text-neutral-300 font-sans">
            <li><strong>web-dev</strong>: Node.js, pnpm, Docker, VS Code integrations</li>
            <li><strong>data-science</strong>: Python, Jupyter Notebooks, pandas, NumPy, Docker</li>
            <li><strong>devops</strong>: Terraform, kubectl, helm, Docker, AWS/GCP CLIs</li>
            <li><strong>mobile-dev</strong>: Flutter SDK, Android Studio tools, CocoaPods</li>
          </ul>
        </>
      ),
    },
    {
      id: "lockfiles",
      category: "Advanced",
      title: "Reproducible Lockfiles",
      content: (
        <>
          <h1 className="text-4xl font-black text-white mb-6 uppercase tracking-tight">Reproducible Lockfiles</h1>
          <p className="text-neutral-300 leading-relaxed mb-4">
            To lock environment tooling configurations across your engineering team, you can generate a lockfile config.
          </p>
          <p className="text-neutral-300 mb-2">Export your current environment config using:</p>
          {renderCodeBlock("adv-export", "autodev export")}
          <p className="text-neutral-300 leading-relaxed mt-4">
            This creates a <code className="font-mono text-[#FFD700] bg-[#111] px-1 py-0.5 rounded">.autodev.yaml</code> file at your project root. When other developers run <code className="font-mono text-neutral-300">autodev setup</code>, the CLI reads this file to guarantee identical compiler and tool versions across everyone&apos;s workstation.
          </p>
        </>
      ),
    },
    {
      id: "docker-usage",
      category: "Advanced",
      title: "Docker Integration",
      content: (
        <>
          <h1 className="text-4xl font-black text-white mb-6 uppercase tracking-tight">Running in Docker</h1>
          <p className="text-neutral-300 leading-relaxed mb-4">
            If you want to keep your host machine clean, you can run AutoDev scanner inside a Docker sandbox.
          </p>
          {renderCodeBlock("adv-dock", "docker run --rm -v $(pwd):/workspace ghcr.io/heetmehta18/autodev scan")}
          <p className="text-neutral-300 leading-relaxed mt-4">
            This will mount your local repository into the container workspace, run the dependency scanners, and print the resulting plan without altering files on your host OS.
          </p>
        </>
      ),
    },
  ];

  // Group sections by category
  const categories = ["Getting Started", "Commands", "Advanced"];

  const currentSection = sections.find((s) => s.id === activeSection) || sections[0];

  return (
    <main className="min-h-screen bg-black flex flex-col">
      <Navbar />

      {/* Main Documentation Wrapper */}
      <div className="flex-1 flex pt-16 max-w-7xl mx-auto w-full px-6">
        
        {/* Left Sidebar (Vercel Style) */}
        <aside className="w-64 shrink-0 border-r border-[#1F1F1F] pr-8 pt-10 hidden md:block sticky top-16 h-[calc(100vh-4rem)] overflow-y-auto">
          <div className="space-y-8">
            {categories.map((cat) => (
              <div key={cat}>
                <h4 className="text-xs font-bold text-neutral-500 uppercase tracking-widest mb-3 font-mono">
                  {cat}
                </h4>
                <ul className="space-y-1.5">
                  {sections
                    .filter((s) => s.category === cat)
                    .map((s) => {
                      const isActive = s.id === activeSection;
                      return (
                        <li key={s.id}>
                          <button
                            onClick={() => setActiveSection(s.id)}
                            className={`w-full text-left px-3 py-1.5 rounded text-sm transition-all font-mono font-medium ${
                              isActive
                                ? "bg-[#111] text-[#FFD700] border-l-2 border-[#FFD700] pl-4"
                                : "text-neutral-400 hover:text-white hover:bg-neutral-900/50"
                            }`}
                          >
                            {s.title}
                          </button>
                        </li>
                      );
                    })}
                </ul>
              </div>
            ))}
          </div>
        </aside>

        {/* Right Content Area */}
        <article className="flex-1 min-w-0 pl-0 md:pl-12 pt-10 pb-20">
          <div className="max-w-3xl">
            
            {/* Mobile Navigation Selector */}
            <div className="md:hidden mb-8 border-2 border-[#2A2A2A] bg-[#111] p-3 rounded flex flex-col gap-2">
              <label className="text-xs font-bold text-[#FFD700] uppercase font-mono tracking-wider">
                Select Page:
              </label>
              <select
                value={activeSection}
                onChange={(e) => setActiveSection(e.target.value)}
                className="bg-black text-white border border-[#2A2A2A] p-2 font-mono text-sm outline-none rounded"
              >
                {sections.map((s) => (
                  <option key={s.id} value={s.id}>
                    {s.category} → {s.title}
                  </option>
                ))}
              </select>
            </div>

            {/* Render Active Content */}
            <div className="prose prose-invert prose-yellow max-w-none">
              {currentSection.content}
            </div>

            {/* Navigation Footer Buttons */}
            <div className="mt-16 pt-8 border-t border-[#1F1F1F] flex justify-between gap-4">
              {(() => {
                const curIdx = sections.findIndex((s) => s.id === activeSection);
                const prevSec = curIdx > 0 ? sections[curIdx - 1] : null;
                const nextSec = curIdx < sections.length - 1 ? sections[curIdx + 1] : null;

                return (
                  <>
                    {prevSec ? (
                      <button
                        onClick={() => setActiveSection(prevSec.id)}
                        className="nb-btn-outline px-4 py-2.5 text-xs text-left flex flex-col font-mono"
                      >
                        <span className="text-neutral-500 font-sans uppercase font-bold text-[10px] tracking-widest">Previous</span>
                        <span>← {prevSec.title}</span>
                      </button>
                    ) : (
                      <div />
                    )}

                    {nextSec ? (
                      <button
                        onClick={() => setActiveSection(nextSec.id)}
                        className="nb-btn px-4 py-2.5 text-xs text-right flex flex-col font-mono"
                      >
                        <span className="text-black/50 font-sans uppercase font-bold text-[10px] tracking-widest">Next</span>
                        <span>{nextSec.title} →</span>
                      </button>
                    ) : (
                      <div />
                    )}
                  </>
                );
              })()}
            </div>

          </div>
        </article>

      </div>

      <Footer />
    </main>
  );
}
