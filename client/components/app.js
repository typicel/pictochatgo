import { LitElement, html, css } from 'lit';
import { Chatbox } from './chat-box'
import { DrawingCanvas } from './drawing-canvas'

export class App extends LitElement {
    static properties = {
        socket: {},
        messages: { type: Array }
    }

    static styles = css`
        display: flex;
        align-items: center;
    `;

    constructor() {
        super();
        this.messages = []
        this.socket = null;
    }

    async createRoom() {
        const options = {
            method: 'POST',
            mode: 'no-cors'
        };

        console.log(this.roomCodeInput.value);

        await fetch(`http://localhost:8080/create/${this.roomCodeInput.value}`, options)
        this.joinRoom()
    }

    joinRoom() {
        if (this.socket) return;

        socket = new WebSocket(`ws://127.0.0.1:8080/ws/${this.roomCodeInput.value}`);

        socket.onmessage = (event) => {
            let incomingMessages = event.data.split('\n');
            for (let i = 0; i < messages.length; i++) {
                messages.push(incomingMessages[i])
            }
        }

        socket.onopen = (event) => {
            console.log("WebSocket Open")
            // document.getElementById('alertBox').innerHTML = "Connected to room " + roomCode
        }

        document.getElementById('form').onsubmit = (event) => {
            console.log("sending message")
            event.preventDefault()

            if (!socket) return;

            let message = document.getElementById('message').value
            if (!message) return;

            socket.send(message)
        }

    }

    sendMessage(event) {
        event.preventDefault()

        if (!this.socket) return;

        let message = this.messageInput.value
        this.socket.send(message)
    }

    get roomCodeInput() {
        return this.renderRoot?.querySelector('#roomcode') ?? null;
    }

    get messageInput() {
        return this.renderRoot?.querySelector('#messageinput') ?? null;
    }

    render() {

        let messageBox = new Chatbox();
        messageBox.messages = this.messages
        return html`
            <header>
                <h1>Gorble Chat</h1>
            </header>
            

            <main>

                ${messageBox}

                <input type="text" id="roomcode"></input>
                <button @click=${this.createRoom}>Create Room</button>
                <button @click=${this.joinRoom}>Join Room</button>


                <form @submit=${this.sendMessage}>
                    <input type="text" id="message">
                    <button type="submit">Gorble</button>
                </form>

                ${DrawingCanvas}
            </main>
        `
    }
}

customElements.define('my-app', App);



