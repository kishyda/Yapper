"use client";

import { useState } from "react";
import styles from './styles.module.css'
import { getServerSession, Session } from "next-auth";

// NOTES TO RAYBO:
// You can change any and every part of the code below, but I set up a couple backend server url endpoints that I want you to follow
// 1. GET request to http://localhost:8080/ws/?userID=${userID} to connect to the chat
// 2. POST request to http://localhost:8080/createChat needs a json object with chatID
// 3. POST request to http://localhost:8080/getMessages needs a json object with chatID
// 4. POST request to http://localhost:8080/getChats needs a json object with userID
// 5. POST request to http://localhost:8080/addUser needs a json object with chatID and userID

export default function Chat({ session }: { session: Session | null }) {
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
            <button onClick={() => console.log(fetch('http://localhost:8080/getMessages'))}>Disconnect</button >
            <input className={styles.textInput} type="text" onChange={(e) => setTextContent(e.target.value)} value={textContent} />
            <button onClick={sendMessage}>SEND</button>
            <div>{errorMessage}</div>
            <div>{wsMessage}</div>
        </>
    )
}