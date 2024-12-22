import type { Metadata } from "next";
import localFont from "next/font/local";
import "./globals.css";
import { Session } from "inspector/promises";
import { SessionProvider } from "next-auth/react";


export default function RootLayout({
    children,
}: Readonly<{
    children: React.ReactNode;
}>) {
    return (
        <html lang="en">
            <body>
                {children}
            </body>
        </html>
    );
}
