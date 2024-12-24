"use client";

import { signIn, signOut, } from "next-auth/react";

export default function LogIn({ session }: { session: any }) {

    if (session) {
        return (
            <>
                Signed in as {session.user?.email} <br />
                <button onClick={() => signOut({ callbackUrl: "/" })}>Sign out</button>
            </>
        )
    }
    return (
        <>
            Not signed in <br />
            <button onClick={() => signIn("google")}>Sign in</button>
        </>
    )
}