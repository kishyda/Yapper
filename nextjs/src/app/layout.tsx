import "./globals.css";
import TopBar from "../components/topBar.tsx";

export default function RootLayout({
    children,
}: Readonly<{
    children: React.ReactNode;
}>) {
    return (
        <html lang="en">
            <body>
                <TopBar></TopBar>
                {children}
            </body>
        </html>
    );
}
