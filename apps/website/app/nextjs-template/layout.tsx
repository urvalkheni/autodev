import type { Metadata } from "next";

export const metadata: Metadata = {
  title: "Next.js Standalone Production Boilerplate — AutoDevs",
  description:
    "Boost your product startup using Next.js App Router, fully prepared for containerized orchestration (Docker/K8s) and continuous integrations pipeline checks.",
  alternates: {
    canonical: "https://autodevs.dev/nextjs-template",
  },
};

export default function NextJsTemplateLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return <>{children}</>;
}
