import { NextRequest } from 'next/server'
import { Server } from 'socket.io'

const SocketHandler = (req: any, res: any) => {
    // Check if the WebSocket server is already running
    if (res.socket.server.io) {
        console.log('Socket is already running')
    } else {
        console.log('Socket is initializing')
        // Create a new Socket.IO server
        const io = new Server(res.socket.server)

        // Attach the Socket.IO server to the HTTP server
        res.socket.server.io = io

        // Set up event listeners for new WebSocket connections
        io.on('connection', (socket) => {
            console.log('A client connected')

            // Handle 'input-change' events from clients
            socket.on('input-change', (msg) => {
                console.log('Received input-change:', msg)
                // Broadcast the message to all other connected clients
                socket.broadcast.emit('update-input', msg)
            })

            // Handle disconnection
            socket.on('disconnect', () => {
                console.log('A client disconnected')
            })
        })
    }

    // End the response (WebSocket logic does not need to interfere with HTTP response)
    res.end()
}

export default SocketHandler
