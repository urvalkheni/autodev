import type { Metadata } from "next";

export const metadata: Metadata = {
  title: "Flutter Clean Architecture Template — AutoDevs",
  description:
    "Initialize cross-platform Flutter mobile, desktop, and web applications instantly, structured according to clean architecture directory conventions.",
  alternates: {
    canonical: "https://autodevs.dev/flutter-template",
  },
};

export default function FlutterTemplateLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return <>{children}</>;
}
