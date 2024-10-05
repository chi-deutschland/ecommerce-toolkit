"use client";

import localFont from "next/font/local";
import "./globals.css";
import { AuthProvider } from '@/context/AuthContext';
import { ThemeProvider } from "@/components/theme-provider"
import { GlobalProvider } from '@/context/GlobalContext';

const geistSans = localFont({
  src: "./fonts/GeistVF.woff",
  variable: "--font-geist-sans",
  weight: "100 900",
});
const geistMono = localFont({
  src: "./fonts/GeistMonoVF.woff",
  variable: "--font-geist-mono",
  weight: "100 900",
});


export default function RootLayout({ children }) {
  return (
    <html lang="en">
      <body
        className={`${geistSans.variable} ${geistMono.variable} antialiased`}
      >
        <AuthProvider>
          <ThemeProvider
              attribute="class"
              defaultTheme="dark"
              enableSystem
              disableTransitionOnChange
            >
            <GlobalProvider>
              {children}
            </GlobalProvider>
          </ThemeProvider>
        </AuthProvider>
      </body>
    </html>
  );
}