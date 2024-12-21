import sqlite3 from 'sqlite3';

export const db = new sqlite3.Database('./src/db/database.sqlite', (err: any) => {
    if (err) {
        console.log("Error occured starting database: " + err);
    }
    else {
        console.log("Database started");
    }
})
