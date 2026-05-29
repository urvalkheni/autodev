export default function Footer() {
  return (
    <footer className="border-t-2 border-[#2A2A2A] py-16 px-6">
      <div className="max-w-7xl mx-auto">
        <div className="grid grid-cols-1 md:grid-cols-4 gap-10 mb-16">
          {/* Brand */}
          <div className="md:col-span-2">
            <div className="text-2xl font-black text-[#FFD700] mb-3">⚡ AUTODEV</div>
            <p className="text-[#666] text-sm leading-relaxed max-w-xs">
              The App Store for Developers. Open-source, cross-platform, and built for thousands of contributors.
            </p>
            <div className="flex gap-3 mt-5">
              <a href="https://github.com/HEETMEHTA18/autodev" target="_blank" rel="noreferrer"
                className="nb-btn px-4 py-2 text-xs">GitHub →</a>
            </div>
          </div>

          {/* Product */}
          <div>
            <h4 className="font-black text-white text-sm uppercase tracking-wider mb-4">Product</h4>
            <ul className="space-y-2 text-sm text-[#666]">
              {["Features", "Profiles", "GitHub Scanner", "Skills Roadmap", "Changelog"].map((l) => (
                <li key={l}><a href="#" className="hover:text-[#FFD700] transition-colors">{l}</a></li>
              ))}
            </ul>
          </div>

          {/* Open Source */}
          <div>
            <h4 className="font-black text-white text-sm uppercase tracking-wider mb-4">Open Source</h4>
            <ul className="space-y-2 text-sm text-[#666]">
              {["Contributing", "Code of Conduct", "Security", "Roadmap", "MIT License"].map((l) => (
                <li key={l}><a href="#" className="hover:text-[#FFD700] transition-colors">{l}</a></li>
              ))}
            </ul>
          </div>
        </div>

        <div className="border-t border-[#1A1A1A] pt-8 flex flex-col sm:flex-row justify-between items-center gap-4">
          <p className="text-xs text-[#444]">© 2026 AutoDev Contributors — MIT License</p>
          <p className="text-xs text-[#444] font-mono">Clone. Scan. Install. Build.</p>
        </div>
      </div>
    </footer>
  );
}
