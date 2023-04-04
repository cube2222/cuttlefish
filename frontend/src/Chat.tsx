import {Messages, ResetConversation, SendMessage} from "../wailsjs/go/main/App";
import React, {useEffect, useRef, useState} from "react";
import {database} from "../wailsjs/go/models";
import Message = database.Message;
import {EventsOn} from "../wailsjs/runtime";
import ReactMarkdown from "react-markdown";
import {Light as SyntaxHighlighter} from "react-syntax-highlighter";
import {dracula} from "react-syntax-highlighter/dist/esm/styles/hljs";

interface Props {

}

const Chat = ({}: Props) => {
    const [messages, setMessages] = useState<Array<Message>>([]);
    const [inputText, setInputText] = useState("");
    const messagesContainerRef = useRef<HTMLDivElement>(null);

    useEffect(() => {
        Messages(42).then((messages) => {
            setMessages(messages);
        });
    }, []);

    useEffect(() => {
        return EventsOn("conversation-42-updated", (data: any) => {
            Messages(42).then((messages) => {
                setMessages(messages);
            });
        })
    }, []);

    useEffect(() => {
        if (messagesContainerRef.current) {
            messagesContainerRef.current.scrollTop = messagesContainerRef.current.scrollHeight;
        }
    }, [messages]);

    const handleKeyDown = (event: React.KeyboardEvent<HTMLTextAreaElement>) => {
        if (event.key === "Enter" && !event.shiftKey) {
            event.preventDefault();
            handleSubmit();
        }
    };

    const handleSubmit = (debugMode: boolean = false) => {
        if (inputText.trim() !== "") {
            SendMessage(42, inputText);
            setInputText("");
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
                className={`${message.sentBySelf ? "bg-gray-600" : "bg-gray-700"} py-2 px-4 rounded-md ${message.sentBySelf ? "text-gray-200" : "text-gray-300"} inline-block max-w-full`}
            />
        );
    };

    return <div className="w-3/4 bg-gray-800">
        <div className="flex flex-col h-full">
            <div
                ref={messagesContainerRef}
                className="flex-1 overflow-y-auto px-4 py-2"
            >
                {messages.map((message, index) => (
                    <div key={index}
                         className={`flex flex-col ${message.sentBySelf ? "items-end" : "items-start"}`}>
                        {/*<div className="text-gray-500 mb-1">*/}
                        {/*    {message.date.toLocaleString()}*/}
                        {/*</div>*/}
                        <div
                            className={`py-2 px-4 rounded-md text-white inline-block relative max-w-full`}
                            style={{wordWrap: "break-word"}}
                        >
                            {renderMarkdown(message)}
                        </div>
                    </div>
                ))}
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
                <div className="flex justify-between">
                    <button
                        type="button"
                        onClick={() => handleSubmit(false)}
                        className="bg-blue-500 text-white p-2 rounded-md mt-2"
                    >
                        Send
                    </button>
                    <button
                        type="button"
                        onClick={() => ResetConversation(42)}
                        className="bg-red-500 text-white p-2 rounded-md mt-2"
                    >
                        Reset
                    </button>
                    <button
                        type="button"
                        onClick={() => handleSubmit(true)}
                        className="bg-orange-500 text-white p-2 rounded-md mt-2"
                    >
                        Send (Other, Debug)
                    </button>
                </div>
            </form>
        </div>
    </div>
}

export default Chat;
