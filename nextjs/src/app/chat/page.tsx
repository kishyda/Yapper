import { getServerSession } from "next-auth";
import Chat from "./chat";

// When sending a websocket message to the server, send the message in the following format
interface Message {
    accountID: string;
    chatID: string;
    message: string | null;
    time: string | null;
    close: boolean
}

export default async function Page() {
    const session = await getServerSession();

    if (!session) {
        return <div>Unauthorized</div>
    }
    return (
        <>
            <Chat session={session} />
        </>
    )
}