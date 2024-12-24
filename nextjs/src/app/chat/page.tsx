"use client";
import { useState } from "react";
import styles from './styles.module.css'

export default function Chat() {

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

    const sendMessage = async () => {
        const response = await fetch("/api/chat", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({ message: textContent }),
        })
        if (!response.ok) {
            console.error(response.statusText);
            setErrorMessage(response.statusText);
        }
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