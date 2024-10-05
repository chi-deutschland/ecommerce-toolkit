"use client"

import * as React from "react"
import { ThemeProvider as NextThemesProvider } from "next-themes"
import { useEffect } from "react"

export function ThemeProvider({ children, ...props }) {
  useEffect(() => {
    document.documentElement.classList.add('dark');
  }, []);

  return <NextThemesProvider defaultTheme="dark" {...props}>{children}</NextThemesProvider>
}