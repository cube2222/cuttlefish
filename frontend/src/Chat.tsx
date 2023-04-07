import {CancelGeneration, GetConversation, Messages, SendMessage} from "../wailsjs/go/main/App";
import React, {useEffect, useRef, useState} from "react";
import {database} from "../wailsjs/go/models";
import {EventsOn} from "../wailsjs/runtime";
import ReactMarkdown from "react-markdown";
import {Light as SyntaxHighlighter} from "react-syntax-highlighter";
import {dracula} from "react-syntax-highlighter/dist/esm/styles/hljs";
import Message = database.Message;
import Conversation = database.Conversation;
import {MinusCircle} from "iconoir-react";

interface Props {
    conversationID: number | null;
    setConversationID: (conversationID: number) => void;
}

const Chat = ({conversationID, setConversationID}: Props) => {
    // TODO: ConversationID null if this is a new chat. Then a message should create the conversation.
    //       Remember you'll want stuff like conversation templates passed too.

    const [messages, setMessages] = useState<Array<Message>>([]);
    const [curConversation, setCurConversation] = useState<database.Conversation | null>(null);
    const [inputText, setInputText] = useState("");
    const messagesContainerRef = useRef<HTMLDivElement>(null);

    useEffect(() => {
        if (conversationID === null) {
            setMessages([]);
            return;
        }
        Messages(conversationID).then((messages) => {
            setMessages(messages);
        });
        GetConversation(conversationID).then((conversation: Conversation) => {
            setCurConversation(conversation);
        })
    }, [conversationID]);

    useEffect(() => {
        if (conversationID === null) {
            setMessages([]);
            return;
        }
        return EventsOn(`conversation-${conversationID}-updated`, (data: any) => {
            Messages(conversationID).then((messages) => {
                setMessages(messages);
            });
            GetConversation(conversationID).then((conversation: Conversation) => {
                setCurConversation(conversation);
            })
        })
    }, [conversationID]);

    useEffect(() => {
        if (messagesContainerRef.current) {
            messagesContainerRef.current.scrollTop = messagesContainerRef.current.scrollHeight;
        }
    }, [messages]); // TODO: Change dependency from messages to something like "sent message" so it doesn't happen on assistant responses.

    const handleKeyDown = async (event: React.KeyboardEvent<HTMLTextAreaElement>) => {
        if (event.key === "Enter" && !event.shiftKey) {
            event.preventDefault();
            await handleSubmit();
        }
    };

    const handleSubmit = async (debugMode: boolean = false) => {
        if (inputText.trim() !== "") {
            let message = await SendMessage(conversationID !== null ? conversationID : -1, inputText);
            setInputText("");
            setConversationID(message.conversationID);
        }
    };

    const copyToClipboard = (text: string) => {
        navigator.clipboard.writeText(text);
        alert("Code copied to clipboard!");
    };

    const renderMarkdown = (message: Message) => {
        return (
            <ReactMarkdown
                children={message.content}
                components={{
                    code({node, inline, className, children, ...props}) {
                        // TODO: Render tool use in a special way. I.e. Python should print the python code nicely.
                        const match = /language-(\w+)/.exec(className || "");
                        if (inline) {
                            return (
                                <code className={className} {...props}>
                                    {children}
                                </code>
                            );
                        }
                        return (
                            <div
                                className="relative py-2"
                                onMouseEnter={(e) => {
                                    const button = e.currentTarget.querySelector("button");
                                    if (button !== null) {
                                        button.style.opacity = "100";
                                    }
                                }}
                                onMouseLeave={(e) => {
                                    const button = e.currentTarget.querySelector("button");
                                    if (button !== null) {
                                        button.style.opacity = "0";
                                    }
                                }}
                            >
                                <button
                                    className="absolute top-0 right-0 bg-gray-700 text-white py-1 px-2 rounded-md opacity-0 hover:opacity-100 transition-opacity"
                                    onClick={() => copyToClipboard(children as string)}
                                >
                                    Copy
                                </button>
                                {match ? (
                                    <SyntaxHighlighter
                                        className="rounded-md"
                                        children={String(children).replace(/\n$/, "")}
                                        style={dracula as any}
                                        language={match[1]}
                                        PreTag="div"
                                        {...props}
                                    />
                                ) : (
                                    <SyntaxHighlighter
                                        className="rounded-md"
                                        children={String(children).replace(/\n$/, "")}
                                        style={dracula as any}
                                        PreTag="div"
                                        {...props}
                                    />
                                )}
                            </div>
                        );
                    },
                }}
                className={`${message.author == 'user' ? "bg-gray-600" : "bg-gray-700"} py-2 px-4 rounded-md ${message.author == 'user' ? "text-gray-200" : "text-gray-300"} inline-block max-w-full`}
            />
        );
    };

    return <div className="w-3/4 bg-gray-800">
        <div className="flex flex-col h-full">
            <div className="relative flex-1 overflow-hidden">
                <div
                    ref={messagesContainerRef}
                    className="flex-1 h-full overflow-y-auto px-4 py-2"
                >
                    {messages.map((message, index) => (
                        <div key={index}
                             className={`flex flex-col ${message.author == 'user' ? "items-end" : "items-start"}`}>
                            {/*<div className="text-gray-500 mb-1">*/}
                            {/*    {message.date.toLocaleString()}*/}
                            {/*</div>*/}
                            {/*TODO: Display author.*/}
                            <div
                                className={`py-2 px-4 rounded-md text-white inline-block relative max-w-full`}
                                style={{wordWrap: "break-word"}}
                            >
                                {renderMarkdown(message)}
                            </div>
                        </div>
                    ))}
                </div>
                {curConversation?.generating && <MinusCircle className="absolute left-4 bottom-1 hover:text-gray-400" onClick={async () => {
                    if (curConversation !== null) {
                        await CancelGeneration(curConversation.id)
                    }
                }}/>}
            </div>
            <form
                onSubmit={(e) => e.preventDefault()}
                className="flex flex-col h-48 px-4 py-2"
            >
                        <textarea
                            value={inputText}
                            onChange={(event) => setInputText(event.target.value)}
                            onKeyDown={handleKeyDown}
                            className="border border-gray-300 border-opacity-50 p-2 w-full h-32 bg-gray-900 text-white resize-none rounded-md"
                        />
                <div className="flex justify-end">
                    <button
                        type="button"
                        onClick={() => handleSubmit(false)}
                        className="bg-blue-500 text-white p-2 rounded-md mt-2"
                    >
                        Send
                    </button>
                </div>
            </form>
        </div>
    </div>
}

export default Chat;
