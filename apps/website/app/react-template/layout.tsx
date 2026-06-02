import type { Metadata } from "next";

export const metadata: Metadata = {
  title: "React + TypeScript + Tailwind CSS Starter Boilerplate — AutoDevs",
  description:
    "Get a production-grade React environment set up in 3 seconds. Pre-configured with modular folders, code styles, lint rules, Docker orchestration, and CI/CD pipelines.",
  alternates: {
    canonical: "https://autodevs.dev/react-template",
  },
};

export default function ReactTemplateLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return <>{children}</>;
}
