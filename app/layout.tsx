import type { Metadata } from "next";
import { GeistSans } from "geist/font/sans";
import "./globals.css";
import { Toaster } from "sonner";

export const metadata: Metadata = {
  title: {
    template: "%s • Task Management",
    default: "Task Management",
  },
  description: "Task Management Website",
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en">
      <body className={`${GeistSans.className} antialiased`}>
        <main>{children}</main>
        <Toaster richColors />
      </body>
    </html>
  );
}
