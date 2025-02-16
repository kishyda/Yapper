'use client';

import { useState } from "react"
import styles from './styles.module.css'

let webSocket: WebSocket;

export default function Page() {

    const [userID, setUserID] = useState('');
    const [message, setMessage] = useState('');
    let timeout: NodeJS.Timeout;

    const startWebSocket = () => {
        webSocket = new WebSocket(`ws://localhost:8080/ws/?userID=${userID}`);

        webSocket.onopen = () => {
            console.log("Connected to server");
        }

        webSocket.onmessage = async (event) => {
            const str = message + String(event.data);
            console.log(str)
            setMessage(str);
        }
    }

    const closeWebSocket = () => {
        webSocket.close();
    }

    const sendMessage = () => {
        const messageObject = {
            accountID: userID,
            accountName: "JohnDoe",
            chatID: "example",
            chatName: "example",
            message: "Hello, world!",
            time: "2023-10-01T12:34:56Z",
            close: false
        };
        webSocket.send(JSON.stringify(messageObject));
    }

    const search = () => {
        clearTimeout(timeout);
        timeout = setTimeout(() => {
            console.log("searching...");
        }, 1000)
    }

    return (
        <>
            <div>
                <>WebSocket</>
                <>
                    <input placeholder="place id" className={styles.input} type="text" value={userID} onChange={(e) => setUserID(e.target.value)} />
                    <div>thing {userID}</div>
                </>
                <button className={styles.button} onClick={() => startWebSocket()}>Start web socket</button >
                <button className={styles.button} onClick={() => closeWebSocket()}>Close web socket</button >
                <button className={styles.button} onClick={() => sendMessage()}>Send message</button >
                <button className={styles.button} onClick={() => setMessage('')}>Clear</button >
                <div>{message}</div>
            </div>
            <div>
                <>some text</>
                <button className={styles.button}>something</button>
                <button className={styles.button}>else</button>
            </div>
            <div>
                <>Search functionality</>
                <input className={styles.input} type="text" placeholder="Search..." onChange={() => search()}/>
            </div>
            <div>
                <>Video and audio</>
                <input placeholder="hi"></input>
            </div>
        </>
    );
}
