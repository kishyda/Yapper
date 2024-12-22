import { getServerSession } from "next-auth";
import { signIn } from "next-auth/react";

export default function Login() {
    const session = getServerSession();
    if (!session) {
        return (
            <div>
                <h1>Login</h1>
                <button onClick={() => signIn("google")}>Login with Google</button>
            </div>
        );
    } else {
        return (
            <div>
                <h1>Logged in</h1>
            </div>
        );
    }
}