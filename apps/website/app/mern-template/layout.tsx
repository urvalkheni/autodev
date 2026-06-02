import type { Metadata } from "next";

export const metadata: Metadata = {
  title: "MERN Stack Complete Docker Orchestration — AutoDevs",
  description:
    "Generate a full stack Mongo-Express-React-Node environment instantly. Fully pre-linked via docker-compose for unified local launches and simplified configurations.",
  alternates: {
    canonical: "https://autodevs.dev/mern-template",
  },
};

export default function MernTemplateLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return <>{children}</>;
}
