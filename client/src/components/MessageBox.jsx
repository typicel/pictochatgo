
const MessageBox = (props) => {
    return (
        <div>
            {props.messages.map((messages) => {
                return <p>{message}</p>
            })}
        </div>
    )
}

export default MessageBox;