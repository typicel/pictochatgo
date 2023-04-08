import { useState } from 'preact/hooks'
import './app.css'

export function App() {
  const [roomCode, setRoomCode] = useState('')
  const [socket, setSocket] = useState(null)
  const [messages, setMessages] = useState([])
  const [messageTextBox, setMessageTextBox] = useState('')
  const [alert, setAlert] = useState('')

  const [showingroomCodeInput, setShowingRoomCodeInput] = useState(false)
  const [roomCodeInput, setRoomCodeInput] = useState('')

  const joinRoom = () => {
    if (socket) return;

    let newSocket = new WebSocket(`wss://localhost:8080/join/${roomCodeInput}`)

    newSocket.onopen = (event) => {
      console.log('Connected to server')
    }

    newSocket.onmessage = (event) => {
      console.log(event.data)
      // set the messages to allow for duplicate messages
      setMessages(messages => [...messages, event.data])
    }

    newSocket.onerror = (event) => {
      console.log('Error: ' + event)
    }


    setSocket(newSocket);
    setAlert("Room Code: " + roomCode)
    setShowingRoomCodeInput(false)
  }

  const createRoom = async () => {
    const options = {
      method: 'POST',
      mode: 'no-cors'
    };

    // Generate a random room code. 6 characters long only uppercase letters A-Z
    let newRoomCode = Math.random().toString(36).substring(2, 8).toUpperCase()
    setRoomCode(newRoomCode)

    await fetch(`http://localhost:8080/create/${newRoomCode}`, options)
    joinRoom()
  }

  const sendMessage = (event) => {
    event.preventDefault()
    if (!socket) return;

    socket.send(messageTextBox)
    setMessageTextBox('')
  }

  return (
    <>
      <h1>Screeble chat!</h1>
      <h4>{alert}</h4>

      {socket === null ?
        <>
          <button onClick={createRoom}>Create room</button>
          <button onClick={() => setShowingRoomCodeInput(true)}>Join room</button>
        </>

        :

        <>
          {messages.map(message => {
            return <p>{message}</p>
          })}

          <input type="text" value={messageTextBox} onInput={(event) => setMessageTextBox(event.target.value)} />
          <button onClick={sendMessage}>Send</button>
        </>
      }

      {showingroomCodeInput ?
        <div>
          <input type="text" value={roomCodeInput} onInput={(event) => setRoomCodeInput(event.target.value)} />
          <button onClick={joinRoom}>Join room</button>
        </div>
        : null
      }

    </>
  )
}
