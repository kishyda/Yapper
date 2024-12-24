"use client";
import { useState } from "react";
import styles from './styles.module.css'
import { headers } from "next/headers";

// When sending a websocket message to the server, send the message in the following format
interface Message {
    accountID: string;
    chatID: string;
    message: string | null;
    time: string | null;
    close: boolean
}

export default function Chat() {

    // Also create queries into the URL to get the accountID, userName, and chatID
    const webSocket = new WebSocket("ws://localhost:8080/ws");

    webSocket.onopen = () => {
        console.log("Connected to server");
    }
    webSocket.onmessage = async (event) => {
        setwsMessage(await event.data)
    }

    const [textContent, setTextContent] = useState("");
    const [errorMessage, setErrorMessage] = useState("");
    const [wsMessage, setwsMessage] = useState("");

    // Change sendMessage to send to the golang server
    const sendMessage = async () => {
        webSocket.send(textContent);
    }
    return (
        <>
            <h1>Chat</h1>
            <input className={styles.textInput} type="text" onChange={(e) => setTextContent(e.target.value)} value={textContent} />
            <button onClick={sendMessage}>SEND</button>
            <div>{errorMessage}</div>
            <div>{wsMessage}</div>
        </>
    )
}