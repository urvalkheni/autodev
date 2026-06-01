import type { Metadata } from "next";
import DocsClient from "./DocsClient";

export const metadata: Metadata = {
  title: "Documentation — AutoDev",
  description:
    "Comprehensive documentation for AutoDev: learn how to scan repositories, install runtimes, automate environments, and boot roles.",
  keywords: [
    "autodev documentation",
    "developer environments",
    "package managers",
    "autodev CLI commands",
    "autodev setup",
    "autodev doctor",
  ],
  openGraph: {
    title: "Documentation — AutoDev",
    description:
      "Learn how to use AutoDev to bootstrap your developer workspaces.",
    url: "https://github.com/HEETMEHTA18/autodev/tree/main/apps/website",
    type: "article",
  },
  twitter: {
    card: "summary_large_image",
    title: "Documentation — AutoDev",
    description: "Bootstrap developer workstations with one command.",
  },
};

export default function DocsPage() {
  return <DocsClient />;
}
