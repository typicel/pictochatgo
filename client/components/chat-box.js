import { LitElement, html } from "lit";


export class Chatbox extends LitElement {
    static properties = {
        messages: {}
    }
    constructor() {
        super();
    }

    render() {
        return html`
        <ul>
            ${this.messages.map(
            (item) => html`
                    <li>${item}</li>
                `
        )}
        </ul>
        `
    }
}

customElements.define('chat-box', Chatbox);