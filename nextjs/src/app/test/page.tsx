'use client';

import { useState } from "react"
import styles from './styles.module.css'

export default function something() {

    const [message, setMessage] = useState('');

    const sendMessage = () => {
        console.log('sending ' + message);
        fetch('/api/textmessage', {
            method: "POST",
            body: JSON.stringify({
                message: message,
            })
        }
        )
    }
    return (
        <>
            <div>HI</div>
            <input onChange={(e) => setMessage(message + e.target.value)}></input>
            <button className={styles.something} onClick={sendMessage}>PRESS ME</button>
        </>
    );
}