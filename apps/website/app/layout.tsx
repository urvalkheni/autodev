import type { Metadata } from "next";
import { Space_Grotesk, JetBrains_Mono } from "next/font/google";
import { ThemeProvider } from "next-themes";
import "./globals.css";

const spaceGrotesk = Space_Grotesk({
  subsets: ["latin"],
  variable: "--font-space",
  display: "swap",
});

const jetbrainsMono = JetBrains_Mono({
  subsets: ["latin"],
  variable: "--font-mono",
  display: "swap",
});

export const metadata: Metadata = {
  title: "AutoDev — The App Store for Developers",
  description:
    "Clone. Scan. Install. Build. An open-source cross-platform developer environment bootstrapper. Install any language, framework, or tool with a single command.",
  keywords: ["developer tools", "CLI", "package manager", "environment setup", "autodev"],
  icons: {
    icon: "/favicon.ico",
    shortcut: "/favicon.ico",
    apple: "/apple-touch-icon.png",
  },
  openGraph: {
    title: "AutoDev — The App Store for Developers",
    description: "Clone. Scan. Install. Build.",
    url: "https://github.com/HEETMEHTA18/autodev",
    siteName: "AutoDev",
    type: "website",
  },
  twitter: {
    card: "summary_large_image",
    title: "AutoDev — The App Store for Developers",
    description: "Clone. Scan. Install. Build.",
  },
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html
      lang="en"
      suppressHydrationWarning
      className={`${spaceGrotesk.variable} ${jetbrainsMono.variable}`}
    >
      <body className="font-space antialiased">
        <ThemeProvider
          attribute="class"
          defaultTheme="dark"
          enableSystem={false}
          disableTransitionOnChange={false}
        >
          {children}
        </ThemeProvider>
      </body>
    </html>
  );
}
