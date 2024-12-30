import NextAuth from "next-auth";
import Google from "next-auth/providers/google";

const googleClientId = process.env.GOOGLE_CLIENT_ID;
const googleClientSecret = process.env.GOOGLE_CLIENT_SECRET;
const nextAuthSecret = process.env.NEXTAUTH_SECRET;

if (!googleClientId) {
    throw new Error("GOOGLE_CLIENT_ID is not set");
}

if (!googleClientSecret) {
    throw new Error("GOOGLE_CLIENT_SECRET is not set");
}

const handler = NextAuth({
    providers: [
        Google({
            clientId: googleClientId,
            clientSecret: googleClientSecret
        })
    ],
    callbacks: {
    },
    secret: nextAuthSecret
});

export { handler as GET, handler as POST };
