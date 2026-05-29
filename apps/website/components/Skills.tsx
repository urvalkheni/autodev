"use client";
import { motion } from "framer-motion";

const roadmap = [
  { from: "React + TypeScript", to: "Next.js + Docker", level: "Frontend" },
  { from: "Next.js + Docker", to: "PostgreSQL + Redis", level: "Full Stack" },
  { from: "PostgreSQL + Redis", to: "Kubernetes + Terraform", level: "DevOps" },
  { from: "Kubernetes + Terraform", to: "Cloud Architect", level: "Cloud" },
];

const currentSkills = ["React", "TypeScript", "Node.js", "Firebase"];
const nextSkills = ["Docker", "CI/CD", "PostgreSQL", "Redis"];
const longTerm = ["Kubernetes", "Terraform", "Go", "Cloud Architecture"];

export default function Skills() {
  return (
    <section id="skills" className="py-24 px-6 bg-[#0D0D0D] border-y-2 border-[#2A2A2A]">
      <div className="max-w-7xl mx-auto">
        <div className="mb-16">
          <span className="text-xs text-[#FFD700] font-bold uppercase tracking-widest">Skills.sh Integration</span>
          <h2 className="text-5xl font-black text-white mt-2 mb-4">
            YOUR LEARNING ROADMAP
          </h2>
          <p className="text-[#888] max-w-xl">
            After scanning your repos, AutoDev generates a personalized roadmap
            powered by <a href="https://skills.sh" className="text-[#FFD700] underline">skills.sh</a>.
          </p>
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-2 gap-10 items-start">
          {/* Skill columns */}
          <div className="space-y-6">
            {[
              { label: "✅ Current Skills", skills: currentSkills, color: "#00FF87" },
              { label: "🚀 Next Steps",     skills: nextSkills,    color: "#FFD700" },
              { label: "🌟 Long-Term Goals",skills: longTerm,       color: "#4A90E2" },
            ].map(({ label, skills, color }) => (
              <div key={label} className="border-2 border-[#2A2A2A] p-5">
                <h4 className="font-bold mb-3 text-sm" style={{ color }}>{label}</h4>
                <div className="flex flex-wrap gap-2">
                  {skills.map((s) => (
                    <span
                      key={s}
                      className="px-3 py-1 border-2 font-semibold text-sm"
                      style={{ borderColor: color, color }}
                    >
                      {s}
                    </span>
                  ))}
                </div>
              </div>
            ))}
          </div>

          {/* Visual roadmap */}
          <div className="border-2 border-[#2A2A2A] p-6">
            <h4 className="font-bold text-white mb-6 text-sm uppercase tracking-wider">Career Progression</h4>
            <div className="space-y-0">
              {roadmap.map((step, i) => (
                <div key={i}>
                  <motion.div
                    initial={{ opacity: 0, x: -20 }}
                    whileInView={{ opacity: 1, x: 0 }}
                    viewport={{ once: true }}
                    transition={{ delay: i * 0.15 }}
                    className="flex items-center gap-4"
                  >
                    <div
                      className="w-3 h-3 border-2 flex-shrink-0"
                      style={{ borderColor: i === 0 ? "#FFD700" : "#444", background: i === 0 ? "#FFD700" : "transparent" }}
                    />
                    <div>
                      <div className="text-xs text-[#FFD700] font-bold uppercase tracking-wider">{step.level}</div>
                      <div className="text-sm text-[#888] font-mono">{step.from}</div>
                    </div>
                  </motion.div>
                  {i < roadmap.length - 1 && (
                    <div className="ml-[5px] w-px h-8 bg-[#333] ml-[6px]" />
                  )}
                </div>
              ))}
              {/* Final destination */}
              <div className="flex items-center gap-4 mt-0">
                <div className="w-3 h-3 border-2 border-[#4A90E2] flex-shrink-0" />
                <div className="text-sm text-[#4A90E2] font-bold">Cloud Architect</div>
              </div>
            </div>

            <div className="mt-8 border-t border-[#2A2A2A] pt-5">
              <div className="terminal">
                <div className="px-4 py-3 font-mono text-xs text-[#00FF87]">
                  $ autodev skills<br />
                  <span className="text-[#555]">  Powered by skills.sh — https://skills.sh</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>
  );
}
