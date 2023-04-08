import { LitElement, css, html } from "lit";

class DrawingCanvas extends LitElement {

    static properties = {
        isPainting: {},
        ctx: {},
        canvasOffsetX: {},
        canvasOffsetY: {},
        startX: {},
        startY: {},
        lineWidth: {}
    }

    static styles = css`
        #canvas {
            border: 1px solid black;
        }
    `

    constructor() {
        super()
        this.isPainting = false;
        this.ctx = canvas().getContext('2d');
        this.canvasOffsetX = canvas.offsetLeft;
        this.canvasOffsetY = canvas.offsetTop;
    }

    get canvas() {
        return this.renderRoot?.querySelector('#canvas') ?? null;
    }

    _handleStartPainting(e) {
        this.isPainting = true;
        this.startX = e.clientX;
        this.startY = e.clientY;
    }

    _handleEndPainting(e) {
        this.isPainting = false;
        this.ctx.stroke()
        this.ctx.beginPath();
    }

    _draw(e) {
        if (!this.isPainting) return;

        this.ctx.lineWidth = this.lineWidth;
        this.ctx.lineCap = 'round'

        this.ctx.lineTo(e.clientX - this.canvasOffsetX, e.clientY);



    }

    render() {

        return html`
            <canvas id="canvas"
                 @mousedown=${this._handleStartPainting}
                 @mouseup=${this._handleEndPainting}
                 @mousemove=${this._draw}
                 >
            </canvas>
        `
    }

}
customElements.define('drawing-canvas', DrawingCanvas)