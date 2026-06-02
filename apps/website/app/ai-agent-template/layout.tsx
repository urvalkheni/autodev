import type { Metadata } from "next";

export const metadata: Metadata = {
  title: "Google Gemini API AI Chatbot Boilerplate — AutoDevs",
  description:
    "Build your AI Agent products immediately using this pre-set monorepo. It features a React chatbot interface connected to an Express.js backend utilizing Google's official Gemini 2.5 Flash SDK.",
  alternates: {
    canonical: "https://autodevs.dev/ai-agent-template",
  },
};

export default function AiAgentTemplateLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return <>{children}</>;
}
