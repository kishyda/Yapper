import { getServerSession } from "next-auth";
import { useSession, signIn, signOut } from "next-auth/react"
import LogIn from "./Login";

export default async function Page() {
    // const { data: session } = useSession()
    const session = await getServerSession();
    return (
        <LogIn session={session} />
    )
}