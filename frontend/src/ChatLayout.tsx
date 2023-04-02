import React, { useState } from "react";
import ReactMarkdown from "react-markdown";
import hljs from "highlight.js";
import { Light as SyntaxHighlighter } from 'react-syntax-highlighter'
import { dracula } from 'react-syntax-highlighter/dist/esm/styles/hljs'
import "highlight.js/styles/github-dark-dimmed.css";

const ChatLayout = () => {
    const [messages, setMessages] = useState<Array<string>>([]);
    const [inputText, setInputText] = useState("");

    const handleKeyDown = (event: React.KeyboardEvent<HTMLTextAreaElement>) => {
        if (event.key === "Enter" && !event.shiftKey) {
            event.preventDefault();
            handleSubmit();
        }
    };

    const handleSubmit = () => {
        if (inputText.trim() !== "") {
            setMessages([...messages, inputText]);
            setInputText("");
        }
    };

    const conversations = [
        {
            id: 1,
            message: "Hey, how's it going?",
            timestamp: "10:23 AM"
        },
        {
            id: 2,
            message: "Can you send me the report?",
            timestamp: "Yesterday"
        },
        {
            id: 3,
            message: "I'll be late to the meeting",
            timestamp: "Tuesday"
        },
    ];

    const renderMarkdown = (message: string) => {
        const renderer = {
            code: (code: string, language: string) => {
                const highlighted = hljs.highlightAuto(code).value;
                return <pre><code className={`hljs ${language}`}>{highlighted}</code></pre>;
            }
        };
        return (
            <ReactMarkdown
                children={message}
                components={{
                    code({ node, inline, className, children, ...props }) {
                        const match = /language-(\w+)/.exec(className || '')
                        if (inline) {
                            return <code className={className} {...props}>
                                {children}
                            </code>
                        }
                        return match ? (
                            <SyntaxHighlighter
                                children={String(children).replace(/\n$/, '')}
                                style={dracula as any}
                                language={match[1]}
                                PreTag="div"
                                {...props}
                            />
                        ) : (
                            <SyntaxHighlighter
                                children={String(children).replace(/\n$/, '')}
                                style={dracula as any}
                                PreTag="div"
                                {...props}
                            />
                        )
                    }
                }}
                className="bg-gray-700 py-2 px-4 rounded-md text-white inline-block"
            />
        );
    };

    return (
        <div className="flex h-screen overflow-hidden">
            <div className="w-1/3 border-r border-gray-300 border-opacity-50 p-4 bg-gray-900">
                <h2 className="font-bold text-lg mb-4">Conversations</h2>
                {conversations.map((conversation, index) => (
                    <div
                        key={conversation.id}
                        className="flex items-center mb-4 cursor-pointer border-b border-gray-400 py-2"
                    >
                        <div className="w-10 h-10 rounded-full bg-gray-300 mr-2"></div>
                        <div className="flex-1 text-gray-500">
                            <div className="flex justify-between">
                                <p className="text-sm">
                                    {conversation.timestamp}
                                </p>
                            </div>
                            <p className="text-gray-500">{conversation.message}</p>
                        </div>
                    </div>
                ))}
            </div>
            <div className="w-2/3 bg-gray-800 rounded-r-lg">
                <div className="flex flex-col h-full">
                    <div className="flex-1 overflow-y-auto px-4 py-2">
                        {messages.map((message, index) => (
                            <div key={index} className="mb-4">
                                <div className="text-gray-500 mb-1">{new Date().toLocaleString()}</div>
                                <div className="bg-gray-700 py-2 px-4 rounded-md text-white inline-block" style={{ wordWrap: "break-word" }}>
                                    {renderMarkdown(message)}
                                </div>
                            </div>
                        ))}
                    </div>
                    <form onSubmit={(e) => e.preventDefault()} className="flex flex-col h-48 px-4 py-2">
                        <textarea
                            value={inputText}
                            onChange={(event) => setInputText(event.target.value)}
                            onKeyDown={handleKeyDown}
                            className="border border-gray-300 border-opacity-50 p-2 w-full h-32 bg-gray-900 text-white resize-none rounded-md"
                        />
                        <button type="button" onClick={handleSubmit} className="bg-blue-500 text-white p-2 rounded-md mt-2">
                            Send
                        </button>
                    </form>
                </div>
            </div>
        </div>
    );
};

export default ChatLayout;
