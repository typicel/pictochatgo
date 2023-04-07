import { LitElement, html, css } from 'lit';

export class HelloWorld extends LitElement {
    static properties = {
        name: {},
    }

    static styles = css`
        :host {
            color: blue;
        }
    `;

    constructor() {
        super();
        this.name = "World"
    }


    render() {
        return html`<p>Hello, ${this.name}!</p>`;
    }
}

customElements.define('hello-world', HelloWorld);


