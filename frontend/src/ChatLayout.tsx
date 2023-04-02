import React, {useState, useRef, useEffect} from "react";
import ReactMarkdown from "react-markdown";
import {Light as SyntaxHighlighter} from "react-syntax-highlighter";
import {dracula} from "react-syntax-highlighter/dist/esm/styles/hljs";
import "highlight.js/styles/github-dark-dimmed.css";
import { Settings } from 'iconoir-react';

type ChatMessage = {
    message: string;
    isSentByMe: boolean;
};

const ChatLayout = () => {
    const [messages, setMessages] = useState<Array<ChatMessage>>([]);
    const [inputText, setInputText] = useState("");
    const messagesContainerRef = useRef<HTMLDivElement>(null);

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
            setMessages([
                ...messages,
                {message: inputText, isSentByMe: !debugMode},
            ]);
            setInputText("");
        }
    };

    const conversations = [
        {
            id: 1,
            message: "Hey, how's it going?",
            timestamp: "10:23 AM",
        },
        {
            id: 2,
            message: "Can you send me the report?",
            timestamp: "Yesterday",
        },
        {
            id: 3,
            message: "I'll be late to the meeting",
            timestamp: "Tuesday",
        },
    ];

    const copyToClipboard = (text: string) => {
        navigator.clipboard.writeText(text);
        alert("Code copied to clipboard!");
    };

    const renderMarkdown = (message: ChatMessage) => {
        return (
            <ReactMarkdown
                children={message.message}
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
                                className="relative"
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
                className={`bg-${message.isSentByMe ? "green" : "gray"}-700 py-2 px-4 rounded-md text-white inline-block max-w-full`}
            />
        );
    };

    return (
        <div className="flex h-screen overflow-hidden">
            <div className="w-1/4 border-r border-gray-300 border-opacity-50 p-4 bg-gray-900">
                <h2 className="font-bold text-lg mb-4">Conversations</h2>
                {conversations.map((conversation, index) => (
                    <div
                        key={conversation.id}
                        className="flex items-center mb-4 cursor-pointer border-b border-gray-400 py-2"
                    >
                        <div className="w-10 h-10 rounded-full bg-gray-300 mr-2"></div>
                        <div className="flex-1 text-gray-500">
                            <div className="flex justify-between">
                                <p className="text-sm">{conversation.timestamp}</p>
                            </div>
                            <p className="text-gray-500">{conversation.message}</p>
                        </div>
                    </div>
                ))}
                <div className="absolute bottom-4 left-4 cursor-pointer">
                    <Settings className="text-gray-500"/>
                </div>
            </div>
            <div className="w-3/4 bg-gray-800">
                <div className="flex flex-col h-full">
                    <div
                        ref={messagesContainerRef}
                        className="flex-1 overflow-y-auto px-4 py-2"
                    >
                        {messages.map((message, index) => (
                            <div key={index} className={`mb-4 flex flex-col ${message.isSentByMe ? "items-end" : "items-start"}`}>
                                <div className="text-gray-500 mb-1">
                                    {new Date().toLocaleString()}
                                </div>
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
                                onClick={() => handleSubmit(true)}
                                className="bg-red-500 text-white p-2 rounded-md mt-2"
                            >
                                Send (Other, Debug)
                            </button>
                        </div>
                    </form>
                </div>
            </div>
        </div>
    );
};

export default ChatLayout;
